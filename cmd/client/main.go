package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"port-server/pkg/client"
	"syscall"
	"time"
)

var (
	clientPort int
	clientHost string
	serverPort int
	serverHost string
	filePath   string
)

func main() {
	flag.IntVar(&clientPort, "clientPort", 9091, "Port of client service")
	flag.StringVar(&clientHost, "clientHost", "localhost", "client hostname")

	flag.IntVar(&serverPort, "serverPort", 9090, "Port of domain service")
	flag.StringVar(&serverHost, "serverHost", "localhost", "Server hostname")
	flag.StringVar(&filePath, "filePath", "./data/ports.json", "path to files containing ports")
	flag.Parse()

	portClient, err := client.NewClient(fmt.Sprintf("%s:%d", serverHost, serverPort),
		fmt.Sprintf("%s:%d", clientHost, clientPort), filePath)

	if err != nil {
		log.Fatalf("error creating client %v", err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	// Root context for terminating outgoing grpc stream
	rootCtx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-c
		log.Printf("Signal %v received shutdown", sig)
		cancel()
		// Context with 10 second timeout to let http server process pending requests
		ctx, cancelShutDown := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelShutDown()
		err := portClient.Shutdown(ctx)
		if err != nil {
			log.Printf("error when shutdown %v", err)
		}
	}()
	log.Fatal(portClient.Start(rootCtx))
}
