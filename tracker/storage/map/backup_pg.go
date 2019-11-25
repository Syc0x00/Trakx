package gomap

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	"github.com/crimist/trakx/tracker/storage"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	// Maximum retention for entries. Rows older than this will be removed
	// "off" to disable
	maxdate = "7 days"

	// Maximum number of rows. Rows exceeding this will be removed by timestamp
	// -1 for unlimited
	maxrows = "10"
)

type PgBackup struct {
	pg *sql.DB
	db *Memory
}

func (bck *PgBackup) Init(db storage.Database) error {
	var err error

	bck.db = db.(*Memory)
	if bck.db == nil {
		panic("db nil on backup init")
	}

	bck.pg, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	err = bck.pg.Ping()
	if err != nil {
		bck.db.logger.Error("postgres ping() failed", zap.Error(err))
		return err
	}

	_, err = bck.pg.Exec("CREATE TABLE IF NOT EXISTS trakx (ts TIMESTAMP DEFAULT now(), bytes BYTEA)")
	if err != nil {
		bck.db.logger.Error("postgres table create failed", zap.Error(err))
		return err
	}

	return nil
}

func (bck PgBackup) save() error {
	data, err := bck.db.encode()
	if err != nil {
		bck.db.logger.Error("Failed to encode", zap.Error(err))
		return err
	}

	_, err = bck.pg.Query("INSERT INTO trakx(bytes) VALUES($1)", data)
	if err != nil {
		bck.db.logger.Error("postgres insert failed", zap.Error(err))
		return errors.New("postgres insert failed")
	}

	rm, err := bck.trim()
	if err != nil {
		bck.db.logger.Error("failed to trim backups", zap.Error(err))
	} else {
		bck.db.logger.Info("Deleted expired postgres records", zap.Int64("deleted", rm))
	}

	bck.db.logger.Info("Deleted expired postgres records", zap.Int64("deleted", rm))

	return nil
}

func (bck PgBackup) Save() error {
	return bck.save()
}

func (bck PgBackup) load() error {
	var data []byte

	err := bck.pg.QueryRow("SELECT bytes FROM trakx ORDER BY ts DESC LIMIT 1").Scan(&data)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// empty postgres table
			bck.db.logger.Info("No stored database found")
			bck.db.make()
			return nil
		}
		return errors.New("postgres SELECT query failed: " + err.Error())
	}

	bck.db.logger.Info("Loading stored database", zap.Int("size", len(data)))
	if err := bck.db.decode(data); err != nil {
		bck.db.logger.Error("Error decoding stored database", zap.Error(err))
		return err
	}
	bck.db.logger.Info("Loaded stored database")

	return nil
}

func (bck PgBackup) Load() error {
	return bck.load()
}

func (bck PgBackup) trim() (int64, error) {
	var trimmed int64

	if maxdate != "off" {
		result, err := bck.pg.Exec("DELETE FROM trakx WHERE ts < NOW() - INTERVAL '" + maxdate + "'")
		if err != nil {
			return -1, err
		}

		trimmed, err = result.RowsAffected()
		if err != nil {
			return -1, err
		}
	}

	if maxrows != "-1" {
		result, err := bck.pg.Exec("DELETE FROM trakx WHERE ctid IN (SELECT ctid FROM trakx ORDER BY ctid DESC OFFSET " + maxrows + ")")
		if err != nil {
			return -1, err
		}
		removedRows, err := result.RowsAffected()
		if err != nil {
			return -1, err
		}

		trimmed += removedRows
	}

	return trimmed, nil
}
