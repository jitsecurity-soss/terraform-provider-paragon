terraform {
  required_providers {
    paragon = {
      source = "arielb135/paragon"
    }
  }
}

provider "paragon" {
  username = "XX@xx.xx"
  password = "XX"
}

# Organizations Data Source
data "paragon_organizations" "all" {}

output "organizations" {
  value = data.paragon_organizations.all.organizations
}

# Teams Data Source
data "paragon_teams" "all" {}

output "teams" {
  value = data.paragon_teams.all.teams
}

# Team Data Source
data "paragon_team" "example" {
  id = "team-id"
}

output "team" {
  value = data.paragon_team.example
}

# Project Resource
resource "paragon_project" "example" {
  organization_id = "org-id"
  name            = "Example Project"
}

output "project" {
  value = paragon_project.example
}

# SDK Keys Resource
resource "paragon_sdk_keys" "example" {
  project_id = paragon_project.example.id
  version    = "1.0.0"
}

output "sdk_keys" {
  value     = paragon_sdk_keys.example
  sensitive = true
}

# Environment Secret Resource
resource "paragon_environment_secret" "example" {
  project_id = paragon_project.example.id
  key        = "SECRET_KEY"
  value      = "secret_value"
}

output "environment_secret" {
  value     = paragon_environment_secret.example
  sensitive = true
}

# Team Member Resource
resource "paragon_team_member" "example" {
  team_id = "team-id"
  email   = "user@example.com"
  role    = "MEMBER"
}

output "team_member" {
  value = paragon_team_member.example
}

# CLI Key Resource
resource "paragon_cli_key" "example" {
  organization_id = "org-id"
  name            = "Example CLI Key"
}

output "cli_key" {
  value     = paragon_cli_key.example
  sensitive = true
}