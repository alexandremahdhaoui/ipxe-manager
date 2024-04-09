package resolverserver

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/alexandremahdhaoui/ipxer/internal/util/certutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

type Expectation = func(echo.Context, ResolveParams) error

type Mock struct {
	t            *testing.T
	expectations []Expectation
	counter      int
	server       http.Server

	CA   *certutil.CA
	Echo *echo.Echo
}

func (m *Mock) Resolve(ctx echo.Context, _ AnyRoutes, params ResolveParams) error {
	m.t.Helper()

	counter := m.counter
	m.counter += 1

	return m.expectations[counter](ctx, params)
}

func (m *Mock) Start() *Mock {
	go func() {
		if err := m.server.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			require.NoError(m.t, err)
		}
	}()

	return m
}

func (m *Mock) PrependExpectation(f Expectation) *Mock {
	m.expectations = append([]Expectation{f}, m.expectations...)

	return m
}

func (m *Mock) AppendExpectation(f Expectation) *Mock {
	m.expectations = append(m.expectations, f)

	return m
}

func (m *Mock) AssertExpectationsAndShutdown() *Mock {
	m.t.Helper()

	ctx := context.Background()

	assert.Equal(m.t, m.counter, len(m.expectations))
	require.NoError(m.t, m.server.Shutdown(ctx))

	return m
}

func NewMock(t *testing.T, addr string) *Mock {
	t.Helper()

	serverName := strings.SplitN(addr, ":", 2)[0] // a bit hacky

	// generate mTLS certs
	ca, err := certutil.NewCA()
	require.NoError(t, err)

	serverKey, serverCrt, err := ca.NewCertifiedKeyPEM(serverName)
	require.NoError(t, err)

	echoServer := echo.New()
	mock := &Mock{
		t:            t,
		expectations: make([]Expectation, 0),
		counter:      0,

		CA:   ca,
		Echo: echoServer,
	}

	RegisterHandlers(echoServer, mock)

	tlsKeyPair, err := tls.X509KeyPair(serverCrt, serverKey)
	require.NoError(t, err)

	mock.server = http.Server{
		Addr:    addr,
		Handler: echoServer, // set Echo as handler
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsKeyPair},
			RootCAs:      ca.Pool(),
			ServerName:   serverName,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    ca.Pool(),
			//TODO: Parameterize InsecureSkipVerify to test use cases where use would allow self-signed certs.
			//      We may also have to update the RootCAs var.
			InsecureSkipVerify: false,
		},
	}

	mock.Start()

	return mock
}
