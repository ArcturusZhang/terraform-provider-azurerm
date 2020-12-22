package client

import (
	"github.com/Azure/azure-sdk-for-go/services/hybridcompute/mgmt/2020-08-02/hybridcompute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	MachineClient *hybridcompute.MachinesClient
}

func NewClient(o *common.ClientOptions) *Client {
	machineClient := hybridcompute.NewMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&machineClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MachineClient: &machineClient,
	}
}
