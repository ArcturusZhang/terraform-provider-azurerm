package hybridcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hybridcompute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HybridComputeMachineResource struct{}

func TestAccHybridComputeMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine", "test")
	r := HybridComputeMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHybridComputeMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine", "test")
	r := HybridComputeMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccHybridComputeMachine_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine", "test")
	r := HybridComputeMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMhybridcomputeMachine_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine", "test")
	r := HybridComputeMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHybridComputeMachine_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine", "test")
	r := HybridComputeMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHybridComputeMachine_updateLocationData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine", "test")
	r := HybridComputeMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateLocationData(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r HybridComputeMachineResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.HybridComputeMachineID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.HybridCompute.MachineClient.Get(ctx, id.ResourceGroup, id.MachineName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Hybridcompute Machine %q (Resource Group %q): %+v", id.MachineName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r HybridComputeMachineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-hybridcompute-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r HybridComputeMachineResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybrid_compute_machine" "test" {
  name = "acctest-hm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (r HybridComputeMachineResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybrid_compute_machine" "import" {
  name = azurerm_hybrid_compute_machine.test.name
  resource_group_name = azurerm_hybrid_compute_machine.test.resource_group_name
  location = azurerm_hybrid_compute_machine.test.location
}
`, config)
}

func (r HybridComputeMachineResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybrid_compute_machine" "test" {
  name = "acctest-hm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  client_public_key = "string"
  virtual_machine_id = "b7a098cc-b0b8-46e8-a205-62f301a62a8f"
  identity {
    type = "SystemAssigned"
  }

  location_datum {
    name = "Redmond"
    city = ""
    country_or_region = ""
    district = ""
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r HybridComputeMachineResource) updateIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybrid_compute_machine" "test" {
  name = "acctest-hm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  client_public_key = "string"
  virtual_machine_id = "b7a098cc-b0b8-46e8-a205-62f301a62a8f"
  identity {
    type = "SystemAssigned"
  }

  location_datum {
    name = "Redmond"
    city = ""
    country_or_region = ""
    district = ""
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r HybridComputeMachineResource) updateLocationData(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybrid_compute_machine" "test" {
  name = "acctest-hm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  client_public_key = "string"
  virtual_machine_id = "b7a098cc-b0b8-46e8-a205-62f301a62a8f"
  identity {
    type = "SystemAssigned"
  }

  location_datum {
    name = "Redmond"
    city = ""
    country_or_region = ""
    district = ""
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
