package compute

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/client"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type managedDiskUpdateMetaData struct {
	// Should we shutdown the VM?
	ShouldShutDown bool

	// Should we deallocate the VM first?
	ShouldDeallocate bool

	// Should we turn back on?
	ShouldTurnBackOn bool

	// The ID of VM that is managing this managed disk
	ManagedBy *string

	Client *client.Client
	ID     *parse.ManagedDiskId

	virtualMachineID *parse.VirtualMachineId
	scaleSetVMID     *parse.VirtualMachineScaleSetVMId
}

func (metadata *managedDiskUpdateMetaData) parseVMID() error {
	if vmID, err := parse.VirtualMachineID(*metadata.ManagedBy); err == nil {
		metadata.virtualMachineID = vmID
		return nil
	}
	if vmID, err := parse.VirtualMachineScaleSetVMID(*metadata.ManagedBy); err == nil {
		metadata.scaleSetVMID = vmID
		return nil
	}

	return fmt.Errorf("cannot parse %q as a Virtual Machine ID or Virtual Machine Scale Set VM ID", *metadata.ManagedBy)
}

func (metadata *managedDiskUpdateMetaData) lock() {
	name := ""
	resourceType := ""
	if metadata.virtualMachineID != nil {
		name = metadata.virtualMachineID.Name
		resourceType = virtualMachineResourceName
	} else {
		name = metadata.scaleSetVMID.Name
		resourceType = "azurerm_virtual_machine_scale_set_instance"
	}
	locks.ByName(name, resourceType)
}

func (metadata *managedDiskUpdateMetaData) unlock() {
	name := ""
	resourceType := ""
	if metadata.virtualMachineID != nil {
		name = metadata.virtualMachineID.Name
		resourceType = virtualMachineResourceName
	} else {
		name = metadata.scaleSetVMID.Name
		resourceType = "azurerm_virtual_machine_scale_set_instance"
	}
	locks.UnlockByName(name, resourceType)
}

func (metadata *managedDiskUpdateMetaData) performUpdate(ctx context.Context, update compute.DiskUpdate) error {
	if err := metadata.parseVMID(); err != nil {
		return err
	}
	if metadata.ShouldShutDown {
		metadata.lock()
		defer metadata.unlock()
		if err := metadata.shutDown(ctx); err != nil {
			return err
		}
	}

	if err := metadata.updateManagedDisk(ctx, update); err != nil {
		return err
	}

	if metadata.ShouldTurnBackOn {
		if err := metadata.turnBackOn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (metadata *managedDiskUpdateMetaData) updateManagedDisk(ctx context.Context, update compute.DiskUpdate) error {
	diskClient := metadata.Client.DisksClient
	id := metadata.ID

	log.Printf("[DEBUG] Updating Managed Disk %q (Resource Group %q)...", id.Name, id.ResourceGroup)
	future, err := diskClient.Update(ctx, id.ResourceGroup, id.Name, update)
	if err != nil {
		return fmt.Errorf("Error updating Managed Disk %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for the update of Managed Disk %q (Resource Group %q)...", id.Name, id.ResourceGroup)
	if err := future.WaitForCompletionRef(ctx, diskClient.Client); err != nil {
		return fmt.Errorf("Error waiting for update operation on Managed Disk %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Updated Managed Disk %q (Resource Group %q)", id.Name, id.ResourceGroup)

	return nil
}

func (metadata *managedDiskUpdateMetaData) shutDown(ctx context.Context) error {
	if vmID, err := parse.VirtualMachineID(*metadata.ManagedBy); err == nil {
		return metadata.shutDownVM(ctx, vmID)
	}
	if vmID, err := parse.VirtualMachineScaleSetVMID(*metadata.ManagedBy); err == nil {
		return metadata.shutDownScaleSetVM(ctx, vmID)
	}
	return fmt.Errorf("cannot parse %q as a Virtual Machine ID or Virtual Machine Scale Set VM ID", *metadata.ManagedBy)
}

func (metadata *managedDiskUpdateMetaData) turnBackOn(ctx context.Context) error {
	if vmID, err := parse.VirtualMachineID(*metadata.ManagedBy); err == nil {
		return metadata.turnBackOnVM(ctx, vmID)
	}
	if vmID, err := parse.VirtualMachineScaleSetVMID(*metadata.ManagedBy); err == nil {
		return metadata.turnBackOnScaleSetVM(ctx, vmID)
	}
	return fmt.Errorf("cannot parse %q as a Virtual Machine ID or Virtual Machine Scale Set VM ID", *metadata.ManagedBy)
}

func (metadata *managedDiskUpdateMetaData) shutDownVM(ctx context.Context, virtualMachine *parse.VirtualMachineId) error {
	vmClient := metadata.Client.VMClient

	instanceView, err := vmClient.InstanceView(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving InstanceView for Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
	}

	metadata.ShouldTurnBackOn = true
	metadata.ShouldDeallocate = true

	if instanceView.Statuses != nil {
		for _, status := range *instanceView.Statuses {
			if status.Code == nil {
				continue
			}

			// could also be the provisioning state which we are not bothered with here
			state := strings.ToLower(*status.Code)
			if !strings.HasPrefix(state, "powerstate/") {
				continue
			}

			state = strings.TrimPrefix(state, "powerstate/")
			switch strings.ToLower(state) {
			case "deallocated":
			case "deallocating":
				metadata.ShouldTurnBackOn = false
				metadata.ShouldShutDown = false
				metadata.ShouldDeallocate = false
			case "stopping":
			case "stopped":
				metadata.ShouldShutDown = false
				metadata.ShouldTurnBackOn = false
			}
		}
	}

	// Shutdown
	if metadata.ShouldShutDown {
		log.Printf("[DEBUG] Shutting Down Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		future, err := vmClient.PowerOff(ctx, virtualMachine.ResourceGroup, virtualMachine.Name, utils.Bool(false))
		if err != nil {
			return fmt.Errorf("Error sending Power Off to Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
			return fmt.Errorf("Error waiting for Power Off of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Shut Down Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
	}

	// De-allocate
	if metadata.ShouldDeallocate {
		log.Printf("[DEBUG] Deallocating Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		deAllocFuture, err := vmClient.Deallocate(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
		if err != nil {
			return fmt.Errorf("Error Deallocating to Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		if err := deAllocFuture.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
			return fmt.Errorf("Error waiting for Deallocation of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Deallocated Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
	}

	return nil
}

func (metadata *managedDiskUpdateMetaData) shutDownScaleSetVM(ctx context.Context, vmID *parse.VirtualMachineScaleSetVMId) error {
	return nil
}

func (metadata *managedDiskUpdateMetaData) turnBackOnVM(ctx context.Context, virtualMachine *parse.VirtualMachineId) error {
	vmClient := metadata.Client.VMClient

	log.Printf("[DEBUG] Starting Linux Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
	future, err := vmClient.Start(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
	if err != nil {
		return fmt.Errorf("Error starting Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
		return fmt.Errorf("Error waiting for start of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Started Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
	return nil
}

func (metadata *managedDiskUpdateMetaData) turnBackOnScaleSetVM(ctx context.Context, vmID *parse.VirtualMachineScaleSetVMId) error {
	return nil
}
