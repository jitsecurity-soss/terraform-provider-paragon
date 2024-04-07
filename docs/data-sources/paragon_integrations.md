---
page_title: "paragon_integrations Data Source - paragon"
subcategory: ""
description: |-
  Fetches a list of integrations currently configured
---

# paragon_integrations (Data Source)

Fetches a list of integrations currently configured

## Example Usage

```terraform
# Read the list of teams
data "paragon_integrations" "example" {
  project_id = "your_project_id"
}
```

## Schema

### Argument Reference

- `project_id` (Required, String) The ID of the project.

### Attributes Reference

- `integrations` (Attributes List) The list of integrations.

The `integration` block contains:
- `id` (String) Identifier for the integration.
- `authentication_type` (String) ***Only for custom integration - for regular ones, it will be empty.*** `oauth` for oauth2, `oauth_client_credential` or `basic` for api keys.
- `type` (String) The type of the integration, such as `jira`, `slack`. 
  - **For custom integrations** - it will always be `custom.<integration_name>`
- `is_active` (Boolean) Indicates if the integration is active for the customers
- `connected_user_count` (Number) The number of users connected to the integration