package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"port-server/api"
)

type PortDomainService struct {
	protocol string
	addr     string
	// Storage for ports
	m     sync.RWMutex
	ports map[string]*api.Port
}

var (
	maxSendSize = 5 * 1024 * 1024 * 1024 * 1024
	maxRecvSize = 5 * 1024 * 1024 * 1024 * 1024
	kaep        = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
	kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     30 * time.Second, // If a client is idle for 30 seconds, send a GOAWAY
		MaxConnectionAge:      5 * time.Minute,  // If any connection is alive for more than 5 minutes, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
)

// NewPortDomainService creates new port server
func NewPortDomainService(protocol, addr string) (*PortDomainService, error) {
	return &PortDomainService{
		protocol: protocol,
		addr:     addr,
		ports:    map[string]*api.Port{},
	}, nil
}

// Start is used to start the server.
// It will return an error if the server is already running.
func (s *PortDomainService) Start() error {
	l, err := net.Listen(s.protocol, s.addr)
	if err != nil {
		return fmt.Errorf("unable to setup server: %s", err.Error())
	}
	// Start the gRPC Server
	grpcServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.MaxSendMsgSize(maxSendSize),
		grpc.MaxRecvMsgSize(maxRecvSize))
	api.RegisterPortServiceServer(grpcServer, s)
	log.Printf("Start server on %s", s.addr)
	return grpcServer.Serve(l)
}

// Save port to storage
func (s *PortDomainService) addPort(port *api.Port) {
	s.m.Lock()
	defer s.m.Unlock()
	// update max and send it to stream
	s.ports[port.Id] = port
}
