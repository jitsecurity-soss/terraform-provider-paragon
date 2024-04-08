---
page_title: "integration_credentials Resource - paragon"
subcategory: ""
description: |-
  Manages credentials for an integration.
---

# paragon_integration_credentials (Resource)

Manages credentials for an integration.

~> **IMPORTANT:** 
The credentials should be stored securely and not exposed in any public repositories.

-> **NOTE:** Currently only oauth creds are supported, Trying to update non-oauth creds will result in an error

-> **NOTE:** For regular non-custom integration, there's no way verifying what type of authentication they required, so there's no restriction updating them.

## Example Usage

Use `paragon_integrations` data source to find out the relevant `integration_id`.

```terraform
# Create credentials for integrating a service
resource "paragon_integration_credentials" "example" {
  integration_id = "your_integration_id"
  project_id = "your_project_id"
  oauth = {
    client_id = "client_id"
    client_secret = "secret"
  }
}
```

## Schema

### Argument Reference

- `integration_id` (String, Required) Identifier of the integration for which to create credentials.
- `project_id` (String, Required) Identifier of the project for which to create credentials.
- `oauth` (Object, Required) OAuth credentials for the relevant OAuth service.
  - `client_id` (String, Required) Client ID for the OAuth service.
  - `client_secret` (String, Required) Client secret for the OAuth service.

### Attributes Reference

- `id` (String) The unique identifier of the credentials resource.
- `creds_provider` (String) Provider of the credentials (e.g., "custom" for custom integration, "jira").
- `scheme` (String) The scheme used for authentication (e.g., "oauth_app").

## JSON State Structure Example

Here's a state sample, Please make sure you keep the `client_secret' attribute secured

```json
{
    "creds_provider": "jira",
    "id": "credentials_id",
    "integration_id": "your_integration_id",
    "oauth": {
      "client_id": "your_client_id",
      "client_secret": "your_client_secret"
    },
    "project_id": "your_project_id",
    "scheme": "oauth_app"
}
```
