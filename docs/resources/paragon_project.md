---
page_title: "paragon_project Resource - paragon"
subcategory: ""
description: |-
  Manages a paragon project.
---

# paragon_project (Resource)

Manages a paragon [project](https://docs-prod.useparagon.com/deploying-integrations/projects).

-> **NOTE:** When creating a project, behind the hood a team is created, and another "older" type of projects (automate) - its ID is saved as reference.

-> **NOTE:** Only the owner of the project can delete it.

## Example Usage

```terraform
# Create a new project
resource "paragon_project" "example" {
  organization_id = "caad9cc6-2914-429d-b6e4-5150e2efb981"
  title           = "Example Project"
}
```

## Schema

### Argument Reference

- `organization_id` (String, Required) Identifier of the organization.
- `title` (String, Required) Name of the project.
- `duplicate_name_allowed` (String, Optional) Indicates whether creating another project with the same name is allowed. (Default = False)

### Attributes Reference

- `id` (String) Identifier of the project.
- `automate_project_id` (String) A hidden older project is created, this ID is for reference to delete it as well if needed.
- `owner_id` (String) Identifier of the project owner.
- `team_id` (String) Identifier of the team associated with the project.
- `is_connect_project` (Boolean) Indicates if the project is a Connect project - This is always true.
- `is_hidden` (Boolean) Indicates if the project is hidden.


## JSON State Structure Example

Here's a state sample:

```json
{
  "id": "40a0685f-ca69-4b1e-8468-a895b2cc0f94",
  "automate_project_id": "df234c5f-d7f4-4667-8838-4aa6701197db",
  "duplicate_name_allowed": true,
  "is_connect_project": true,
  "is_hidden": false,
  "organization_id": "caad9cc6-2914-429d-b6e4-5150e2efb981",
  "owner_id": "user_id_that_created_the_project",
  "team_id": "236fab2b-f92f-459b-98c1-aa676b943681",
  "title": "project_title"
}
```
