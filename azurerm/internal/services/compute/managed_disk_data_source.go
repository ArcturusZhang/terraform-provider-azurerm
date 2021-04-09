package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceManagedDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceManagedDiskRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"create_option": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_encryption_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_iops_read_write": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"disk_mbps_read_write": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"disk_size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"image_reference_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_account_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),

			"zones": azure.SchemaZonesComputed(),
		},
	}
}

func dataSourceManagedDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(err) {
			return fmt.Errorf("Error: Managed Disk %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("[ERROR] Error making Read request on Azure Managed Disk %q (Resource Group %q): %s", name, resGroup, err)
	}

	disk := *resp.Disk
	d.SetId(*disk.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	if sku := disk.SKU; sku != nil {
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
			d.Set("os_type", string(*props.OSType))
		}

		diskEncryptionSetId := ""
		if props.Encryption != nil && props.Encryption.DiskEncryptionSetID != nil {
			diskEncryptionSetId = *props.Encryption.DiskEncryptionSetID
		}
		d.Set("disk_encryption_set_id", diskEncryptionSetId)
	}

	d.Set("zones", utils.FlattenStringPtrSlice(disk.Zones))

	return tags.Track2FlattenAndSet(d, disk.Tags)
}
