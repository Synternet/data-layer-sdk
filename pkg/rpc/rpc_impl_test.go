package rpc_test

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/synternet/data-layer-sdk/pkg/rpc"
	rpctypes "github.com/synternet/data-layer-sdk/types/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ rpctypes.TestServiceServer = (*Test)(nil)

type Test struct {
}

// Test implements rpc.TestServiceServer.
func (t *Test) Test(ctx context.Context, r *rpctypes.TestRequest) (*rpctypes.TestResponse, error) {
	subject, _ := rpc.GetSubject(ctx)
	headers, _ := rpc.GetHeaders(ctx)
	slog.Info("Test", "r", r, "subject", subject, "header", headers)
	if r.A < 0 {
		return nil, fmt.Errorf("negative value: %v", r)
	}

	hdr := make(map[string]string)
	for k, v := range headers {
		hdr[k] = strings.Join(v, ",")
	}

	return &rpctypes.TestResponse{
		Ab:      r.A + r.B,
		Subject: subject.String(),
		Header:  hdr,
	}, nil
}

// TestStream implements rpc.TestServiceServer.
func (t *Test) TestStream(r *rpctypes.TestRequest, srv rpctypes.TestService_TestStreamServer) error {
	var i float32 = 1
	slog.Info("TestStream", "r", r)
	for {
		if err := srv.Send(&rpctypes.TestResponse{
			Ab: (r.A + r.B) * i,
		}); err != nil {
			return err
		}
		i += 1
	}
}

// TestStreamOnly implements rpc.TestServiceServer.
func (t *Test) TestStreamOnly(_ *emptypb.Empty, srv rpctypes.TestService_TestStreamOnlyServer) error {
	var i float32 = 1
	slog.Info("TestStreamOnly")
	for {
		if err := srv.Send(&rpctypes.TestResponse{
			Ab: i,
		}); err != nil {
			return err
		}
		i += 1
	}
}

// TestStreamBidirectional implements rpc.TestServiceServer.
func (t *Test) TestStreamBidirectional(srv rpctypes.TestService_TestStreamBidirectionalServer) error {
	for {
		var msg rpctypes.TestRequest
		if err := srv.RecvMsg(&msg); err != nil {
			return err
		}

		if err := srv.Send(&rpctypes.TestResponse{
			Ab: msg.A + msg.B,
		}); err != nil {
			return err
		}
	}
}
