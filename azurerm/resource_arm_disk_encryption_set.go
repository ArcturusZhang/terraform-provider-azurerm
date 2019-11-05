package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDiskEncryptionSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDiskEncryptionSetCreateUpdate,
		Read:   resourceArmDiskEncryptionSetRead,
		Update: resourceArmDiskEncryptionSetCreateUpdate,
		Delete: resourceArmDiskEncryptionSetDelete,

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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"active_key": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"source_vault_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"identity_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.SystemAssigned),
				}, false),
				Default: string(compute.SystemAssigned),
			},

			"previous_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"source_vault_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDiskEncryptionSetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Compute.DiskEncryptionSetsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_disk_encryption_set", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	activeKey := d.Get("active_key").([]interface{})
	identityType := d.Get("identity_type").(string)
	t := d.Get("tags").(map[string]interface{})

	diskEncryptionSet := compute.DiskEncryptionSet{
		Identity: &compute.EncryptionSetIdentity{
			Type: compute.DiskEncryptionSetIdentityType(identityType),
		},
		Location: utils.String(location),
		EncryptionSetProperties: &compute.EncryptionSetProperties{
			ActiveKey: expandArmDiskEncryptionSetKeyVaultAndKeyReference(activeKey),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, diskEncryptionSet)
	if err != nil {
		return fmt.Errorf("Error creating Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Disk Encryption Set %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmDiskEncryptionSetRead(d, meta)
}

func resourceArmDiskEncryptionSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Compute.DiskEncryptionSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["diskEncryptionSets"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Disk Encryption Set %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

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
		d.Set("identity_type", string(identity.Type))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDiskEncryptionSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Compute.DiskEncryptionSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["diskEncryptionSets"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmDiskEncryptionSetKeyVaultAndKeyReference(input []interface{}) *compute.KeyVaultAndKeyReference {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	sourceVaultId := v["source_vault_id"].(string)
	keyUrl := v["key_url"].(string)

	result := compute.KeyVaultAndKeyReference{
		KeyURL: utils.String(keyUrl),
		SourceVault: &compute.SourceVault{
			ID: utils.String(sourceVaultId),
		},
	}
	return &result
}

func flattenArmDiskEncryptionSetKeyVaultAndKeyReference(input *compute.KeyVaultAndKeyReference) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if keyUrl := input.KeyURL; keyUrl != nil {
		result["key_url"] = *keyUrl
	}
	if sourceVault := input.SourceVault; sourceVault != nil {
		if sourceVaultId := sourceVault.ID; sourceVaultId != nil {
			result["source_vault_id"] = *sourceVaultId
		}
	}

	return []interface{}{result}
}

func flattenArmDiskEncryptionSetKeyVaultAndKeyReferenceArray(input *[]compute.KeyVaultAndKeyReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if sourceVault := item.SourceVault; sourceVault != nil {
			if sourceVaultId := sourceVault.ID; sourceVaultId != nil {
				v["source_vault_id"] = *sourceVaultId
			}
		}

		results = append(results, v)
	}

	return results
}
