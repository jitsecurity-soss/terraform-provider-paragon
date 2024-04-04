---
page_title: "paragon_team Data Source - paragon"
subcategory: ""
description: |-
  Fetches a team by its ID.
---

# paragon_team (Data Source)

Fetches a team by its ID.

-> **NOTE:** In the console, there's not really a concept of teams. Whenever a new project is created, a team is created with the same name as the project. This data source is used to fetch the list of the teams if a team_id is required at any part. If you change the project name, team name will not change.

## Example Usage

```terraform
# Read a specific team by ID
data "paragon_team" "example" {
  id = "your_team_id"
}
```

## Schema

### Argument Reference

- `id` (String, Required) Identifier for the team.

### Attributes Reference

- `date_created` (String) The creation date of the team.
- `date_updated` (String) The last update date of the team.
- `name` (String) The name of the team.
- `website` (String) The website of the team.
- `organization_id` (String) The ID of the organization the team belongs to.