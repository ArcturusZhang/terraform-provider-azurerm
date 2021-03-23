package msi

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/msi/2018-11-30/armmsi"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmUserAssignedIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmUserAssignedIdentityCreateUpdate,
		Read:   resourceArmUserAssignedIdentityRead,
		Update: resourceArmUserAssignedIdentityCreateUpdate,
		Delete: resourceArmUserAssignedIdentityDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.UserAssignedIdentityID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			migration.UserAssignedIdentityV0ToV1(),
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
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),

			"principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmUserAssignedIdentityCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSI.UserAssignedIdentitiesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for User Assigned Identity create/update.")

	location := d.Get("location").(string)
	t := d.Get("tags").(map[string]interface{})

	resourceId := parse.NewUserAssignedIdentityID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.Name, nil)
		log.Printf("[DAPZHANG] %+v", err)
		if err != nil {
			if !utils.Track2ResponseWasNotFound(err) {
				return fmt.Errorf("checking for presence of existing User Assigned Identity %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
			}
		}

		if existing.Identity != nil && existing.Identity.ID != nil && *existing.Identity.ID != "" {
			return tf.ImportAsExistsError("azurerm_user_assigned_identity", resourceId.ID())
		}
	}

	identity := armmsi.Identity{
		TrackedResource: armmsi.TrackedResource{
			Location: utils.String(location),
			Tags:     tags.Track2Expand(t),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.Name, identity, nil); err != nil {
		return fmt.Errorf("creating/updating User Assigned Identity %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceArmUserAssignedIdentityRead(d, meta)
}

func resourceArmUserAssignedIdentityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSI.UserAssignedIdentitiesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.UserAssignedIdentityID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil)
	if err != nil {
		if utils.Track2ResponseWasNotFound(resp.RawResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving User Assigned Identity %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	identity := resp.Identity
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(identity.Location))

	if props := identity.Properties; props != nil {
		if principalId := props.PrincipalID; principalId != nil {
			d.Set("principal_id", principalId)
		}

		if clientId := props.ClientID; clientId != nil {
			d.Set("client_id", clientId)
		}
	}

	return tags.Track2FlattenAndSet(d, identity.Tags)
}

func resourceArmUserAssignedIdentityDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSI.UserAssignedIdentitiesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.UserAssignedIdentityID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.Name, nil); err != nil {
		return fmt.Errorf("deleting User Assigned Identity %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
