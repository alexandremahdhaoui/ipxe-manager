package main

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8s struct {
	client    client.Client
	namespace string
}

func NewK8s(cl client.Client, namespace string) *K8s {
	return &K8s{
		client:    cl,
		namespace: namespace,
	}
}

func (k *K8s) Setup(ctx context.Context, kubeconfigPath string) error {
	// 1. create the local-container-registry namespace
	ns := corev1.Namespace{}
	ns.Name = k.namespace

	if err := k.client.Create(ctx, &ns); !apierrors.IsAlreadyExists(err) && err != nil {
		return err // TODO: wrap err
	}

	// 2. set kubeconfig for the kubectl subprocess which applies the cert manager config
	if err := os.Setenv("KUBECONFIG", kubeconfigPath); err != nil {
		return err // TODO: wrap err
	}

	return nil
}

func (k *K8s) Teardown(ctx context.Context, kubeconfigPath string) error {
	ns := &corev1.Namespace{}
	ns.Name = k.namespace

	if err := k.client.Delete(ctx, ns); !apierrors.IsNotFound(err) && err != nil {
		return err // TODO: wrap err
	}

	// 2. set kubeconfig for the kubectl subprocess which applies the cert manager config
	if err := os.Setenv("KUBECONFIG", kubeconfigPath); err != nil {
		return err // TODO: wrap err
	}

	return nil
}
