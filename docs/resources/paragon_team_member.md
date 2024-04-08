---
page_title: "paragon_team_member Resource - paragon"
subcategory: ""
description: |-
  Manages a team member.
---

# paragon_team_member (Resource)

Manages a [team member](https://docs-prod.useparagon.com/managing-account/teams), This resource handles either if the member already accepted the invitation or not.

-> **NOTE:** `team_id` can be retrieved from the `paragon_project` resource, or `paragon_teams`/`paragon_team` data source.

-> **NOTE:** It appears even though team members are creted under a team - only the invites are "per team (project)", but once a user had accepted the invitation - it has the same permissions to all other projects.

## Example Usage

```terraform
# Option 1 - by team name, note that if you changed project name - the team name will stay the original team name.
data "paragon_team" "team" {
  name = "my_team_name"
}

resource "paragon_team_member" "team_member" {
  team_id = data.paragon_team.team.id
  email   = "example@example.com"
  role    = "MEMBER"
}

# option 2 - from the project you created
resource "paragon_project" "my_proj" {
  organization_id = "org-id"
  name            = "Example Project"
}

resource "paragon_team_member" "team_member" {
  team_id = paragon_project.my_proj.team_id
  email   = "example@example.com"
  role    = "MEMBER"
}
```

## Errors
Email must be unique, This resource blocks the option to create 2 team members with the same email.

## Schema

### Argument Reference

- `team_id` (String) Identifier of the team, Can be retrieved from `paragon_teams` data source or `paragon_project` resource.
- `email` (String) Email address of the team member.
- `role` (String) Role of the team member (ADMIN, MEMBER, SUPPORT).

### Attributes Reference

- `id` (String) Identifier of the team member. This can change after the user accepts the invitation.

## JSON State Structure Example

Here's a state sample:

```json
{
  "email": "example@example.com",
  "id": "user_id",
  "role": "MEMBER",
  "team_id": "team_id"
}
```