// codec package implements various codecs to be used with the publisher to encode and decode messages.
package codec

import (
	"google.golang.org/protobuf/proto"
)

type PbCodec struct {
}

func NewPbCodec() *PbCodec {
	return &PbCodec{}
}

func (c *PbCodec) Encode(buf []byte, msg proto.Message) ([]byte, error) {
	return proto.Marshal(msg)
}

func (c *PbCodec) Decode(buf []byte, msg proto.Message) error {
	return proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(buf, msg)
}
