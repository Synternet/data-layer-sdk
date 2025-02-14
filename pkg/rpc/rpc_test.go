package rpc_test

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/synternet/data-layer-sdk/pkg/rpc"
	rpctypes "github.com/synternet/data-layer-sdk/types/rpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	once      sync.Once
	publisher rpc.Publisher
)

func makeServer(t *testing.T, ctx context.Context) rpc.Publisher {

	once.Do(func() {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
		grp, ctx := errgroup.WithContext(ctx)
		pub := NewPublisher(ctx, "test_prefix")
		srv := rpc.NewServiceRegistrar(grp, pub)

		rpctypes.RegisterTestServiceServer(srv, &Test{})
		go srv.Start(ctx)
		time.Sleep(time.Millisecond * 10) // Allow go routines to run before we start the client
		publisher = pub
	})

	return publisher
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
}

func TestRequestReplyWithError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sub := makeServer(t, ctx)
	clt := rpc.NewClientConn(ctx, sub, "test_prefix")
	client := rpctypes.NewTestServiceClient(clt)

	ctx1, cancel1 := context.WithTimeout(ctx, time.Millisecond*50)
	defer cancel1()
	res, err := client.Test(ctx1, &rpctypes.TestRequest{A: -123, B: -321})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "types.rpc.TestService.Test: negative value: a:-123 b:-321", res.Error)
	assert.Equal(t, "", res.Subject, "")
	assert.Equal(t, map[string]string(nil), res.Header)
}

func _TestRequestStream(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sub := makeServer(t, ctx)
	clt := rpc.NewClientConn(ctx, sub, "test_prefix")
	client := rpctypes.NewTestServiceClient(clt)

	ctx1, cancel1 := context.WithTimeout(ctx, time.Millisecond*50)
	defer cancel1()
	str, err := client.TestStream(ctx1, &rpctypes.TestRequest{A: 123, B: 321})
	require.NoError(t, err)
	require.NotNil(t, str)

	for i := 0; i < 10; i++ {
		var msg rpctypes.TestResponse
		err := str.RecvMsg(&msg)
		require.NoError(t, err)
		slog.Info("RecvMsg", "i", i, "msg", &msg)
		assert.Equal(t, msg.Ab, float32(123.0+321.0)*float32(i))
	}
}

func _TestRequestStreamOnly(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sub := makeServer(t, ctx)
	clt := rpc.NewClientConn(ctx, sub, "test_prefix")
	client := rpctypes.NewTestServiceClient(clt)

	ctx1, cancel1 := context.WithTimeout(ctx, time.Millisecond*50)
	defer cancel1()
	str, err := client.TestStreamOnly(ctx1, &emptypb.Empty{})
	require.NoError(t, err)
	require.NotNil(t, str)

	for i := 0; i < 10; i++ {
		var msg rpctypes.TestResponse
		err := str.RecvMsg(&msg)
		require.NoError(t, err)
		slog.Info("RecvMsg", "i", i, "msg", &msg)
		assert.Equal(t, msg.Ab, float32(123.0+321.0)*float32(i))
	}
}
