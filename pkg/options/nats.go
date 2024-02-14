package options

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

// MakeNats is a convenience function to create NATS connection from various parameters.
func MakeNats(name, urls, userCreds, nkey, jwt, caCert, clientCert, clientKey string) (*nats.Conn, error) {
	if urls == "" {
		slog.Info("Will use NATS stub")
		return nil, nil
	}

	// Connect Options.
	opts := []nats.Option{}

	// Use TLS client authentication
	if clientCert != "" && clientKey != "" {
		opts = append(opts, nats.ClientCert(clientCert, clientKey))
	}

	// Use specific CA certificate
	if caCert != "" {
		opts = append(opts, nats.RootCAs(caCert))
	}

	// Use UserCredentials
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}

	// Use Nkey + JWT authentication.
	if nkey != "" && jwt != "" {
		opts = append(opts, nats.UserJWTAndSeed(jwt, nkey))
	}

	// Connect to NATS
	opts = append(opts, nats.Name(name))

	var err error
	natsConnection, err := nats.Connect(urls, opts...)
	if err != nil {
		return nil, fmt.Errorf("connection to NATS failed: %w", err)
	}
	return natsConnection, nil
}
