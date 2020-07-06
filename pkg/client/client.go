package client

import (
	"context"
	"google.golang.org/grpc/keepalive"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"

	"port-server/api"
)

type Client struct {
	serverUrl string
	clientUrl string
	filePath  string

	httpServer *http.Server
	client     api.PortServiceClient

	getReader func(string) (io.ReadCloser, error)
	loadFile  func(context.Context, api.PortService_UploadPortsClient, io.Reader) error
}

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

// NewClient creates new client object pointing to server
func NewClient(serverUrl, clientUrl, filePath string) (*Client, error) {
	// Create http server for rest API
	portClient := &Client{
		serverUrl: serverUrl,
		clientUrl: clientUrl,
		filePath:  filePath,
		getReader: func(fileName string) (io.ReadCloser, error) {
			f, err := os.Open(fileName)
			if err != nil {
				return nil, err
			}
			return f, nil
		},
		loadFile: uploadFile,
	}

	router := mux.NewRouter()
	router.HandleFunc("/port", portClient.getPorts)
	router.HandleFunc("/port/{id}", portClient.getPort)
	server := &http.Server{
		Addr:    clientUrl,
		Handler: router,
	}
	portClient.httpServer = server
	return portClient, nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	return c.httpServer.Shutdown(ctx)
}

// Start launches HTTP server for client and start uploading file to PortDomainService
func (c *Client) Start(ctx context.Context) (err error) {
	// Establish connection to server
	conn, err := grpc.Dial(c.serverUrl, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return err
	}
	defer func() {
		cerr := conn.Close()
		if err == nil {
			err = cerr
		}
	}()
	c.client = api.NewPortServiceClient(conn)
	go func() {
		// Upload file to server
		log.Printf("upload file %s", c.filePath)
		if err := c.uploadFile(ctx, c.filePath); err != nil {
			log.Printf("error uploading file %v", err)
		}
	}()
	// Start client server
	log.Printf("Start client  on %s pointing to server %s", c.clientUrl, c.serverUrl)
	return c.httpServer.ListenAndServe()
}

// uploadFile sends json file with ports to server
func (c *Client) uploadFile(ctx context.Context, fileName string) (err error) {
	var (
		sender api.PortService_UploadPortsClient
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	if sender, err = c.client.UploadPorts(ctx); err != nil {
		return err
	}
	reader, err := c.getReader(fileName)
	if err != nil {
		return err
	}
	defer func() {
		cerr := reader.Close()
		if err == nil {
			err = cerr
		}
	}()
	if err := c.loadFile(ctx, sender, reader); err != nil {
		return err
	}
	return err
}
