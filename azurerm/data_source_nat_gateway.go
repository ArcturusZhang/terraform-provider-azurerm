package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmNatGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNatGatewayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"idle_timeout_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"public_ip_address_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"public_ip_prefix_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// "public_ip_addresses": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 		},
			// 	},
			// },

			// "public_ip_prefixes": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 		},
			// 	},
			// },

			"resource_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.NatGatewaysClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Nat Gateway %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("sku", resp.Sku.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if natGatewayPropertiesFormat := resp.NatGatewayPropertiesFormat; natGatewayPropertiesFormat != nil {
		d.Set("idle_timeout_in_minutes", natGatewayPropertiesFormat.IdleTimeoutInMinutes)
		if err := d.Set("public_ip_address_ids", flattenArmNatGatewayIPSubResource(natGatewayPropertiesFormat.PublicIPAddresses)); err != nil {
			return fmt.Errorf("Error setting `public_ip_address_ids`: %+v", err)
		}
		if err := d.Set("public_ip_prefix_ids", flattenArmNatGatewayIPSubResource(natGatewayPropertiesFormat.PublicIPPrefixes)); err != nil {
			return fmt.Errorf("Error setting `public_ip_prefix_ids`: %+v", err)
		}
		d.Set("resource_guid", natGatewayPropertiesFormat.ResourceGUID)
		if err := d.Set("subnets", flattenArmNatGatewaySubResource(natGatewayPropertiesFormat.Subnets)); err != nil {
			return fmt.Errorf("Error setting `subnets`: %+v", err)
		}
	}
	d.Set("type", resp.Type)
	d.Set("zones", utils.FlattenStringSlice(resp.Zones))
	tags.FlattenAndSet(d, resp.Tags)

	return nil
}
