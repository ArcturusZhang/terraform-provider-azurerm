---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
sidebar_current: "docs-azurerm-resource-nat-gateway"
description: |-
  Manage Azure NatGateway instance.
---

# azurerm_nat_gateway

Manage Azure NatGateway instance.


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the nat gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `idle_timeout_in_minutes` - (Optional) The idle timeout of the nat gateway.

* `public_ip_address_ids` - (Optional) A list of IDs of existing `azurerm_public_ip` resources.

* `public_ip_prefix_ids` - (Optional) A list of IDs of existing `azurerm_public_ip_prefix` resources.

* `resource_guid` - (Optional) The resource GUID property of the NAT gateway resource.

* `sku_name` - (Optional) The nat gateway SKU, supported values: Standard. Defaults to `Standard`.

* `zones` - (Optional) A list of availability zones denoting the zone in which Nat Gateway should be deployed. Changing this forces a new resource to be created.

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `subnet_ids` - A list of IDs of existing `azurerm_subnet` resources.
