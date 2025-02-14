package service

import (
	"strings"

	"google.golang.org/protobuf/proto"
)

// Publish will sign the message and publish it to a subject constructed from "{prefix}.{name}.{suffixes}".
// Publish will use PubNats connection.
func (b *Service) Publish(msg proto.Message, suffixes ...string) error {
	return b.PublishTo(msg, b.Subject(suffixes...))
}

// PublishBuf is the same as Publish, but will publish the raw bytes.
func (b *Service) PublishBuf(buf []byte, suffixes ...string) error {
	return b.PublishBufTo(buf, b.Subject(suffixes...))
}

// PublishTo will sign the message and publish it to a specific subject constructed from subject tokens.
// PublishTo will use PubNats connection.
func (b *Service) PublishTo(msg proto.Message, tokens ...string) error {
	payload, err := b.Codec.Encode(nil, msg)
	if err != nil {
		return err
	}
	return b.PublishBufTo(payload, tokens...)
}

// PublishBufTo is the same as PublishTo, but for raw bytes.
func (b *Service) PublishBufTo(buf []byte, tokens ...string) error {
	if b.PubNats == nil {
		return ErrPubConnection
	}
	msg, err := b.makeMsg(buf, "", strings.Join(tokens, "."))
	if err != nil {
		return err
	}

	select {
	case <-b.Context.Done():
		b.Logger.Info("PublishBufTo cancelled", "err", b.Context.Err(), "queue_size", len(b.publishCh))
		return b.Context.Err()
	case b.publishCh <- msg:
	}
	return nil
}
