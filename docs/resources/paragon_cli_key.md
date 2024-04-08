---
page_title: "paragon_cli_key Resource - paragon"
subcategory: ""
description: |-
  Manages a CLI key.
---

# paragon_cli_key (Resource)

Manages a CLI key. The CLI key can be used to work with [Paragraph](https://docs.useparagon.com/paragraph/getting-started).

~> **IMPORTANT:** 
The key that is created should be stored securely.

-> **NOTE:** This resource will prohibit creating several keys with the same name and the same user used in the provider authentication.

-> **NOTE:** CLI keys are organization-wide, and can be used to interact with all projects within the organization.

## Example Usage

```terraform
# Create a new CLI key
resource "paragon_cli_key" "example" {
  organization_id = "your_organization_id"
  name            = "example_cli_key"
}
```

## Schema

### Argument Reference

- `organization_id` (String) Identifier of the organization.
- `name` (String) Name of the CLI key.

### Attributes Reference

- `id` (String) Identifier of the CLI key.
- `key` (String, Sensitive) The CLI key.


## JSON State Structure Example

Here's a state sample, Please make sure you keep the `key' attribute secured

```json
{
  "id": "cli_key_id",
  "key": "cli_key.XXXXXXXXXXX",
  "name": "key_name",
  "organization_id": "your_org_id"
}
```
