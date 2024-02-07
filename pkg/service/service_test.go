package service

import (
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/syntropynet/data-layer-sdk/pkg/options"
)

const (
	testSeed = "SUAAFVOGJFQB4RI5DWNGCV4SNIYQDS77XEN523SP3A32MOR6E2XJGWTZUM"
	testPeer = "BmaYfPNTt8mjUhxvYD1XT2wHuFQx3YxYP5Zki4ESjCWX"
)

func TestBase_MakeMsgVerify(t *testing.T) {
	publisher := Service{}
	publisher.Configure(
		WithName("test"),
		WithPrefix("prefix"),
		WithNKeySeed(testSeed),
	)

	if publisher.Identity == "" {
		t.Error("identity is empty")
	}

	msg, err := publisher.makeMsg([]byte("lore ipsum"), "test.test")
	if err != nil {
		t.Error("failure: ", err.Error())
		t.Fail()
	}

	var (
		msg_out_counter   atomic.Uint64
		bytes_out_counter atomic.Uint64
	)

	wrapped := wrapMessage(publisher.Codec, &msg_out_counter, &bytes_out_counter, publisher.makeMsg, msg)
	err = publisher.Verify(wrapped)
	if err != nil {
		t.Error("failure: ", err.Error())
		t.Fail()
	}
}

func TestPublisher_getUserId(t *testing.T) {
	tests := []struct {
		name string
		subj string
		suff []string
		want []string
	}{
		{
			"simple case",
			"foo.bar.req.a.b",
			[]string{"req", ">"},
			[]string{"req", "a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Service{
				Options: options.Options{
					Prefix: "foo",
					Name:   "bar",
				},
			}
			if got := p.GetStreamIdParts(&nats.Msg{Subject: tt.subj}, tt.suff...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Publisher.getUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase_Verify(t *testing.T) {
	tests := []struct {
		name    string
		peers   []string
		makeMsg func(t *testing.T, b *Service) *nats.Msg
		wantErr bool
	}{
		{
			"message signed and empty peers",
			[]string{},
			func(t *testing.T, b *Service) *nats.Msg {
				msg, err := b.makeMsg([]byte("lore ipsum"), "test.test")
				if err != nil {
					t.Error("failure: ", err.Error())
					t.Fail()
				}
				return msg
			},
			false,
		},
		{
			"message w/o identity and empty peers",
			[]string{},
			func(t *testing.T, b *Service) *nats.Msg {
				msg, err := b.makeMsg([]byte("lore ipsum"), "test.test")
				if err != nil {
					t.Error("failure: ", err.Error())
					t.Fail()
				}
				msg.Header.Del("identity")
				return msg
			},
			false,
		},
		{
			"message w/o signature and empty peers",
			[]string{},
			func(t *testing.T, b *Service) *nats.Msg {
				msg, err := b.makeMsg([]byte("lore ipsum"), "test.test")
				if err != nil {
					t.Error("failure: ", err.Error())
					t.Fail()
				}
				msg.Header.Del("signature")
				return msg
			},
			false,
		},
		{
			"modified message signed and empty peers",
			[]string{},
			func(t *testing.T, b *Service) *nats.Msg {
				msg, err := b.makeMsg([]byte("lore ipsum"), "test.test")
				if err != nil {
					t.Error("failure: ", err.Error())
					t.Fail()
				}
				msg.Data[3] = 1
				return msg
			},
			true,
		},
		{
			"message signed and good peers",
			[]string{testPeer},
			func(t *testing.T, b *Service) *nats.Msg {
				msg, err := b.makeMsg([]byte("lore ipsum"), "test.test")
				if err != nil {
					t.Error("failure: ", err.Error())
					t.Fail()
				}
				return msg
			},
			false,
		},
		{
			"message signed and wrong peers",
			[]string{"CspLUFkBWXdimZ2b1sEo4nYpCRDdDpnCEhjWbNVKPmEN"},
			func(t *testing.T, b *Service) *nats.Msg {
				msg, err := b.makeMsg([]byte("lore ipsum"), "test.test")
				if err != nil {
					t.Error("failure: ", err.Error())
					t.Fail()
				}
				return msg
			},
			true,
		},
		{
			"message w/o identity signed with peers",
			[]string{testPeer},
			func(t *testing.T, b *Service) *nats.Msg {
				msg, err := b.makeMsg([]byte("lore ipsum"), "test.test")
				if err != nil {
					t.Error("failure: ", err.Error())
					t.Fail()
				}
				msg.Header.Del("identity")
				return msg
			},
			true,
		},
		{
			"message w/o signature with peers",
			[]string{testPeer},
			func(t *testing.T, b *Service) *nats.Msg {
				msg, err := b.makeMsg([]byte("lore ipsum"), "test.test")
				if err != nil {
					t.Error("failure: ", err.Error())
					t.Fail()
				}
				msg.Header.Del("signature")
				return msg
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Service{}
			b.Configure(
				WithName("bar"),
				WithPrefix("foo"),
				WithNKeySeed(testSeed),
				WithKnownIdentities(tt.peers...),
			)
			msg := tt.makeMsg(t, b)

			var (
				msg_out_counter   atomic.Uint64
				bytes_out_counter atomic.Uint64
			)

			wrapped := wrapMessage(b.Codec, &msg_out_counter, &bytes_out_counter, b.makeMsg, msg)
			if err := b.Verify(wrapped); (err != nil) != tt.wantErr {
				t.Errorf("Base.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
