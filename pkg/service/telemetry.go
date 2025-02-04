package service

import (
	"strconv"
	"time"

	"github.com/synternet/data-layer-sdk/types/telemetry"
	"google.golang.org/protobuf/proto"
)

func (b *Service) handleTelemetryPing(nmsg Message) (proto.Message, error) {
	now := time.Now().UnixNano()
	var ping telemetry.Ping
	_, err := b.Unmarshal(nmsg, &ping)
	if err != nil {
		b.Logger.Error("ping unmarshal", "err", err)
		return nil, err
	}

	return &telemetry.Pong{
		Nonce:     ping.Nonce,
		Timestamp: now,
		Owl:       now - ping.Timestamp,
	}, nil
}

func (b *Service) reportTelemetry() {
	nonce := b.nonce.Add(1)
	now := time.Now()
	b.Publish(
		&telemetry.Telemetry{
			Nonce:  strconv.FormatUint(nonce, 16),
			Status: b.collectStatus(),
		},
		"telemetry",
	)
	b.prevTelemetry = now
}
