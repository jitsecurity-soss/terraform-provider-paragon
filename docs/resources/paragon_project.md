---
page_title: "paragon_project Resource - paragon"
subcategory: ""
description: |-
  Manages a paragon project.
---

# paragon_project (Resource)

Manages a paragon [project](https://docs-prod.useparagon.com/deploying-integrations/projects).

-> **NOTE:** When creating a project, behind the hood a team is created, and another "older" type of projects (automate) - its ID is saved as reference.


## Example Usage

```terraform
# Create a new project
resource "paragon_project" "example" {
  organization_id = "your_organization_id"
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
  "id": "project_id",
  "automate_project_id": "older_project_id",
  "duplicate_name_allowed": true,
  "is_connect_project": true,
  "is_hidden": false,
  "organization_id": "your_org_id",
  "owner_id": "your_user_id",
  "team_id": "your_team_id",
  "title": "project_title"
}
```
