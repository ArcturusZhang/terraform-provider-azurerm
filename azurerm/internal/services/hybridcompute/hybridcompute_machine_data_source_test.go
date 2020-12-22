package hybridcompute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type HybridComputeMachineDataSource struct{}

func TestAccDataSourceHybridComputeMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hybrid_compute_machine", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: HybridComputeMachineDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("tags").Exists(),
			),
		},
	})
}

func (d HybridComputeMachineDataSource) basic(data acceptance.TestData) string {
	config := HybridComputeMachineResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hybrid_compute_machine" "test" {
  name = azurerm_hybrid_compute_machine.test.name
  resource_group_name = azurerm_hybrid_compute_machine.test.resource_group_name
}
`, config)
}
