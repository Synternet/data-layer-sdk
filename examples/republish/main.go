package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/syntropynet/data-layer-sdk/pkg/options"
	"github.com/syntropynet/data-layer-sdk/pkg/service"
)

func main() {
	urls := flag.String("urls", "", "NATS urls")
	prefix := flag.String("prefix", "syntropy", "Subject prefix. The subject will be {prefix}.my_publisher")
	source := flag.String("source", "syntropy.ethereum.tx", "Source Subject republish messages from.")
	creds := flag.String("nats-creds", "", "NATS credentials file")
	nkey := flag.String("nats-nkey", "", "NATS NKey string")
	jwt := flag.String("nats-jwt", "", "NATS JWT string")
	credsPub := flag.String("nats-creds-pub", "", "NATS credentials file for publishing")
	nkeyPub := flag.String("nats-nkey-pub", "", "NATS NKey string for publishing")
	jwtPub := flag.String("nats-jwt-pub", "", "NATS JWT string for publishing")
	verbose := flag.Bool("verbose", false, "Verbose logs")

	flag.Parse()

	conn, err := options.MakeNats("Republish Publisher", *urls, *creds, *nkey, *jwt, "", "", "")
	if err != nil {
		log.Fatal("Failed creating NATS connection: ", err.Error())
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
			log.Fatal("Failed creating publishing NATS connection: ", err.Error())
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
		log.Fatal("Failed creating the republisher: ", err.Error())
	}

	pubCtx := publisher.Start()
	defer publisher.Close()

	select {
	case <-ctx.Done():
		log.Println("Shutdown")
	case <-pubCtx.Done():
		log.Println("Publisher stopped with cause: ", context.Cause(pubCtx).Error())
	}
}
