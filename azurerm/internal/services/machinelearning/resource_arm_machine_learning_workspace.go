package machinelearning

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2019-11-01/machinelearningservices"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"application_insights_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"container_registry_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"discovery_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Basic",
					"Enterprise",
				}, true),
				Default: "Basic",
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
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	storageAccount := d.Get("storage_account_id").(string)
	keyVault := d.Get("key_vault_id").(string)
	applicationInsights := d.Get("application_insights_id").(string)
	skuName := d.Get("sku_name").(string)
	t := d.Get("tags").(map[string]interface{})

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for existing AML Workspace %q (Resource Group %q): %s", name, resourceGroup, err)
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_machine_learning_workspace", *existing.ID)
		}
	}

	// TODO -- it turns out that you cannot use premium storage accounts
	workspace := machinelearningservices.Workspace{
		Name:     &name,
		Location: &location,
		Tags:     tags.Expand(t),
		WorkspaceProperties: &machinelearningservices.WorkspaceProperties{
			StorageAccount:      &storageAccount,
			ApplicationInsights: &applicationInsights,
			KeyVault:            &keyVault,
		},
		Identity: expandAmlIdentity(d),
		Sku:      &machinelearningservices.Sku{Name: utils.String(skuName)},
	}

	if v, ok := d.GetOk("description"); ok {
		workspace.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("friendly_name"); ok {
		workspace.FriendlyName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("discovery_url"); ok {
		workspace.DiscoveryURL = utils.String(v.(string))
	}

	if v, ok := d.GetOk("container_registry_id"); ok {
		workspace.ContainerRegistry = utils.String(v.(string))
	}

	if err := validateStorageAccount(storageAccount, resourceGroup, d, meta); err != nil {
		return fmt.Errorf("Error creating Machine Learning Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := validateContainerRegistry(workspace.ContainerRegistry, resourceGroup, d, meta); err != nil {
		return fmt.Errorf("Error creating Machine Learning Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	result, err := client.CreateOrUpdate(ctx, resourceGroup, name, workspace)
	if err != nil {
		return fmt.Errorf("Error during workspace creation %q in resource group (%q): %+v", name, resourceGroup, err)
	}

	fmt.Printf("created AML Workspace %q", result.Name)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read workspace %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmAmlWorkspaceRead(d, meta)
}

func validateStorageAccount(accountID string, resourceGroup string, d *schema.ResourceData, meta interface{}) error {
	if accountID == "" {
		return fmt.Errorf("Error validating Storage Account: Empty ID")
	}
	id, err := azure.ParseAzureResourceID(accountID)
	if err != nil {
		return fmt.Errorf("Error validating Storage Account: %+v", err)
	}
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	accountName := id.Path["storageAccounts"]
	account, err := client.GetProperties(ctx, resourceGroup, accountName, "")
	if err != nil {
		return fmt.Errorf("Error validating Storage Account %q (Resource Group %q): %+v", accountName, resourceGroup, err)
	}
	if sku := account.Sku; sku != nil {
		if sku.Tier == storage.Premium {
			return fmt.Errorf("Error validating Storage Account %q (Resource Group %q): The associated Storage Account must not be Premium", accountName, resourceGroup)
		}
	}
	return nil
}

func validateContainerRegistry(acrID *string, resourceGroup string, d *schema.ResourceData, meta interface{}) error {
	if acrID == nil {
		return nil
	}
	id, err := azure.ParseAzureResourceID(*acrID)
	if err != nil {
		return fmt.Errorf("Error validating Container Registry: %+v", err)
	}
	client := meta.(*clients.Client).Containers.RegistriesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	acrName := id.Path["registries"]
	acr, err := client.Get(ctx, resourceGroup, acrName)
	if err != nil {
		return fmt.Errorf("Error validating Container Registry %q (Resource Group %q): %+v", acrName, resourceGroup, err)
	}
	if acr.AdminUserEnabled == nil || !*acr.AdminUserEnabled {
		return fmt.Errorf("Error validating Container Registry%q (Resource Group %q): The associated Container Registry must set `admin_enabled` to true", acrName, resourceGroup)
	}
	return nil
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
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

func flattenAmlIdentity(identity *machinelearningservices.Identity) []interface{} {
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
	return []interface{}{result}
}

func expandAmlIdentity(d *schema.ResourceData) *machinelearningservices.Identity {
	v := d.Get("identity").([]interface{})
	// Rest api will return an error if you did not set the identity field, if not set return default
	if len(v) == 0 {
		return &machinelearningservices.Identity{Type: machinelearningservices.SystemAssigned}
	}
	identity := v[0].(map[string]interface{})
	identityType := machinelearningservices.ResourceIdentityType(identity["type"].(string))

	amlIdentity := machinelearningservices.Identity{
		Type: identityType,
	}

	return &amlIdentity
}
