package ipfs

import (
	"context"
	"log"

	"github.com/ipfs/boxo/gateway"
	shell "github.com/ipfs/go-ipfs-api"
)

type IPFS struct {
	Shell         *shell.Shell
	GatewayConfig *gateway.Config
	Backend       *gateway.IPFSBackend
}

// TODO: determine the localhost port for the IPFS daemon
// TODO: and fetch it from the environment variables or configuration file
func NewIPFSShell() *shell.Shell {
	return shell.NewShell("localhost:5001")
}

func NewIPFSGatewayConfig() *gateway.Config {
	return &gateway.Config{}
}

func NewIPFSGateway() *IPFS {
	return &IPFS{
		Shell:         NewIPFSShell(),
		GatewayConfig: NewIPFSGatewayConfig(),
		Backend:       gateway.NewIPFSBackend(),
	}
}

func (ipfs *IPFS) Publish(ctx context.Context, data []byte) (string, error) {

	log.Println("publish: start")

	log.Println("publish: done")

	return "", nil

}
