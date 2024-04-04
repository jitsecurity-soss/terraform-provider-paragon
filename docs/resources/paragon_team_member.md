---
page_title: "paragon_team_member Resource - paragon"
subcategory: ""
description: |-
  Manages a team member.
---

# paragon_team_member (Resource)

Manages a [team member](https://docs-prod.useparagon.com/managing-account/teams), This resource handles either if the member already accepted the invitation or not.

-> **NOTE:** `key` argument cannot be updated, it will cause recreation of the resource.

-> **NOTE:** `ADMIN` roles administer the entire organziation, while `MEMBER` and `SUPPORT` roles belong to a specific project/team.

## Example Usage

```terraform
# Create a new team member
resource "paragon_team_member" "example" {
  team_id = "your_team_id"
  email   = "team_member@example.com"
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