package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func init() {
	SchemeBuilder.Register(&Assignment{}, &AssignmentList{})
}

// apiVersion: ipxe.cloud.alexandre.mahdhaoui.com/v1alpha1
// kind: Assignment
// metadata:
//   name: your-assignment
// spec:
//   # subjectSelectors map[string]string
//   # the specified labels selects subjects that can iPXE boot the selected profile below.
//   subjectSelectors:
//     serialNumber:
//       - c4a94672-05a1-4eda-a186-b4aa4544b146
//     uuid:
//       - 47c6da67-7477-4970-aa03-84e48ff4f6ad
//   # profileSelectors map[string]string
//   # the specified labels selects which profile should be used.
//   profileSelectors:
//     ipxe/buildarch: aarch64
//     profile-uuid: 819f1859-a669-410b-adfc-d0bc128e2d7a
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
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`

		Items []AssignmentList `json:"items"`
	}

	AssignmentSpec struct {
		SubjectSelectors map[string][]string `json:"subjectSelectors"`
		ProfileSelectors map[string]string   `json:"profileSelectors"`
	}

	AssignmentStatus struct{}
)
