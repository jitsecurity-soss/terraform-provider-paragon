---
page_title: "paragon_team Data Source - paragon"
subcategory: ""
description: |-
  Fetches a team by its ID.
---

# paragon_team (Data Source)

Fetches a team by its name.

-> **NOTE:** In the console, there's not really a concept of teams. Whenever a new project is created, a team is created with the same name as the project. This data source is used to fetch the list of the teams if a team_id is required at any part. If you change the project name, team name will not change - so always use original team name to fetch a team.

-> **NOTE:** There's no restriction to create 2 teams with the same name (in paragon) - so this data source will return the first team it finds with the given name.

## Example Usage

```terraform
# Read a specific team by ID
data "paragon_team" "example" {
  name = "your_team_name"
}
```

## Schema

### Argument Reference

- `name` (String, Required) The name of the team to search.

### Attributes Reference
- `id` (String) Identifier for the team.
- `date_created` (String) The creation date of the team.
- `date_updated` (String) The last update date of the team.
- `website` (String) The website of the team.
- `organization_id` (String) The ID of the organization the team belongs to.

## JSON State Structure Example

Here's a state sample:

```json
{
  "date_created": "2024-03-21T17:37:39.902Z",
  "date_updated": "2024-03-21T17:37:39.902Z",
  "id": "your_team_id",
  "name": "tean_name",
  "organization_id": "your_org_id",
  "website": ""
}
```
