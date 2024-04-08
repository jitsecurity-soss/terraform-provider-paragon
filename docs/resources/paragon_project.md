---
page_title: "paragon_project Resource - paragon"
subcategory: ""
description: |-
  Manages a paragon project.
---

# paragon_project (Resource)

Manages a paragon [project](https://docs-prod.useparagon.com/deploying-integrations/projects).


## Example Usage

```terraform
# Create a new project
resource "paragon_project" "example" {
  organization_id = "your_organization_id"
  name            = "Example Project"
}
```

## Schema

### Argument Reference

- `organization_id` (String, Required) Identifier of the organization.
- `name` (String, Required) Name of the project.
- `duplicate_name_allowed` (String, Optional) Indicates whether creating another project with the same name is allowed. (Default = False)

### Attributes Reference

- `id` (String) Identifier of the project.
- `title` (String) Title of the project.
- `owner_id` (String) Identifier of the project owner.
- `team_id` (String) Identifier of the team associated with the project.
- `is_connect_project` (Boolean) Indicates if the project is a Connect project - This is always true.
- `is_hidden` (Boolean) Indicates if the project is hidden.
- `date_created` (String) Date when the project was created.
- `date_updated` (String) Date when the project was last updated.


## JSON State Structure Example

Here's a state sample:

```json
{
    "duplicate_name_allowed": false,
    "id": "your_project_id",
    "is_connect_project": true,
    "is_hidden": false,
    "name": "dev",
    "organization_id": "your_organization_id",
    "owner_id": "your_user_id",
    "team_id": "your_team_id",
    "title": "project_name"
}
```
