package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"io"
	"log"

	"port-server/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UploadPorts client-streaming method for uploading ports
func (s *PortDomainService) UploadPorts(stream api.PortService_UploadPortsServer) error {
	ctx := stream.Context()
	for {
		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// receive data from stream
			port, err := stream.Recv()
			if err == io.EOF {
				// return will close stream from server side
				return nil
			}
			log.Println(port)
			if err != nil {
				log.Printf("receive error %v", err)
				continue
			}
			s.addPort(port)
		}
	}
}

// GetPort unary method return port by id
func (s *PortDomainService) GetPort(_ context.Context, req *api.PortRequest) (*api.Port, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	log.Printf("Total ports %d\n", len(s.ports))
	port, ok := s.ports[req.PortID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Port %s not found", req.PortID)
	}
	log.Printf("Send port %v\n", port)
	return port, nil
}

// GetPorts server streaming method returns ports
func (s *PortDomainService) GetPorts(_ *empty.Empty, stream api.PortService_GetPortsServer) error {
	s.m.RLock()
	defer s.m.RUnlock()
	log.Printf("Total ports %d\n", len(s.ports))
	for _, port := range s.ports {
		log.Printf("Send port %v\n", port)
		if err := stream.Send(port); err != nil {
			return err
		}
	}
	return nil
}
