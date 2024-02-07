package service

import (
	"log"
	"strconv"
	"time"
)

func (b *Service) handleTelemetryPing(nmsg Message) {
	now := time.Now().UnixNano()
	var ping TelemetryPing
	_, err := b.Unmarshal(nmsg, &ping)
	if err != nil {
		log.Println("ERR: ping: ", err.Error())
		return
	}

	timestamp, err := strconv.ParseInt(ping.Timestamp, 10, 64)
	if err != nil {
		log.Println("ERR: not a number in timestamp")
		return
	}

	nmsg.Respond(
		&TelemetryPong{
			Nonce:     ping.Nonce,
			Timestamp: strconv.FormatInt(now, 10),
			Owl:       strconv.FormatInt(now-timestamp, 10),
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
