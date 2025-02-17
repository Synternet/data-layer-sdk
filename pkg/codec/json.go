// codec package implements various codecs to be used with the publisher to encode and decode messages.
package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type JsonCodec struct {
}

func NewJsonCodec() *JsonCodec {
	return &JsonCodec{}
}

func (c *JsonCodec) Encode(buf []byte, msg proto.Message) ([]byte, error) {
	return protojson.Marshal(msg)
}

func (c *JsonCodec) Decode(buf []byte, msg proto.Message) error {
	return protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(buf, msg)
}
