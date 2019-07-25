package acr

import (
	"context"

	//"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-06-01-preview/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/jananiv/acroperator/resourcehelper/config"
	"github.com/jananiv/acroperator/resourcehelper/iam"
)

func getacrclient() containerregistry.RegistriesClient {
	registryClient := containerregistry.NewRegistriesClient("7060bca0-7a3c-44bd-b54c-4bb1e9facfac")
	auth, _ := iam.GetResourceManagementAuthorizer()
	registryClient.Authorizer = auth
	registryClient.AddToUserAgent(config.UserAgent())
	return registryClient
}

// CreateRegistry creates an Event Hubs hub in a namespace
// Parameters:
// resourceGroupName - name of the resource group within the azure subscription.
// registryName - the Container Registry name
func CreateRegistry(ctx context.Context, resourceGroupName string, registryName string, location string, sku string, adminUserEnabled bool) (containerregistry.RegistriesCreateFuture, error) {
	registryclient := getacrclient()
	skuName := containerregistry.SkuName(sku)
	return registryclient.Create(
		ctx,
		resourceGroupName,
		registryName,
		containerregistry.Registry{
			Sku: &containerregistry.Sku{
				Name: skuName,
			},
			Location: to.StringPtr(location),
			RegistryProperties: &containerregistry.RegistryProperties{
				AdminUserEnabled: to.BoolPtr(adminUserEnabled),
			},
		},
	)
}

// DeleteRegistry deletes a container registry from the resource group.
// Parameters:
// resourceGroupName - name of the resource group within the azure subscription.
// registryName - the name of Container Registry
func DeleteRegistry(ctx context.Context, resourceGroupName string, registryName string) (result containerregistry.RegistriesDeleteFuture, err error) {
	registryclient := getacrclient()
	return registryclient.Delete(ctx,
		resourceGroupName,
		registryName)
}

//GetRegistry gets info abput registry
func GetRegistry(ctx context.Context, resourceGroupName string, registryName string) (containerregistry.Registry, error) {
	registryclient := getacrclient()
	return registryclient.Get(ctx, resourceGroupName, registryName)
}
