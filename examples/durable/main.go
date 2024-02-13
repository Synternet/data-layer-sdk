package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	_ "github.com/syntropynet/data-layer-sdk/pkg/dotenv"

	"github.com/syntropynet/data-layer-sdk/pkg/options"
	"github.com/syntropynet/data-layer-sdk/pkg/service"
)

func main() {
	urls := flag.String("urls", os.Getenv("NATS_URL"), "NATS urls")
	source := flag.String("source", "syntropy_defi.price.single.ATOM", "Source Subject to stream from.")
	creds := flag.String("nats-creds", os.Getenv("NATS_CREDS"), "NATS credentials file")
	nkey := flag.String("nats-nkey", os.Getenv("NATS_NKEY"), "NATS NKey string")
	jwt := flag.String("nats-jwt", os.Getenv("NATS_JWT"), "NATS JWT string")
	verbose := flag.Bool("verbose", false, "Verbose logs")

	flag.Parse()

	conn, err := options.MakeNats("Streaming consumer", *urls, *creds, *nkey, *jwt, "", "", "")
	if err != nil {
		log.Fatal("Failed creating NATS connection: ", err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var opts = []options.Option{
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
		log.Fatal("Failed creating the consumer: ", err.Error())
	}

	pubCtx := publisher.Start()
	defer publisher.Close()

	select {
	case <-ctx.Done():
		log.Println("Shutdown")
	case <-pubCtx.Done():
		log.Println("Service stopped with cause: ", context.Cause(pubCtx).Error())
	}
}
