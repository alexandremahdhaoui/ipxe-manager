package certutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// Just copied this over from:
// https://github.com/madflojo/testcerts/blob/main/testcerts.go

type CA struct {
	rootCert       *x509.Certificate
	selfSignedCert *x509.Certificate
	pool           *x509.CertPool
	privateKey     *ecdsa.PrivateKey
}

func NewCA() (*CA, error) {
	// 1. create a ca cert.
	caCert := &x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Use in test only!"},
		},
		SerialNumber:          big.NewInt(123),
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(1 * time.Hour),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	// 2. create private key.
	caKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err //TOOO: wrap err
	}

	// 3. self sign a certificate with the CA's own cert and the previously generated private key.
	selfSignedCertRaw, err := x509.CreateCertificate(rand.Reader, caCert, caCert, caKey.Public(), caKey)
	if err != nil {
		return nil, err //TODO: wrap err
	}

	// 4. Add self-signed cert to pool
	certPool := x509.NewCertPool()
	//certPool.AppendCertsFromPEM(selfSignedCertRaw)
	//TODO: check if the code below cannot be replaced with above code.

	selfSignedCert, err := x509.ParseCertificate(selfSignedCertRaw)
	if err != nil {
		return nil, err
	}

	certPool.AddCert(selfSignedCert)

	return &CA{
		rootCert:       caCert,
		pool:           certPool,
		selfSignedCert: selfSignedCert,
		privateKey:     caKey,
	}, nil
}
