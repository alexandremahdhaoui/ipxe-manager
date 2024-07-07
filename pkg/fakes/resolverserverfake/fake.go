package resolverserverfake

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/alexandremahdhaoui/ipxer/internal/util/certutil"
	"github.com/alexandremahdhaoui/ipxer/pkg/generated/resolverserver"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Expectation = func(echo.Context, resolverserver.ResolveParams) error

type Fake struct {
	t            *testing.T
	expectations []Expectation
	counter      int
	server       http.Server

	CA   *certutil.CA
	Echo *echo.Echo
}

func (f *Fake) Resolve(ctx echo.Context, _ resolverserver.AnyRoutes, params resolverserver.ResolveParams) error {
	f.t.Helper()

	counter := f.counter
	f.counter += 1

	return f.expectations[counter](ctx, params)
}

func (f *Fake) Start() *Fake {
	go func() {
		if err := f.server.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			require.NoError(f.t, err)
		}
	}()

	return f
}

func (f *Fake) PrependExpectation(expectation Expectation) *Fake {
	f.expectations = append([]Expectation{expectation}, f.expectations...)

	return f
}

func (f *Fake) AppendExpectation(expectation Expectation) *Fake {
	f.expectations = append(f.expectations, expectation)

	return f
}

func (f *Fake) AssertExpectationsAndShutdown() *Fake {
	f.t.Helper()

	ctx := context.Background()

	assert.Equal(f.t, f.counter, len(f.expectations))
	require.NoError(f.t, f.server.Shutdown(ctx))

	return f
}

func New(t *testing.T, addr string) *Fake {
	t.Helper()

	serverName := strings.SplitN(addr, ":", 2)[0] // a bit hacky

	// generate mTLS certs
	ca, err := certutil.NewCA()
	require.NoError(t, err)

	serverKey, serverCrt, err := ca.NewCertifiedKeyPEM(serverName)
	require.NoError(t, err)

	echoServer := echo.New()
	fake := &Fake{
		t:            t,
		expectations: make([]Expectation, 0),
		counter:      0,

		CA:   ca,
		Echo: echoServer,
	}

	resolverserver.RegisterHandlers(echoServer, fake)

	tlsKeyPair, err := tls.X509KeyPair(serverCrt, serverKey)
	require.NoError(t, err)

	fake.server = http.Server{
		Addr:    addr,
		Handler: echoServer, // set Echo as handler
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsKeyPair},
			RootCAs:      ca.Pool(),
			ServerName:   serverName,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    ca.Pool(),
			// TODO: Parameterize InsecureSkipVerify to test use cases where use would allow self-signed certs.
			//      We may also have to update the RootCAs var.
			InsecureSkipVerify: false,
		},
	}

	fake.Start()

	return fake
}
