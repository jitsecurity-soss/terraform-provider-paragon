---
page_title: "paragon_teams Data Source - paragon"
subcategory: ""
description: |-
  Fetches a list of paragon teams.
---

# paragon_teams (Data Source)

Fetches a list of paragon teams.

-> **NOTE:** In the console, there's not really a concept of teams. Whenever a new project is created, a team is created with the same name as the project. This data source is used to fetch the list of the teams if a team_id is required at any part. If you change the project name, team name will not change.


## Example Usage

```terraform
# Read the list of teams
data "paragon_teams" "example" {}
```

## Schema

### Attributes Reference

- `teams` (Attributes List) The list of teams.

The `teams` block contains:

- `id` (String) Identifier for the team.
- `date_created` (String) The creation date of the team.
- `date_updated` (String) The last update date of the team.
- `name` (String) The name of the team.
- `website` (String) The website of the team.
- `organization_id` (String) The ID of the organization the team belongs to.


## JSON State Structure Example

Here's a state sample:

```json
{
  "teams": [
    {
      "date_created": "2024-03-21T17:37:39.902Z",
      "date_updated": "2024-03-21T17:37:39.902Z",
      "id": "c8fbefd4-6d54-4c82-9951-78aa1d92bd50",
      "name": "team_name",
      "organization_id": "c1dbaa21-bf20-4131-a1b9-5072a4c78f7e",
      "website": ""
    },
    {
      "date_created": "2024-04-21T17:37:39.902Z",
      "date_updated": "2024-04-21T17:37:39.902Z",
      "id": "b1f9827e-e880-4a29-930c-f0eeb1f142d1",
      "name": "team_name2",
      "organization_id": "c1dbaa21-bf20-4131-a1b9-5072a4c78f7e",
      "website": ""
    }    
  ]
}
```
