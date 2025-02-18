package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	v1 "github.com/synternet/data-layer-sdk/examples/rpc/types/example/v1"
	"github.com/synternet/data-layer-sdk/pkg/codec"
	_ "github.com/synternet/data-layer-sdk/pkg/dotenv"
	"github.com/synternet/data-layer-sdk/pkg/rpc"

	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/data-layer-sdk/pkg/service"
)

func main() {
	urls := flag.String("urls", os.Getenv("NATS_URL"), "NATS urls")
	source := flag.String("source", "synternet.rpc-example", "Source Subject to stream from.")
	creds := flag.String("nats-creds", os.Getenv("SUBSCRIBER_CREDS"), "NATS credentials file")
	nkey := flag.String("nats-nkey", os.Getenv("SUBSCRIBER_NKEY"), "NATS NKey string")
	jwt := flag.String("nats-jwt", os.Getenv("SUBSCRIBER_JWT"), "NATS JWT string")
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
		service.WithName("client"),
		service.WithNats(conn),
		service.WithVerbose(*verbose),
		service.WithCodec(codec.NewProtoJsonCodec()),
	}

	opts = append(
		opts,
		service.WithUserCreds(*creds),
		service.WithNKeySeed(*nkey),
	)

	var subscriber service.Service

	subscriber.Configure(opts...)
	if err != nil {
		panic(fmt.Errorf("Failed creating the consumer: %w", err))
	}
	users := v1.NewUserServiceClient(rpc.NewClientConn(ctx, &subscriber, *source))

	subscriber.SubscribeTo(func(msg service.Message) {
		var user v1.RegistrationsResponse
		_, err := subscriber.Unmarshal(msg, &user)
		if err != nil {
			panic(err)
		}
		slog.Info("New user just registered", "user", &user)
	}, *source, "users.registrations")

	pubCtx := subscriber.Start()
	defer subscriber.Close()

	_, err = users.Add(ctx, &v1.AddRequest{Name: "John", Surname: "Smith"})
	if err != nil {
		panic(fmt.Errorf("Failed adding John: %w", err))
	}
	_, err = users.Add(ctx, &v1.AddRequest{Name: "Mr.", Surname: "Andersen"})
	if err != nil {
		panic(fmt.Errorf("Failed adding Andersen: %w", err))
	}

	user, err := users.Get(ctx, &v1.GetRequest{Id: 1})
	if err != nil {
		panic(fmt.Errorf("Failed adding Andersen: %w", err))
	}
	slog.Info("Retrieved user", "user", user)

	_, err = users.Add(ctx, &v1.AddRequest{Name: "Trinity", Surname: "Tiff"})
	if err != nil {
		panic(fmt.Errorf("Failed adding Andersen: %w", err))
	}

	select {
	case <-ctx.Done():
		slog.Info("Shutdown")
	case <-pubCtx.Done():
		slog.Info("subscriber stopped", "cause", context.Cause(pubCtx).Error())
	}
}
