package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-07-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNatGatewayCreateUpdate,
		Read:   resourceArmNatGatewayRead,
		Update: resourceArmNatGatewayCreateUpdate,
		Delete: resourceArmNatGatewayDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"idle_timeout_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
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
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// },

			// "public_ip_prefixes": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// },

			"resource_guid": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Standard),
				}, false),
				Default: string(network.Standard),
			},

			"tags": tags.Schema(),

			"zones": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmNatGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.NatGatewaysClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_nat_gateway", *resp.ID)
		}
	}

	// id := d.Get("id").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	idleTimeoutInMinutes := d.Get("idle_timeout_in_minutes").(int)
	publicIpAddressIds := d.Get("public_ip_address_ids").([]interface{})
	publicIpPrefixIds := d.Get("public_ip_prefix_ids").([]interface{})
	resourceGuid := d.Get("resource_guid").(string)
	sku := d.Get("sku").(string)
	zones := d.Get("zones").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	parameters := network.NatGateway{
		// ID:       utils.String(id),
		Location: utils.String(location),
		NatGatewayPropertiesFormat: &network.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: utils.Int32(int32(idleTimeoutInMinutes)),
			PublicIPAddresses:    expandArmNatGatewayIPSubResource(publicIpAddressIds),
			PublicIPPrefixes:     expandArmNatGatewayIPSubResource(publicIpPrefixIds),
			ResourceGUID:         utils.String(resourceGuid),
		},
		Sku: &network.NatGatewaySku{
			Name: network.NatGatewaySkuName(sku),
		},
		Tags:  tags.Expand(t),
		Zones: utils.ExpandStringSlice(zones),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Nat Gateway %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNatGatewayRead(d, meta)
}

func resourceArmNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.NatGatewaysClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["natGateways"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Nat Gateway %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

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

func resourceArmNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.NatGatewaysClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["natGateways"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmNatGatewaySubResource(input []interface{}) *[]network.SubResource {
	results := make([]network.SubResource, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		id := v["id"].(string)

		result := network.SubResource{
			ID: utils.String(id),
		}

		results = append(results, result)
	}
	return &results
}

func expandArmNatGatewayIPSubResource(input []interface{}) *[]network.SubResource {
	results := make([]network.SubResource, 0)
	for _, item := range input {
		id := item.(string)

		result := network.SubResource{
			ID: utils.String(id),
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmNatGatewaySubResource(input *[]network.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if id := item.ID; id != nil {
			v["id"] = *id
		}

		results = append(results, v)
	}

	return results
}

func flattenArmNatGatewayIPSubResource(input *[]network.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, *item.ID)
	}

	return results
}

// func flattenArmNatGatewaySubResource(input *[]network.SubResource) []interface{} {
// 	results := make([]interface{}, 0)
// 	if input == nil {
// 		return results
// 	}

// 	for _, item := range *input {
// 		v := make(map[string]interface{})

// 		results = append(results, v)
// 	}

// 	return results
// }
