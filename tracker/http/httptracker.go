package http

import (
	"fmt"
	"net"

	"github.com/crimist/trakx/config"
	"github.com/crimist/trakx/tracker/storage"
	"github.com/pkg/errors"
)

const (
	httpRequestMax = 2600 // enough for scrapes up to 40 info_hashes
)

type Tracker struct {
	peerdb   storage.Database
	workers  workers
	shutdown chan struct{}
}

// Init sets up the HTTPTracker.
func (tracker *Tracker) Init(peerdb storage.Database) {
	tracker.peerdb = peerdb
	tracker.shutdown = make(chan struct{})
}

// Serve begins listening and serving clients.
func (tracker *Tracker) Serve() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port))
	if err != nil {
		return errors.Wrap(err, "Failed to open TCP listen socket")
	}

	cache, err := config.GenerateEmbeddedCache()
	if err != nil {
		return errors.Wrap(err, "failed to generate embedded cache")
	}

	tracker.workers = workers{
		tracker:   tracker,
		listener:  ln,
		fileCache: cache,
	}

	tracker.workers.startWorkers(config.Config.HTTP.Threads)

	<-tracker.shutdown
	if err := ln.Close(); err != nil {
		return errors.Wrap(err, "Failed to close tcp listen socket")
	}

	return nil
}

// Shutdown stops the HTTP tracker server by closing the socket.
func (t *Tracker) Shutdown() {
	if t == nil || t.shutdown == nil {
		return
	}
	var die struct{}
	t.shutdown <- die
}
