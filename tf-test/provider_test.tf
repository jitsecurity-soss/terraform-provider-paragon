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

# Get all organizations
data "paragon_organizations" "orgs" {}