package udp

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/Syc0x00/Trakx/tracker/shared"
	"go.uber.org/zap"
)

type udpTracker struct {
	conn    *net.UDPConn
	avgResp time.Time
}

// Run runs the UDP tracker
func Run(trimInterval time.Duration) {
	u := udpTracker{}
	loadConnDB()
	rand.Seed(time.Now().UnixNano() * time.Now().Unix())

	go shared.RunOn(trimInterval, connDB.trim)
	u.listen()
}

func (u *udpTracker) listen() {
	var err error

	u.conn, err = net.ListenUDP("udp4", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: shared.Config.UDPPort, Zone: ""})
	if err != nil {
		panic(err)
	}
	defer u.conn.Close()

	var pool sync.Pool
	pool.New = func() interface{} {
		return make([]byte, 1496, 1496)
	}

	for {
		b := pool.Get().([]byte)
		len, remote, err := u.conn.ReadFromUDP(b)
		if err != nil {
			shared.Logger.Error("ReadFromUDP()", zap.Error(err))
			pool.Put(b)
			continue
		}
		go func() {
			u.process(b[:len], remote)

			// optimized zero
			b = b[:cap(b)]
			for i := range b {
				b[i] = 0
			}
			pool.Put(b)
		}()
	}
}

func (u *udpTracker) process(data []byte, remote *net.UDPAddr) {
	base := connect{}
	var addr [4]byte
	ip := remote.IP.To4()
	copy(addr[:], ip)

	if ip == nil {
		u.conn.WriteToUDP(newClientError("IPv6?", base.TransactionID, zap.String("ip", remote.IP.String())), remote)
		return
	}

	err := base.unmarshall(data)
	if err != nil {
		u.conn.WriteToUDP(newServerError("base.unmarshall()", err, base.TransactionID), remote)
	}

	if base.Action == 0 {
		u.connect(&base, remote, addr)
		return
	}

	if dbID, ok := connDB.check(base.ConnectionID, addr); !ok && shared.UDPCheckConnID {
		u.conn.WriteToUDP(newClientError("bad connid", base.TransactionID, zap.Int64("dbID", dbID), zap.Int64("clientID", base.ConnectionID), zap.Int32("action", base.Action), zap.Any("addr", addr)), remote)
		return
	}

	switch base.Action {
	case 1:
		announce := announce{}
		if err := announce.unmarshall(data); err != nil {
			u.conn.WriteToUDP(newServerError("announce.unmarshall()", err, base.TransactionID), remote)
			return
		}
		u.announce(&announce, remote, addr)

	case 2:
		scrape := scrape{}
		if err := scrape.unmarshall(data); err != nil {
			u.conn.WriteToUDP(newServerError("scrape.unmarshall()", err, base.TransactionID), remote)
			return
		}
		u.scrape(&scrape, remote)
	default:
		u.conn.WriteToUDP(newClientError("bad action", base.TransactionID, zap.Int32("action", base.Action), zap.Any("addr", addr)), remote)
	}
}
