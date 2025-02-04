package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/types/telemetry"
	"google.golang.org/protobuf/proto"
)

// Serve is a convenience method to serve a service subject. It acts the same as Subscribe, but takes `ServiceHandler` instead, and will respond
// either with Error type or response from the handler. Serve will use ReqNats connection.
func (b *Service) Serve(handler ServiceHandler, suffixes ...string) (*nats.Subscription, error) {
	return b.subscribeTo(
		b.ReqNats,
		func(msg Message) {
			resp, err := handler(msg)
			if err != nil {
				b.Logger.Error("service handler failed", "err", err, "suffixes", suffixes)
				err1 := msg.Respond(&telemetry.Error{Error: err.Error()})
				if err1 != nil {
					b.Logger.Error("service handler failed during error", "err", err, "err1", err1, "suffixes", suffixes)
				}
				return
			}
			err = msg.Respond(resp)
			if err != nil {
				b.Logger.Error("service handler failed", "err", err, "suffixes", suffixes)
			}
		},
		b.Subject(suffixes...),
	)
}

// RequestFrom requests a reply from a subject using ReqNats connection. The subject will be constructed from tokens.
// This a synchronous operation that does not involve publisher queue.
func (b *Service) RequestFrom(ctx context.Context, msg proto.Message, resp proto.Message, tokens ...string) (Message, error) {
	payload, err := b.Codec.Encode(nil, msg)
	if err != nil {
		return nil, err
	}
	response, err := b.RequestBufFrom(ctx, payload, tokens...)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		_, err := b.Unmarshal(response, resp)
		if err != nil {
			return response, fmt.Errorf("unmarshal failed: %w", err)
		}
	}

	return response, err
}

// RequestBufFrom requests a reply from a subject using ReqNats connection.
// This a synchronous operation that does not involve publisher queue.
func (b *Service) RequestBufFrom(ctx context.Context, buf []byte, tokens ...string) (Message, error) {
	if b.ReqNats == nil {
		return nil, ErrReqConnection
	}

	msg, err := b.makeMsg(buf, strings.Join(tokens, "."))
	if err != nil {
		return nil, err
	}

	ret, err := b.ReqNats.RequestMsgWithContext(ctx, msg)
	if err != nil {
		return nil, err
	}

	return wrapMessage(b.Codec, &b.msg_out_counter, &b.bytes_out_counter, b.makeMsg, ret), nil
}

// Respond will respond to a message sent as a request.
// This is a helper function and Message.Respond should be used instead.
func (b *Service) Respond(nmsg Message, msg proto.Message) error {
	payload, err := b.Codec.Encode(nil, msg)
	if err != nil {
		return err
	}
	return b.RespondBuf(nmsg, payload)
}

// RespondBuf is the same as Respond, but will respond with raw bytes.
func (b *Service) RespondBuf(msg Message, buf []byte) error {
	reply, err := b.makeMsg(buf, "")
	if err != nil {
		return err
	}

	return msg.Message().RespondMsg(reply)
}
