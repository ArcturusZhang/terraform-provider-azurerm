package parse

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineScaleSetVMId struct {
	ResourceGroup              string
	VirtualMachineScaleSetName string
	Name                       string
}

func VirtualMachineScaleSetVMID(input string) (*VirtualMachineScaleSetVMId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Virtual Machine Scale Set VM ID %q: %+v", input, err)
	}

	vm := VirtualMachineScaleSetVMId{
		ResourceGroup: id.ResourceGroup,
	}

	if vm.VirtualMachineScaleSetName, err = id.PopSegment("virtualMachineScaleSets"); err != nil {
		return nil, err
	}

	if vm.Name, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vm, nil
}
