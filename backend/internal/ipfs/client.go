package ipfs

import (
	"context"
	"lgtm/internal/test"
	"log"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

type IPFS struct {
	Shell *shell.Shell
}

// TODO: determine the localhost port for the IPFS daemon
// TODO: and fetch it from the environment variables or configuration file
func NewIPFSShell() *shell.Shell {
	return shell.NewShell("localhost:5001")
}

func NewIPFSClient() *IPFS {
	return &IPFS{
		Shell: NewIPFSShell(),
	}
}

func (ipfs *IPFS) Publish(ctx context.Context, data []byte) (string, error) {

	log.Println("publish: start")

	if err := test.SleepOrCancel(ctx, 2*time.Second, "publish"); err != nil {
		return "", err
	}

	log.Println("publish: done")

	return "", nil

}
