package publisher

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ipfs/kubo/client/rpc"
	iface "github.com/ipfs/kubo/core/coreiface"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ResponseCID struct {
	Source string `json:"source"`
	Stdout string `json:"stdout"`
}

type IPFSPublisher struct {
	node   iface.CoreAPI
	tracer trace.Tracer
	logger *slog.Logger
}

func NewIPFSPublisher(t trace.Tracer, l *slog.Logger, daemonURL string) *IPFSPublisher {
	node, err := rpc.NewURLApiWithClient(daemonURL, http.DefaultClient)
	if err != nil {
		return nil
	}

	return &IPFSPublisher{
		node:   node,
		tracer: t,
		logger: l,
	}
}

func (i *IPFSPublisher) Publish(ctx context.Context, source []byte, stdout []byte) (ResponseCID, error) {

	i.logger.InfoContext(ctx, "publish: start")
	ctx, span := i.tracer.Start(ctx, "publisher.publish")
	defer span.End()

	response := ResponseCID{}

	sourceReader := bytes.NewReader(source)
	sourceBlock, err := i.node.Block().Put(ctx, sourceReader)
	if err != nil {
		span.RecordError(err)
		i.logger.ErrorContext(ctx, "publish: error adding source block", "error", err)
		return response, fmt.Errorf("ipfs add source: %w", err)
	}
	response.Source = sourceBlock.Path().String()
	response.Source = response.Source[6:] // Remove the "/ipfs/" prefix
	span.AddEvent("source block published", trace.WithAttributes(
		attribute.String("source_cid", response.Source),
	))
	i.logger.InfoContext(ctx, "publish: block published", "source", response.Source)

	outputReader := bytes.NewReader(stdout)
	outputBlock, err := i.node.Block().Put(ctx, outputReader)
	if err != nil {
		span.RecordError(err)
		i.logger.ErrorContext(ctx, "publish: error adding stdout block", "error", err)
		return response, fmt.Errorf("ipfs add stdout: %w", err)
	}
	response.Stdout = outputBlock.Path().String()
	response.Stdout = response.Stdout[6:]
	span.AddEvent("stdout block published", trace.WithAttributes(
		attribute.String("stdout_cid", response.Stdout),
	))
	i.logger.InfoContext(ctx, "publish: block published", "stdout", response.Stdout)

	span.AddEvent("publish: done")
	i.logger.InfoContext(ctx, "publish: done")

	return response, nil

}
