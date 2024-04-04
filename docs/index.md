---
page_title: "Provider: Paragon"
description: |-
  The Paragon provider is an interface to the Paragon service API.
---

# Paragon Provider

The Paragon provider is designed to work with the IPaaS vendor [Paragon](https://www.useparagon.com/).

Use the navigation to the left to read about the available resources.

~> **IMPORTANT:** 
This is not an official provider and is not supported by Paragon. This provider is maintained by the community and is not officially supported by HashiCorp. The APIs used by this provider are not officially supported by Paragon and may change at any time.


## Example Usage

```terraform

# Configure the connection details for the Paragon service
# Please create an "Admin" user for your organization.
provider "paragon" {
  username = "your_email"
  password = "your_password"
}

# Read your organization to get it by name.
data "paragon_organization" "my_org" {
  name = "my_paragon_organization"
}

# Create a new project.
resource "paragon_project" "main_pargon_project" {
  organization_id = data.paragon_organization.my_org.organization.id
  name            = "my_paragon_project"
}

# Create a new Paragon environment secret
resource "paragon_environment_secret" "my_secret" {
  project_id = paragon_project.main_pargon_project.id
  key        = "SECRET_KEY"
  value      = "SECRET_VALUE"
}

```
## Schema

### Required

- `username` (String) The email address of the paragon admin user.
- `password` (String) The password of the paragon admin user.

### Optional

- `base_url` (String) The base URL of the Paragon service. Default: `https://zeus.useparagon.com`.