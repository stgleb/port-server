package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"port-server/api"
)

// uploadFile sends  file to server
func uploadFile(ctx context.Context, client api.PortService_UploadPortsClient, reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	// We expect an object
	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected object")
	}
	// Read props
	for decoder.More() {
		// Read items (large objects)
		for decoder.More() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Skip key
				t, _ := decoder.Token()
				portKey, ok := t.(string)
				if !ok {
					return fmt.Errorf("error converting port key to string")
				}
				// Read next item (large object)
				port := &api.Port{}
				if err := decoder.Decode(port); err != nil {
					return err
				}
				port.Id = portKey
				if err := client.Send(port); err != nil {
					if err == io.EOF {
						break
					}
					return fmt.Errorf("error send port %v", port)
				}
			}
		}
	}
	// Object closing delim
	t, err = decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		log.Fatal("expected object closing")
	}
	// Close stream
	_, err = client.CloseAndRecv()
	if err != nil && err != io.EOF {
		return fmt.Errorf("error closing %v", err)
	}
	log.Println("File has been uploaded")
	return nil
}
