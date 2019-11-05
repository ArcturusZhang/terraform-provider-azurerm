package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMDiskEncryptionSet_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDiskEncryptionSet_basic(ri, location),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}
func TestAccDataSourceAzureRMDiskEncryptionSet_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDiskEncryptionSet_complete(ri, location),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccDataSourceDiskEncryptionSet_basic(rInt int, location string) string {
	config := testAccAzureRMDiskEncryptionSet_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_disk_encryption_set" "test" {
  resource_group_name = "${azurerm_disk_encryption_set.test.resource_group_name}"
  name                = "${azurerm_disk_encryption_set.test.name}"
}
`, config)
}

func testAccDataSourceDiskEncryptionSet_complete(rInt int, location string) string {
	config := testAccAzureRMDiskEncryptionSet_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_disk_encryption_set" "test" {
  resource_group_name = "${azurerm_disk_encryption_set.test.resource_group_name}"
  name                = "${azurerm_disk_encryption_set.test.name}"
}
`, config)
}
