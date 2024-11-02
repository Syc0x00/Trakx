package udp

import (
	"net"
	"net/netip"

	"github.com/crimist/trakx/tracker/udp/udpprotocol"
	"go.uber.org/zap"
)

const maximumScrapeHashes = 74

func (tracker *Tracker) scrape(udpAddr *net.UDPAddr, addrPort netip.AddrPort, transactionID int32, data []byte) {
	if tracker.stats != nil {
		tracker.stats.Scrapes.Add(1)
	}

	scrape, err := udpprotocol.NewScrapeRequest(data)
	if err != nil {
		tracker.fatal(udpAddr, []byte("failed to parse scrape"), transactionID)
		zap.L().Info("failed to parse clients scrape packet", zap.Binary("packet", data), zap.Error(err), zap.Any("remote", addrPort))
		return
	}

	if len(scrape.InfoHashes) > maximumScrapeHashes {
		tracker.fatal(udpAddr, []byte("exceeded 74 hashes"), scrape.TransactionID)
		zap.L().Debug("client sent over sized scrape request (> 74 hashes)", zap.Int("hashes", len(scrape.InfoHashes)), zap.Any("scrape", scrape), zap.Any("remote", udpAddr))
		return
	}

	resp := udpprotocol.ScrapeResponse{
		Action:        udpprotocol.ActionScrape,
		TransactionID: scrape.TransactionID,
	}

	for _, hash := range scrape.InfoHashes {
		if len(hash) != 20 {
			tracker.fatal(udpAddr, append([]byte("missized hash "), hash[0:7]...), scrape.TransactionID)
			zap.L().Debug("client sent scrape with missized hash", zap.Any("hash", hash), zap.Any("scrape", scrape), zap.Any("remote", udpAddr))
			return
		}

		seeds, leeches := tracker.peerDB.TorrentStats(hash)
		info := udpprotocol.ScrapeResponseInfo{
			Complete:   int32(seeds),
			Incomplete: int32(leeches),
			Downloaded: -1,
		}
		resp.Info = append(resp.Info, info)
	}

	marshalledResp, err := resp.Marshal()
	if err != nil {
		tracker.fatal(udpAddr, []byte("failed to marshall scrape response"), scrape.TransactionID)
		zap.L().Error("failed to marshall scrape response", zap.Error(err), zap.Any("scrape", scrape), zap.Any("remote", udpAddr))
		return
	}

	tracker.socket.WriteToUDP(marshalledResp, udpAddr)
}
