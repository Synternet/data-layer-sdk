package rpc_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/synternet/data-layer-sdk/pkg/rpc"
	rpctypes "github.com/synternet/data-layer-sdk/x/synternet/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ rpctypes.TestServiceServer = (*Test)(nil)

type Test struct {
	t *testing.T
}

// Test implements rpc.TestServiceServer.
func (t *Test) Test(ctx context.Context, r *rpctypes.TestRequest) (*rpctypes.TestResponse, error) {
	subject, _ := rpc.GetSubject(ctx)
	headers, _ := rpc.GetHeaders(ctx)
	t.t.Log("Test", "r=", r, "subject=", subject, "header=", headers)
	if r.A < 0 {
		return nil, fmt.Errorf("negative value: a:%v b:%v", r.A, r.B)
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

// TestVars implements rpc.TestServiceServer.
func (t *Test) TestVars(ctx context.Context, r *rpctypes.TestRequest) (*rpctypes.TestResponse, error) {
	subject, _ := rpc.GetSubject(ctx)
	headers, _ := rpc.GetHeaders(ctx)
	t.t.Log("Test", "r=", r, "subject=", subject, "header=", headers)
	if r.A < 0 {
		return nil, fmt.Errorf("negative value: a:%v b:%v", r.A, r.B)
	}

	assert.Equal(t.t, "123456", subject.Tokens()[len(subject.Tokens())-1])

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
	var i float32
	t.t.Log("TestStream", "r=", r)
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
	var i float32
	t.t.Log("TestStreamOnly")
	for {
		if err := srv.Send(&rpctypes.TestResponse{
			Ab: i,
		}); err != nil {
			t.t.Log("TestStreamOnly", "err=", err)
			return err
		}
		i += 1
	}
}

// TestStreamBidirectional implements rpc.TestServiceServer.
func (t *Test) TestStreamBidirectional(srv rpctypes.TestService_TestStreamBidirectionalServer) error {
	t.t.Log("TestStreamBidirectional")
	for {
		var msg rpctypes.TestRequest
		if err := srv.RecvMsg(&msg); err != nil {
			t.t.Log("TestStreamBidirectional recv", "err=", err)
			return err
		}

		if err := srv.Send(&rpctypes.TestResponse{
			Ab: msg.A + msg.B,
		}); err != nil {
			t.t.Log("TestStreamBidirectional send", "err=", err)
			return err
		}
	}
}
