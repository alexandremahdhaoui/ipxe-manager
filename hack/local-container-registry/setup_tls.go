package main

import (
	"context"
	"github.com/alexandremahdhaoui/ipxer/hack/internal"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"os/exec"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	certManagerManifests = `https://github.com/cert-manager/cert-manager/releases/download/v1.15.1/cert-manager.yaml`
)

type TLS struct {
	client client.Client

	caCrtPath           string
	registryNamespace   string
	registryServiceFQDN string
}

func NewTLS(cl client.Client, caCrtPath, registryNamespace, registryServiceFQDN string) *TLS {
	return &TLS{
		client: cl,

		caCrtPath:           caCrtPath,
		registryNamespace:   registryNamespace,
		registryServiceFQDN: registryServiceFQDN,
	}
}

func (t *TLS) Setup(ctx context.Context) error {
	// 1. Install cert-manager.
	helmRepoAdd := exec.Command("helm", strings.Split(
		"repo add jetstack https://charts.jetstack.io --force-update", " ")...)
	if err := internal.RunCmdWithStdPipes(helmRepoAdd); err != nil {
		return err // TODO: wrap err
	}

	helmInstall := exec.Command("helm", strings.Split(
		"install cert-manager jetstack/cert-manager "+
			"--namespace cert-manager "+
			"--create-namespace "+
			"--version v1.15.1 "+
			"--set crds.enabled=true", " ")...)
	if err := internal.RunCmdWithStdPipes(helmInstall); err != nil {
		return err // TODO: wrap err
	}

	// TODO: await cert-manager pods are running

	// 2. Create self signed issuer.
	issuer := &certmanagerv1.Issuer{}

	issuer.Name = t.ResourceName()
	issuer.Namespace = t.registryNamespace
	issuer.Spec.SelfSigned = &certmanagerv1.SelfSignedIssuer{}

	if err := t.client.Create(ctx, issuer); err != nil {
		return err // TODO: wrap err
	}

	// 3. Create certificate in registry namespace.
	cert := &certmanagerv1.Certificate{}

	cert.Name = t.ResourceName()
	cert.Namespace = t.registryNamespace

	cert.Spec.DNSNames = []string{t.registryServiceFQDN}
	cert.Spec.SecretName = t.ResourceName()
	cert.Spec.IssuerRef.Name = t.ResourceName()

	if err := t.client.Create(ctx, cert); err != nil {
		return err // TODO: wrap err
	}

	return nil
}

func (t *TLS) ResourceName() string {
	return "local-container-registry"
}

func (t *TLS) Teardown() error {
	cmd := exec.Command("kubectl", "delete", "-f", certManagerManifests)

	if err := internal.RunCmdWithStdPipes(cmd); err != nil {
		return err // TODO: wrap err
	}

	return nil
}
