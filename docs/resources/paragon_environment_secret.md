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
  project_id = paragon_project.example.id
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