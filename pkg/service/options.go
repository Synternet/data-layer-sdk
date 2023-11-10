package service

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cosmos/btcutil/base58"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/syntropynet/data-layer-sdk/pkg/options"
)

// WithParam adds a key-value pair to the configuration of the publisher.
func WithParam(key string, val any) options.Option {
	return func(o *options.Options) {
		o.Params[key] = val
	}
}

// WithContext sets the context for the publisher.
func WithContext(ctx context.Context) options.Option {
	return func(o *options.Options) {
		o.SetContext(ctx)
	}
}

// WithNats sets up preconfigured NATS connector for publishing and subscribing.
func WithNats(nc *nats.Conn) options.Option {
	return func(o *options.Options) {
		if nc == nil {
			return
		}
		o.PubNats = nc
		o.SubNats = nc
	}
}

// WithPubNats sets up preconfigured NATS connector specifically for publishing.
func WithPubNats(nc *nats.Conn) options.Option {
	return func(o *options.Options) {
		if nc == nil {
			return
		}
		o.PubNats = nc
	}
}

// WithSubNats sets up preconfigured NATS connector speficivally for subscribing.
func WithSubNats(nc *nats.Conn) options.Option {
	return func(o *options.Options) {
		if nc == nil {
			return
		}
		o.SubNats = nc
	}
}

// WithName sets name of the publisher. The subject is in the form of {prefix}.{name}.
// Subscriptions and publishing will use the subject constructed from prefix and name.
func WithName(name string) options.Option {
	return func(o *options.Options) {
		if name == "" {
			log.Fatal("name must be provided")
		}
		o.Name = name
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
			log.Fatal(errors.New("private key must implement Signer interface"))
		}
		if _, ok := pkey.(ed25519.PrivateKey); !ok {
			log.Fatal(errors.New("private key must be ED25519"))
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
			log.Fatal(err)
		}

		if len(privateKey) > 1024 {
			log.Fatal(errors.New("private key exceeds 1 KB"))
		}
		if !strings.HasPrefix(string(privateKey), "-----BEGIN PRIVATE KEY-----") {
			log.Fatal(errors.New("does not contain private key section"))
		}

		var block *pem.Block
		if block, _ = pem.Decode(privateKey); block == nil {
			log.Fatal(errors.New("not a valid pem private key"))
		}

		// Parse the key
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		pkey, ok := parsedKey.(ed25519.PrivateKey)
		if !ok {
			log.Fatal(errors.New("private key must be ed25519 key"))
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
			log.Fatalln("credentials file:", err.Error())
		}
		defer f.Close()
		contents, err := io.ReadAll(f)
		if err != nil {
			log.Fatalln("credentials file:", err.Error())
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
			log.Fatalln("credentials nkey:", err.Error())
		}
		npkey, err := key.PrivateKey()
		if err != nil {
			log.Fatalln("credentials private key:", err.Error())
		}

		pkey, err := nkeys.Decode(nkeys.PrefixBytePrivate, []byte(npkey))
		if err != nil {
			log.Fatalln("nkey decode:", err.Error(), "nkey:", npkey)
		}
		if len(pkey) != ed25519.PrivateKeySize {
			log.Fatalln("NKey: key size mismatch: ", len(pkey), ed25519.PrivateKeySize)
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
