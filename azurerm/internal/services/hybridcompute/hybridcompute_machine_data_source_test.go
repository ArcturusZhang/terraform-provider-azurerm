package hybridcompute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMhybridcomputeMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hybridcompute_machine", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMhybridcomputeMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcehybridcomputeMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMhybridcomputeMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "tags"),
				),
			},
		},
	})
}

func testAccDataSourcehybridcomputeMachine_basic(data acceptance.TestData) string {
	config := testAccAzureRMhybridcomputeMachine_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hybridcompute_machine" "test" {
  name = azurerm_hybridcompute_machine.test.name
  resource_group_name = azurerm_hybridcompute_machine.test.resource_group_name
}
`, config)
}
