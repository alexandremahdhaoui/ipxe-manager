package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func init() {
	SchemeBuilder.Register(&Profile{}, &ProfileList{})
}

// apiVersion: ipxe.cloud.alexandre.mahdhaoui.com/v1alpha1
// kind: Profile
// metadata:
//   name: your-profile
//   labels:
//     assign/ipxe/buildarch: aarch64
//     assign/extrinsic/region: us-cal
// spec:
//   # ipxe string
//   ipxe: |
//     # ipxe
//     command ... --with-config {{ config0 }} --ignition-url {{ ignitionFile }} --or-cloud-init {{ cloudInit }}
//   # additionalConfig map[string]string
//   additionalConfig:
//     config0: |
//       YOUR CONFIG HERE
//     ignitionFile: |
//       YOUR IGNITION CONFIG HERE
//     cloudInit: |
//       YOUR CLOUD INIT CONFIG HERE
// status:
//   # UUIDs that are used to fetch
//   additionalConfig:
//     config0: 89952e35-2a85-4f03-a6b2-7f9526bfafc0
//     ignitionFile: 445a4753-3d59-4429-8cea-7db9febdeca

type (
	//+kubebuilder:object:root=true
	//+kubebuilder:subresources:status

	Profile struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`

		Spec   ProfileSpec   `json:"spec,omitempty"`
		Status ProfileStatus `json:"status,omitempty"`
	}

	//+kubebuilder:object:root=true

	ProfileList struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`

		Items []ProfileList `json:"items"`
	}

	ProfileSpec struct {
		IPXE             string            `json:"ipxe"`
		AdditionalConfig map[string]string `json:"additionalConfig,omitempty"`
	}

	ProfileStatus struct{}
)
