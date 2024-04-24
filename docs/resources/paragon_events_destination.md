---
page_title: "paragon_events_destination Resource - paragon"
subcategory: ""
description: |-
  Manages an events destination (Alerts)
---

# paragon_events_destination (Resource)

Set up a destination to send events about your Paragon project to your logging, incident management, or analytics services.

Manages an [events destination](https://docs.useparagon.com/monitoring/event-destinations).

### Supported destinations
1. Webhook - which sends a webhook with a custom payload to a URL, optionally adding headers, This supports variable substitution from the event itself.
2. Email - sends the incident details to a specific email address.

### Supported events
- `workflow_failure` - Any unhandled error in a workflow.

-> **NOTE:** For webhooks - Payload structure is not verified in this provider, It's recommended to verify it after creation 

-> **NOTE:** This resource does not manage whether the destination is enabled or disabled, It's enabled by default when created.

## Example Usage

### Webhook Destination

```terraform
resource "paragon_events_destination" "webhook_example" {
  project_id = "a7321f97-9c6a-437d-b51e-bd4ce549635f"
  events     = ["workflow_failure"]

  webhook {
    url = "https://example.com/webhook"
    headers = {
      "Content-Type" = "application/json"
      "Authorization" = "Bearer my-auth-token"
    }
    body = <<EOF
{
  "message": "Workflow failed: {{$.event.workflow.name}}",
  "timestamp": "{{$.event.timestamp}}"
}
EOF
  }
}
```

### Email Destination
```terraform
resource "paragon_events_destination" "email_example" {
  project_id = "a7321f97-9c6a-437d-b51e-bd4ce549635f"
  events     = ["workflow_failure"]

  email {
    address = "notifications@example.com"
  }
}
```

## Webhook templates

### Datadog
``` terraform
resource "paragon_events_destination" "webhook_example" {
  project_id = "a7321f97-9c6a-437d-b51e-bd4ce549635f"
  events     = ["workflow_failure"]

  webhook = {
    url = "https://http-intake.logs.<DD-SITE>/api/v2/logs"
    headers = {
      "DD-API-KEY" = "<Datadog API Key>"
    }
    body = <<EOF
[
  {
    "hostname": "paragon",
    "service": "[Paragon] {{$.event.project.name}}",
    "ddsource": "paragon",
    "message": {{$.event}}
  }
]
EOF
  }
}
```

### New Relic
``` terraform
resource "paragon_events_destination" "webhook_example" {
  project_id = "a7321f97-9c6a-437d-b51e-bd4ce549635f"
  events     = ["workflow_failure"]

  webhook = {
    url = "https://log-api.newrelic.com/log/v1"
    headers = {
      "Api-Key" = "<New Relic API Key>"
    }
    body = <<EOF
[{
  "common": {
      "attributes": {
        "logtype": "{{$.event.type}}",
        "service": "[Paragon] {{$.event.project.name}}",
        "hostname": "paragon"
      }
    },
  "logs": [{
    "timestamp": {{$.event.timestamp}},
    "message": "{{$.event.message}}",
    "attributes": {{$.event}}
  }]
}]
EOF
  }
}
```

### Sentry
``` terraform
resource "paragon_events_destination" "webhook_example" {
  project_id = "a7321f97-9c6a-437d-b51e-bd4ce549635f"
  events     = ["workflow_failure"]

  webhook = {
    url = "<sentry URL>"
    headers = {
      "X-Sentry-Auth" = "Sentry sentry_version=7, sentry_key=<Sentry Key>"
    }
    body = <<EOF
{
  "event_id": "{{$.event.data.workflowExecution.id}}",
  "timestamp": "{{$.event.timestampISO}}",
  "exception": {
    "values": [{
      "type": "Paragon {{$.event.type}}",
      "value": "{{$.event.message}}",
      "mechanism": {
        "type": "paragon",
        "description": "{{$.event.data.error}}",
        "data": {{$.event}}
      }
    }]
  }
}
EOF
  }
}
```


### Slack
``` terraform
resource "paragon_events_destination" "webhook_example" {
  project_id = "a7321f97-9c6a-437d-b51e-bd4ce549635f"
  events     = ["workflow_failure"]

  webhook = {
    url = "<Slack webhook url>"
    body = <<EOF
{
      "text": "*{{$.event.message}}* (<https://dashboard.useparagon.com/connect/projects/{{$.event.project.id}}/history/workflows/{{$.event.workflow.id}}/executions/{{$.event.data.workflowExecution.id}}|Link to execution>)\n\n```{{$.event.data.error}}```"
}
EOF
  }
}
```


## Schema
### Argument Reference
* `project_id` (String, Required) Identifier of the project.
* `events` (List of String, Required) List of events to subscribe to, Currently only `workflow_failure` is supported.
* `webhook` (Block, Optional) Webhook destination configuration. Cannot be used with `email`.
  * `url` (String, Required) URL to send webhook notifications to.
  * `headers` (Map of String, Sensitive, Optional) Headers to include in the webhook request.
  * `body` (String, Required) Body to send with the webhook, Supports variable substitution from the event.
* `email` (Block, Optional) Email destination configuration. Cannot be used with `webhook`.
  * `address` (String, Required) Email address to send notifications to.

### Attributes Reference
- `id` (String) Identifier of the event destination.

## JSON State Structure Example

Here's a state sample

```json
{
  "email": null,
  "events": [
    "workflow_failure"
  ],
  "id": "ab86fd8f-4e52-433c-82bd-1dd968103256",
  "project_id": "a7321f97-9c6a-437d-b51e-bd4ce549635f",
  "webhook": {
    "body": "[\n  {\n    \"hostname\": \"paragon\",\n    \"service\": \"[Paragon] {{$.event.project.name}}\",\n    \"ddsource\": \"paragon\",\n    \"message\": \"{{$.event}}\",\n    \"some more\": \"{{$.event.timestamp}}\"\n  }\n]\n",
    "headers": {
      "key": "value",
      "key2": "value"
    },
    "url": "https://example.com/webhook"
  }
}
```