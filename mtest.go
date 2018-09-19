package mtest

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type Cmd struct {
	Cmd *exec.Cmd

	mainfn func()
	hash   string

	mu     sync.Mutex
	sigErr error
}

func New(m func(), args ...string) Cmd {
	c := Cmd{}

	c.mainfn = m
	c.Cmd = exec.Command(os.Args[0], args...)
	c.Cmd.Stdout = os.Stdout
	c.Cmd.Stderr = os.Stderr

	h := sha256.New()
	h.Write([]byte(os.Args[0] + strings.Join(args, "")))
	c.hash = fmt.Sprintf("FOO_%x", h.Sum(nil))
	c.Cmd.Env = []string{c.hash + "=1"}

	return c
}

func (c *Cmd) Start(ctx context.Context) error {
	if os.Getenv(c.hash) != "" {
		c.mainfn()
		os.Exit(0)
	}

	for _, k := range os.Environ() {
		if strings.HasPrefix(k, "FOO_") {
			return nil
		}
	}

	err := c.Cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		c.mu.Lock()
		defer c.mu.Unlock()
		c.sigErr = c.Cmd.Process.Signal(syscall.SIGINT)
	}()

	return nil
}

func (c *Cmd) Wait() error {
	err := c.Cmd.Wait()
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sigErr != nil {
		return errors.New(c.sigErr.Error())
	}

	return nil
}
