package service

import (
	"strconv"
	"time"
)

//go:generate protoc -I ../../proto --proto_path=../../proto --go_opt=paths=source_relative --go_out=./ ../../proto/telemetry.proto

func (b *Service) handleTelemetryPing(nmsg Message) {
	now := time.Now().UnixNano()
	var ping Ping
	_, err := b.Unmarshal(nmsg, &ping)
	if err != nil {
		b.Logger.Error("ping unmarshal", "err", err)
		return
	}

	nmsg.Respond(
		&Pong{
			Nonce:     ping.Nonce,
			Timestamp: now,
			Owl:       now - ping.Timestamp,
		},
	)
}

func (b *Service) reportTelemetry() {
	nonce := b.nonce.Add(1)
	now := time.Now()
	b.Publish(
		&Telemetry{
			Nonce:  strconv.FormatUint(nonce, 16),
			Status: b.collectStatus(),
		},
		"telemetry",
	)
	b.prevTelemetry = now
}
