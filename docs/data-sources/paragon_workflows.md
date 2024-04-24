---
page_title: "paragon_workflows Data Source - paragon"
subcategory: ""
description: |-
  Returns list of workflows associated with a project and integration.
---

# paragon_workflow (Data Source)

Returns list of workflows associated with a project and integration.

## Example Usage

```terraform
# Read a specific workflow by project ID, integration ID, and description
data "paragon_workflows" "example" {
  project_id     = "c555a650-cd0b-4782-ae66-674517a12fb0"
  integration_id = "461a6e87-0cd5-4eb2-b2c8-6585f7077fdb"
}
```

## Errors
Error will be thrown if the workflow is not found.

## Schema

### Argument Reference

- `project_id` (String, Required) The ID of the project.
- `integration_id` (String, Required) The ID of the integration.


### Attributes Reference

- `workflows` (Attributes List) The list of teams.

The `workflows` block contains:

- `id` (String)  The ID of the workflow.
- `project_id` (String) The ID of the project.
- `integration_id` (String) The ID of the integration.- 
- `description` (String) The description of the workflow to search for.
- `date_created` (String) The creation date of the workflow.
- `date_updated` (String) The last update date of the workflow.
- `tags` (List of String) The tags associated with the workflow.
- `workflow_version` (Number) The version of the workflow.

## JSON State Structure Example

Here's a state sample:

```json
{
  "workflows": [
    {
      "date_created": "2024-04-15T10:59:44.207Z",
      "date_updated": "2024-04-17T07:47:10.659Z",
      "description": "Create tickets from issues",
      "id": "461a6e87-0cd5-4eb2-b2c8-6585f7077fdb",
      "integration_id": "461a6e87-0cd5-4eb2-b2c8-6585f7077fdb",
      "project_id": "c555a650-cd0b-4782-ae66-674517a12fb0",      
      "tags": [],
      "workflow_version": 0
    },
    {
      "date_created": "2024-04-24T15:53:46.039Z",
      "date_updated": "2024-04-24T15:53:50.544Z",
      "description": "New Workflow",
      "id": "6cdad43e-3090-4d48-83bb-cb1563fb7789",
      "integration_id": "461a6e87-0cd5-4eb2-b2c8-6585f7077fdb",
      "project_id": "c555a650-cd0b-4782-ae66-674517a12fb0",
      "tags": [],
      "workflow_version": 1
    }
  ]
}
```
