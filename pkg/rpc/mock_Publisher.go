// Code generated by mockery v2.42.3. DO NOT EDIT.

package rpc

import (
	context "context"

	nats "github.com/nats-io/nats.go"
	mock "github.com/stretchr/testify/mock"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	service "github.com/synternet/data-layer-sdk/pkg/service"
)

// MockPublisher is an autogenerated mock type for the Publisher type
type MockPublisher struct {
	mock.Mock
}

type MockPublisher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPublisher) EXPECT() *MockPublisher_Expecter {
	return &MockPublisher_Expecter{mock: &_m.Mock}
}

// PublishTo provides a mock function with given fields: msg, tokens
func (_m *MockPublisher) PublishTo(msg protoreflect.ProtoMessage, tokens ...string) error {
	_va := make([]interface{}, len(tokens))
	for _i := range tokens {
		_va[_i] = tokens[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PublishTo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(protoreflect.ProtoMessage, ...string) error); ok {
		r0 = rf(msg, tokens...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPublisher_PublishTo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishTo'
type MockPublisher_PublishTo_Call struct {
	*mock.Call
}

// PublishTo is a helper method to define mock.On call
//   - msg protoreflect.ProtoMessage
//   - tokens ...string
func (_e *MockPublisher_Expecter) PublishTo(msg interface{}, tokens ...interface{}) *MockPublisher_PublishTo_Call {
	return &MockPublisher_PublishTo_Call{Call: _e.mock.On("PublishTo",
		append([]interface{}{msg}, tokens...)...)}
}

func (_c *MockPublisher_PublishTo_Call) Run(run func(msg protoreflect.ProtoMessage, tokens ...string)) *MockPublisher_PublishTo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(protoreflect.ProtoMessage), variadicArgs...)
	})
	return _c
}

func (_c *MockPublisher_PublishTo_Call) Return(_a0 error) *MockPublisher_PublishTo_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPublisher_PublishTo_Call) RunAndReturn(run func(protoreflect.ProtoMessage, ...string) error) *MockPublisher_PublishTo_Call {
	_c.Call.Return(run)
	return _c
}

// RequestFrom provides a mock function with given fields: ctx, msg, resp, tokens
func (_m *MockPublisher) RequestFrom(ctx context.Context, msg protoreflect.ProtoMessage, resp protoreflect.ProtoMessage, tokens ...string) (service.Message, error) {
	_va := make([]interface{}, len(tokens))
	for _i := range tokens {
		_va[_i] = tokens[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, msg, resp)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for RequestFrom")
	}

	var r0 service.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, protoreflect.ProtoMessage, protoreflect.ProtoMessage, ...string) (service.Message, error)); ok {
		return rf(ctx, msg, resp, tokens...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, protoreflect.ProtoMessage, protoreflect.ProtoMessage, ...string) service.Message); ok {
		r0 = rf(ctx, msg, resp, tokens...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, protoreflect.ProtoMessage, protoreflect.ProtoMessage, ...string) error); ok {
		r1 = rf(ctx, msg, resp, tokens...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPublisher_RequestFrom_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestFrom'
type MockPublisher_RequestFrom_Call struct {
	*mock.Call
}

// RequestFrom is a helper method to define mock.On call
//   - ctx context.Context
//   - msg protoreflect.ProtoMessage
//   - resp protoreflect.ProtoMessage
//   - tokens ...string
func (_e *MockPublisher_Expecter) RequestFrom(ctx interface{}, msg interface{}, resp interface{}, tokens ...interface{}) *MockPublisher_RequestFrom_Call {
	return &MockPublisher_RequestFrom_Call{Call: _e.mock.On("RequestFrom",
		append([]interface{}{ctx, msg, resp}, tokens...)...)}
}

func (_c *MockPublisher_RequestFrom_Call) Run(run func(ctx context.Context, msg protoreflect.ProtoMessage, resp protoreflect.ProtoMessage, tokens ...string)) *MockPublisher_RequestFrom_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), args[1].(protoreflect.ProtoMessage), args[2].(protoreflect.ProtoMessage), variadicArgs...)
	})
	return _c
}

func (_c *MockPublisher_RequestFrom_Call) Return(_a0 service.Message, _a1 error) *MockPublisher_RequestFrom_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPublisher_RequestFrom_Call) RunAndReturn(run func(context.Context, protoreflect.ProtoMessage, protoreflect.ProtoMessage, ...string) (service.Message, error)) *MockPublisher_RequestFrom_Call {
	_c.Call.Return(run)
	return _c
}

// Serve provides a mock function with given fields: handler, suffixes
func (_m *MockPublisher) Serve(handler service.ServiceHandler, suffixes ...string) (*nats.Subscription, error) {
	_va := make([]interface{}, len(suffixes))
	for _i := range suffixes {
		_va[_i] = suffixes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handler)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Serve")
	}

	var r0 *nats.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(service.ServiceHandler, ...string) (*nats.Subscription, error)); ok {
		return rf(handler, suffixes...)
	}
	if rf, ok := ret.Get(0).(func(service.ServiceHandler, ...string) *nats.Subscription); ok {
		r0 = rf(handler, suffixes...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(service.ServiceHandler, ...string) error); ok {
		r1 = rf(handler, suffixes...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPublisher_Serve_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Serve'
type MockPublisher_Serve_Call struct {
	*mock.Call
}

// Serve is a helper method to define mock.On call
//   - handler service.ServiceHandler
//   - suffixes ...string
func (_e *MockPublisher_Expecter) Serve(handler interface{}, suffixes ...interface{}) *MockPublisher_Serve_Call {
	return &MockPublisher_Serve_Call{Call: _e.mock.On("Serve",
		append([]interface{}{handler}, suffixes...)...)}
}

func (_c *MockPublisher_Serve_Call) Run(run func(handler service.ServiceHandler, suffixes ...string)) *MockPublisher_Serve_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(service.ServiceHandler), variadicArgs...)
	})
	return _c
}

