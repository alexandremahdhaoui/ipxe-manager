// Package v1alpha1 contains API Schema definitions for the v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=ipxe.cloud.alexandre.mahdhaoui.com
package v1alpha1

import (
	"fmt"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	Group   = "ipxe.cloud.alexandre.mahdhaoui.com"
	Version = "v1alpha1"
)

var (
	GroupVersion  = schema.GroupVersion{Group: Group, Version: Version}
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
	AddToScheme   = SchemeBuilder.AddToScheme
)

func UUIDLabelSelector(id uuid.UUID) string {
	return fmt.Sprintf("%s/%s", Group, id.String())
}

func LabelSelector(key string) string {
	return fmt.Sprintf("%s/%s", Group, key)
}
