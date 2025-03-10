---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sentry_issue_alert Resource - terraform-provider-sentry"
subcategory: ""
description: |-
  Create an Issue Alert Rule for a Project. See the Sentry Documentation https://docs.sentry.io/api/alerts/create-an-issue-alert-rule-for-a-project/ for more information.
  Please note the following changes since v0.12.0:
  - The attributes conditions, filters, and actions are in JSON string format. The types must match the Sentry API, otherwise Terraform will incorrectly detect a drift. Use parseint("string", 10) to convert a string to an integer. Avoid using jsonencode() as it is unable to distinguish between an integer and a float.
  - The attribute internal_id has been removed. Use id instead.
  - The attribute id is now the ID of the issue alert. Previously, it was a combination of the organization, project, and issue alert ID.
---

# sentry_issue_alert (Resource)

Create an Issue Alert Rule for a Project. See the [Sentry Documentation](https://docs.sentry.io/api/alerts/create-an-issue-alert-rule-for-a-project/) for more information.

Please note the following changes since v0.12.0:
- The attributes `conditions`, `filters`, and `actions` are in JSON string format. The types must match the Sentry API, otherwise Terraform will incorrectly detect a drift. Use `parseint("string", 10)` to convert a string to an integer. Avoid using `jsonencode()` as it is unable to distinguish between an integer and a float.
- The attribute `internal_id` has been removed. Use `id` instead.
- The attribute `id` is now the ID of the issue alert. Previously, it was a combination of the organization, project, and issue alert ID.

## Example Usage

```terraform
resource "sentry_issue_alert" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id
  name         = "My issue alert"

  action_match = "any"
  filter_match = "any"
  frequency    = 30

  conditions = <<EOT
[
  {
    "id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
  },
  {
    "id": "sentry.rules.conditions.regression_event.RegressionEventCondition"
  },
  {
    "id": "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
    "value": 500,
    "interval": "1h"
  },
  {
    "id": "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition",
    "value": 1000,
    "interval": "15m"
  },
  {
    "id": "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition",
    "value": 50,
    "interval": "10m"
  }
]
EOT

  actions = "[]" # Please see below for examples

  filters = <<EOT
[
  {
    "id": "sentry.rules.filters.age_comparison.AgeComparisonFilter",
    "comparison_type": "older",
    "value": 3,
    "time": "week"
  },
  {
    "id": "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
    "value": 120
  },
  {
    "id": "sentry.rules.filters.assigned_to.AssignedToFilter",
    "targetType": "Unassigned"
  },
  {
    "id": "sentry.rules.filters.assigned_to.AssignedToFilter",
    "targetType": "Member",
    "targetIdentifier": 895329789
  },
  {
    "id": "sentry.rules.filters.latest_release.LatestReleaseFilter"
  },
  {
    "id": "sentry.rules.filters.issue_category.IssueCategoryFilter",
    "value": 2
  },
  {
    "id": "sentry.rules.conditions.event_attribute.EventAttributeCondition",
    "attribute": "http.url",
    "match": "nc",
    "value": "localhost"
  },
  {
    "id": "sentry.rules.filters.tagged_event.TaggedEventFilter",
    "key": "level",
    "match": "eq",
    "value": "error"
  },
  {
    "id": "sentry.rules.filters.level.LevelFilter",
    "match": "gte",
    "level": "50"
  }
]
EOT
}

#
# Send a notification to Suggested Assignees
#

resource "sentry_issue_alert" "member_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.mail.actions.NotifyEmailAction",
    "targetType": "IssueOwners",
    "fallthroughType": "ActiveMembers"
  }
]
EOT
  // ...
}

#
# Send a notification to a Member
#

data "sentry_organization_member" "member" {
  organization = data.sentry_organization.test.id
  email        = "test@example.com"
}

resource "sentry_issue_alert" "member_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.mail.actions.NotifyEmailAction",
    "targetType": "Member",
    "fallthroughType": "AllMembers",
    "targetIdentifier": ${parseint(data.sentry_organization_member.member.id, 10)}
  }
]
EOT
  // ...
}

#
# Send a notification to a Team
#

data "sentry_team" "team" {
  organization = sentry_project.test.organization
  slug         = "my-team"
}

resource "sentry_issue_alert" "team_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.mail.actions.NotifyEmailAction",
    "targetType": "Team",
    "fallthroughType": "AllMembers",
    "targetIdentifier": ${parseint(data.sentry_team.team.internal_id, 10)}
  }
]
EOT
  // ...
}

#
# Send a Slack notification
#

# Retrieve a Slack integration
data "sentry_organization_integration" "slack" {
  organization = sentry_project.test.organization

  provider_key = "slack"
  name         = "Slack Workspace" # Name of your Slack workspace
}

resource "sentry_issue_alert" "slack_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
    "workspace": ${parseint(data.sentry_organization_integration.slack.id, 10)},
    "channel": "#warning",
    "tags": "environment,level"
  }
]
EOT
  // ...
}

#
# Send a Microsoft Teams notification
#

data "sentry_organization_integration" "msteams" {
  organization = sentry_project.test.organization

  provider_key = "msteams"
  name         = "My Team" # Name of your Microsoft Teams team
}

resource "sentry_issue_alert" "msteams_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.msteams.notify_action.MsTeamsNotifyServiceAction",
    "team": ${parseint(data.sentry_organization_integration.msteams.id, 10)},
    "channel": "General"
  }
]
EOT
  // ...
}

#
# Send a Discord notification
#

data "sentry_organization_integration" "discord" {
  organization = sentry_project.test.organization

  provider_key = "discord"
  name         = "Discord Server" # Name of your Discord server
}

resource "sentry_issue_alert" "discord_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.discord.notify_action.DiscordNotifyServiceAction",
    "server": ${parseint(data.sentry_organization_integration.discord.id, 10)},
    "channel_id": 94732897,
    "tags": "browser,user"
  }
]
EOT
  // ...
}

#
# Create a Jira Ticket
#

data "sentry_organization_integration" "jira" {
  organization = sentry_project.test.organization

  provider_key = "jira"
  name         = "JIRA" # Name of your Jira server
}

resource "sentry_issue_alert" "jira_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.jira.notify_action.JiraCreateTicketAction",
    "integration": ${parseint(data.sentry_organization_integration.jira.id, 10)},
    "project": "349719",
    "issueType": "1"
  }
]
EOT
  // ...
}

#
# Create a Jira Server Ticket
#

data "sentry_organization_integration" "jira_server" {
  organization = sentry_project.test.organization

  provider_key = "jira_server"
  name         = "JIRA" # Name of your Jira server
}

resource "sentry_issue_alert" "jira_server_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.jira_server.notify_action.JiraServerCreateTicketAction",
    "integration": ${parseint(data.sentry_organization_integration.jira_server.id, 10)},
    "project": "349719",
    "issueType": "1"
  }
]
EOT
  // ...
}

#
# Create a GitHub Issue
#

data "sentry_organization_integration" "github" {
  organization = sentry_project.test.organization

  provider_key = "github"
  name         = "GitHub"
}

resource "sentry_issue_alert" "github_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.github.notify_action.GitHubCreateTicketAction",
    "integration": ${parseint(data.sentry_organization_integration.github.id, 10)},
    "repo": default,
    "title": "My Test Issue",
    "assignee": "Baxter the Hacker",
    "labels": ["bug", "p1"]
  }
]
EOT
  // ...
}

#
# Create an Azure DevOps work item
#

data "sentry_organization_integration" "vsts" {
  organization = sentry_project.test.organization

  provider_key = "vsts"
  name         = "Azure DevOps"
}

resource "sentry_issue_alert" "vsts_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.vsts.notify_action.AzureDevopsCreateTicketAction",
    "integration": ${parseint(data.sentry_organization_integration.vsts.id, 10)},
    "project": "0389485",
    "work_item_type": "Microsoft.VSTS.WorkItemTypes.Task"
  }
]
EOT
  // ...
}

#
# Send a PagerDuty notification
#

data "sentry_organization_integration" "pagerduty" {
  organization = sentry_project.test.organization

  provider_key = "pagerduty"
  name         = "PagerDuty"
}

resource "sentry_issue_alert" "pagerduty_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.pagerduty.notify_action.PagerDutyNotifyServiceAction",
    "account": ${parseint(data.sentry_organization_integration.pagerduty.id, 10)},
    "service": 9823924
  }
]
EOT
  // ...
}

#
# Send an Opsgenie notification
#

data "sentry_organization_integration" "opsgenie" {
  organization = sentry_project.test.organization

  provider_key = "opsgenie"
  name         = "Opsgenie"
}

resource "sentry_issue_alert" "opsgenie_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.opsgenie.notify_action.OpsgenieNotifyTeamAction",
    "account": ${parseint(data.sentry_organization_integration.opsgenie.id, 10)},
    "team": "9438930258-fairy"
  }
]
EOT
  // ...
}

#
# Send a notification to a service
#

resource "sentry_issue_alert" "notification_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.rules.actions.notify_event_service.NotifyEventServiceAction",
    "service": "mail"
  }
]
EOT
  // ...
}

#
# Send a notification to a Sentry app with a custom webhook payload
#

resource "sentry_issue_alert" "notification_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.rules.actions.notify_event_sentry_app.NotifyEventSentryAppAction",
    "settings": [
        {"name": "title", "value": "Team Rocket"},
        {"name": "summary", "value": "We're blasting off again."}
    ],
    "sentryAppInstallationUuid": 643522,
    "hasSchemaFormConfig": true
  }
]
EOT
  // ...
}

#
# Send a notification (for all legacy integrations)
#

resource "sentry_issue_alert" "notification_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.rules.actions.notify_event.NotifyEventAction"
  }
]
EOT
  // ...
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action_match` (String) Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen.
- `actions` (String) List of actions. In JSON string format.
- `conditions` (String) List of conditions. In JSON string format.
- `frequency` (Number) Perform actions at most once every `X` minutes for this issue.
- `name` (String) The issue alert name.
- `organization` (String) The slug of the organization the resource belongs to.
- `project` (String) The slug of the project the resource belongs to.

### Optional

- `environment` (String) Perform issue alert in a specific environment.
- `filter_match` (String) A string determining which filters need to be true before any actions take place. Required when a value is provided for `filters`.
- `filters` (String) A list of filters that determine if a rule fires after the necessary conditions have been met. In JSON string format.
- `owner` (String) The ID of the team or user that owns the rule.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# import using the organization, project slugs and rule id from the URL:
# https://sentry.io/organizations/[org-slug]/alerts/rules/[project-slug]/[rule-id]/details/
terraform import sentry_issue_alert.default org-slug/project-slug/rule-id
```