func (_c *MockPublisher_Serve_Call) Return(_a0 *nats.Subscription, _a1 error) *MockPublisher_Serve_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPublisher_Serve_Call) RunAndReturn(run func(service.ServiceHandler, ...string) (*nats.Subscription, error)) *MockPublisher_Serve_Call {
	_c.Call.Return(run)
	return _c
}

// Subject provides a mock function with given fields: suffixes
func (_m *MockPublisher) Subject(suffixes ...string) string {
	_va := make([]interface{}, len(suffixes))
	for _i := range suffixes {
		_va[_i] = suffixes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Subject")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(...string) string); ok {
		r0 = rf(suffixes...)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockPublisher_Subject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Subject'
type MockPublisher_Subject_Call struct {
	*mock.Call
}

// Subject is a helper method to define mock.On call
//   - suffixes ...string
func (_e *MockPublisher_Expecter) Subject(suffixes ...interface{}) *MockPublisher_Subject_Call {
	return &MockPublisher_Subject_Call{Call: _e.mock.On("Subject",
		append([]interface{}{}, suffixes...)...)}
}

func (_c *MockPublisher_Subject_Call) Run(run func(suffixes ...string)) *MockPublisher_Subject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *MockPublisher_Subject_Call) Return(_a0 string) *MockPublisher_Subject_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPublisher_Subject_Call) RunAndReturn(run func(...string) string) *MockPublisher_Subject_Call {
	_c.Call.Return(run)
	return _c
}

// SubscribeTo provides a mock function with given fields: handler, tokens
func (_m *MockPublisher) SubscribeTo(handler service.MessageHandler, tokens ...string) (*nats.Subscription, error) {
	_va := make([]interface{}, len(tokens))
	for _i := range tokens {
		_va[_i] = tokens[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handler)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeTo")
	}

	var r0 *nats.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(service.MessageHandler, ...string) (*nats.Subscription, error)); ok {
		return rf(handler, tokens...)
	}
	if rf, ok := ret.Get(0).(func(service.MessageHandler, ...string) *nats.Subscription); ok {
		r0 = rf(handler, tokens...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(service.MessageHandler, ...string) error); ok {
		r1 = rf(handler, tokens...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPublisher_SubscribeTo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscribeTo'
type MockPublisher_SubscribeTo_Call struct {
	*mock.Call
}

// SubscribeTo is a helper method to define mock.On call
//   - handler service.MessageHandler
//   - tokens ...string
func (_e *MockPublisher_Expecter) SubscribeTo(handler interface{}, tokens ...interface{}) *MockPublisher_SubscribeTo_Call {
	return &MockPublisher_SubscribeTo_Call{Call: _e.mock.On("SubscribeTo",
		append([]interface{}{handler}, tokens...)...)}
}

func (_c *MockPublisher_SubscribeTo_Call) Run(run func(handler service.MessageHandler, tokens ...string)) *MockPublisher_SubscribeTo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(service.MessageHandler), variadicArgs...)
	})
	return _c
}

func (_c *MockPublisher_SubscribeTo_Call) Return(_a0 *nats.Subscription, _a1 error) *MockPublisher_SubscribeTo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPublisher_SubscribeTo_Call) RunAndReturn(run func(service.MessageHandler, ...string) (*nats.Subscription, error)) *MockPublisher_SubscribeTo_Call {
	_c.Call.Return(run)
	return _c
}

// Unmarshal provides a mock function with given fields: nmsg, msg
func (_m *MockPublisher) Unmarshal(nmsg service.Message, msg protoreflect.ProtoMessage) (nats.Header, error) {
	ret := _m.Called(nmsg, msg)

	if len(ret) == 0 {
		panic("no return value specified for Unmarshal")
	}

	var r0 nats.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(service.Message, protoreflect.ProtoMessage) (nats.Header, error)); ok {
		return rf(nmsg, msg)
	}
	if rf, ok := ret.Get(0).(func(service.Message, protoreflect.ProtoMessage) nats.Header); ok {
		r0 = rf(nmsg, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(nats.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(service.Message, protoreflect.ProtoMessage) error); ok {
		r1 = rf(nmsg, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPublisher_Unmarshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unmarshal'
type MockPublisher_Unmarshal_Call struct {
	*mock.Call
}

// Unmarshal is a helper method to define mock.On call
//   - nmsg service.Message
//   - msg protoreflect.ProtoMessage
func (_e *MockPublisher_Expecter) Unmarshal(nmsg interface{}, msg interface{}) *MockPublisher_Unmarshal_Call {
	return &MockPublisher_Unmarshal_Call{Call: _e.mock.On("Unmarshal", nmsg, msg)}
}

func (_c *MockPublisher_Unmarshal_Call) Run(run func(nmsg service.Message, msg protoreflect.ProtoMessage)) *MockPublisher_Unmarshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(service.Message), args[1].(protoreflect.ProtoMessage))
	})
	return _c
}

func (_c *MockPublisher_Unmarshal_Call) Return(_a0 nats.Header, _a1 error) *MockPublisher_Unmarshal_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPublisher_Unmarshal_Call) RunAndReturn(run func(service.Message, protoreflect.ProtoMessage) (nats.Header, error)) *MockPublisher_Unmarshal_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPublisher creates a new instance of MockPublisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPublisher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPublisher {
	mock := &MockPublisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
