// Package v1alpha1 contains API Schema definitions for the v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=ipxe.cloud.alexandre.mahdhaoui.com
package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"sigs.k8s.io/controller-runtime/pkg/client"

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

func NewUUIDLabelSelector(id uuid.UUID) string {
	return LabelSelector(id.String(), UUIDPrefix)
}

func SetUUIDLabelSelector(obj client.Object, id uuid.UUID, value string) {
	obj.GetLabels()[NewUUIDLabelSelector(id)] = value
}

func IsUUIDLabelSelector(key string) bool {
	return strings.Contains(key, LabelSelector("", UUIDPrefix))
}

func IsInternalLabel(key string) bool {
	return strings.Contains(key, Group)
}

func UUIDLabelSelectors(labels map[string]string) (idNameMap map[uuid.UUID]string, reverse map[string]uuid.UUID, err error) {
	idNameMap = make(map[uuid.UUID]string)
	reverse = make(map[string]uuid.UUID)
	for k, v := range labels {
		if !IsUUIDLabelSelector(k) {
			continue
		}

		id, err := uuid.Parse(strings.TrimPrefix(k, LabelSelector("", UUIDPrefix)))
		if err != nil {
			return nil, nil, err // TODO: wrap err
		}

		idNameMap[id] = v
		reverse[k] = id
	}

	return idNameMap, reverse, nil
}
