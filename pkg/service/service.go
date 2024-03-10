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

type JetStreamer interface {
	JetStream(opts ...nats.JSOpt) (nats.JetStreamContext, error)
}

var _ JetStreamer = &nats.Conn{}

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
	status["messages"] = map[string]any{
		"out_queue":     len(b.publishCh),
		"out_queue_cap": cap(b.publishCh),
		"in":            b.msg_in_counter.Swap(0),
		"out":           b.msg_out_counter.Swap(0),
		"bytes_in":      b.bytes_in_counter.Swap(0),
		"bytes_out":     b.bytes_out_counter.Swap(0),
	}

	return status
}

// Subscribe will subscribe to a subject constructed from {prefix}.{name}.{...suffixes}, where
// suffixes are joined using ".".
func (b *Service) Subscribe(handler MessageHandler, suffixes ...string) (*nats.Subscription, error) {
	return b.SubscribeTo(handler, b.Subject(suffixes...))
}

// Serve is a convenience method to serve a service subject. It acts the same as Subscribe, but takes `ServiceHandler` instead.
func (b *Service) Serve(handler ServiceHandler, suffixes ...string) (*nats.Subscription, error) {
	return b.SubscribeTo(
		func(msg Message) {
			resp, err := handler(msg)
			if err != nil {
				b.Logger.Error("service handler failed", "err", err, "suffixes", suffixes)
				err1 := msg.Respond(fmt.Sprintf("error:%s", err.Error()))
				if err != nil {
					b.Logger.Error("service handler failed during error", "err", err, "err1", err1, "suffixes", suffixes)
				}
				return
			}
			err = msg.Respond(resp)
			if err != nil {
				b.Logger.Error("service handler failed", "err", err, "suffixes", suffixes)
			}
		},
		b.Subject(suffixes...),
	)
}

// SubscribeTo will subscribe to a subject constructed {...suffixes}, where
// suffixes are joined using ".".
//
// Experimental: When a stream was registered with AddStream SubscribeTo will use durable stream instead of realtime.
func (b *Service) SubscribeTo(handler MessageHandler, suffixes ...string) (*nats.Subscription, error) {
	if b.SubNats == nil {
		return nil, fmt.Errorf("subscribing NATS connection is nil")
	}
	if b.VerboseLog {
		b.Logger.Debug("SubscribeTo", "suffixes", suffixes)
	}
	natsHandler := func(msg *nats.Msg) {
		b.msg_in_counter.Add(1)
		b.bytes_in_counter.Add(uint64(len(msg.Data)))
		wrapped := wrapMessage(b.Codec, &b.msg_out_counter, &b.bytes_out_counter, b.makeMsg, msg)
		handler(wrapped)
	}

	subject := strings.Join(suffixes, ".")

	if sub, err := b.attemptJSConsume(natsHandler, subject); err == nil {
		return sub, nil
	}

	if b.QueueName != "" {
		return b.SubNats.QueueSubscribe(subject, b.QueueName, natsHandler)
	}
	return b.SubNats.Subscribe(subject, natsHandler)
}

// Close should be closed to clean-up the publisher.
func (b *Service) Close() error {
	b.Cancel(nil)
	return b.Group.Wait()
}

// Fail is a convenience function that allows to asynchronously propagate errors.
func (b *Service) Fail(err error) {
	b.Logger.Error("Publisher failed", err)
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

func (b *Service) Respond(msg Message, buf []byte, suffixes ...string) error {
	payload, err := b.Codec.Encode(nil, msg.Message())
	if err != nil {
		return err
	}
	return b.RespondBuf(msg, payload, suffixes...)
}

func (b *Service) RespondBuf(msg Message, buf []byte, suffixes ...string) error {
	reply, err := b.makeMsg(buf, strings.Join(suffixes, "."))
	if err != nil {
		return err
	}

	return msg.Message().RespondMsg(reply)
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
		b.Logger.Info("PublishBufTo cancelled", "err", b.Context.Err(), "queue_size", len(b.publishCh))
		return b.Context.Err()
	case b.publishCh <- msg:
	}
	return nil
}

// RequestFrom requests a reply or a stream from a subject using subscribing NATS connection.
// This a synchronous operation that does not involve publisher queue.
func (b *Service) RequestFrom(ctx context.Context, msg any, resp any, suffixes ...string) (Message, error) {
	payload, err := b.Codec.Encode(nil, msg)
	if err != nil {
		return nil, err
	}
	response, err := b.RequestBufFrom(ctx, payload, suffixes...)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		_, err := b.Unmarshal(response, resp)
		if err != nil {
			return response, fmt.Errorf("unmarshal failed: %w", err)
		}
	}

	return response, err
}

// RequestBufFrom requests a reply or a stream from a subject using subscribing NATS connection.
// This a synchronous operation that does not involve publisher queue.
func (b *Service) RequestBufFrom(ctx context.Context, buf []byte, suffixes ...string) (Message, error) {
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

	return wrapMessage(b.Codec, &b.msg_out_counter, &b.bytes_out_counter, b.makeMsg, ret), nil
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
