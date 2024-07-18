package main

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	credentialsSecretName = "local-container-registry-credentials"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Credential struct {
	client          client.Client
	credentials     Credentials
	credentialsPath string
	namespace       string
}

func NewCredential(cl client.Client, credentialsPath, namespace string) *Credential {
	return &Credential{
		client: cl,
		credentials: Credentials{
			Username: generateRandomString(32),
			Password: generateRandomString(32),
		},
		credentialsPath: credentialsPath,
		namespace:       namespace,
	}
}

func (c *Credential) Setup() error {
	// 1. write credentials.
	if err := c.writeCredentials(); err != nil {
		return err // TODO: wrap err
	}

	// 2. create secret.
	ctx := context.Background()

	// 3. credential secret.
	credentialsSecret := &corev1.Secret{} //nolint:exhaustruct

	credentialsSecret.Name = credentialsSecretName
	credentialsSecret.Namespace = c.namespace
	credentialsSecret.Type = corev1.SecretTypeOpaque

	credentialsSecret.Data = map[string][]byte{
		"username": []byte(c.credentials.Username),
		"password": []byte(c.credentials.Password),
	}

	if err := c.client.Create(ctx, credentialsSecret); err != nil {
		return err // TODO: wrap err
	}

	return nil
}

func (c *Credential) writeCredentials() error {
	b, err := json.Marshal(c.credentials)
	if err != nil {
		return err // TODO: wrap err
	}

	if err := os.WriteFile(c.credentialsPath, b, 0600); err != nil {
		return err // TODO: wrap err
	}

	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
