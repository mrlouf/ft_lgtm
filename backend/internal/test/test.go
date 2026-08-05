package test

import (
	"context"
	"fmt"
	"time"
)

// sleepOrCancel sleeps for the specified duration or returns early if the context is canceled.
// This is useful to simulate long-running operations in the different stages
// of the backend pipeline (compile, execute, publish) and test cancellation.
func SleepOrCancel(ctx context.Context, d time.Duration, stage string) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: interrupted: %w", stage, ctx.Err())
	}
}

func ReturnTestsSources() map[string]string {
	return map[string]string{
		"fs_read": `package main

import "os"

func main() {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		println("BLOCKED:", err.Error())
		return
	}
	println("LEAK:", string(data))
}`,
		"path_traversal": `package main

import "os"

func main() {
	data, err := os.ReadFile("/tmp/../../etc/passwd")
	if err != nil {
		println("BLOCKED:", err.Error())
		return
	}
	println("LEAK:", string(data))
}`,
		"network_access": `package main

import (
	"net/http"
	"io/ioutil"
)

func main() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		println("BLOCKED:", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("BLOCKED:", err.Error())
		return
	}
	println("LEAK:", string(body))
}`,
		"fs_write": `package main

import "os"

func main() {
	err := os.WriteFile("/tmp/testfile", []byte("test"), 0644)
	if err != nil {
		println("BLOCKED:", err.Error())
		return
	}
	println("LEAK: wrote to /tmp/testfile")
}`,
		"infinite_loop": `package main

func main() {
	for {
		// Infinite loop
	}
}`,
	}
}
