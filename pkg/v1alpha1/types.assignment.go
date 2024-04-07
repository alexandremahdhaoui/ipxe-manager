package v1alpha1

import (
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func init() {
	SchemeBuilder.Register(&Assignment{}, &AssignmentList{})
}

var (
	// DefaultAssignmentLabel is used to query default assignments.
	DefaultAssignmentLabel = LabelSelector("default-assignment")

	// BuildarchList Label Selector

	I386BuildarchLabelSelector  = LabelSelector(I386.String(), BuildarchPrefix)
	X8664BuildarchLabelSelector = LabelSelector(X8664.String(), BuildarchPrefix)
	Arm32BuildarchLabelSelector = LabelSelector(Arm32.String(), BuildarchPrefix)
	Arm64BuildarchLabelSelector = LabelSelector(Arm64.String(), BuildarchPrefix)
)

type Buildarch string

func (b Buildarch) String() string {
	return string(b)
}

const (
	// BuildarchList

	// I386 - i386	32-bit x86 CPU
	I386 Buildarch = "i386"
	// X8664 - x86_64	64-bit x86 CPU
	X8664 Buildarch = "x86_64"
	// Arm32 - arm32	32-bit ARM CPU
	Arm32 Buildarch = "arm32"
	// Arm64 - arm64	64-bit ARM CPU
	Arm64 Buildarch = "arm64"
)

// apiVersion: ipxer.cloud.alexandre.mahdhaoui.com/v1alpha1
// kind: Assignment
// metadata:
//   name: your-assignment
//   labels:
//     ipxer.cloud.alexandre.mahdhaoui.com/buildarch: arm64
//     uuid.ipxer.cloud.alexandre.mahdhaoui.com/c4a94672-05a1-4eda-a186-b4aa4544b146: ""
//     uuid.ipxer.cloud.alexandre.mahdhaoui.com/3f5f3c39-584e-4c7c-b6ff-137e1aaa7175: ""
// spec:
//   # subjectSelectors map[string][]string
//   # the specified labels selects subjects that can iPXE boot the selected profile below.
//   subjectSelectors:
//     buildarch: # please note only 1 buildarch mat be specified at a time.
//       - arm64
//     serialNumber:
//       - c4a94672-05a1-4eda-a186-b4aa4544b146
//     uuid:
//       - 47c6da67-7477-4970-aa03-84e48ff4f6ad
//       - 3f5f3c39-584e-4c7c-b6ff-137e1aaa7175
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
		SubjectSelectors SubjectSelectors `json:"subjectSelectors"`
		ProfileName      string           `json:"profileName"`
		IsDefault        bool             `json:"isDefault"`
	}

	AssignmentStatus struct{}

	SubjectSelectors struct {
		BuildarchList []Buildarch `json:"buildarch"`
		UUIDList      []string    `json:"uuidList"`
	}
)

func (a *Assignment) GetBuildarchList() []Buildarch {
	out := make([]Buildarch, 0)

	if _, ok := a.Labels[Arm32BuildarchLabelSelector]; ok {
		out = append(out, Arm32)
	}

	if _, ok := a.Labels[Arm64BuildarchLabelSelector]; ok {
		out = append(out, Arm64)
	}

	if _, ok := a.Labels[I386BuildarchLabelSelector]; ok {
		out = append(out, I386)
	}

	if _, ok := a.Labels[X8664BuildarchLabelSelector]; ok {
		out = append(out, X8664)
	}

	return out
}

func (a *Assignment) SetBuildarch(buildarch Buildarch) {
	switch buildarch {
	case Arm32:
		a.Labels[Arm32BuildarchLabelSelector] = ""
	case Arm64:
		a.Labels[Arm64BuildarchLabelSelector] = ""
	case I386:
		a.Labels[I386BuildarchLabelSelector] = ""
	case X8664:
		a.Labels[X8664BuildarchLabelSelector] = ""
	}
}

func (a *Assignment) SetUUIDLabelSelector(id uuid.UUID) {
	a.Labels[NewUUIDLabelSelector(id)] = ""
}

func NewUUIDLabelSelector(id uuid.UUID) string {
	return LabelSelector(id.String(), UUIDPrefix)
}

func IsInternalLabel(key string) bool {
	return strings.Contains(key, Group)
}
