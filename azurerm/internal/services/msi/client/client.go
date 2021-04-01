package client

import (
	"github.com/Azure/azure-sdk-for-go/sdk/arm/msi/2018-11-30/armmsi"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	UserAssignedIdentitiesClient *armmsi.UserAssignedIdentitiesClient
}

func NewClient(o *common.ClientOptions) *Client {
	UserAssignedIdentitiesClient := armmsi.NewUserAssignedIdentitiesClient(o.ResourceManagerConnection, o.SubscriptionId)

	return &Client{
		UserAssignedIdentitiesClient: UserAssignedIdentitiesClient,
	}
}
