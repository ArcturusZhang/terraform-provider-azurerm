package compute

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-12-01/armcompute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceSnapshotCreateUpdate,
		Read:   resourceSnapshotRead,
		Update: resourceSnapshotCreateUpdate,
		Delete: resourceSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateSnapshotName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"create_option": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Copy),
					string(compute.Import),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"source_uri": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"storage_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"disk_size_gb": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"encryption_settings": encryptionSettingsSchema(),

			"tags": tags.Schema(),
		},
	}
}

func resourceSnapshotCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	createOption := d.Get("create_option").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, nil)
		if err != nil {
			if !utils.Track2ResponseWasNotFound(err) {
				return fmt.Errorf("Error checking for presence of existing Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.Snapshot != nil && existing.Snapshot.ID != nil && *existing.Snapshot.ID != "" {
			return tf.ImportAsExistsError("azurerm_snapshot", *existing.Snapshot.ID)
		}
	}

	snapshot := armcompute.Snapshot{
		Resource: armcompute.Resource{
			Location: utils.String(location),
			Tags: tags.Track2Expand(t),
		},
		Properties: &armcompute.SnapshotProperties{
			CreationData: &armcompute.CreationData{
				CreateOption: armcompute.DiskCreateOption(createOption).ToPtr(),
			},
		},
	}

	if v, ok := d.GetOk("source_uri"); ok {
		snapshot.Properties.CreationData.SourceURI = utils.String(v.(string))
	}

	if v, ok := d.GetOk("source_resource_id"); ok {
		snapshot.Properties.CreationData.SourceResourceID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		snapshot.Properties.CreationData.StorageAccountID = utils.String(v.(string))
	}

	diskSizeGB := d.Get("disk_size_gb").(int)
	if diskSizeGB > 0 {
		snapshot.Properties.DiskSizeGB = utils.Int32(int32(diskSizeGB))
	}

	if v, ok := d.GetOk("encryption_settings"); ok {
		encryptionSettings := v.([]interface{})
		settings := encryptionSettings[0].(map[string]interface{})
		snapshot.Properties.EncryptionSettingsCollection = expandManagedDiskEncryptionSettings(settings)
	}

	future, err := client.BeginCreateOrUpdate(ctx, resourceGroup, name, snapshot, nil)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if _, err := future.PollUntilDone(ctx, common.DefaultPollingInterval); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		return fmt.Errorf("Error issuing get request for Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.Snapshot.ID)

	return resourceSnapshotRead(d, meta)
}

func resourceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["snapshots"]

	resp, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(err) {
			log.Printf("[INFO] Error reading Snapshot %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Snapshot %q: %+v", name, err)
	}
	
	snapshot := resp.Snapshot
	d.Set("name", snapshot.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := snapshot.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := snapshot.Properties; props != nil {
		if data := props.CreationData; data != nil {
			if data.CreateOption != nil {
				d.Set("create_option", string(*data.CreateOption))
			}

			if accountId := data.StorageAccountID; accountId != nil {
				d.Set("storage_account_id", accountId)
			}
		}

		if props.DiskSizeGB != nil {
			d.Set("disk_size_gb", int(*props.DiskSizeGB))
		}

		if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
			return fmt.Errorf("Error setting `encryption_settings`: %+v", err)
		}
	}

	return tags.Track2FlattenAndSet(d, snapshot.Tags)
}

func resourceSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["snapshots"]

	future, err := client.BeginDelete(ctx, resourceGroup, name, nil)
	if err != nil {
		return fmt.Errorf("Error deleting Snapshot: %+v", err)
	}

	if _, err := future.PollUntilDone(ctx, common.DefaultPollingInterval); err != nil {
		return fmt.Errorf("Error deleting Snapshot: %+v", err)
	}

	return nil
}

func ValidateSnapshotName(v interface{}, _ string) (warnings []string, errors []error) {
	// a-z, A-Z, 0-9, _ and -. The max name length is 80
	value := v.(string)

	if !regexp.MustCompile("^[A-Za-z0-9_-]+$").MatchString(value) {
		errors = append(errors, fmt.Errorf("Snapshot Names can only contain alphanumeric characters and underscores."))
	}

	length := len(value)
	if length > 80 {
		errors = append(errors, fmt.Errorf("Snapshot Name can be up to 80 characters, currently %d.", length))
	}

	return warnings, errors
}
