---
page_title: "paragon_organization Data Source - paragon"
subcategory: ""
description: |-
  Fetches a specific organization by its name.
---

# paragon_organization (Data Source)

Fetches a specific organization by its name.

## Example Usage

```terraform
# Read a specific organization by name
data "paragon_organization" "example" {
  name = "your organization name"
}
```

## Errors
Error will be thrown if the organization is not found.

## Schema

### Argument Reference

- `name` (String, Required) The name of the organization to search.

### Attributes Reference

- `organization` (Attributes) The organization details.

The `organization` block contains:

- `id` (String) Identifier for the organization.
- `date_created` (String) The creation date of the organization.
- `date_updated` (String) The last update date of the organization.
- `name` (String) The name of the organization.
- `website` (String) The website of the organization.
- `type` (String) The type of the organization. (e.g. - BUSINESS)
- `purpose` (String) The purpose of the organization.
- `referral` (String) The referral of the organization. (e.g. - email)
- `size` (String) The size of the organization (e.g. - 1-10).
- `role` (String) The role of the owner of the organization.
- `completed_qualification` (Boolean) Indicates if the organization has completed qualification.

## JSON State Structure Example

Here's a state sample:

```json
{
    "name": "your_org_name",
    "organization": {
      "completed_qualification": true,
      "date_created": "2024-03-06T13:18:11.762Z",
      "date_updated": "2024-03-06T13:18:39.943Z",
      "id": "c1dbaa21-bf20-4131-a1b9-5072a4c78f7e",
      "name": "your org name",
      "purpose": "let our customers use integrations",
      "referral": "email",
      "role": "cto",
      "size": "10-49",
      "type": "BUSINESS",
      "website": "https://example.com"
    }
}
```
