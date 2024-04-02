package transformerserver

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type Expectation = func(echo.Context) error

type Mock struct {
	t            *testing.T
	expectations []Expectation
	counter      uint
	Echo         *echo.Echo
}

func (w *Mock) Transform(ctx echo.Context, _ AnyRoutes) error {
	w.t.Helper()

	counter := w.counter
	w.counter += 1

	return w.expectations[counter](ctx)
}

func (w *Mock) PrependExpectation(f Expectation) {
	w.expectations = append([]Expectation{f}, w.expectations...)
}

func (w *Mock) AppendExpectation(f Expectation) {
	w.expectations = append(w.expectations, f)
}

func (w *Mock) AssertExpectationsAndShutdown() {
	w.t.Helper()

	assert.Equal(w.t, w.counter, len(w.expectations))
	require.NoError(w.t, w.Echo.Shutdown(context.Background()))
}

func NewMock(t *testing.T, address string) *Mock {
	t.Helper()

	echoServer := echo.New()
	mock := &Mock{
		t:            t,
		expectations: make([]Expectation, 0),
		counter:      0,
		Echo:         echoServer,
	}

	RegisterHandlers(echoServer, mock)

	go func() {
		require.NoError(t, echoServer.Start(address))
	}()

	return mock
}
