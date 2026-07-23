package ipfs

import (
	shell "github.com/ipfs/go-ipfs-api"
)

func NewIPFSShell() *shell.Shell {
	return shell.NewShell("localhost:5001")
}
