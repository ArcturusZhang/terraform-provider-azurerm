package hybridcompute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hybridcompute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMhybridcomputeMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybridcompute_machine", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMhybridcomputeMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMhybridcomputeMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMhybridcomputeMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybridcompute_machine", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMhybridcomputeMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMhybridcomputeMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMhybridcomputeMachine_requiresImport),
		},
	})
}

func TestAccAzureRMhybridcomputeMachine_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybridcompute_machine", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMhybridcomputeMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMhybridcomputeMachine_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMhybridcomputeMachine_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybridcompute_machine", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMhybridcomputeMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMhybridcomputeMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMhybridcomputeMachine_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMhybridcomputeMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMhybridcomputeMachine_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybridcompute_machine", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMhybridcomputeMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMhybridcomputeMachine_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMhybridcomputeMachine_updateIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMhybridcomputeMachine_updateLocationData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybridcompute_machine", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMhybridcomputeMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMhybridcomputeMachine_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMhybridcomputeMachine_updateLocationData(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMhybridcomputeMachineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).HybridCompute.MachineClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("hybridcompute Machine not found: %s", resourceName)
		}
		id, err := parse.HybridComputeMachineID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.MachineName, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Hybridcompute Machine %q does not exist", id.MachineName)
			}
			return fmt.Errorf("bad: Get on Hybridcompute.MachineClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMhybridcomputeMachineDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).HybridCompute.MachineClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hybridcompute_machine" {
			continue
		}
		id, err := parse.HybridComputeMachineID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.MachineName, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Hybridcompute.MachineClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMhybridcomputeMachine_template(data acceptance.TestData) string {
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

func testAccAzureRMhybridcomputeMachine_basic(data acceptance.TestData) string {
	template := testAccAzureRMhybridcomputeMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybridcompute_machine" "test" {
  name = "acctest-hm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMhybridcomputeMachine_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMhybridcomputeMachine_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybridcompute_machine" "import" {
  name = azurerm_hybridcompute_machine.test.name
  resource_group_name = azurerm_hybridcompute_machine.test.resource_group_name
  location = azurerm_hybridcompute_machine.test.location
}
`, config)
}

func testAccAzureRMhybridcomputeMachine_complete(data acceptance.TestData) string {
	template := testAccAzureRMhybridcomputeMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybridcompute_machine" "test" {
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

func testAccAzureRMhybridcomputeMachine_updateIdentity(data acceptance.TestData) string {
	template := testAccAzureRMhybridcomputeMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybridcompute_machine" "test" {
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

func testAccAzureRMhybridcomputeMachine_updateLocationData(data acceptance.TestData) string {
	template := testAccAzureRMhybridcomputeMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hybridcompute_machine" "test" {
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
