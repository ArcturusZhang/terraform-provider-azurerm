package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMNatGateway_basic(t *testing.T) {
	dataSourceName := "data.azurerm_nat_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNatGateway_basic(ri, location),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccDataSourceNatGateway_basic(rInt int, location string) string {
	config := testAccAzureRMNatGateway_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_nat_gateway" "test" {
  resource_group_name = "${azurerm_nat_gateway.test.resource_group_name}"
  name                = "${azurerm_nat_gateway.test.name}"
}
`, config)
}
