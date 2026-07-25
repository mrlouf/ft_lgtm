package ipfs

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ipfs/kubo/client/rpc"
	iface "github.com/ipfs/kubo/core/coreiface"
)

type IPFSPublisher struct {
	node iface.CoreAPI
}

func NewIPFSPublisher(daemonURL string) (*IPFSPublisher, error) {
	node, err := rpc.NewURLApiWithClient(daemonURL, http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("ipfs client: %w", err)
	}
	return &IPFSPublisher{node: node}, nil
}

func (ipfs *IPFSPublisher) Publish(ctx context.Context, data []byte) (string, error) {

	log.Println("publish: start")

	log.Println("publish: done")

	return "", nil

}
