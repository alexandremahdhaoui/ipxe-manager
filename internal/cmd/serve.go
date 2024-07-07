package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/internal/util/gracefulshutdown"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Serve(portServerMap map[int]*echo.Echo, gs *gracefulshutdown.GracefulShutdown) {
	// 1. Run the servers.
	for port, server := range portServerMap {
		port, server := port, server
		addr := fmt.Sprintf(":%d", port)

		// sets the base context to be the GracefulShutdown's context.
		server.Server.BaseContext = func(_ net.Listener) context.Context {
			return gs.Context()
		}

		go func() {
			gs.WaitGroup().Add(1)

			if err := server.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
				server.Logger.Errorf("❌ received error: %s", err.Error())
				gs.WaitGroup().Done()
				gs.Shutdown(1) // Initiate a graceful shutdown.
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
	for _, server := range portServerMap {
		server := server

		go func() {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
			defer cancel()

			if err := server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				server.Logger.Errorf("❌ received error while shutting down server: %s", err.Error())
				return
			}

			server.Logger.Info("✅ server shut down successfully")
		}()
	}
}
