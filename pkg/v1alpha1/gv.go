// Package v1alpha1 contains API Schema definitions for the v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=ipxe.cloud.alexandre.mahdhaoui.com
package v1alpha1

import (
	"fmt"
	"github.com/google/uuid"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	Group   = "ipxer.cloud.alexandre.mahdhaoui.com"
	Version = "v1alpha1"

	UUIDPrefix      = "uuid"
	BuildarchPrefix = "buildarch"
)

var (
	GroupVersion  = schema.GroupVersion{Group: Group, Version: Version}
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
	AddToScheme   = SchemeBuilder.AddToScheme
)

func LabelSelector(key string, prefixes ...string) string {
	label := fmt.Sprintf("%s/%s", Group, key)

	if len(prefixes) > 0 {
		label = fmt.Sprintf("%s.%s", strings.Join(prefixes, "."), label)
	}

	return label
}

func IsInternalLabel(key string) bool {
	return strings.Contains(key, Group)
}

func SetUUIDLabelSelector(obj client.Object, id uuid.UUID, value string) {
	obj.GetLabels()[NewUUIDLabelSelector(id)] = value
}

func NewUUIDLabelSelector(id uuid.UUID) string {
	return LabelSelector(id.String(), UUIDPrefix)
}

func IsUUIDLabelSelector(key string) bool {
	return strings.Contains(key, LabelSelector("", UUIDPrefix))
}
