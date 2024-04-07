---
page_title: "integration_status Resource - paragon"
subcategory: ""
description: |-
  Activates/deactivates an integration.

# integration_status (Resource)

Controls the enablement of an integration.

## Example Usage

Use `paragon_integrations` data source to find out the relevant `integration_id`.

```terraform
# Create credentials for integrating a service
resource "paragon_integration_status" "example" {
  integration_id = "your_integration_id"
  project_id = "your_project_id"
  active = true
}
```

## Schema

### Argument Reference

- `integration_id` (String, Required) Identifier of the integration to enable.
- `project_id` (String, Required) Identifier of the project of the integration to enable.
- `active` (Object, Required) weather the integration is active or not.

### Attributes Reference

- `id` (String) Identifier of the integration. Same as `integration_id`.