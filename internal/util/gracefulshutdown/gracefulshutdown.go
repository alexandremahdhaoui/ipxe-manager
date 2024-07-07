package gracefulshutdown

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type GracefulShutdown struct {
	ctx    context.Context
	cancel context.CancelFunc
	name   string

	mu sync.Mutex
	wg *sync.WaitGroup
}

// New creates a new GracefulShutdown struct initializing a sync.WaitGroup and a new context.Context cancelable by a
// CancelFunc, a SIGTERM, SIGINT or SIGKILL.
func New(name string) *GracefulShutdown {
	// 1. initialize a new cancelable context.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)

	// 2. initialize a new wait group.
	wg := &sync.WaitGroup{}

	// 3. create the GracefulShutdown struct.
	gs := &GracefulShutdown{
		ctx:    ctx,
		cancel: cancel,
		name:   name,
		wg:     wg,
	}

	// 4. Ensure gs.Shutdown is always called at least once when the context is done.
	go func() {
		<-ctx.Done()
		gs.Shutdown(0)
	}()

	return gs
}

func (s *GracefulShutdown) Shutdown(exitCode int) {
	// 1. Try to lock the GracefulShutdown struct. This oneliner ensures Shutdown is idempotent.
	if !s.mu.TryLock() {
		return
	}

	defer s.mu.Unlock() // NB: there isn't really a point to release the lock.

	// 2. Print a log line.
	slog.InfoContext(s.ctx, fmt.Sprintf("⌛ gracefully shutting down %s", s.name))

	// 3. Cancel the context.
	s.cancel()

	// 4. Wait until all goroutines which incremented the wait group are done.
	s.wg.Wait()

	// 5. Exit.
	os.Exit(exitCode)
}

func (s *GracefulShutdown) Context() context.Context {
	return s.ctx
}

func (s *GracefulShutdown) CancelFunc() context.CancelFunc {
	return s.cancel
}

func (s *GracefulShutdown) WaitGroup() *sync.WaitGroup {
	return s.wg
}
