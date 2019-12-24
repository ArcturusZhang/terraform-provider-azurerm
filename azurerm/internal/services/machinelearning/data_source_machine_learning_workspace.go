package machinelearning

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAMLWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAMLWorkspaceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_vault_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"application_insights_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"container_registry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"discovery_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
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
		},
	}
}

func dataSourceArmAMLWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Machine Learning Workspace %q (Resource Group %q) was not found: %+v", name, resourceGroup, err)
		}
		return fmt.Errorf("Error reading Machine Learning Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Machine Learning Workspace %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("description", resp.Description)
	d.Set("friendly_name", resp.FriendlyName)
	d.Set("key_vault_id", resp.KeyVault)
	d.Set("application_insights_id", resp.ApplicationInsights)
	d.Set("container_registry_id", resp.ContainerRegistry)
	d.Set("storage_account_id", resp.StorageAccount)
	d.Set("discovery_url", resp.DiscoveryURL)
	if err := d.Set("identity", flattenAmlIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
