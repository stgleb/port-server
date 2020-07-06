package client

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"port-server/api"
)

var (
	jsonStream = `
{
  "AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "Abu Dhabi",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  }
}`

	jsonStreamErr = `{`
)

type mockSender struct {
	mock.Mock
	grpc.ClientStream
}

func (m *mockSender) Send(port *api.Port) error {
	args := m.Called(port)
	val, ok := args.Get(0).(error)
	if !ok {
		return args.Error(0)
	}
	return val
}

func (m *mockSender) CloseAndRecv() (*empty.Empty, error) {
	args := m.Called()
	val, ok := args.Get(0).(*empty.Empty)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestUploadFile(t *testing.T) {
	testCases := []struct {
		expectedErr bool
		sendErr     error
		closeErr    error
	}{
		{
			expectedErr: false,
		},
		{
			expectedErr: true,
			sendErr:     errors.New("send error"),
		},
		{
			expectedErr: true,
			closeErr:    errors.New("close error"),
		},
	}

	for _, testCase := range testCases {
		sender := &mockSender{}
		reader := strings.NewReader(jsonStream)
		sender.On("Send", mock.Anything).Return(testCase.sendErr)
		sender.On("Send", mock.Anything).Return(testCase.sendErr)
		sender.On("CloseAndRecv").Return(mock.Anything, testCase.closeErr)

		err := uploadFile(context.Background(), sender, reader)

		if !testCase.expectedErr && err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if testCase.expectedErr && err == nil {
			t.Errorf("error must not be  nil")
		}
	}
}

func TestUploadFileJsonErr(t *testing.T) {
	sender := &mockSender{}
	reader := strings.NewReader(jsonStreamErr)
	err := uploadFile(context.Background(), sender, reader)

	if err == nil {
		t.Errorf("error expected")
	}
}
