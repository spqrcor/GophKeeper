package main

import (
	"syscall"
	"testing"
	"time"
)

func TestMains(t *testing.T) {
	go func() {
		time.Sleep(1 * time.Second)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	main()
}
