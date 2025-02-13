// Code generated by mockery v2.42.3. DO NOT EDIT.

package service

import (
	nats "github.com/nats-io/nats.go"
	mock "github.com/stretchr/testify/mock"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	time "time"
)

// MockMessage is an autogenerated mock type for the Message type
type MockMessage struct {
	mock.Mock
}

type MockMessage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMessage) EXPECT() *MockMessage_Expecter {
	return &MockMessage_Expecter{mock: &_m.Mock}
}

// Ack provides a mock function with given fields: opts
func (_m *MockMessage) Ack(opts ...nats.AckOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Ack")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...nats.AckOpt) error); ok {
		r0 = rf(opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessage_Ack_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ack'
type MockMessage_Ack_Call struct {
	*mock.Call
}

// Ack is a helper method to define mock.On call
//   - opts ...nats.AckOpt
func (_e *MockMessage_Expecter) Ack(opts ...interface{}) *MockMessage_Ack_Call {
	return &MockMessage_Ack_Call{Call: _e.mock.On("Ack",
		append([]interface{}{}, opts...)...)}
}

func (_c *MockMessage_Ack_Call) Run(run func(opts ...nats.AckOpt)) *MockMessage_Ack_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.AckOpt, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(nats.AckOpt)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *MockMessage_Ack_Call) Return(_a0 error) *MockMessage_Ack_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Ack_Call) RunAndReturn(run func(...nats.AckOpt) error) *MockMessage_Ack_Call {
	_c.Call.Return(run)
	return _c
}

// AckSync provides a mock function with given fields: opts
func (_m *MockMessage) AckSync(opts ...nats.AckOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AckSync")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...nats.AckOpt) error); ok {
		r0 = rf(opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessage_AckSync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AckSync'
type MockMessage_AckSync_Call struct {
	*mock.Call
}

// AckSync is a helper method to define mock.On call
//   - opts ...nats.AckOpt
func (_e *MockMessage_Expecter) AckSync(opts ...interface{}) *MockMessage_AckSync_Call {
	return &MockMessage_AckSync_Call{Call: _e.mock.On("AckSync",
		append([]interface{}{}, opts...)...)}
}

func (_c *MockMessage_AckSync_Call) Run(run func(opts ...nats.AckOpt)) *MockMessage_AckSync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.AckOpt, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(nats.AckOpt)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *MockMessage_AckSync_Call) Return(_a0 error) *MockMessage_AckSync_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_AckSync_Call) RunAndReturn(run func(...nats.AckOpt) error) *MockMessage_AckSync_Call {
	_c.Call.Return(run)
	return _c
}

// Data provides a mock function with given fields:
func (_m *MockMessage) Data() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Data")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// MockMessage_Data_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Data'
type MockMessage_Data_Call struct {
	*mock.Call
}

// Data is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Data() *MockMessage_Data_Call {
	return &MockMessage_Data_Call{Call: _e.mock.On("Data")}
}

func (_c *MockMessage_Data_Call) Run(run func()) *MockMessage_Data_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Data_Call) Return(_a0 []byte) *MockMessage_Data_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Data_Call) RunAndReturn(run func() []byte) *MockMessage_Data_Call {
	_c.Call.Return(run)
	return _c
}

// Equal provides a mock function with given fields: msg
func (_m *MockMessage) Equal(msg Message) bool {
	ret := _m.Called(msg)

	if len(ret) == 0 {
		panic("no return value specified for Equal")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(Message) bool); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockMessage_Equal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Equal'
type MockMessage_Equal_Call struct {
	*mock.Call
}

// Equal is a helper method to define mock.On call
//   - msg Message
func (_e *MockMessage_Expecter) Equal(msg interface{}) *MockMessage_Equal_Call {
	return &MockMessage_Equal_Call{Call: _e.mock.On("Equal", msg)}
}

func (_c *MockMessage_Equal_Call) Run(run func(msg Message)) *MockMessage_Equal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Message))
	})
	return _c
}

func (_c *MockMessage_Equal_Call) Return(_a0 bool) *MockMessage_Equal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Equal_Call) RunAndReturn(run func(Message) bool) *MockMessage_Equal_Call {
	_c.Call.Return(run)
	return _c
}

// Header provides a mock function with given fields:
func (_m *MockMessage) Header() nats.Header {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Header")
	}

	var r0 nats.Header
	if rf, ok := ret.Get(0).(func() nats.Header); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(nats.Header)
		}
	}

	return r0
}

// MockMessage_Header_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Header'
type MockMessage_Header_Call struct {
	*mock.Call
}

// Header is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Header() *MockMessage_Header_Call {
	return &MockMessage_Header_Call{Call: _e.mock.On("Header")}
}

func (_c *MockMessage_Header_Call) Run(run func()) *MockMessage_Header_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Header_Call) Return(_a0 nats.Header) *MockMessage_Header_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Header_Call) RunAndReturn(run func() nats.Header) *MockMessage_Header_Call {
	_c.Call.Return(run)
	return _c
}

