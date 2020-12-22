package hybridcompute

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/hybridcompute/mgmt/2020-08-02/hybridcompute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hybridcompute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hybridcompute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceHybridComputeMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceHybridComputeMachineCreate,
		Read:   resourceHybridComputeMachineRead,
		Update: resourceHybridComputeMachineUpdate,
		Delete: resourceHybridComputeMachineDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.HybridComputeMachineID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.HybridComputeMachineName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"client_public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
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

			"location_data": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"city": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"country_or_region": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"district": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"agent_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_profile": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"computer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"os_sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vm_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceHybridComputeMachineCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).HybridCompute.MachineClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewHybridComputeMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.MachineName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Hybridcompute Machine %q (Resource Group %q): %+v", id.MachineName, id.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_hybrid_compute_machine", id.ID())
	}

	identity := expandArmMachineIdentity(d.Get("identity").([]interface{}))

	parameter := hybridcompute.Machine{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		MachinePropertiesModel: &hybridcompute.MachinePropertiesModel{
			ClientPublicKey: utils.String(d.Get("client_public_key").(string)),
			VMID:            utils.String(d.Get("virtual_machine_id").(string)),
			LocationData:    expandArmMachineLocationData(d.Get("location_data").([]interface{})),
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MachineName, parameter); err != nil {
		return fmt.Errorf("creating Hybridcompute Machine %q (Resource Group %q): %+v", id.MachineName, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MachineName, "")
	if err != nil {
		return fmt.Errorf("retrieving Hybridcompute Machine %q (Resource Group %q): %+v", id.MachineName, id.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Hybridcompute Machine %q (Resource Group %q) ID", id.MachineName, id.ResourceGroup)
	}

	d.SetId(id.ID())

	return resourceHybridComputeMachineRead(d, meta)
}

func resourceHybridComputeMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HybridCompute.MachineClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridComputeMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MachineName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] hybridcompute %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Hybridcompute Machine %q (Resource Group %q): %+v", id.MachineName, id.ResourceGroup, err)
	}
	d.Set("name", id.MachineName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.MachinePropertiesModel; props != nil {
		d.Set("client_public_key", props.ClientPublicKey)
		d.Set("virtual_machine_id", props.VMID)
		if err := d.Set("location_data", flattenArmMachineLocationData(props.LocationData)); err != nil {
			return fmt.Errorf("setting `.location_data`: %+v", err)
		}
		d.Set("agent_version", props.AgentVersion)
		d.Set("display_name", props.DisplayName)
		d.Set("machine_fqdn", props.MachineFqdn)
		d.Set("dns_fqdn", props.DNSFqdn)
		d.Set("domain_name", props.DomainName)
		d.Set("os_name", props.OsName)
		if err := d.Set("os_profile", flattenArmMachineOsProfile(props.OsProfile)); err != nil {
			return fmt.Errorf("setting `os_profile`: %+v", err)
		}
		d.Set("os_sku", props.OsSku)
		d.Set("os_version", props.OsVersion)
		d.Set("vm_uuid", props.VMUUID)
		d.Set("ad_fqdn", props.AdFqdn)
	}
	if err := d.Set("identity", flattenArmMachineIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceHybridComputeMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HybridCompute.MachineClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridComputeMachineID(d.Id())
	if err != nil {
		return err
	}

	parameter := hybridcompute.MachineUpdate{
		MachineUpdateIdentity:                &hybridcompute.MachineUpdateIdentity{},
		MachineUpdatePropertiesModel: &hybridcompute.MachineUpdatePropertiesModel{},
	}
	if d.HasChange("location_data") {
		parameter.MachineUpdatePropertiesModel.LocationData = expandArmMachineLocationData(d.Get("location_data").([]interface{}))
	}
	if d.HasChange("tags") {
		parameter.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.MachineName, parameter); err != nil {
		return fmt.Errorf("updating Hybridcompute Machine %q (Resource Group %q): %+v", id.MachineName, id.ResourceGroup, err)
	}
	return resourceHybridComputeMachineRead(d, meta)
}

func resourceHybridComputeMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HybridCompute.MachineClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridComputeMachineID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.MachineName); err != nil {
		return fmt.Errorf("deleting Hybridcompute Machine %q (Resource Group %q): %+v", id.MachineName, id.ResourceGroup, err)
	}
	return nil
}

func expandArmMachineIdentity(input []interface{}) *hybridcompute.MachineIdentity {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &hybridcompute.MachineIdentity{
		Type: utils.String(v["type"].(string)),
	}
}

func expandArmMachineLocationData(input []interface{}) *hybridcompute.LocationData {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &hybridcompute.LocationData{
		Name:            utils.String(v["name"].(string)),
		City:            utils.String(v["city"].(string)),
		District:        utils.String(v["district"].(string)),
		CountryOrRegion: utils.String(v["country_or_region"].(string)),
	}
}

func flattenArmMachineIdentity(input *hybridcompute.MachineIdentity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var t string
	if input.Type != nil {
		t = *input.Type
	}
	var principalId string
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}
	var tenantId string
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}
	return []interface{}{
		map[string]interface{}{
			"type":         t,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
}

func flattenArmMachineLocationData(input *hybridcompute.LocationData) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}
	var city string
	if input.City != nil {
		city = *input.City
	}
	var countryOrRegion string
	if input.CountryOrRegion != nil {
		countryOrRegion = *input.CountryOrRegion
	}
	var district string
	if input.District != nil {
		district = *input.District
	}
	return []interface{}{
		map[string]interface{}{
			"name":              name,
			"city":              city,
			"country_or_region": countryOrRegion,
			"district":          district,
		},
	}
}

func flattenArmMachineOsProfile(input *hybridcompute.MachinePropertiesOsProfile) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var computerName string
	if input.ComputerName != nil {
		computerName = *input.ComputerName
	}
	return []interface{}{
		map[string]interface{}{
			"computer_name": computerName,
		},
	}
}
