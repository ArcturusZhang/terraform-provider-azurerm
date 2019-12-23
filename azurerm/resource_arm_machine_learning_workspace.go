package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2019-11-01/machinelearningservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAmlWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAmlWorkspaceCreateUpdate,
		Read:   resourceArmAmlWorkspaceRead,
		Update: resourceArmAmlWorkspaceCreateUpdate,
		Delete: resourceArmAmlWorkspaceDelete,

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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"application_insights_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"container_registry_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"discovery_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tags.Schema(),

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(machinelearningservices.SystemAssigned),
							}, false),
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
		},
	}
}

func resourceArmAmlWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	description := d.Get("description").(string)
	friendlyName := d.Get("friendly_name").(string)
	storageAccount := d.Get("storage_account_id").(string)
	keyVault := d.Get("key_vault_id").(string)
	containerRegistry := d.Get("container_registry_id").(string)
	applicationInsights := d.Get("application_insights_id").(string)
	discoveryUrl := d.Get("discovery_url").(string)
	t := d.Get("tags").(map[string]interface{})

	existing, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for existing AML Workspace %q (Resource Group %q): %s", name, resGroup, err)
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_machine_learning_workspace", *existing.ID)
		}
	}

	// TODO -- should validate container registry enable_admin is enabled.
	workspace := machinelearningservices.Workspace{
		Name:     &name,
		Location: &location,
		Tags:     tags.Expand(t),
		WorkspaceProperties: &machinelearningservices.WorkspaceProperties{
			Description:         &description,
			FriendlyName:        &friendlyName,
			StorageAccount:      &storageAccount,
			DiscoveryURL:        &discoveryUrl,
			ContainerRegistry:   &containerRegistry,
			ApplicationInsights: &applicationInsights,
			KeyVault:            &keyVault,
		},
		Identity: expandAmlIdentity(d),
	}

	result, err := client.CreateOrUpdate(ctx, resGroup, name, workspace)
	if err != nil {
		return fmt.Errorf("Error during workspace creation %q in resource group (%q): %+v", name, resGroup, err)
	}

	fmt.Printf("created AML Workspace %q", result.Name)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read workspace %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmAmlWorkspaceRead(d, meta)
}

func resourceArmAmlWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["machineLearningServices"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.WorkspaceProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("storage_account_id", props.StorageAccount)
		d.Set("discovery_url", props.DiscoveryURL)
		d.Set("container_registry_id", props.ContainerRegistry)
		d.Set("application_insights_id", props.ApplicationInsights)
		d.Set("key_vault_id", props.KeyVault)
	}
	if err := d.Set("identity", flattenAmlIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error flattening identity on Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAmlWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx := meta.(*clients.Client).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["machineLearningServices"]

	_, err = client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func flattenAmlIdentity(identity *machinelearningservices.Identity) interface{} {
	if identity == nil {
		return nil
	}
	result := make(map[string]interface{})
	result["type"] = string(identity.Type)
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}
	if identity.TenantID != nil {
		result["tenant_id"] = *identity.TenantID
	}
	return result
}

func expandAmlIdentity(d *schema.ResourceData) *machinelearningservices.Identity {
	v := d.Get("identity")
	// Rest api will return an error if you did not set the identity field.
	if v == nil {
		return &machinelearningservices.Identity{Type: machinelearningservices.SystemAssigned}
	}
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := machinelearningservices.ResourceIdentityType(identity["type"].(string))

	amlIdentity := machinelearningservices.Identity{
		Type: identityType,
	}

	return &amlIdentity
}
