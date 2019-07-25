// Package config manages loading configuration from environment and command-line params
package config

import (
	"fmt"
	"os"

	"github.com/Azure/go-autorest/autorest/azure"
)

var (
	// these are our *global* config settings, to be shared by all packages.
	// each has corresponding public accessors below.
	// if anything requires a `Set` accessor, that indicates it perhaps
	// shouldn't be set here, because mutable vars shouldn't be global.
	clientID               string
	clientSecret           string
	tenantID               string
	subscriptionID         string
	authorizationServerURL string
	cloudName              string = "AzurePublicCloud"
	useDeviceFlow          bool

	userAgent   string
	environment *azure.Environment
)

// ClientID is the OAuth client ID.
func ClientID() string {
	return clientID
}

// ClientSecret is the OAuth client secret.
func ClientSecret() string {
	return clientSecret
}

// TenantID is the AAD tenant to which this client belongs.
func TenantID() string {
	return tenantID
}

// SubscriptionID is a target subscription for Azure resources.
func SubscriptionID() string {
	return subscriptionID
}

// AuthorizationServerURL is the OAuth authorization server URL.
// Q: Can this be gotten from the `azure.Environment` in `Environment()`?
func AuthorizationServerURL() string {
	return authorizationServerURL
}

//UseDeviceFlow specifies if interactive auth should be used. Interactive
// auth uses the OAuth Device Flow grant type.
func UseDeviceFlow() bool {
	return useDeviceFlow
}

//UserAgent specifies a string to append to the agent identifier.
func UserAgent() string {
	if len(userAgent) > 0 {
		return userAgent
	}
	return "AzureOperator-resourcemanager"
}

//Environment returns an `azure.Environment{...}` for the current cloud.
func Environment() *azure.Environment {
	if environment != nil {
		return environment
	}
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		// TODO: move to initialization of var
		panic(fmt.Sprintf(
			"invalid cloud name '%s' specified, cannot continue\n", cloudName))
	}
	environment = &env
	LoadSettings()
	return environment
}

// LoadSettings loads blah
func LoadSettings() error {

	fmt.Println()

	azureEnv, _ := azure.EnvironmentFromName("AzurePublicCloud") // shouldn't fail
	authorizationServerURL = azureEnv.ActiveDirectoryEndpoint

	// these must be provided by environment
	// clientID
	clientID = os.Getenv("AZURE_CLIENT_ID")
	if len(clientID) == 0 {
		return fmt.Errorf("expected env vars not provided")
	}
	fmt.Println("Client ID: ", clientID)
	// clientSecret
	clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	if len(clientSecret) == 0 { // don't need a secret for device flow
		return fmt.Errorf("expected env vars not provided")
	}
	fmt.Println("Client secret: ", clientSecret)
	// tenantID (AAD)
	tenantID = os.Getenv("AZURE_TENANT_ID")
	if len(tenantID) == 0 {
		return fmt.Errorf("expected env vars not provided")
	}

	fmt.Println("tenantID: ", tenantID)
	// subscriptionID (ARM)
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		return fmt.Errorf("expected env vars not provided")
	}

	return nil
}