// InProgress provides a mock function with given fields: opts
func (_m *MockMessage) InProgress(opts ...nats.AckOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for InProgress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...nats.AckOpt) error); ok {
		r0 = rf(opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessage_InProgress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InProgress'
type MockMessage_InProgress_Call struct {
	*mock.Call
}

// InProgress is a helper method to define mock.On call
//   - opts ...nats.AckOpt
func (_e *MockMessage_Expecter) InProgress(opts ...interface{}) *MockMessage_InProgress_Call {
	return &MockMessage_InProgress_Call{Call: _e.mock.On("InProgress",
		append([]interface{}{}, opts...)...)}
}

func (_c *MockMessage_InProgress_Call) Run(run func(opts ...nats.AckOpt)) *MockMessage_InProgress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.AckOpt, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(nats.AckOpt)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *MockMessage_InProgress_Call) Return(_a0 error) *MockMessage_InProgress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_InProgress_Call) RunAndReturn(run func(...nats.AckOpt) error) *MockMessage_InProgress_Call {
	_c.Call.Return(run)
	return _c
}

// Message provides a mock function with given fields:
func (_m *MockMessage) Message() *nats.Msg {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Message")
	}

	var r0 *nats.Msg
	if rf, ok := ret.Get(0).(func() *nats.Msg); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.Msg)
		}
	}

	return r0
}

// MockMessage_Message_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Message'
type MockMessage_Message_Call struct {
	*mock.Call
}

// Message is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Message() *MockMessage_Message_Call {
	return &MockMessage_Message_Call{Call: _e.mock.On("Message")}
}

func (_c *MockMessage_Message_Call) Run(run func()) *MockMessage_Message_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Message_Call) Return(_a0 *nats.Msg) *MockMessage_Message_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Message_Call) RunAndReturn(run func() *nats.Msg) *MockMessage_Message_Call {
	_c.Call.Return(run)
	return _c
}

// Metadata provides a mock function with given fields:
func (_m *MockMessage) Metadata() (*nats.MsgMetadata, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Metadata")
	}

	var r0 *nats.MsgMetadata
	var r1 error
	if rf, ok := ret.Get(0).(func() (*nats.MsgMetadata, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *nats.MsgMetadata); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.MsgMetadata)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMessage_Metadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Metadata'
type MockMessage_Metadata_Call struct {
	*mock.Call
}

// Metadata is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Metadata() *MockMessage_Metadata_Call {
	return &MockMessage_Metadata_Call{Call: _e.mock.On("Metadata")}
}

func (_c *MockMessage_Metadata_Call) Run(run func()) *MockMessage_Metadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Metadata_Call) Return(_a0 *nats.MsgMetadata, _a1 error) *MockMessage_Metadata_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMessage_Metadata_Call) RunAndReturn(run func() (*nats.MsgMetadata, error)) *MockMessage_Metadata_Call {
	_c.Call.Return(run)
	return _c
}

// Nak provides a mock function with given fields: opts
func (_m *MockMessage) Nak(opts ...nats.AckOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Nak")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...nats.AckOpt) error); ok {
		r0 = rf(opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessage_Nak_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Nak'
type MockMessage_Nak_Call struct {
	*mock.Call
}

// Nak is a helper method to define mock.On call
//   - opts ...nats.AckOpt
func (_e *MockMessage_Expecter) Nak(opts ...interface{}) *MockMessage_Nak_Call {
	return &MockMessage_Nak_Call{Call: _e.mock.On("Nak",
		append([]interface{}{}, opts...)...)}
}

func (_c *MockMessage_Nak_Call) Run(run func(opts ...nats.AckOpt)) *MockMessage_Nak_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.AckOpt, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(nats.AckOpt)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *MockMessage_Nak_Call) Return(_a0 error) *MockMessage_Nak_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Nak_Call) RunAndReturn(run func(...nats.AckOpt) error) *MockMessage_Nak_Call {
	_c.Call.Return(run)
	return _c
}

// NakWithDelay provides a mock function with given fields: delay, opts
func (_m *MockMessage) NakWithDelay(delay time.Duration, opts ...nats.AckOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, delay)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for NakWithDelay")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Duration, ...nats.AckOpt) error); ok {
		r0 = rf(delay, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessage_NakWithDelay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NakWithDelay'
type MockMessage_NakWithDelay_Call struct {
	*mock.Call
}

// NakWithDelay is a helper method to define mock.On call
//   - delay time.Duration
//   - opts ...nats.AckOpt
func (_e *MockMessage_Expecter) NakWithDelay(delay interface{}, opts ...interface{}) *MockMessage_NakWithDelay_Call {
	return &MockMessage_NakWithDelay_Call{Call: _e.mock.On("NakWithDelay",
		append([]interface{}{delay}, opts...)...)}
}

func (_c *MockMessage_NakWithDelay_Call) Run(run func(delay time.Duration, opts ...nats.AckOpt)) *MockMessage_NakWithDelay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.AckOpt, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(nats.AckOpt)
			}
		}
		run(args[0].(time.Duration), variadicArgs...)
	})
	return _c
}

