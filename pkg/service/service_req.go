package service

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
)

// RequestFrom requests a reply or a stream from a subject using ReqNats connection. The subject will be constructed from tokens.
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

// RequestBufFrom requests a reply or a stream from a subject using ReqNats connection.
// This a synchronous operation that does not involve publisher queue.
func (b *Service) RequestBufFrom(ctx context.Context, buf []byte, tokens ...string) (Message, error) {
	if b.ReqNats == nil {
		return nil, fmt.Errorf("request NATS connection is nil")
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
