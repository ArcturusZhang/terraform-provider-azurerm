package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDiskEncryptionSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDiskEncryptionSetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"active_key": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"previous_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDiskEncryptionSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Compute.DiskEncryptionSetsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Disk Encryption Set %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if encryptionSetProperties := resp.EncryptionSetProperties; encryptionSetProperties != nil {
		if err := d.Set("active_key", flattenArmDiskEncryptionSetKeyVaultAndKeyReference(encryptionSetProperties.ActiveKey)); err != nil {
			return fmt.Errorf("Error setting `active_key`: %+v", err)
		}
		if err := d.Set("previous_keys", flattenArmDiskEncryptionSetKeyVaultAndKeyReferenceArray(encryptionSetProperties.PreviousKeys)); err != nil {
			return fmt.Errorf("Error setting `previous_keys`: %+v", err)
		}
	}
	if identity := resp.Identity; identity != nil {
		if err := d.Set("identity", flattenArmDiskEncryptionSetIdentity(identity)); err != nil {
			return fmt.Errorf("Error setting `identity`: %+v", err)
		}
	}

	return nil
}
