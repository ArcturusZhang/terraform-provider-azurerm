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
func TestAccDataSourceAzureRMNatGateway_complete(t *testing.T) {
	dataSourceName := "data.azurerm_nat_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNatGateway_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_address_ids.0", "${azurerm_public_ip.<%= resource_id_hint -%>.id}"),
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_prefix_ids.0", "${azurerm_public_ip_prefix.<%= resource_id_hint -%>.id}"),
					resource.TestCheckResourceAttr(dataSourceName, "sku", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "idle_timeout_in_minutes", "10"),
				),
			},
		},
	})
}
func TestAccDataSourceAzureRMNatGateway_update(t *testing.T) {
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
			{
				Config: testAccDataSourceNatGateway_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_address_ids.0", "${azurerm_public_ip.<%= resource_id_hint -%>.id}"),
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_prefix_ids.0", "${azurerm_public_ip_prefix.<%= resource_id_hint -%>.id}"),
					resource.TestCheckResourceAttr(dataSourceName, "sku", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "idle_timeout_in_minutes", "10"),
				),
			},
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

func testAccDataSourceNatGateway_complete(rInt int, location string) string {
	config := testAccAzureRMNatGateway_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_nat_gateway" "test" {
  resource_group_name = "${azurerm_nat_gateway.test.resource_group_name}"
  name                = "${azurerm_nat_gateway.test.name}"
}
`, config)
}
