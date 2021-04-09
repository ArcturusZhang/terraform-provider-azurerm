package compute

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-12-01/armcompute"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAvailabilitySet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAvailabilitySetCreateUpdate,
		Read:   resourceAvailabilitySetRead,
		Update: resourceAvailabilitySetCreateUpdate,
		Delete: resourceAvailabilitySetDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,78}[a-zA-Z0-9_])?$"),
					"The Availability set name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 80 characters, and it must begin a letter or number and end with a letter, number or underscore.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"platform_update_domain_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 20),
			},

			"platform_fault_domain_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},

			"managed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"proximity_placement_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,

				// We have to ignore case due to incorrect capitalisation of resource group name in
				// proximity placement group ID in the response we get from the API request
				//
				// todo can be removed when https://github.com/Azure/azure-sdk-for-go/issues/5699 is fixed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAvailabilitySetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.AvailabilitySetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Availability Set creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, nil)
		if err != nil {
			if !utils.Track2ResponseWasNotFound(err) {
				return fmt.Errorf("Error checking for presence of existing Availability Set %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.AvailabilitySet != nil && existing.AvailabilitySet.ID != nil && *existing.AvailabilitySet.ID != "" {
			return tf.ImportAsExistsError("azurerm_availability_set", *existing.AvailabilitySet.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	updateDomainCount := d.Get("platform_update_domain_count").(int)
	faultDomainCount := d.Get("platform_fault_domain_count").(int)
	managed := d.Get("managed").(bool)
	t := d.Get("tags").(map[string]interface{})

	availSet := armcompute.AvailabilitySet{
		Resource: armcompute.Resource{
			Location: &location,
			Tags:     tags.Track2Expand(t),
		},
		Properties: &armcompute.AvailabilitySetProperties{
			PlatformFaultDomainCount:  utils.Int32(int32(faultDomainCount)),
			PlatformUpdateDomainCount: utils.Int32(int32(updateDomainCount)),
		},
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		availSet.Properties.ProximityPlacementGroup = &armcompute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if managed {
		n := "Aligned"
		availSet.SKU = &armcompute.SKU{
			Name: &n,
		}
	}

	resp, err := client.CreateOrUpdate(ctx, resGroup, name, availSet, nil)
	if err != nil {
		return err
	}

	d.SetId(*resp.AvailabilitySet.ID)

	return resourceAvailabilitySetRead(d, meta)
}

func resourceAvailabilitySetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.AvailabilitySetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AvailabilitySetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Availability Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	set := resp.AvailabilitySet
	d.Set("name", set.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := set.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if set.SKU != nil && set.SKU.Name != nil {
		d.Set("managed", strings.EqualFold(*set.SKU.Name, "Aligned"))
	}

	if props := set.Properties; props != nil {
		d.Set("platform_update_domain_count", props.PlatformUpdateDomainCount)
		d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)

		if proximityPlacementGroup := props.ProximityPlacementGroup; proximityPlacementGroup != nil {
			d.Set("proximity_placement_group_id", proximityPlacementGroup.ID)
		}
	}

	return tags.Track2FlattenAndSet(d, set.Tags)
}

func resourceAvailabilitySetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.AvailabilitySetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AvailabilitySetID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, id.Name, nil)
	return err
}
