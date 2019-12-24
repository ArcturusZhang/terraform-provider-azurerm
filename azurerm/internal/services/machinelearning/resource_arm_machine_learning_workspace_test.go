package machinelearning

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMMachineLearningWorkspace_basic(t *testing.T) {
	resourceName := "azurerm_machine_learning_workspace.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspace_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(resourceName, "application_insights_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMMachineLearningWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Machine Learning Workspace not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).MachineLearning.WorkspacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Machine Learning Workspace %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on machinelearningservices.WorkspacesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMachineLearningWorkspaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MachineLearning.WorkspacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_machine_learning_workspace" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on machinelearningservices.WorkspacesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMMachineLearningWorkspace_basic(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mlservices-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                 = "acctestvault%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  tenant_id            = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "premium"
}

resource "azurerm_storage_account" "sa" {
  name                     = "acctestsa%d"
  location                 = "${azurerm_resource_group.rg.location}"
  resource_group_name      = "${azurerm_resource_group.rg.name}"
  account_tier             = "Premium"
  account_replication_type = "LRS"
}

resource "azurerm_application_insights" "test" {
  name                 = "acctestai-%d"
  location             = "${azurerm_resource_group.rg.location}"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  application_type     = "web"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctestworkspace-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  key_vault_id            = "${azurerm_key_vault.test.id}"
  storage_account_id      = "${azurerm_storage_account.test.id}"
  application_insights_id = "${azurerm_application_insights.test.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
