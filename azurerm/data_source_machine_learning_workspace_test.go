package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMachineLearningWorkspace_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMachineLearningWorkspace_basic(ri, location),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccDataSourceMachineLearningWorkspace_basic(rInt int, location string) string {
	config := testAccAzureRMMachineLearningWorkspace_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_machine_learning_workspace" "test" {
  resource_group_name = "${azurerm_machine_learning_workspace.test.resource_group_name}"
  name                = "${azurerm_machine_learning_workspace.test.name}"
}
`, config)
}
