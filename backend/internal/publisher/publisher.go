package publisher

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ipfs/kubo/client/rpc"
	iface "github.com/ipfs/kubo/core/coreiface"
	"go.opentelemetry.io/otel/trace"
)

type ResponseCID struct {
	Source string `json:"source"`
	Stdout string `json:"stdout"`
}

type IPFSPublisher struct {
	node   iface.CoreAPI
	tracer trace.Tracer
}

func NewIPFSPublisher(t trace.Tracer, daemonURL string) *IPFSPublisher {
	node, err := rpc.NewURLApiWithClient(daemonURL, http.DefaultClient)
	if err != nil {
		return nil
	}

	return &IPFSPublisher{
		node:   node,
		tracer: t,
	}
}

func (i *IPFSPublisher) Publish(ctx context.Context, source []byte, stdout []byte) (ResponseCID, error) {

	log.Println("publish: start")
	ctx, span := i.tracer.Start(ctx, "publisher.publish")
	defer span.End()

	response := ResponseCID{}

	sourceReader := bytes.NewReader(source)
	sourceBlock, err := i.node.Block().Put(ctx, sourceReader)
	if err != nil {
		span.RecordError(err)
		return response, fmt.Errorf("ipfs add source: %w", err)
	}

	outputReader := bytes.NewReader(stdout)
	outputBlock, err := i.node.Block().Put(ctx, outputReader)
	if err != nil {
		span.RecordError(err)
		return response, fmt.Errorf("ipfs add stdout: %w", err)
	}

	response.Source = sourceBlock.Path().String()
	response.Stdout = outputBlock.Path().String()

	response.Source = response.Source[6:] // Remove the "/ipfs/" prefix
	response.Stdout = response.Stdout[6:]

	log.Println("publish: block published:", response.Source)
	log.Println("publish: block published:", response.Stdout)

	log.Println("publish: done")

	return response, nil

}
