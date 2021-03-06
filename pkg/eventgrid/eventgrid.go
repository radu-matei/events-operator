package eventgrid

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2018-01-01/eventgrid"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	defaultLocation = "westeurope"

	defaultActiveDirectoryEndpoint = azure.PublicCloud.ActiveDirectoryEndpoint
	defaultResourceManagerEndpoint = azure.PublicCloud.ResourceManagerEndpoint

	subscriptionID = getEnvVarOrExit("AZ_SUBSCRIPTION_ID")
	tenantID       = getEnvVarOrExit("AZ_TENANT_ID")
	clientID       = getEnvVarOrExit("AZ_CLIENT_ID")
	clientSecret   = getEnvVarOrExit("AZ_CLIENT_SECRET")
)

func getEventGridClient() (eventgrid.EventSubscriptionsClient, error) {
	var subscriptionsClient eventgrid.EventSubscriptionsClient

	oAuthConfig, err := adal.NewOAuthConfig(defaultActiveDirectoryEndpoint, tenantID)
	if err != nil {
		return subscriptionsClient, fmt.Errorf("cannot get oauth config: %v", err)
	}
	token, err := adal.NewServicePrincipalToken(*oAuthConfig, clientID, clientSecret, defaultResourceManagerEndpoint)
	if err != nil {
		return subscriptionsClient, fmt.Errorf("cannot get service principal token: %v", err)
	}

	subscriptionsClient = eventgrid.NewEventSubscriptionsClient(subscriptionID)
	subscriptionsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return subscriptionsClient, nil
}

// CheckEventSubscription checks the existence of an event subscription
func CheckEventSubscription(name, resourceGroup, storageAccount, tlsWebHook string) (bool, error) {

	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s", subscriptionID, resourceGroup, storageAccount)

	c, err := getEventGridClient()
	if err != nil {
		log.Fatalf("cannot get eventgrid client: %v", err)
	}

	_, err = c.Get(context.Background(), scope, name)
	if err != nil {
		return false, fmt.Errorf("cannot get event subscription: %v", err)
	}

	return true, nil
}

// CreateOrUpdateEventSubscription creates an Azure Event Grid eventsubscription
func CreateOrUpdateEventSubscription(resourceGroup, storageAccountName, tlsWebhook string) error {
	c, err := getEventGridClient()
	if err != nil {
		log.Fatalf("cannot get eventgrid client: %v", err)
	}

	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s", subscriptionID, resourceGroup, storageAccountName)
	subscriptionName := fmt.Sprintf("%seventsubscription", storageAccountName)

	subscription := eventgrid.EventSubscription{
		EventSubscriptionProperties: &eventgrid.EventSubscriptionProperties{
			Destination: eventgrid.WebHookEventSubscriptionDestination{
				EndpointType: eventgrid.EndpointTypeWebHook,
				WebHookEventSubscriptionDestinationProperties: &eventgrid.WebHookEventSubscriptionDestinationProperties{
					EndpointURL: to.StringPtr(tlsWebhook),
				},
			},
		},
	}

	ctx := context.Background()

	f, err := c.CreateOrUpdate(ctx, scope, subscriptionName, subscription)
	if err != nil {
		return fmt.Errorf("cannot create event subscription: %v", err)
	}

	err = f.WaitForCompletion(ctx, c.Client)
	if err != nil {
		return fmt.Errorf("cannot get the subscription create or update future response: %v", err)
	}

	return nil
}

func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatalf("missing environment variable %s\n", varName)
	}

	return value
}
