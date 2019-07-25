/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AzureContainerRegistrySpec defines the desired state of AzureContainerRegistry
type AzureContainerRegistrySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ResourceGroup    string `json:"resourcegroup"`
	Location         string `json:"location"`
	AdminUserEnabled bool   `json:"adminuserenabled,omitempty"`
	Sku              string `json:"sku,omitempty"`
}

// AzureContainerRegistryStatus defines the observed state of AzureContainerRegistry
type AzureContainerRegistryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Provisioning bool `json:"provisioning,omitempty"`
	Provisioned  bool `json:"provisioned,omitempty"`
}

// +kubebuilder:object:root=true

// AzureContainerRegistry is the Schema for the azurecontainerregistries API
type AzureContainerRegistry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureContainerRegistrySpec   `json:"spec,omitempty"`
	Status AzureContainerRegistryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AzureContainerRegistryList contains a list of AzureContainerRegistry
type AzureContainerRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureContainerRegistry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureContainerRegistry{}, &AzureContainerRegistryList{})
}

// IsBeingDeleted does this
func (acr *AzureContainerRegistry) IsBeingDeleted() bool {
	return !acr.ObjectMeta.DeletionTimestamp.IsZero()
}

// IsSubmitted does this
func (acr *AzureContainerRegistry) IsSubmitted() bool {
	return acr.Status.Provisioning || acr.Status.Provisioned
}

// HasFinalizer does this
func (acr *AzureContainerRegistry) HasFinalizer(finalizerName string) bool {
	return ContainsString(acr.ObjectMeta.Finalizers, finalizerName)
}

// AddFinalizer does this
func (acr *AzureContainerRegistry) AddFinalizer(finalizerName string) {
	acr.ObjectMeta.Finalizers = append(acr.ObjectMeta.Finalizers, finalizerName)
}

// RemoveFinalizer does this
func (acr *AzureContainerRegistry) RemoveFinalizer(finalizerName string) {
	acr.ObjectMeta.Finalizers = RemoveString(acr.ObjectMeta.Finalizers, finalizerName)
}

// ContainsString does this
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// RemoveString does this
func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
