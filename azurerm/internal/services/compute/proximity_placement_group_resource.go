package compute

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-12-01/armcompute"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceProximityPlacementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceProximityPlacementGroupCreateUpdate,
		Read:   resourceProximityPlacementGroupRead,
		Update: resourceProximityPlacementGroupCreateUpdate,
		Delete: resourceProximityPlacementGroupDelete,

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceProximityPlacementGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ProximityPlacementGroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Proximity Placement Group creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, nil)
		if err != nil {
			if !utils.Track2ResponseWasNotFound(err) {
				return fmt.Errorf("Error checking for presence of existing Proximity Placement Group %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ProximityPlacementGroup != nil && existing.ProximityPlacementGroup.ID != nil && *existing.ProximityPlacementGroup.ID != "" {
			return tf.ImportAsExistsError("azurerm_proximity_placement_group", *existing.ProximityPlacementGroup.ID)
		}
	}

	ppg := armcompute.ProximityPlacementGroup{
		Resource: armcompute.Resource{
			Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
			Tags:     tags.Track2Expand(d.Get("tags").(map[string]interface{})),
		},
	}

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, name, ppg, nil)
	if err != nil {
		return err
	}

	d.SetId(*resp.ProximityPlacementGroup.ID)

	return resourceProximityPlacementGroupRead(d, meta)
}

func resourceProximityPlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ProximityPlacementGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["proximityPlacementGroups"]

	resp, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Proximity Placement Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	ppg := resp.ProximityPlacementGroup
	d.Set("name", ppg.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := ppg.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.Track2FlattenAndSet(d, ppg.Tags)
}

func resourceProximityPlacementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ProximityPlacementGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["proximityPlacementGroups"]

	_, err = client.Delete(ctx, resGroup, name, nil)
	return err
}
