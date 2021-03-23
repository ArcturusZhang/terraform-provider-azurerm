package client

import (
	"github.com/Azure/azure-sdk-for-go/sdk/arm/msi/2018-11-30/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	UserAssignedIdentitiesClient *armmsi.UserAssignedIdentitiesClient
}

func NewClient(o *common.ClientOptions) *Client {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	UserAssignedIdentitiesClient := armmsi.NewUserAssignedIdentitiesClient(armcore.NewDefaultConnection(cred, &armcore.ConnectionOptions{
		Logging: azcore.LogOptions{IncludeBody: true},
	}), o.SubscriptionId)

	return &Client{
		UserAssignedIdentitiesClient: UserAssignedIdentitiesClient,
	}
}
