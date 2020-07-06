package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"port-server/api"
)

type Receiver interface {
	Recv() (*api.Port, error)
}

// Encodes stream of ports to json
func encodeStream(w io.Writer, stream api.PortService_GetPortsClient) {
	first := true
	w.Write([]byte("{"))
	for {
		port, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("getPorts %v", err)
			continue
		}
		// Don't put coma if this is first object
		if first {
			first = false
		} else {
			w.Write([]byte("\t,"))
		}
		log.Printf("port %v\n", port)
		w.Write([]byte(fmt.Sprintf("\"%s\":", port.Id)))
		if err := json.NewEncoder(w).Encode(port); err != nil {
			log.Printf("error encoding port %v %v", port, err)
		}
	}
	w.Write([]byte("}"))
}
