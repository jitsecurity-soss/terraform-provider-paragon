---
page_title: "paragon_environment_secret Resource - paragon"
subcategory: ""
description: |-
  Manages an environment secret.
---

# paragon_environment_secret (Resource)

Manages an [environment secret](https://docs-prod.useparagon.com/workflows/environment-secrets). Those values can then be used inside workflows.

-> **NOTE:** `key` argument cannot be updated, it will cause recreation of the resource.

## Example Usage

```terraform
# Create a new environment secret
resource "paragon_environment_secret" "example" {
  project_id = "08ae44e3-d506-4c0e-87b0-a6934aa2f3a1"
  key        = "SECRET_KEY"
  value      = "SECRET_VALUE"
}
```

## Schema

### Argument Reference

- `project_id` (String, Required) Identifier of the project.
- `key` (String, Required) Key of the environment secret.
- `value` (String, Required, Sensitive) Value of the environment secret.

### Attributes Reference

- `id` (String) Identifier of the environment secret.
- `hash` (String) Hash of the environment secret.

## JSON State Structure Example

Here's a **full** state sample, Note that the input value is marked as sensitive attribute.

```json
{
    "hash": "secret_hash",
    "id": "2c24d3db-cc78-48db-b0ec-61c70f25ebc2",
    "key": "secret_name",
    "project_id": "08ae44e3-d506-4c0e-87b0-a6934aa2f3a1",
    "value": "secret_value"
}
```