func (_c *MockMessage_NakWithDelay_Call) Return(_a0 error) *MockMessage_NakWithDelay_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_NakWithDelay_Call) RunAndReturn(run func(time.Duration, ...nats.AckOpt) error) *MockMessage_NakWithDelay_Call {
	_c.Call.Return(run)
	return _c
}

// QueueName provides a mock function with given fields:
func (_m *MockMessage) QueueName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for QueueName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockMessage_QueueName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueueName'
type MockMessage_QueueName_Call struct {
	*mock.Call
}

// QueueName is a helper method to define mock.On call
func (_e *MockMessage_Expecter) QueueName() *MockMessage_QueueName_Call {
	return &MockMessage_QueueName_Call{Call: _e.mock.On("QueueName")}
}

func (_c *MockMessage_QueueName_Call) Run(run func()) *MockMessage_QueueName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_QueueName_Call) Return(_a0 string) *MockMessage_QueueName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_QueueName_Call) RunAndReturn(run func() string) *MockMessage_QueueName_Call {
	_c.Call.Return(run)
	return _c
}

// Reply provides a mock function with given fields:
func (_m *MockMessage) Reply() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Reply")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockMessage_Reply_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Reply'
type MockMessage_Reply_Call struct {
	*mock.Call
}

// Reply is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Reply() *MockMessage_Reply_Call {
	return &MockMessage_Reply_Call{Call: _e.mock.On("Reply")}
}

func (_c *MockMessage_Reply_Call) Run(run func()) *MockMessage_Reply_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Reply_Call) Return(_a0 string) *MockMessage_Reply_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Reply_Call) RunAndReturn(run func() string) *MockMessage_Reply_Call {
	_c.Call.Return(run)
	return _c
}

// Respond provides a mock function with given fields: _a0
func (_m *MockMessage) Respond(_a0 protoreflect.ProtoMessage) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Respond")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(protoreflect.ProtoMessage) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessage_Respond_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Respond'
type MockMessage_Respond_Call struct {
	*mock.Call
}

// Respond is a helper method to define mock.On call
//   - _a0 protoreflect.ProtoMessage
func (_e *MockMessage_Expecter) Respond(_a0 interface{}) *MockMessage_Respond_Call {
	return &MockMessage_Respond_Call{Call: _e.mock.On("Respond", _a0)}
}

func (_c *MockMessage_Respond_Call) Run(run func(_a0 protoreflect.ProtoMessage)) *MockMessage_Respond_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(protoreflect.ProtoMessage))
	})
	return _c
}

func (_c *MockMessage_Respond_Call) Return(_a0 error) *MockMessage_Respond_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Respond_Call) RunAndReturn(run func(protoreflect.ProtoMessage) error) *MockMessage_Respond_Call {
	_c.Call.Return(run)
	return _c
}

// Subject provides a mock function with given fields:
func (_m *MockMessage) Subject() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Subject")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockMessage_Subject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Subject'
type MockMessage_Subject_Call struct {
	*mock.Call
}

// Subject is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Subject() *MockMessage_Subject_Call {
	return &MockMessage_Subject_Call{Call: _e.mock.On("Subject")}
}

func (_c *MockMessage_Subject_Call) Run(run func()) *MockMessage_Subject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Subject_Call) Return(_a0 string) *MockMessage_Subject_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Subject_Call) RunAndReturn(run func() string) *MockMessage_Subject_Call {
	_c.Call.Return(run)
	return _c
}

// Term provides a mock function with given fields: opts
func (_m *MockMessage) Term(opts ...nats.AckOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Term")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...nats.AckOpt) error); ok {
		r0 = rf(opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessage_Term_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Term'
type MockMessage_Term_Call struct {
	*mock.Call
}

// Term is a helper method to define mock.On call
//   - opts ...nats.AckOpt
func (_e *MockMessage_Expecter) Term(opts ...interface{}) *MockMessage_Term_Call {
	return &MockMessage_Term_Call{Call: _e.mock.On("Term",
		append([]interface{}{}, opts...)...)}
}

func (_c *MockMessage_Term_Call) Run(run func(opts ...nats.AckOpt)) *MockMessage_Term_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.AckOpt, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(nats.AckOpt)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *MockMessage_Term_Call) Return(_a0 error) *MockMessage_Term_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Term_Call) RunAndReturn(run func(...nats.AckOpt) error) *MockMessage_Term_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMessage creates a new instance of MockMessage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMessage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMessage {
	mock := &MockMessage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
