package client

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/gorilla/mux"
	"port-server/api"
)

// get port returns port by id
func (c *Client) getPort(w http.ResponseWriter, r *http.Request) {
	portID := mux.Vars(r)["id"]
	req := &api.PortRequest{
		PortID: portID,
	}
	log.Printf("portReq %v\n", req)
	port, err := c.client.GetPort(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("get port %v\n", port)
	if err := json.NewEncoder(w).Encode(port); err != nil {
		log.Printf("error encoding port %v %v", port, err)
	}
}

// get ports return ports from server
func (c *Client) getPorts(w http.ResponseWriter, r *http.Request) {
	stream, err := c.client.GetPorts(r.Context(), &empty.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encodeStream(w, stream)
}
