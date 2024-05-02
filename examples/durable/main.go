package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	_ "github.com/synternet/data-layer-sdk/pkg/dotenv"

	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/data-layer-sdk/pkg/service"
)

func main() {
	urls := flag.String("urls", os.Getenv("NATS_URL"), "NATS urls")
	source := flag.String("source", "synternet.price.single.ATOM", "Source Subject to stream from.")
	creds := flag.String("nats-creds", os.Getenv("NATS_CREDS"), "NATS credentials file")
	nkey := flag.String("nats-nkey", os.Getenv("NATS_NKEY"), "NATS NKey string")
	jwt := flag.String("nats-jwt", os.Getenv("NATS_JWT"), "NATS JWT string")
	verbose := flag.Bool("verbose", false, "Verbose logs")

	flag.Parse()

	conn, err := options.MakeNats("Streaming consumer", *urls, *creds, *nkey, *jwt, "", "", "")
	if err != nil {
		panic(fmt.Errorf("Failed creating NATS connection: %w", err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	opts := []options.Option{
		service.WithContext(ctx),
		service.WithName("consume"),
		service.WithNats(conn),
		service.WithVerbose(*verbose),
		service.WithParam(SourceParam, *source),
	}

	opts = append(
		opts,
		service.WithUserCreds(*creds),
		service.WithNKeySeed(*nkey),
	)

	publisher, err := New(opts...)
	if err != nil {
		panic(fmt.Errorf("Failed creating the consumer: %w", err))
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
