package compute

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-12-01/armcompute"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedHostGroupCreate,
		Read:   resourceDedicatedHostGroupRead,
		Update: resourceDedicatedHostGroupUpdate,
		Delete: resourceDedicatedHostGroupDelete,

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
				ValidateFunc: validateDedicatedHostGroupName(),
			},

			"location": azure.SchemaLocation(),

			// There's a bug in the Azure API where this is returned in upper-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/8068
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"platform_fault_domain_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},

			// Currently only one endpoint is allowed.
			// we'll leave this open to enhancement when they add multiple zones support.
			"zones": azure.SchemaSingleZone(),

			"tags": tags.Schema(),
		},
	}
}

func resourceDedicatedHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name, nil)
		if err != nil {
			if !utils.Track2ResponseWasNotFound(err) {
				return fmt.Errorf("Error checking for present of existing Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.DedicatedHostGroup != nil && existing.DedicatedHostGroup.ID != nil && *existing.DedicatedHostGroup.ID != "" {
			return tf.ImportAsExistsError("azurerm_dedicated_host_group", *existing.DedicatedHostGroup.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	platformFaultDomainCount := d.Get("platform_fault_domain_count").(int)
	t := d.Get("tags").(map[string]interface{})

	parameters := armcompute.DedicatedHostGroup{
		Resource: armcompute.Resource{
			Location: utils.String(location),
			Tags: tags.Track2Expand(t),
		}, 
		Properties: &armcompute.DedicatedHostGroupProperties{
			PlatformFaultDomainCount: utils.Int32(int32(platformFaultDomainCount)),
		},
	}
	if zones, ok := d.GetOk("zones"); ok {
		parameters.Zones = utils.ExpandStringPtrSlice(zones.([]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters, nil); err != nil {
		return fmt.Errorf("Error creating Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, name, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if resp.DedicatedHostGroup == nil || resp.DedicatedHostGroup.ID == nil {
		return fmt.Errorf("Cannot read Dedicated Host Group %q (Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.DedicatedHostGroup.ID)

	return resourceDedicatedHostGroupRead(d, meta)
}

func resourceDedicatedHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["hostGroups"]

	resp, err := client.Get(ctx, resourceGroupName, name, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(err) {
			log.Printf("[INFO] Dedicated Host Group %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	group := *resp.DedicatedHostGroup
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := group.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := group.Properties; props != nil {
		platformFaultDomainCount := 0
		if props.PlatformFaultDomainCount != nil {
			platformFaultDomainCount = int(*props.PlatformFaultDomainCount)
		}
		d.Set("platform_fault_domain_count", platformFaultDomainCount)
	}
	d.Set("zones", utils.FlattenStringPtrSlice(group.Zones))

	return tags.Track2FlattenAndSet(d, group.Tags)
}

func resourceDedicatedHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := armcompute.DedicatedHostGroupUpdate{
		UpdateResource: armcompute.UpdateResource{
			Tags: tags.Track2Expand(t),
		},
	}

	if _, err := client.Update(ctx, resourceGroupName, name, parameters, nil); err != nil {
		return fmt.Errorf("Error updating Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	return resourceDedicatedHostGroupRead(d, meta)
}

func resourceDedicatedHostGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["hostGroups"]

	if _, err := client.Delete(ctx, resourceGroup, name, nil); err != nil {
		return fmt.Errorf("Error deleting Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}
