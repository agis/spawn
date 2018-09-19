package spawn

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"
)

const envPrefix = "SPAWN_"

var envRe = regexp.MustCompile(`\A` + envPrefix + `[[:xdigit:]]{64}=1\z`)

// Cmd represents an external binary being prepared to run.
type Cmd struct {
	Cmd *exec.Cmd

	fn   func()
	hash string

	mu     sync.Mutex
	sigErr error
}

// New returns the command to execute fn() with the given arguments.
// The command will use os.Stdout and os.Stderr.
func New(fn func(), args ...string) Cmd {
	c := Cmd{}

	c.fn = fn
	c.Cmd = exec.Command(os.Args[0], args...)
	c.Cmd.Stdout = os.Stdout
	c.Cmd.Stderr = os.Stderr

	h := sha256.New()
	h.Write([]byte(os.Args[0] + strings.Join(args, "")))
	c.hash = fmt.Sprintf(envPrefix+"%x", h.Sum(nil))
	c.Cmd.Env = []string{c.hash + "=1"}

	return c
}

// Start starts c. It does not wait for it to complete. When ctx is complete,
// a SIGINT will be sent to c.
//
// The Wait method will return the exit code and release associated resources
// once the command exits.
func (c *Cmd) Start(ctx context.Context) error {
	if os.Getenv(c.hash) != "" {
		c.fn()
		os.Exit(0)
	}

	for _, k := range os.Environ() {
		if envRe.MatchString(k) {
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

// Wait waits for the command to exit and waits for any copying to stdin or
// copying from stdout or stderr to complete.
//
// The command must have been started by Start.
//
// The returned error is nil if the command runs, has no problems copying
// stdin, stdout, and stderr, exits with a zero exit status and any signals
// to the command were delivered successfully.
//
// Wait releases any resources associated with the command.
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
