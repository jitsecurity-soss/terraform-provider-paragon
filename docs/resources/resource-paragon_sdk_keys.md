---
page_title: "paragon_sdk_keys Resource - paragon"
subcategory: ""
description: |-
  Manages SDK keys for a project.
---

# paragon_sdk_keys (Resource)

Manages SDK keys for a project. The private key that is created should be used to sign JWT tokens described as [paragon user](https://docs.useparagon.com/billing/connected-users). 

~> **IMPORTANT:** 
The private key that is created should be stored securely.

-> **NOTE:** paragon supports 3 types of sdk keys - signing keys, Auth0 or Firebase - This resource only supports signing keys (`paragon` auth type).

## Example Usage

```terraform
# Create an SDK key for the project
resource "paragon_sdk_keys" "example" {
  project_id = paragon_project.example.id
  version    = "1"
}
```

## Schema

### Argument Reference

- `project_id` (String, Required) Identifier of the project to create keys for.
- `version` (String, Required) Version of the SDK key - Change this value after creation to recreate the keys.

### Attributes Reference

- `id` (String) Identifier of the SDK key.
- `auth_type` (String) Authentication type of the SDK key. (e.g. - paragon)
- `revoked` (Boolean) Indicates if the SDK key is revoked.
- `generated_date` (String) Date when the SDK key was generated.
- `private_key` (String, Sensitive) Private key of the SDK key.