package compute

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-12-01/armcompute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceManagedDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagedDiskCreateUpdate,
		Read:   resourceManagedDiskRead,
		Update: resourceManagedDiskUpdate,
		Delete: resourceManagedDiskDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ManagedDiskID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"zones": azure.SchemaSingleZone(),

			"storage_account_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.StandardLRS),
					string(compute.PremiumLRS),
					string(compute.StandardSSDLRS),
					string(compute.UltraSSDLRS),
				}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"create_option": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Copy),
					string(compute.Empty),
					string(compute.FromImage),
					string(compute.Import),
					string(compute.Restore),
				}, false),
			},

			"source_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"source_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true, // Not supported by disk update
				ValidateFunc: azure.ValidateResourceID,
			},

			"image_reference_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Windows),
					string(compute.Linux),
				}, true),
			},

			"disk_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateManagedDiskSizeGB,
			},

			"disk_iops_read_write": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"disk_mbps_read_write": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"disk_encryption_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
				//       https://github.com/Azure/azure-rest-api-specs/issues/8132
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"encryption_settings": encryptionSettingsSchema(),

			"tags": tags.Schema(),
		},
	}
}

func resourceManagedDiskCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, nil)
		if err != nil {
			if !utils.Track2ResponseWasNotFound(err) {
				return fmt.Errorf("Error checking for presence of existing Managed Disk %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.Disk != nil && existing.Disk.ID != nil && *existing.Disk.ID != "" {
			return tf.ImportAsExistsError("azurerm_managed_disk", *existing.Disk.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	createOption := armcompute.DiskCreateOption(d.Get("create_option").(string))
	storageAccountType := d.Get("storage_account_type").(string)
	osType := d.Get("os_type").(string)
	t := d.Get("tags").(map[string]interface{})
	zones := azure.Track2ExpandZones(d.Get("zones").([]interface{}))
	skuName := armcompute.DiskStorageAccountTypes(storageAccountType)

	props := &armcompute.DiskProperties{
		CreationData: &armcompute.CreationData{
			CreateOption: &createOption,
		},
		OSType: armcompute.OperatingSystemTypes(osType).ToPtr(),
		Encryption: &armcompute.Encryption{
			Type: armcompute.EncryptionTypeEncryptionAtRestWithPlatformKey.ToPtr(),
		},
	}

	if v := d.Get("disk_size_gb"); v != 0 {
		diskSize := int32(v.(int))
		props.DiskSizeGB = &diskSize
	}

	if storageAccountType == string(compute.UltraSSDLRS) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskIOPS := int64(v.(int))
			props.DiskIOPSReadWrite = &diskIOPS
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskMBps := int64(v.(int))
			props.DiskMBpsReadWrite = &diskMBps
		}
	} else if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") {
		return fmt.Errorf("[ERROR] disk_iops_read_write and disk_mbps_read_write are only available for UltraSSD disks")
	}

	if createOption == armcompute.DiskCreateOptionImport {
		sourceUri := d.Get("source_uri").(string)
		if sourceUri == "" {
			return fmt.Errorf("`source_uri` must be specified when `create_option` is set to `Import`")
		}

		storageAccountId := d.Get("storage_account_id").(string)
		if storageAccountId == "" {
			return fmt.Errorf("`storage_account_id` must be specified when `create_option` is set to `Import`")
		}

		props.CreationData.StorageAccountID = utils.String(storageAccountId)
		props.CreationData.SourceURI = utils.String(sourceUri)
	}
	if createOption == armcompute.DiskCreateOptionCopy || createOption == armcompute.DiskCreateOptionRestore {
		sourceResourceId := d.Get("source_resource_id").(string)
		if sourceResourceId == "" {
			return fmt.Errorf("`source_resource_id` must be specified when `create_option` is set to `Copy` or `Restore`")
		}

		props.CreationData.SourceResourceID = utils.String(sourceResourceId)
	}
	if createOption == armcompute.DiskCreateOptionFromImage {
		imageReferenceId := d.Get("image_reference_id").(string)
		if imageReferenceId == "" {
			return fmt.Errorf("`image_reference_id` must be specified when `create_option` is set to `Import`")
		}

		props.CreationData.ImageReference = &armcompute.ImageDiskReference{
			ID: utils.String(imageReferenceId),
		}
	}

	if v, ok := d.GetOk("encryption_settings"); ok {
		encryptionSettings := v.([]interface{})
		settings := encryptionSettings[0].(map[string]interface{})
		props.EncryptionSettingsCollection = expandManagedDiskEncryptionSettings(settings)
	}

	if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
		props.Encryption = &armcompute.Encryption{
			Type:                armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey.ToPtr(),
			DiskEncryptionSetID: utils.String(diskEncryptionSetId),
		}
	}

	createDisk := armcompute.Disk{
		Resource: armcompute.Resource{
			Name:           &name,
			Location:       &location,
			Tags:  tags.Track2Expand(t),
		},
		Properties: props,
		SKU: &armcompute.DiskSKU{
			Name: &skuName,
		},
		Zones: zones,
	}

	future, err := client.BeginCreateOrUpdate(ctx, resourceGroup, name, createDisk, nil)
	if err != nil {
		return fmt.Errorf("Error creating/updating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if _, err = future.PollUntilDone(ctx, common.DefaultPollingInterval); err != nil {
		return fmt.Errorf("Error waiting for create/update of Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.Disk != nil && read.Disk.ID == nil {
		return fmt.Errorf("Error reading Managed Disk %s (Resource Group %q): ID was nil", name, resourceGroup)
	}

	d.SetId(*read.Disk.ID)

	return resourceManagedDiskRead(d, meta)
}

func resourceManagedDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	storageAccountType := d.Get("storage_account_type").(string)
	shouldShutDown := false

	resp, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(err) {
			return fmt.Errorf("Error Managed Disk %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Azure Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	diskUpdate := armcompute.DiskUpdate{
		Properties: &armcompute.DiskUpdateProperties{},
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		diskUpdate.Tags = tags.Track2Expand(t)
	}

	if d.HasChange("storage_account_type") {
		shouldShutDown = true
		var skuName armcompute.DiskStorageAccountTypes
		for _, v := range armcompute.PossibleDiskStorageAccountTypesValues() {
			if strings.EqualFold(storageAccountType, string(v)) {
				skuName = v
			}
		}
		diskUpdate.SKU = &armcompute.DiskSKU{
			Name: &skuName,
		}
	}

	if strings.EqualFold(storageAccountType, string(armcompute.DiskStorageAccountTypesUltraSSDLRS)) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskIOPS := int64(v.(int))
			diskUpdate.Properties.DiskIOPSReadWrite = &diskIOPS
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskMBps := int64(v.(int))
			diskUpdate.Properties.DiskMBpsReadWrite = &diskMBps
		}
	} else if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") {
		return fmt.Errorf("[ERROR] disk_iops_read_write and disk_mbps_read_write are only available for UltraSSD disks")
	}

	if d.HasChange("os_type") {
		diskUpdate.Properties.OSType = armcompute.OperatingSystemTypes(d.Get("os_type").(string)).ToPtr()
	}

	if d.HasChange("disk_size_gb") {
		if old, new := d.GetChange("disk_size_gb"); new.(int) > old.(int) {
			shouldShutDown = true
			diskUpdate.Properties.DiskSizeGB = utils.Int32(int32(new.(int)))
		} else {
			return fmt.Errorf("Error - New size must be greater than original size. Shrinking disks is not supported on Azure")
		}
	}

	if d.HasChange("disk_encryption_set_id") {
		shouldShutDown = true
		if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			diskUpdate.Properties.Encryption = &armcompute.Encryption{
				Type:                armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey.ToPtr(),
				DiskEncryptionSetID: utils.String(diskEncryptionSetId),
			}
		} else {
			return fmt.Errorf("Once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	disk := *resp.Disk
	// whilst we need to shut this down, if we're not attached to anything there's no point
	if shouldShutDown && disk.ManagedBy == nil {
		shouldShutDown = false
	}

	// if we are attached to a VM we bring down the VM as necessary for the operations which are not allowed while it's online
	if shouldShutDown {
		virtualMachine, err := parse.VirtualMachineID(*disk.ManagedBy)
		if err != nil {
			return fmt.Errorf("Error parsing VMID %q for disk attachment: %+v", *disk.ManagedBy, err)
		}
		// check instanceView State
		vmClient := meta.(*clients.Client).Compute.VMClient

		locks.ByName(name, virtualMachineResourceName)
		defer locks.UnlockByName(name, virtualMachineResourceName)

		instanceView, err := vmClient.InstanceView(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
		if err != nil {
			return fmt.Errorf("Error retrieving InstanceView for Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		shouldTurnBackOn := true
		shouldDeallocate := true

		if instanceView.Statuses != nil {
			for _, status := range *instanceView.Statuses {
				if status.Code == nil {
					continue
				}

				// could also be the provisioning state which we're not bothered with here
				state := strings.ToLower(*status.Code)
				if !strings.HasPrefix(state, "powerstate/") {
					continue
				}

				state = strings.TrimPrefix(state, "powerstate/")
				switch strings.ToLower(state) {
				case "deallocated":
				case "deallocating":
					shouldTurnBackOn = false
					shouldShutDown = false
					shouldDeallocate = false
				case "stopping":
				case "stopped":
					shouldShutDown = false
					shouldTurnBackOn = false
				}
			}
		}

		// Shutdown
		if shouldShutDown {
			log.Printf("[DEBUG] Shutting Down Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			forceShutdown := false
			future, err := vmClient.PowerOff(ctx, virtualMachine.ResourceGroup, virtualMachine.Name, utils.Bool(forceShutdown))
			if err != nil {
				return fmt.Errorf("Error sending Power Off to Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
				return fmt.Errorf("Error waiting for Power Off of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Shut Down Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}

		// De-allocate
		if shouldDeallocate {
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

		// Update Disk
		updateFuture, err := client.BeginUpdate(ctx, resourceGroup, name, diskUpdate, nil)
		if err != nil {
			return fmt.Errorf("Error updating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		if _, err := updateFuture.PollUntilDone(ctx, common.DefaultPollingInterval); err != nil {
			return fmt.Errorf("Error waiting for update of Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if shouldTurnBackOn {
			log.Printf("[DEBUG] Starting Linux Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			future, err := vmClient.Start(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
			if err != nil {
				return fmt.Errorf("Error starting Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
				return fmt.Errorf("Error waiting for start of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Started Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}
	} else { // otherwise, just update it
		diskFuture, err := client.BeginUpdate(ctx, resourceGroup, name, diskUpdate, nil)
		if err != nil {
			return fmt.Errorf("Error expanding managed disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		
		if _, err = diskFuture.PollUntilDone(ctx, common.DefaultPollingInterval); err != nil {
			return fmt.Errorf("Error waiting for expand operation on managed disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourceManagedDiskRead(d, meta)
}

func resourceManagedDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DiskName, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(err) {
			log.Printf("[INFO] Disk %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Managed Disk %s (resource group %s): %s", id.DiskName, id.ResourceGroup, err)
	}

	disk := resp.Disk
	d.Set("name", disk.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zones", utils.FlattenStringPtrSlice(disk.Zones))

	if location := disk.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := disk.SKU; sku != nil && sku.Name != nil {
		d.Set("storage_account_type", string(*sku.Name))
	}

	if props := disk.Properties; props != nil {
		if creationData := props.CreationData; creationData != nil {
			if creationData.CreateOption != nil {
				d.Set("create_option", string(*creationData.CreateOption))
			}

			imageReferenceID := ""
			if creationData.ImageReference != nil && creationData.ImageReference.ID != nil {
				imageReferenceID = *creationData.ImageReference.ID
			}
			d.Set("image_reference_id", imageReferenceID)

			d.Set("source_resource_id", creationData.SourceResourceID)
			d.Set("source_uri", creationData.SourceURI)
			d.Set("storage_account_id", creationData.StorageAccountID)
		}

		d.Set("disk_size_gb", props.DiskSizeGB)
		d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
		d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
		if props.OSType != nil {
			d.Set("os_type", props.OSType)
		}

		diskEncryptionSetId := ""
		if props.Encryption != nil && props.Encryption.DiskEncryptionSetID != nil {
			diskEncryptionSetId = *props.Encryption.DiskEncryptionSetID
		}
		d.Set("disk_encryption_set_id", diskEncryptionSetId)

		if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
			return fmt.Errorf("Error setting `encryption_settings`: %+v", err)
		}
	}

	return tags.Track2FlattenAndSet(d, disk.Tags)
}

func resourceManagedDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.BeginDelete(ctx, id.ResourceGroup, id.DiskName, nil)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Disk %q (Resource Group %q): %+v", id.DiskName, id.ResourceGroup, err)
	}

	if _, err = future.PollUntilDone(ctx, common.DefaultPollingInterval); err != nil {
		return fmt.Errorf("Error waiting for deletion of Managed Disk %q (Resource Group %q): %+v", id.DiskName, id.ResourceGroup, err)
	}

	return nil
}
