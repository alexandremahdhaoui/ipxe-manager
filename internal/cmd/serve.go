package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/alexandremahdhaoui/ipxer/internal/util/gracefulshutdown"
	"github.com/alexandremahdhaoui/ipxer/pkg/constants"
)

func Serve(servers map[string]*http.Server, gs *gracefulshutdown.GracefulShutdown) {
	// 1. Run the servers.
	for name, server := range servers {
		ctx := context.WithValue(gs.Context(), constants.ServerNameContextKey, name)

		// sets the base context to be the GracefulShutdown's context.
		server.BaseContext = func(_ net.Listener) context.Context {
			return ctx
		}

		gs.WaitGroup().Add(1)

		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.ErrorContext(ctx, "❌ received error", "error", err)

				// we need to call Done() before requesting the shutdown. Otherwise, the WaitGroup will never decrement.
				gs.WaitGroup().Done()
				gs.Shutdown(1) // Initiate a graceful shutdown. This call is blocking and awaits for wg.

				return
			}

			gs.WaitGroup().Done()

			// The server stopped running without errors, thus we initiate a graceful shutdown if none was previously
			// initiated.
			gs.Shutdown(0)
		}()
	}

	// 2. Await context is done.
	<-gs.Context().Done()

	// 3. Gracefully shutdown each server.
	for name, server := range servers {
		go func() {
			ctx := context.WithValue(context.Background(), constants.ServerNameContextKey, name)

			ctx, cancel := context.WithDeadline(ctx, time.Now().Add(1*time.Minute)) // 1 min deadline.
			defer cancel()

			if err := server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.ErrorContext(ctx, "❌ received error while shutting down server", "error", err)

				return
			}

			slog.Info("✅ gracefully shut down server")
		}()
	}
}
