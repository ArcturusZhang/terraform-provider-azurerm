package compute

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-12-01/armcompute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDiskAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiskAccessCreateUpdate,
		Read:   resourceDiskAccessRead,
		Update: resourceDiskAccessCreateUpdate,
		Delete: resourceDiskAccessDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DiskAccessID(id)
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

			"tags": tags.Schema(),
		},
	}
}

func resourceDiskAccessCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskAccessClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Disk Access creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, nil)
		if err != nil {
			if !utils.Track2ResponseWasNotFound(err) {
				return fmt.Errorf("Error checking for presence of existing Disk Access %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}
		if existing.DiskAccess != nil && existing.DiskAccess.ID != nil && *existing.DiskAccess.ID != "" {
			return tf.ImportAsExistsError("azurerm_disk_access", *existing.DiskAccess.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	createDiskAccess := armcompute.DiskAccess{
		Resource: armcompute.Resource{
			Location: &location,
			Tags:     tags.Track2Expand(t),
		},
	}

	future, err := client.BeginCreateOrUpdate(ctx, resourceGroup, name, createDiskAccess, nil)
	if err != nil {
		return fmt.Errorf("Error creating/updating Disk Access %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if _, err = future.PollUntilDone(ctx, common.DefaultPollingInterval); err != nil {
		return fmt.Errorf("Error waiting for create/update of Disk Access %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving Disk Access %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.DiskAccess == nil || read.DiskAccess.ID == nil {
		return fmt.Errorf("Error reading Disk Access %s (Resource Group %q): ID was nil", name, resourceGroup)
	}

	d.SetId(*read.DiskAccess.ID)

	return resourceDiskAccessRead(d, meta)
}

func resourceDiskAccessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskAccessClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DiskAccessID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(err) {
			log.Printf("[INFO] Disk Access %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Disk Access %s (resource group %s): %s", id.Name, id.ResourceGroup, err)
	}

	access := *resp.DiskAccess
	d.Set("name", access.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := access.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.Track2FlattenAndSet(d, access.Tags)
}

func resourceDiskAccessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskAccessClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DiskAccessID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.BeginDelete(ctx, id.ResourceGroup, id.Name, nil)
	if err != nil {
		return fmt.Errorf("Error deleting Disk Access %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if _, err = future.PollUntilDone(ctx, common.DefaultPollingInterval); err != nil {
		return fmt.Errorf("Error waiting for deletion of Disk Access %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
