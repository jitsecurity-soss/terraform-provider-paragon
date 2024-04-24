---
page_title: "paragon_workflow Data Source - paragon"
subcategory: ""
description: |-
  Fetches a specific workflow by description.
---

# paragon_workflow (Data Source)

Fetches a specific workflow by description.

## Example Usage

```terraform
# Read a specific workflow by project ID, integration ID, and description
data "paragon_workflow" "example" {
  project_id     = "e0da0789-cd90-4ca7-897b-8a89404eb329"
  integration_id = "fb549b70-658b-4a14-9318-4dca3a88bfa7"
  description    = "your workflow description"
}
```

## Errors
Error will be thrown if the workflow is not found.

## Schema

### Argument Reference

- `project_id` (String, Required) The ID of the project.
- `integration_id` (String, Required) The ID of the integration.
- `description` (String, Required) The description of the workflow to search for.

### Attributes Reference
- `id` (String)  The ID of the workflow.
- `date_created` (String) The creation date of the workflow.
- `date_updated` (String) The last update date of the workflow.
- `tags` (List of String) The tags associated with the workflow.
- `workflow_version` (Number) The version of the workflow.

## JSON State Structure Example

Here's a state sample:

```json
{
  "id": "workflow_id"  
  "project_id": "e0da0789-cd90-4ca7-897b-8a89404eb329",
  "integration_id": "fb549b70-658b-4a14-9318-4dca3a88bfa7",
  "description": "Create tickets from issues", 
  "date_created": "2024-04-15T10:59:44.207Z", 
  "date_updated": "2024-04-17T07:47:10.659Z", 
  "tags": [], 
  "workflow_version": 0  
}
```
