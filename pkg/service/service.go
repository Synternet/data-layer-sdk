package service

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/synternet/data-layer-sdk/pkg/options"
	"google.golang.org/protobuf/proto"

	"github.com/cosmos/btcutil/base58"
	"github.com/nats-io/nats.go"
)

// StatusFunc is a type for callback func. This will be called periodically to construct telemetry status.
// "uptime", "goroutines", "period" keys will be overriden internally.
type StatusFunc func() map[string]string

type JetStreamer interface {
	JetStream(opts ...nats.JSOpt) (nats.JetStreamContext, error)
}

var _ JetStreamer = &nats.Conn{}

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrInvalidIdentity  = errors.New("invalid identity")
	ErrUnknownIdentity  = errors.New("unknown identity")
	ErrPubConnection    = errors.New("publishing NATS connection is nil")
	ErrSubConnection    = errors.New("subscribing NATS connection is nil")
	ErrReqConnection    = errors.New("request NATS connection is nil")
)

type jsStream struct {
	cfgStream    *nats.StreamConfig
	cfgConsumer  *nats.ConsumerConfig
	streamInfo   *nats.StreamInfo
	consumerInfo *nats.ConsumerInfo
}

// Service is the base publisher structure. It must be embedded in the publisher to benefit from the
// common implementation. Config method must be called on this embedded struct in order to
// properly set it up.
type Service struct {
	options.Options
	mu                sync.Mutex
	startTime         time.Time
	prevTelemetry     time.Time
	publishCh         chan *nats.Msg // This is out of bounds publish channel
	nonce             atomic.Uint64
	msg_in_counter    atomic.Uint64
	msg_out_counter   atomic.Uint64
	bytes_in_counter  atomic.Uint64
	bytes_out_counter atomic.Uint64
	statusCallback    map[uintptr]StatusFunc

	// Experimental feature
	js             nats.JetStreamContext
	streams        []jsStream
	streamSubjects SubjectMap
}

// Configure must be called by the publisher implementation.
func (b *Service) Configure(opts ...options.Option) error {
	if b.statusCallback == nil {
		b.statusCallback = make(map[uintptr]StatusFunc)
	}
	err := b.Parse(opts...)
	if err != nil {
		return fmt.Errorf("failed parsing options: %w", err)
	}
	b.publishCh = make(chan *nats.Msg, b.PublishQueueSize)

	privKey, ok := b.PrivateKey.(ed25519.PrivateKey)
	if !ok {
		return fmt.Errorf("private key is not ED25519")
	}
	pubkey, ok := privKey.Public().(ed25519.PublicKey)
	if !ok {
		return fmt.Errorf("derived public key is not ED25519")
	}
	b.Identity = base58.Encode(pubkey)
	defer func() {
		b.Logger.Info("Service configured", "identity", b.Identity, "JetStream", b.js != nil)
	}()

	if b.SubNats == nil {
		return nil
	}

	if js, ok := b.SubNats.(JetStreamer); ok {
		b.js, err = js.JetStream(nats.Context(b.Context))
		if err != nil {
			b.Logger.Error("JetStream failed", err)
			return nil
		}
		b.streamSubjects = make(SubjectMap)
	}
	return nil
}

func (b *Service) run() error {
	sub, err := b.Subscribe(b.handleTelemetryPing, "telemetry", "ping")
	if err != nil {
		err = fmt.Errorf("Telemetry subscription failed: %w", err)
		b.Cancel(err)
		return err
	}
	defer sub.Unsubscribe()

	ticker := time.NewTicker(b.TelemetryPeriod)
	for {
		select {
		case <-b.Context.Done():
			b.Logger.Info("Message Publish loop Context closed", "err", b.Context.Err())
			return nil
		case <-ticker.C:
			b.reportTelemetry()
		case msg := <-b.publishCh:
			if b.PubNats == nil {
				b.Logger.Warn("Messages are being published to nil NATS connection")
				continue
			}
			b.PubNats.PublishMsg(msg)
			b.msg_out_counter.Add(1)
			b.bytes_out_counter.Add(uint64(len(msg.Data)))
		}
	}
}

// Start will start the publisher base implementation that includes telemetry and other internal goroutines.
func (b *Service) Start() context.Context {
	b.startTime = time.Now()
	b.prevTelemetry = time.Now()
	b.Group.Go(b.run)
	return b.Context
}

