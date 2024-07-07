package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/labstack/echo/v4"
)

// Serve takes a map of Port to Echo servers, starts the server and gracefully shutdown the servers if any error or
// SIGINT occurs.
func Serve(portServerMap map[int]*echo.Echo) error {
	errChan := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	runBulk(portServerMap, nil, func(port int, e *echo.Echo) {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
			e.Logger.Error(err)
			errChan <- err
		}
	})

	select {
	case err := <-errChan:
		return errors.Join(err, awaitShutdown(portServerMap))
	case <-ctx.Done():
		return awaitShutdown(portServerMap)
	}
}

func awaitShutdown(servers map[int]*echo.Echo) error {
	var errs error

	wg := new(sync.WaitGroup)
	errChan := make(chan error, len(servers))

	runBulk(servers, wg, func(i int, e *echo.Echo) {
		defer wg.Done()

		e.Logger.Info("shutting down server")
		switch err := e.Shutdown(context.Background()); {
		case err == nil:
			e.Logger.Info("successfully shutdown server")
		case errors.Is(err, http.ErrServerClosed):
			e.Logger.Info("server already shut down")
		default:
			e.Logger.Errorf("error while shutting down server: %q", err.Error())
			errChan <- err
		}
	})

	wg.Wait()
	close(errChan)

	for err := range errChan {
		errs = errors.Join(errs, err)
	}

	return errs
}

func runBulk(servers map[int]*echo.Echo, wg *sync.WaitGroup, f func(int, *echo.Echo)) {
	for port, e := range servers {
		if wg != nil {
			wg.Add(1)
		}

		port, e := port, e

		go f(port, e)
	}
}
