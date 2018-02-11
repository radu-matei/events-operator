package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventProvider is a specification for an EventProvider resource
type EventProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EventProviderSpec `json:"spec"`
}

// EventProviderSpec is the spec for an EventProvider resource
type EventProviderSpec struct {
	ProviderName    string `json:"providerName"`
	EventType       string `json:"eventType"`
	StorageAccount  string `json:"storageAccount"`
	ResourceGroup   string `json:"resourceGroup"`
	AzureSecretName string `json:"azureSecretName"`
	Host            string `json:"host"`
	HostImage       string `json:"hostImage"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventProviderList is a list of EventProvider resources
type EventProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []EventProvider `json:"items"`
}
