// Package spawn makes it easy to end-to-end test Go servers. The main idea is
// that you spin up your server in your TestMain(), use it throughout your tests
// and shut it down at the end.
//
// Refer to the examples directory for usage information.
package spawn

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// envVar is the environment variable set on the spawned Cmd and is used
// to hijack TestMain and execute main() instead.
const envVar = "GOSPAWN_EXEC_MAIN"

// Cmd wraps exec.Cmd and represents a binary being prepared or run.
//
// In the typical end-to-end testing scenario, Cmd will end up running
// two times:
//
//  1. from TestMain when the test suite is first executed. At this
//     point it will spawn the already-compiled test binary (itself) again
//     and...
//  2. from the spawned binary, inside TestMain again. But this time it will
//     intercept TestMain and will execute main() instead (i.e. the actual program)
type Cmd struct {
	Cmd  *exec.Cmd
	main func()

	mu     sync.Mutex
	sigErr error
}

// New returns a Cmd that will either execute the passed in main() function, or the parent binary with
// the given arguments. The program's main() function should be passed as f.
func New(main func(), args ...string) *Cmd {
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), envVar+"=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Cmd{
		main: main,
		Cmd:  cmd,
	}
}

// Start starts c until it terminates or ctx is cancelled. It does not wait
// for it to complete. When ctx is cancelled a SIGINT is sent to c.
//
// The Wait method will return the exit code and release associated resources
// once the command exits.
func (c *Cmd) Start(ctx context.Context) error {
	if os.Getenv(envVar) == "1" {
		c.main()

		// we don't want to continue executing in TestMain()
		os.Exit(0)
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
