package main

import (
	"flag"
	"fmt"
	"log"
	"port-server/pkg/server"
)

var (
	host string
	port int
)

func main() {
	flag.IntVar(&port, "port", 9090, "Port of domain service")
	flag.StringVar(&host, "host", "localhost", "Server hostname")
	flag.Parse()

	srv, err := server.NewPortDomainService("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("error creating server %v", err)
	}
	log.Fatal(srv.Start())
}
