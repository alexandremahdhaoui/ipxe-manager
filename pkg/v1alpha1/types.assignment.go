package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// DefaultAssignmentLabel is used to query default assignments.
	DefaultAssignmentLabel = LabelSelector("default-assignment")

	// BuildarchAssignmentLabel is used to query assignment based on cpu architecture
	BuildarchAssignmentLabel = LabelSelector("buildarch")
)

type Buildarch string

const (
	// Buildarch

	// I386 - i386	32-bit x86 CPU
	I386 Buildarch = "i386"
	// X8664 - x86_64	64-bit x86 CPU
	X8664 Buildarch = "x86_64"
	// Arm32 - arm32	32-bit ARM CPU
	Arm32 Buildarch = "arm32"
	// Arm64 - arm64	64-bit ARM CPU
	Arm64 Buildarch = "arm64"
)

func init() {
	SchemeBuilder.Register(&Assignment{}, &AssignmentList{})
}

// apiVersion: ipxe.cloud.alexandre.mahdhaoui.com/v1alpha1
// kind: Assignment
// metadata:
//   name: your-assignment
//   labels:
//     ipxe.cloud.alexandre.mahdhaoui.com/buildarch: arm64
//     ipxe.cloud.alexandre.mahdhaoui.com/c4a94672-05a1-4eda-a186-b4aa4544b146: ""
// spec:
//   # subjectSelectors map[string]string
//   # the specified labels selects subjects that can iPXE boot the selected profile below.
//   subjectSelectors:
//     serialNumber:
//       - c4a94672-05a1-4eda-a186-b4aa4544b146
//     uuid:
//       - 47c6da67-7477-4970-aa03-84e48ff4f6ad
//   # profileName string
//   profileName: 819f1859-a669-410b-adfc-d0bc128e2d7a
// status:
//   conditions: []

type (
	//+kubebuilder:object:root=true
	//+kubebuilder:subresources:status

	Assignment struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`

		Spec   AssignmentSpec   `json:"spec,omitempty"`
		Status AssignmentStatus `json:"status,omitempty"`
	}

	//+kubebuilder:object:root=true

	AssignmentList struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitempty"`

		Items []Assignment `json:"items"`
	}

	AssignmentSpec struct {
		SubjectSelectors map[string][]string `json:"subjectSelectors"`
		ProfileName      string
	}

	AssignmentStatus struct{}
)

func (a *Assignment) SetBuildarch(buildarch Buildarch) {
	a.Labels[BuildarchAssignmentLabel] = string(buildarch)
}
