package udp

import (
	"net"

	"github.com/crimist/trakx/tracker/udp/udpprotocol"
	"go.uber.org/zap"
)

// TODO: rename this to `fatal`
func (tracker *Tracker) sendError(remote *net.UDPAddr, message string, TransactionID int32) {
	if tracker.stats != nil {
		tracker.stats.ServerErrors.Add(1)
	}

	protoError := udpprotocol.ErrorResponse{
		Action:        udpprotocol.ActionError,
		TransactionID: TransactionID,
		ErrorString:   []byte("internal error"),
	}

	data, err := protoError.Marshall()
	if err != nil {
		zap.L().Error("failed to marshal error packet", zap.Error(err))
		tracker.socket.WriteToUDP([]byte("catastrophic failure"), remote)
	} else {
		tracker.socket.WriteToUDP(data, remote)
	}
}