// AddStatusCallback will register a status callback.
func (b *Service) AddStatusCallback(callback StatusFunc) {
	if callback == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.statusCallback[reflect.ValueOf(callback).Pointer()] = callback
}

// RemoveStatusCallback will remove a status callback.
func (b *Service) RemoveStatusCallback(callback StatusFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.statusCallback, reflect.ValueOf(callback).Pointer())
}

// Concatenate status copies (k,v) from items to status optionally prefixing the key with path.
func (b *Service) ConcatenateStatus(path string, status, items map[string]string) map[string]string {
	if items == nil || status == nil {
		return status
	}
	for k, v := range items {
		if path != "" {
			k = fmt.Sprintf("%s.%s", path, k)
		}
		status[k] = v
	}

	return status
}

func (b *Service) collectStatus() map[string]string {
	b.mu.Lock()
	defer b.mu.Unlock()

	status := make(map[string]string)

	for _, cb := range b.statusCallback {
		newStatus := cb()
		b.ConcatenateStatus("", status, newStatus)
	}

	status["uptime"] = time.Since(b.startTime).String()
	status["period"] = time.Since(b.prevTelemetry).String()
	status["goroutines"] = strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
	b.ConcatenateStatus(
		"messages",
		status,
		map[string]string{
			"out_queue":     strconv.FormatInt(int64(len(b.publishCh)), 10),
			"out_queue_cap": strconv.FormatInt(int64(cap(b.publishCh)), 10),
			"in":            strconv.FormatUint(b.msg_in_counter.Swap(0), 10),
			"out":           strconv.FormatUint(b.msg_out_counter.Swap(0), 10),
			"bytes_in":      strconv.FormatUint(b.bytes_in_counter.Swap(0), 10),
			"bytes_out":     strconv.FormatUint(b.bytes_out_counter.Swap(0), 10),
		},
	)

	return status
}

// Close should be closed to clean-up the publisher.
func (b *Service) Close() error {
	b.Cancel(nil)
	return b.Group.Wait()
}

// Fail is a convenience function that allows to asynchronously propagate errors.
func (b *Service) Fail(err error) {
	b.Logger.Error("Publisher failed", "err", err)
	b.Cancel(err)
}

func (b *Service) makeMsg(payload []byte, subject string) (*nats.Msg, error) {
	signature, _, err := b.Sign(payload)
	if err != nil {
		return nil, err
	}

	result := &nats.Msg{
		Subject: subject,
		Data:    payload,
		Header: map[string][]string{
			"identity":  {b.Identity},
			"signature": {base64.StdEncoding.EncodeToString(signature)},
			"timestamp": {strconv.FormatInt(time.Now().UnixNano(), 10)},
		},
	}
	return result, nil
}

// Unmarshal is a convenience function that first verifies any signatures in the message and unmarshals bytes into a message.
func (b *Service) Unmarshal(nmsg Message, msg proto.Message) (nats.Header, error) {
	// TODO check signatures
	if err := b.Verify(nmsg); err != nil {
		return nmsg.Header(), err
	}
	return nmsg.Header(), b.Codec.Decode(nmsg.Data(), msg)
}

// Sign will sign the bytes.
func (b *Service) Sign(msg []byte) (signature []byte, publicKey []byte, err error) {
	if b.PrivateKey == nil {
		return nil, nil, nil
	}

	pkey := b.PrivateKey.(ed25519.PrivateKey)
	pubkey := pkey.Public().(ed25519.PublicKey)

	return ed25519.Sign(pkey, msg), []byte(pubkey), nil
}

// Verify will marshal the unmarshalled payload back to bytes and verifies the signature with that.
func (b *Service) Verify(nmsg Message) error {
	identity := nmsg.Header().Get("identity")
	signature := nmsg.Header().Get("signature")
	var pkey ed25519.PublicKey

	switch {
	case len(b.KnownPublicKeys) != 0:
		if key, ok := b.KnownPublicKeys[identity]; !ok {
			return ErrUnknownIdentity
		} else {
			pkey = key
		}
	case identity == "" || signature == "":
		return nil
	default:
		identityBytes := base58.Decode(identity)
		if len(identityBytes) != ed25519.PublicKeySize {
			return ErrInvalidIdentity
		}
		pkey = ed25519.PublicKey(identityBytes)
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	if !ed25519.Verify(pkey, nmsg.Data(), signatureBytes) {
		return ErrInvalidSignature
	}

	return nil
}
