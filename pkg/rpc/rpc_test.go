package rpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/synternet/data-layer-sdk/pkg/rpc"
	service "github.com/synternet/data-layer-sdk/pkg/service"
	rpctypes "github.com/synternet/data-layer-sdk/x/synternet/rpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
)

func makeServer(t *testing.T, ctx context.Context) rpc.Publisher {
	ctx, _ = context.WithTimeout(ctx, time.Second*30)
	grp, ctx := errgroup.WithContext(ctx)
	pub := NewPublisher(ctx, t, "test_prefix")
	srv := rpc.NewServiceRegistrar(grp, pub)

	rpctypes.RegisterTestServiceServer(srv, &Test{t: t})
	go srv.Start(ctx)
	time.Sleep(time.Millisecond * 10) // Allow go routines to run before we start the client

	return pub
}

// TestPubSub tests the mock publisher implementation
func TestPubSub(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	srv := makeServer(t, ctx)

	received := make(chan struct{})
	srv.SubscribeTo(func(msg service.Message) {
		received <- struct{}{}
	}, "test")

	srv.PublishTo(&rpctypes.TestRequest{}, "test")

	select {
	case <-ctx.Done():
		require.False(t, true)
	case <-received:
	}
	time.Sleep(time.Millisecond * 10)
}

// TestPubRpcSub tests the mock publisher implementation
func TestPubRpcSub(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	srv := makeServer(t, ctx)

	received := make(chan struct{})
	srv.SubscribeTo(func(msg service.Message) {
		require.Equal(t, "reply.To", msg.Reply())
		received <- struct{}{}
	}, "test")

	srv.PublishToRpc(&rpctypes.TestRequest{}, "reply.To", "test")

	select {
	case <-received:
	case <-ctx.Done():
		require.False(t, true)
	}
	time.Sleep(time.Millisecond * 10)
}

func TestRequestReply(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sub := makeServer(t, ctx)
	clt := rpc.NewClientConn(ctx, sub, "test_prefix")
	client := rpctypes.NewTestServiceClient(clt)

	ctx1, cancel1 := context.WithTimeout(ctx, time.Millisecond*50)
	defer cancel1()
	res, err := client.Test(ctx1, &rpctypes.TestRequest{A: 123, B: 321})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, float32(123.0+321.0), res.Ab)
	assert.Equal(t, "", res.Error)
	assert.Equal(t, "test_prefix.override.test.override.test.method", res.Subject)
	assert.Equal(t, map[string]string{"identity": "some,identity"}, res.Header)

	time.Sleep(time.Millisecond * 10)
}

func TestRequestReplyWithError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	sub := makeServer(t, ctx)
	clt := rpc.NewClientConn(ctx, sub, "test_prefix")
	client := rpctypes.NewTestServiceClient(clt)

	ctx1, cancel1 := context.WithTimeout(ctx, time.Millisecond*50)
	defer cancel1()
	res, err := client.Test(ctx1, &rpctypes.TestRequest{A: -123, B: -321})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "synternet.rpc.TestService.Test: negative value: a:-123 b:-321", res.Error)
	assert.Equal(t, "", res.Subject, "")
	assert.Equal(t, map[string]string(nil), res.Header)

	time.Sleep(time.Millisecond * 10)
}

func TestRequestStream(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sub := makeServer(t, ctx)
	clt := rpc.NewClientConn(ctx, sub, "test_prefix")
	client := rpctypes.NewTestServiceClient(clt)

	ctx1, cancel1 := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel1()
	str, err := client.TestStream(ctx1, &rpctypes.TestRequest{A: 123, B: 321})
	require.NoError(t, err)
	require.NotNil(t, str)

	for i := 0; i < 10; i++ {
		var msg rpctypes.TestResponse
		err := str.RecvMsg(&msg)
		require.NoError(t, err)
		t.Log("RecvMsg", "i=", i, "msg=", &msg)
		assert.Equal(t, float32(123.0+321.0)*float32(i), msg.Ab)
	}

	cancel()
	cancel1()
	time.Sleep(time.Millisecond * 10)
}

func TestRequestStreamOnly(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sub := makeServer(t, ctx)
	clt := rpc.NewClientConn(ctx, sub, "test_prefix")
	client := rpctypes.NewTestServiceClient(clt)

	ctx1, cancel1 := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel1()
	str, err := client.TestStreamOnly(ctx1, &emptypb.Empty{})
	require.NoError(t, err)
	require.NotNil(t, str)

	for i := 0; i < 10; i++ {
		var msg rpctypes.TestResponse
		err := str.RecvMsg(&msg)
		require.NoError(t, err)
		t.Log("RecvMsg", "i=", i, "msg=", &msg)
		assert.Equal(t, float32(i), msg.Ab)
	}

	cancel()
	time.Sleep(time.Millisecond * 10)
}
