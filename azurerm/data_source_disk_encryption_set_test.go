package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMDiskEncryptionSet_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	vaultName := fmt.Sprintf("vault%d", ri)
	keyName := fmt.Sprintf("key-%s", rs)
	desName := fmt.Sprintf("acctestdes-%d", ri)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccDataSourceDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName string) string {
	config := testAccAzureRMDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName)
	return fmt.Sprintf(`
%s

data "azurerm_disk_encryption_set" "test" {
  resource_group_name = "${azurerm_disk_encryption_set.test.resource_group_name}"
  name                = "${azurerm_disk_encryption_set.test.name}"
}
`, config)
}
