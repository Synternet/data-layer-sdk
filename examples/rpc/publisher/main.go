package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/synternet/data-layer-sdk/pkg/codec"
	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/data-layer-sdk/pkg/service"
)

func main() {
	urls := flag.String("urls", os.Getenv("NATS_URL"), "NATS urls")
	prefix := flag.String("prefix", "synternet", "Subject prefix. The subject will be {prefix}.my_publisher")
	creds := flag.String("nats-creds", os.Getenv("PUBLISHER_CREDS"), "NATS credentials file")
	nkey := flag.String("nats-nkey", os.Getenv("PUBLISHER_NKEY"), "NATS NKey string")
	jwt := flag.String("nats-jwt", os.Getenv("PUBLISHER_JWT"), "NATS JWT string")
	verbose := flag.Bool("verbose", false, "Verbose logs")

	flag.Parse()

	conn, err := options.MakeNats("RPC Publisher", *urls, *creds, *nkey, *jwt, "", "", "")
	if err != nil {
		panic(fmt.Errorf("Failed creating NATS connection: %w", err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	opts := []options.Option{
		service.WithContext(ctx),
		service.WithName("rpc-example"),
		service.WithPrefix(*prefix),
		service.WithNats(conn),
		service.WithNKeySeed(*nkey),
		service.WithUserCreds(*creds),
		service.WithCodec(codec.NewProtoJsonCodec()),
		service.WithVerbose(*verbose),
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
