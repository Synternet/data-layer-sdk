// codec package implements various codecs to be used with the publisher to encode and decode messages.
package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ProtoJsonCodec struct {
}

func NewProtoJsonCodec() *ProtoJsonCodec {
	return &ProtoJsonCodec{}
}

func (c *ProtoJsonCodec) Encode(buf []byte, msg proto.Message) ([]byte, error) {
	return protojson.Marshal(msg)
}

func (c *ProtoJsonCodec) Decode(buf []byte, msg proto.Message) error {
	return protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(buf, msg)
}
