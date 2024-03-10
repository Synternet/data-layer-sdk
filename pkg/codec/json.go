// codec package implements various codecs to be used with the publisher to encode and decode messages.
package codec

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

type JsonCodec struct {
}

func NewJsonCodec() *JsonCodec {
	return &JsonCodec{}
}

func (c *JsonCodec) Encode(buf []byte, msg proto.Message) ([]byte, error) {
	return json.Marshal(msg)
}

func (c *JsonCodec) Decode(buf []byte, msg proto.Message) error {
	return json.Unmarshal(buf, msg)
}
