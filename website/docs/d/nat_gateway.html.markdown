---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
sidebar_current: "docs-azurerm-datasource-nat-gateway"
description: |-
  Gets information about an existing Nat Gateway
---

# Data Source: azurerm_nat_gateway

Use this data source to access information about an existing Nat Gateway.



## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the nat gateway.

* `resource_group_name` - (Required) The Name of the Resource Group where the App Service exists.


## Attributes Reference

The following attributes are exported:

* `location` - Resource location.

* `idle_timeout_in_minutes` - The idle timeout of the nat gateway.

* `public_ip_address_ids` - A list of IDs of existing `azurerm_public_ip` resources.

* `public_ip_prefix_ids` - A list of IDs of existing `azurerm_public_ip_prefix` resources.

* `resource_guid` - The resource GUID property of the NAT gateway resource.

* `sku_name` - The nat gateway SKU, supported values: Standard.

* `subnet_ids` - A list of IDs of existing `azurerm_subnet` resources.

* `type` - Resource type.

* `zones` - A list of availability zones denoting the zone in which Nat Gateway should be deployed.

* `tags` - Resource tags.
