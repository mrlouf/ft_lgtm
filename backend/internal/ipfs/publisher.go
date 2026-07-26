package ipfs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ipfs/kubo/client/rpc"
	iface "github.com/ipfs/kubo/core/coreiface"
)

type ResponseCID struct {
	Source string `json:"source"`
	Stdout string `json:"stdout"`
}

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

func (i *IPFSPublisher) Publish(ctx context.Context, source []byte, stdout []byte) (ResponseCID, error) {

	log.Println("publish: start")

	response := ResponseCID{}

	sourceReader := bytes.NewReader(source)
	sourceBlock, err := i.node.Block().Put(ctx, sourceReader)
	if err != nil {
		return response, fmt.Errorf("ipfs add source: %w", err)
	}

	outputReader := bytes.NewReader(stdout)
	outputBlock, err := i.node.Block().Put(ctx, outputReader)
	if err != nil {
		return response, fmt.Errorf("ipfs add stdout: %w", err)
	}

	response.Source = sourceBlock.Path().String()
	response.Stdout = outputBlock.Path().String()

	log.Println("publish: block published:", response.Source)
	log.Println("publish: block published:", response.Stdout)

	log.Println("publish: done")

	return response, nil

}
