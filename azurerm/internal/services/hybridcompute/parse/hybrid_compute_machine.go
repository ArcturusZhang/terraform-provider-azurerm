package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HybridComputeMachineId struct {
	SubscriptionId string
	ResourceGroup  string
	MachineName    string
}

func NewHybridComputeMachineID(subscriptionId, resourceGroup, machineName string) HybridComputeMachineId {
	return HybridComputeMachineId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		MachineName:    machineName,
	}
}

func (id HybridComputeMachineId) String() string {
	segments := []string{
		fmt.Sprintf("Machine Name %q", id.MachineName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Hybrid Compute Machine", segmentsStr)
}

func (id HybridComputeMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MachineName)
}

// HybridComputeMachineID parses a HybridComputeMachine ID into an HybridComputeMachineId struct
func HybridComputeMachineID(input string) (*HybridComputeMachineId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HybridComputeMachineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MachineName, err = id.PopSegment("machines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
