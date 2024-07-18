package main

import (
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ContainerRegistry struct {
	client    client.Client
	namespace string
}

func NewContainerRegistry(cl client.Client, namespace string) *ContainerRegistry {
	return &ContainerRegistry{
		client:    cl,
		namespace: namespace,
	}
}

func (r *ContainerRegistry) Setup() error {
	panic("not implemented")
}

func (r *ContainerRegistry) ServiceFQDN() string {
	return fmt.Sprintf("local-container-registry.%s.svc.cluster.local", r.namespace)
}
