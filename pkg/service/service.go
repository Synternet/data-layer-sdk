package service

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/syntropynet/data-layer-sdk/pkg/options"

	"github.com/cosmos/btcutil/base58"
	"github.com/nats-io/nats.go"
)

// StatusFunc is a type for callback func. This will be called periodically to construct telemetry status.
// "uptime", "goroutines", "period" keys will be overriden internally.
type StatusFunc func() map[string]any

// Service is the base publisher structure. It must be embedded in the publisher to benefit from the
// common implementation. Config method must be called on this embedded struct in order to
// properly set it up.
type Service struct {
	options.Options
	mu             sync.Mutex
	startTime      time.Time
	prevTelemetry  time.Time
	publishCh      chan *nats.Msg // This is out of bounds publish channel
	counter        atomic.Uint64
	statusCallback map[uintptr]StatusFunc
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
	log.Println("Identity: ", b.Identity)
	return nil
}

func (b *Service) run() error {
	sub, err := b.Subscribe(b.handleTelemetryPing, "telemetry", "ping")
	if err != nil {
		b.Cancel(err)
		return err
	}
	defer sub.Unsubscribe()

	ticker := time.NewTicker(b.TelemetryPeriod)
	for {
		select {
		case <-b.Context.Done():
			log.Println("Message Publish loop Context closed")
			return nil
		case <-ticker.C:
			b.reportTelemetry()
		case msg := <-b.publishCh:
			if b.PubNats == nil {
				log.Println("Messages are being published to nil NATS connection")
				continue
			}
			b.PubNats.PublishMsg(msg)
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

func (b *Service) collectStatus() map[string]any {
	b.mu.Lock()
	defer b.mu.Unlock()

	status := make(map[string]any)

	for _, cb := range b.statusCallback {
		newStatus := cb()
		if newStatus == nil {
			continue
		}
		for k, v := range newStatus {
			status[k] = v
		}
	}

	status["uptime"] = time.Since(b.startTime).String()
	status["period"] = time.Since(b.prevTelemetry).String()
	status["goroutines"] = runtime.NumGoroutine()
	status["publish_queue"] = len(b.publishCh)
	status["publish_queue_cap"] = cap(b.publishCh)

	return status
}

// Subscribe will subscribe to a subject constructed from {prefix}.{name}.{...suffixes}, where
// suffixes are joined using ".".
func (b *Service) Subscribe(handler MessageHandler, suffixes ...string) (*nats.Subscription, error) {
	return b.SubscribeTo(handler, b.Subject(suffixes...))
}

// SubscribeTo will subscribe to a subject constructed {...suffixes}, where
// suffixes are joined using ".".
func (b *Service) SubscribeTo(handler MessageHandler, suffixes ...string) (*nats.Subscription, error) {
	if b.SubNats == nil {
		return nil, fmt.Errorf("subscribing NATS connection is nil")
	}
	if b.VerboseLog {
		log.Println("SubscribeTo", strings.Join(suffixes, "."))
	}
	natsHandler := func(msg *nats.Msg) {
		wrapped := wrapMessage(b.Codec, b.makeMsg, msg)
		handler(wrapped)
	}

	if b.QueueName != "" {
		return b.SubNats.QueueSubscribe(strings.Join(suffixes, "."), b.QueueName, natsHandler)
	}
	return b.SubNats.Subscribe(strings.Join(suffixes, "."), natsHandler)
}

// Close should be closed to clean-up the publisher.
func (b *Service) Close() error {
	b.Cancel(nil)
	return b.Group.Wait()
}

// Fail is a convenience function that allows to asynchronously propagate errors.
func (b *Service) Fail(err error) {
	log.Println("Publisher failed: ", err.Error())
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

// GetStreamIdParts returns subject parts that are considered as a stream identifier.
// For example, if you receive messages on `prefix.name.query.>`, then every subtopic in place of `>` will be
// considered as a stream id. Suffixes argument in this case will be `query.>`.
// In this example, GetStreamIdParts will return `query.part1.part2.part3` if the message was received on subject `prefix.name.query.part1.part2.part3`.
func (b *Service) GetStreamIdParts(nmsg *nats.Msg, suffixes ...string) []string {
	N := strings.Count(b.Subject(suffixes...), ".")
	if N < 0 {
		return nil
	}
	parts := strings.Split(nmsg.Subject, ".")
	return parts[N-1:]
}

// GetStreamId will return the stream parts as well as joined stream id.
// For example it will return (`[query, part1, part2, part3]`, `part1.part2.part3`) if the message was received on subject `prefix.name.query.part1.part2.part3`.
func (b *Service) GetStreamId(nmsg *nats.Msg, suffixes ...string) ([]string, string) {
	parts := b.GetStreamIdParts(nmsg, suffixes...)
	if parts == nil {
		return nil, ""
	}
	streamId := strings.Join(parts[1:], ".")
	return parts, streamId
}

// Unmarshal is a convenience function that first verifies any signatures in the message and unmarshals bytes into a message.
func (b *Service) Unmarshal(nmsg Message, msg any) (nats.Header, error) {
	// TODO check signatures
	if err := b.Verify(nmsg); err != nil {
		return nmsg.Header(), err
	}
	return nmsg.Header(), b.Codec.Decode(nmsg.Data(), msg)
}

// Publish will sign the message and publish it.
// It will sign the message in place. Also it will update the timestamp and identity fields.
func (b *Service) Publish(msg any, suffixes ...string) error {
	return b.PublishTo(msg, b.Subject(suffixes...))
}

// PublishBuf will publish the raw bytes.
func (b *Service) PublishBuf(buf []byte, suffixes ...string) error {
	return b.PublishBufTo(buf, b.Subject(suffixes...))
}

// PublishTo will sign the message and publish it to a specific subject.
// It will sign the message in place. Also it will update the timestamp and identity fields.
func (b *Service) PublishTo(msg any, suffixes ...string) error {
	payload, err := b.Codec.Encode(nil, msg)
	if err != nil {
		return err
	}
	return b.PublishBufTo(payload, suffixes...)
}

func (b *Service) Respond(msg *nats.Msg, buf []byte, suffixes ...string) error {
	payload, err := b.Codec.Encode(nil, msg)
	if err != nil {
		return err
	}
	return b.RespondBuf(msg, payload, suffixes...)
}

func (b *Service) RespondBuf(msg *nats.Msg, buf []byte, suffixes ...string) error {
	msg, err := b.makeMsg(buf, strings.Join(suffixes, "."))
	if err != nil {
		return err
	}

	return msg.RespondMsg(msg)
}

// PublishBufTo will publish the raw bytes to a specific subject.
func (b *Service) PublishBufTo(buf []byte, suffixes ...string) error {
	if b.PubNats == nil {
		return fmt.Errorf("publishing NATS connection is nil")
	}
	msg, err := b.makeMsg(buf, strings.Join(suffixes, "."))
	if err != nil {
		return err
	}

	select {
	case <-b.Context.Done():
		return b.Context.Err()
	case b.publishCh <- msg:
	}
	return nil
}

// RequestFrom requests a reply or a stream from a subject using subscribing NATS connection.
// This a synchronous operation that does not involve publisher queue.
func (b *Service) RequestFrom(ctx context.Context, msg any, suffixes ...string) (*nats.Msg, error) {
	payload, err := b.Codec.Encode(nil, msg)
	if err != nil {
		return nil, err
	}
	return b.RequestBufFrom(ctx, payload, suffixes...)
}

// RequestBufFrom requests a reply or a stream from a subject using subscribing NATS connection.
// This a synchronous operation that does not involve publisher queue.
func (b *Service) RequestBufFrom(ctx context.Context, buf []byte, suffixes ...string) (*nats.Msg, error) {
	if b.SubNats == nil {
		return nil, fmt.Errorf("subscribing NATS connection is nil")
	}

	msg, err := b.makeMsg(buf, strings.Join(suffixes, "."))
	if err != nil {
		return nil, err
	}

	ret, err := b.SubNats.RequestMsgWithContext(ctx, msg)
	if err != nil {
		return nil, err
	}

	return ret, nil
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
			return errors.New("unknown identity")
		} else {
			pkey = key
		}
	case identity == "" || signature == "":
		return nil
	default:
		identityBytes := base58.Decode(identity)
		if len(identityBytes) != ed25519.PublicKeySize {
			return errors.New("invalid identity")
		}
		pkey = ed25519.PublicKey(identityBytes)
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	if !ed25519.Verify(pkey, nmsg.Data(), signatureBytes) {
		return errors.New("invalid signature")
	}

	return nil
}
