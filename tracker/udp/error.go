package udp

import (
	"net"

	"github.com/crimist/trakx/tracker/udp/udpprotocol"
	"go.uber.org/zap"
)

func (tracker *Tracker) fatal(remote *net.UDPAddr, message []byte, TransactionID int32) {
	if tracker.stats != nil {
		// TODO: this isn't right
		tracker.stats.ServerErrors.Add(1)
	}

	protoError := udpprotocol.ErrorResponse{
		Action:        udpprotocol.ActionError,
		TransactionID: TransactionID,
		ErrorString:   message,
	}

	data, err := protoError.Marshal()
	if err != nil {
		zap.L().Error("failed to marshal error packet", zap.Error(err))
		tracker.socket.WriteToUDP([]byte("catastrophic failure"), remote)
	} else {
		tracker.socket.WriteToUDP(data, remote)
	}
}
