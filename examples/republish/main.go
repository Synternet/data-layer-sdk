package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	_ "github.com/syntropynet/data-layer-sdk/pkg/dotenv"

	"github.com/syntropynet/data-layer-sdk/pkg/options"
	"github.com/syntropynet/data-layer-sdk/pkg/service"
)

func main() {
	urls := flag.String("urls", os.Getenv("NATS_URL"), "NATS urls")
	prefix := flag.String("prefix", "syntropy", "Subject prefix. The subject will be {prefix}.my_publisher")
	source := flag.String("source", "syntropy.ethereum.tx", "Source Subject republish messages from.")
	creds := flag.String("nats-creds", os.Getenv("NATS_CREDS"), "NATS credentials file")
	nkey := flag.String("nats-nkey", os.Getenv("NATS_NKEY"), "NATS NKey string")
	jwt := flag.String("nats-jwt", os.Getenv("NATS_JWT"), "NATS JWT string")
	credsPub := flag.String("nats-creds-pub", os.Getenv("NATS_PUB_CREDS"), "NATS credentials file for publishing")
	nkeyPub := flag.String("nats-nkey-pub", os.Getenv("NATS_PUB_NKEY"), "NATS NKey string for publishing")
	jwtPub := flag.String("nats-jwt-pub", os.Getenv("NATS_PUB_JWT"), "NATS JWT string for publishing")
	verbose := flag.Bool("verbose", false, "Verbose logs")

	flag.Parse()

	conn, err := options.MakeNats("Republish Publisher", *urls, *creds, *nkey, *jwt, "", "", "")
	if err != nil {
		panic(fmt.Errorf("Failed creating NATS connection: %w", err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var opts = []options.Option{
		service.WithContext(ctx),
		service.WithName("republish"),
		service.WithPrefix(*prefix),
		service.WithNats(conn),
		service.WithVerbose(*verbose),
		service.WithParam(SourceParam, *source),
	}

	if *credsPub != "" || *nkeyPub != "" || *jwtPub != "" {
		conn, err := options.MakeNats("Republish Publisher", *urls, *credsPub, *nkeyPub, *jwtPub, "", "", "")
		if err != nil {
			panic(fmt.Errorf("Failed creating publishing NATS connection: %w", err))
		}

		// NOTE: Using publisher's side credentials for identity.
		// Better use service.WithPrivateKey() or service.WithPemPrivateKey()
		// to decouple NATS credentials and service credentials.
		opts = append(
			opts,
			service.WithPubNats(conn),
			service.WithUserCreds(*credsPub),
			service.WithNKeySeed(*nkeyPub),
		)
	} else {
		opts = append(
			opts,
			service.WithUserCreds(*creds),
			service.WithNKeySeed(*nkey),
		)
	}

	publisher, err := New(opts...)
	if err != nil {
		panic(fmt.Errorf("Failed creating the republisher: %w", err))
	}

	pubCtx := publisher.Start()
	defer publisher.Close()

	select {
	case <-ctx.Done():
		slog.Info("Shutdown")
	case <-pubCtx.Done():
		slog.Info("Publisher stopped", "cause", context.Cause(pubCtx).Error())
	}
}
