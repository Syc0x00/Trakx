package tracker

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/crimist/trakx/tracker/storage"
	"github.com/crimist/trakx/tracker/udp"

	"go.uber.org/zap"
)

var (
	SigStop     = os.Interrupt
	exitSuccess = 0
)

func sigHandler(peerdb storage.Database, udptracker *udp.UDPTracker) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR1)

	for {
		sig := <-c

		switch sig {
		case os.Interrupt, os.Kill, syscall.SIGTERM:
			logger.Info("Got exit signal", zap.Any("sig", sig))

			peerdb.Backup().Save()
			if udptracker != nil {
				udptracker.WriteConns()
			}

			os.Exit(exitSuccess)
		default:
			logger.Info("Got unknown sig", zap.Any("Signal", sig))
		}
		// os.Exit(128 + int(sig.(syscall.Signal)))
	}
}
