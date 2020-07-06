package client

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"port-server/api"
)

type mockReceiver struct {
	mock.Mock
	grpc.ClientStream
}

func (m *mockReceiver) Recv() (*api.Port, error) {
	args := m.Called()
	val, ok := args.Get(0).(*api.Port)
	if !ok {
		return nil, args.Error(1)
	}
	return val, nil
}

// Verify that stream is encoded in valid json
func TestEncodeStream(t *testing.T) {
	ports := []*api.Port{
		{
			Id:   "id-1",
			Name: "name-1",
			City: "Amsterdam",
		},
		{
			Id:   "id-2",
			Name: "name-2",
			City: "Roterdam",
		},
	}
	receiver := &mockReceiver{}
	for index := range ports {
		receiver.On("Recv").Return(ports[index], nil).Once()
	}
	receiver.On("Recv").Return(nil, io.EOF)
	writer := &bytes.Buffer{}
	encodeStream(writer, receiver)
	output := map[string]*api.Port{}
	err := json.Unmarshal(writer.Bytes(), &output)
	if err != nil {
		t.Errorf("unmarshalling error %v", err)
	}
}
