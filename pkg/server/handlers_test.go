package server

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"port-server/api"
)

type mockSender struct {
	mock.Mock
	grpc.ServerStream
}

func (m *mockSender) Send(port *api.Port) error {
	args := m.Called(port)
	_, ok := args.Get(0).(*api.Port)
	if !ok {
		return args.Error(0)
	}
	return nil
}

type mockReceiver struct {
	mock.Mock
	grpc.ServerStream
}

func (m *mockReceiver) Recv() (*api.Port, error) {
	args := m.Called()
	val, ok := args.Get(0).(*api.Port)
	if !ok {
		return nil, args.Error(1)
	}
	return val, nil
}

func (m *mockReceiver) Context() context.Context {
	args := m.Called()
	val, ok := args.Get(0).(context.Context)
	if !ok {
		return nil
	}
	return val
}

func (m *mockReceiver) SendAndClose(e *empty.Empty) error {
	args := m.Called(e)
	_, ok := args.Get(0).(*empty.Empty)
	if !ok {
		return args.Error(1)
	}
	return nil
}

func TestGetPort(t *testing.T) {
	port1 := &api.Port{
		Id:   "port-1",
		Name: "Amsterdam",
	}
	srv := PortDomainService{
		ports: map[string]*api.Port{
			port1.Id: port1,
		},
	}
	req := &api.PortRequest{
		PortID: port1.Id,
	}
	port, err := srv.GetPort(nil, req)
	if err != nil {
		t.Errorf("unexpected error when get port %s", err)
	}
	if port.Id != port1.Id {
		t.Errorf("wrong id expected %s actual %s", port1.Id, port.Id)
	}
}

func TestGetPorts(t *testing.T) {
	port1 := &api.Port{
		Id:   "port-1",
		Name: "Amsterdam",
	}
	testCases := []struct {
		caseName string
		port     *api.Port
		err      error
	}{
		{
			caseName: "success",
			port:     port1,
			err:      nil,
		},
		{
			caseName: "send error",
			port:     nil,
			err:      errors.New("send error"),
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.caseName)
		srv := PortDomainService{
			ports: map[string]*api.Port{
				port1.Id: port1,
			},
		}
		sender := &mockSender{}
		sender.On("Send", port1).Return(testCase.err)
		sender.On("SendAndClose", mock.Anything).Return(nil)
		err := srv.GetPorts(nil, sender)
		if testCase.err == nil && err != nil {
			t.Errorf("unexpected error %v", err)
		}
		if testCase.err != nil && err == nil {
			t.Errorf("error must not be nil")
		}
	}
}

func TestUploadPort(t *testing.T) {
	port1 := &api.Port{
		Id:   "port-1",
		Name: "Amsterdam",
	}
	srv := PortDomainService{
		ports: map[string]*api.Port{},
	}
	receiver := &mockReceiver{}
	receiver.On("Recv").Return(port1, nil).Once()
	receiver.On("Recv").Return(nil, io.EOF).Once()
	receiver.On("Context").Return(context.Background())
	err := srv.UploadPorts(receiver)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if _, ok := srv.ports[port1.Id]; !ok {
		t.Errorf("port %s not found", port1.Id)
	}
}

func TestUploadPortCancelled(t *testing.T) {
	port1 := &api.Port{
		Id:   "port-1",
		Name: "Amsterdam",
	}
	srv := PortDomainService{
		ports: map[string]*api.Port{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	receiver := &mockReceiver{}
	receiver.On("Recv").Return(port1, nil).Once()
	receiver.On("Recv").Return(nil, io.EOF).Once()
	receiver.On("Context").Return(ctx)
	cancel()
	err := srv.UploadPorts(receiver)
	if err != context.Canceled {
		t.Errorf("unexpected error %v", err)
	}
}
