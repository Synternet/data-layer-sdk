package service

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/cosmos/btcutil/base58"
	"github.com/nats-io/nkeys"
	"github.com/synternet/data-layer-sdk/pkg/options"
)

// WithParam adds a key-value pair to the configuration of the publisher.
func WithParam(key string, val any) options.Option {
	return func(o *options.Options) {
		o.Params[key] = val
	}
}

// WithLogger sets the logger for the publisher.
func WithLogger(logger *slog.Logger) options.Option {
	return func(o *options.Options) {
		if logger == nil {
			panic(errors.New("logger must not be nil"))
		}
		o.Logger = logger
	}
}

// WithContext sets the context for the publisher.
func WithContext(ctx context.Context) options.Option {
	return func(o *options.Options) {
		o.SetContext(ctx)
	}
}

// WithNats sets up preconfigured NATS connector for publishing, subscribing, and request/reply.
func WithNats(nc options.NatsConn) options.Option {
	return func(o *options.Options) {
		WithPubNats(nc)(o)
		WithSubNats(nc)(o)
		WithReqNats(nc)(o)
	}
}

// WithPubNats sets up preconfigured NATS connector specifically for publishing.
func WithPubNats(nc options.NatsConn) options.Option {
	return func(o *options.Options) {
		if reflect.ValueOf(nc).IsNil() {
			return
		}
		o.PubNats = nc
	}
}

// WithSubNats sets up preconfigured NATS connector specifically for subscribing.
//
// NOTE: This will also set ReqNats for compatibility.
func WithSubNats(nc options.NatsConn) options.Option {
	return func(o *options.Options) {
		if reflect.ValueOf(nc).IsNil() {
			return
		}
		o.SubNats = nc
		WithReqNats(nc)(o)
	}
}

// WithReqNats sets up preconfigured NATS connector specifically for request/reply.
func WithReqNats(nc options.NatsConn) options.Option {
	return func(o *options.Options) {
		if reflect.ValueOf(nc).IsNil() {
			return
		}
		o.ReqNats = nc
	}
}

// WithName sets name of the publisher. The subject is in the form of {prefix}.{name}.
// Subscriptions and publishing will use the subject constructed from prefix and name.
func WithName(name string) options.Option {
	return func(o *options.Options) {
		if name == "" {
			panic(errors.New("name must be provided"))
		}
		o.Name = name
	}
}

// WithStreamName sets name of the JetStream stream.
// If an empty string is used, then stream name will be "{prefix}-{name}".
func WithStreamName(name string) options.Option {
	return func(o *options.Options) {
		o.StreamName = name
	}
}

// WithPrefix sets the prefix for the subject. The subject is in the form of {prefix}.{name}.
// Subscriptions and publishing will use the subject constructed from prefix and name.
func WithPrefix(prefix string) options.Option {
	return func(o *options.Options) {
		if prefix == "" {
			return
		}
		o.Prefix = prefix
	}
}

// WithPrivateKey will configure identity using ED25519 crypto.PrivateKey.
func WithPrivateKey(pkey crypto.PrivateKey) options.Option {
	return func(o *options.Options) {
		if pkey == nil {
			return
		}
		if _, ok := pkey.(crypto.Signer); !ok {
			panic(errors.New("private key must implement Signer interface"))
		}
		if _, ok := pkey.(ed25519.PrivateKey); !ok {
			panic(errors.New("private key must be ED25519"))
		}
		o.PrivateKey = pkey
	}
}

// WithPemPrivateKey will load ED25519 private key from a PEM file and use it for identity.
func WithPemPrivateKey(keyFile string) options.Option {
	return func(o *options.Options) {
		if keyFile == "" {
			return
		}

		privateKey, err := os.ReadFile(keyFile)
		if err != nil {
			panic(err)
		}

		if len(privateKey) > 1024 {
			panic(errors.New("private key exceeds 1 KB"))
		}
		if !strings.HasPrefix(string(privateKey), "-----BEGIN PRIVATE KEY-----") {
			panic(errors.New("does not contain private key section"))
		}

		var block *pem.Block
		if block, _ = pem.Decode(privateKey); block == nil {
			panic(errors.New("not a valid pem private key"))
		}

		// Parse the key
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			panic(err)
		}
		pkey, ok := parsedKey.(ed25519.PrivateKey)
		if !ok {
			panic(errors.New("private key must be ed25519 key"))
		}

		WithPrivateKey(pkey)(o)
	}
}

// WithUserCreds will load NKey from NATS User Credentials file and use it for identity.
func WithUserCreds(path string) options.Option {
	return func(o *options.Options) {
		if path == "" {
			return
		}
		f, err := os.Open(path)
		if err != nil {
			panic(fmt.Errorf("credentials file: %w", err))
		}
		defer f.Close()
		contents, err := io.ReadAll(f)
		if err != nil {
			panic(fmt.Errorf("credentials file: %w", err))
		}
		WithNKeySeed(string(contents))(o)
	}
}

// WithNKeySeed will decode decorated NATS NKey and use it for identity.
func WithNKeySeed(contents string) options.Option {
	if contents == "" {
		return func(o *options.Options) {}
	}
	return func(o *options.Options) {
		key, err := nkeys.ParseDecoratedNKey([]byte(contents))
		if err != nil {
			panic(fmt.Errorf("credentials nkey: %w", err))
		}
		npkey, err := key.PrivateKey()
		if err != nil {
			panic(fmt.Errorf("credentials private key: %w", err))
		}

		pkey, err := nkeys.Decode(nkeys.PrefixBytePrivate, []byte(npkey))
		if err != nil {
			panic(fmt.Errorf("nkey decode: %w nkey=%s", err, npkey))
		}
		if len(pkey) != ed25519.PrivateKeySize {
			panic(fmt.Errorf("NKey: key size mismatch: %d != %d", len(pkey), ed25519.PrivateKeySize))
		}
		o.PrivateKey = ed25519.PrivateKey(pkey)
	}
}

// WithVerbose will configure verbosity level.
func WithVerbose(v bool) options.Option {
	return func(o *options.Options) {
		o.VerboseLog = v
	}
}

// WithTelemetryPeriod will configure verbosity level.
func WithTelemetryPeriod(p time.Duration) options.Option {
	return func(o *options.Options) {
		o.TelemetryPeriod = p
	}
}

// WithPublishQueueSize will configure the size of the publish queue.
func WithPublishQueueSize(n int) options.Option {
	return func(o *options.Options) {
		o.PublishQueueSize = n
	}
}

// WithCodec will configure the codec.
func WithCodec(c options.Codec) options.Option {
	return func(o *options.Options) {
		o.Codec = c
	}
}

// WithKnownPublicKey will add known public keys.
// Only Messages from publishers identified by these keys will be accepted.
func WithKnownPublicKey(keys ...ed25519.PublicKey) options.Option {
	return func(o *options.Options) {
		for _, key := range keys {
			o.KnownPublicKeys[base58.Encode(key)] = key
		}
	}
}

// WithKnownIdentities will add known identities as base58 encoded ed25519 public keys.
// Only Messages from publishers identified by these keys will be accepted.
func WithKnownIdentities(ids ...string) options.Option {
	return func(o *options.Options) {
		for _, id := range ids {
			bytes := base58.Decode(id)
			if bytes == nil {
				continue
			}
			if len(bytes) != ed25519.PublicKeySize {
				continue
			}
			o.KnownPublicKeys[id] = ed25519.PublicKey(bytes)
		}
	}
}
