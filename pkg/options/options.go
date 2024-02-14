// options package is a package that stores publisher's configuration and has convenience functions to configure it.
package options

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/syntropynet/data-layer-sdk/pkg/codec"
	"golang.org/x/sync/errgroup"
)

type ContextKeyType string

const (
	CancelFuncKey ContextKeyType = "cancel"
	GroupFuncKey  ContextKeyType = "group"
)

type Option func(*Options)

type NatsConn interface {
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	QueueSubscribe(subj, queue string, cb nats.MsgHandler) (*nats.Subscription, error)
	PublishMsg(m *nats.Msg) error
	RequestMsgWithContext(ctx context.Context, msg *nats.Msg) (*nats.Msg, error)
	Flush() error
}

// Codec represents a message encoder and decoder
type Codec interface {
	Encode(nmsg []byte, msg any) ([]byte, error)
	Decode(nmsg []byte, msg any) error
}

var _ NatsConn = &nats.Conn{}

// Options struct used by the publisher.
// This structure shouldn't be modified by the user, however, it is open for modification if needed.
type Options struct {
	// The global context
	Context context.Context
	// The global error group that should be used to start internal goroutines.
	// .Close() methods will wait on this group until all goroutines have exited.
	// It will cancel the Context if any of the goroutines exit with an error.
	Group *errgroup.Group
	// Cancel function is used to initiate the shutdown of the service.
	Cancel context.CancelCauseFunc

	// Include verbose logging
	VerboseLog bool
	Logger     *slog.Logger
	// Determines how often to send telemetry messages.
	TelemetryPeriod time.Duration

	// The size of the publish queue
	PublishQueueSize int
	// The size of the subscriber's queue
	SubscribeQueueSize int

	// Publishing NATS connection
	PubNats NatsConn
	// Subscribing NATS connection
	SubNats NatsConn

	// The private key for this publisher.
	// This key is used to sign the messages and is used to derive the identity.
	PrivateKey crypto.PrivateKey
	// Identity is used internally by publishers and is generated from the private key
	Identity string

	// A collection of known public keys from which we are allowed to accept messages from.
	KnownPublicKeys map[string]ed25519.PublicKey

	// Subject prefix for publishing.
	Prefix string
	// Name of the NATS queue
	QueueName string
	// Name of the publisher
	Name string
	// Codec is used to marshal published messages to whire format.
	// Can also be used to unmarshal received messages.
	Codec Codec

	// Other generic parameters that can be obtained using generic Param function.
	Params map[string]any
}

func (o *Options) setDefaults() {
	_, pkey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(fmt.Errorf("failed generating ed25519 key: %w", err))
	}
	o.Logger = slog.Default()
	o.SetContext(context.Background())
	o.PubNats = &natsStub{Verbose: &o.VerboseLog, logger: o.Logger}
	o.SubNats = &natsStub{Verbose: &o.VerboseLog, logger: o.Logger}
	o.PrivateKey = pkey
	o.Prefix = "syntropy"
	o.Name = "rnd"
	o.TelemetryPeriod = time.Second * 10
	o.Params = make(map[string]any)
	o.KnownPublicKeys = make(map[string]ed25519.PublicKey)
	o.Codec = codec.NewJsonCodec()
	o.PublishQueueSize = 1000
}

func (o *Options) SetContext(ctxMain context.Context) {
	ctx, cancel := context.WithCancelCause(ctxMain)
	o.Cancel = cancel
	o.Group, ctx = errgroup.WithContext(ctx)
	ctx = context.WithValue(ctx, CancelFuncKey, cancel)
	ctx = context.WithValue(ctx, GroupFuncKey, o.Group)
	o.Context = ctx
}

// Parse initializes the Options with default values and
// applies any functional option modifiers.
// Generally this must be called in `New` function that creates a publisher.
func (o *Options) Parse(opts ...Option) error {
	o.setDefaults()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(o)
	}

	if o.Codec == nil {
		return errors.New("codec is nil")
	}
	if o.PrivateKey == nil {
		return errors.New("private key is nil")
	}
	if o.PubNats == nil && o.SubNats == nil {
		return errors.New("nats is nil")
	}
	if o.Name == "" {
		return errors.New("empty name")
	}
	if o.Prefix == "" {
		return errors.New("empty prefix")
	}

	return nil
}

// subject returns a service subject made up of a prefix, service subject parts, and streamId.
// The subject follows the following format: `{prefix}.{parts joined}.streamId`.
// If the last element of parts contains special character ">", then it will be replaced by streamId, unless streamId is empty.
// If the prefix is "." - o.Subject() will be used instead.
func (o Options) subject(prefix string, parts []string, streamId string) string {
	var builder strings.Builder

	if prefix == "." {
		builder.WriteString(o.Subject())
	} else if prefix != "" {
		builder.WriteString(prefix)
	}
	if builder.Len() > 0 {
		builder.WriteString(".")
	}

	N := len(parts)
	for i, item := range parts {
		if i == N-1 && item == ">" && streamId != "" {
			continue
		}
		builder.WriteString(item)
		if i < N-1 {
			builder.WriteString(".")
		}
	}
	if streamId != "" {
		builder.WriteString(streamId)
	}
	return builder.String()
}

// Subject returns a constructed subject from base prefix and suffixes (subtopics) like so:
// prefix := {Prefix}.{Name}
// subject := {prefix}.{suffix1}.{suffix2}...
func (o Options) Subject(suffixes ...string) string {
	parts := make([]string, 0, len(suffixes)+3)
	if o.Prefix != "" {
		parts = append(parts, o.Prefix)
	}
	parts = append(parts, o.Name)
	for _, s := range suffixes {
		if s == "" {
			continue
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ".")
}

// Param will return a parameter if it was stored by the provided key.
// Otherwise it will return a default(first optional argument) or nil.
func (o Options) Param(key string, defaults ...any) any {
	if val, ok := o.Params[key]; ok {
		return val
	}
	if len(defaults) == 0 {
		return nil
	}
	return defaults[0]
}

// Param is a generic helper function that will return a parameter stored in Options.
// It will return a default if no key is found. Default is also used for type inference.
func Param[T any](o Options, key string, dflt T) T {
	result := o.Param(key, dflt)
	if ret, ok := result.(T); ok {
		return ret
	}
	var zero T
	return zero
}
