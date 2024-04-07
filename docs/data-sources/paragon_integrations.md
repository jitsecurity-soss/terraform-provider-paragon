---
page_title: "paragon_integrations Data Source - paragon"
subcategory: ""
description: |-
  Fetches a list of integrations currently configured
---


# paragon_integrations Data Source - paragon

The `paragon_integrations` data source retrieves a list of currently configured integrations, presented in a key-value format. Each key represents the name of an integration, and the value provides detailed information about that integration.

## Example Usage

This example demonstrates how to retrieve the list of integrations for a specified project. The output is a structured map where each key is an integration name, and its value is an object containing details about the integration.

```terraform
data "paragon_integrations" "example" {
  project_id = "your_project_id"
}

// Output for example the id of jira integration
output "jira_integration_id" {
  value = data.paragon_integrations.example.integrations["jira"].id
}
```

## Schema

### Argument Reference

- `project_id` (Required, String): The ID of the project for which to fetch integrations.

### Attributes Reference

The following attributes are exported:

- `integrations` (Map of Objects): A map where each key is an integration name, and its value is an object with the following keys:
  - `id` (String): The unique identifier for the integration.
  - `authentication_type` (String): Specifies the type of authentication used. For custom integrations, this can be `oauth`, `oauth_client_credential`, or `basic`. For regular integrations, this will be an empty string.
  - `type` (String): The type of the integration (e.g., `jira`, `slack`). For custom integrations, this will be formatted as `custom.<integration_name>`.
  - `is_active` (Boolean): Indicates whether the integration is currently active.
  - `connected_user_count` (Number): The number of users connected to this integration.
  - `custom_integration_id` (String, Optional): The unique identifier for custom integrations. This will be present only for custom integrations.

## JSON State Structure Example

Here's a state sample:

```json
{
  "schema_version": 0,
  "attributes": {
    "integrations": {
      "custom.wiremock": {
        "authentication_type": "oauth",
        "connected_user_count": 0,
        "custom_integration_id": "adcc040f-46f0-4ee2-8d14-ef822258aeb8",
        "id": "7152a676-97a5-4ee6-b50c-b988d8d8f41c",
        "is_active": false,
        "type": "custom.wiremock"
      },
      "jira": {
        "authentication_type": "",
        "connected_user_count": 0,
        "custom_integration_id": "",
        "id": "e14c4be7d-a811-49c0-9a76-e717d2e5c10e",
        "is_active": false,
        "type": "jira"
      }
    },
    "project_id": "ce6c7134-98c9-48d6-a19c-30cd54dba758"
  },
  "sensitive_attributes": []
}
```
