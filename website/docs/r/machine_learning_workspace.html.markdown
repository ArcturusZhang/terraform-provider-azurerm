---
subcategory: "MachineLearning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_workspace"
sidebar_current: "docs-azurerm-resource-machine-learning-workspace"
description: |-
  Manages a Azure Machine Learning Workspace.
---
# azurerm_machine_learning_workspace

Manages a Azure Machine Learning Workspace

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Machine Learning Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Machine Learning Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Optional) Specifies the supported Azure location where the Machine Learning Workspace should exist. Changing this forces a new resource to be created.

* `description` - (Optional) The description of this Machine Learning Workspace.

* `friendly_name` - (Optional) Friendly name for this Machine Learning Workspace.

* `key_vault_id` - (Required) The ID of key vault associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `application_insights_id` - (Required) The ID of the Application Insights associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) The ID of the Storage Account associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `container_registry_id` - (Optional) The ID of the container registry associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

-> **NOTE:** The `admin_enabled` should be `true` in order to associate the Container Registry to this Machine Learning Workspace.

* `discovery_url` - (Optional) The URL for the discovery service to identify regional endpoints for machine learning experimentation services.

* `sku_name` - (Optional) SKU/edition of the Machine Learning Workspace, possible values are `Basic` for a basic workspace or `Enterprise` for a feature rich workspace. Default to `Basic`.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Machine Learning Workspace.

## Import

Machine Learning Workspace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_workspace.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/workspace1
```
