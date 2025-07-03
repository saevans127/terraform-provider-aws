---
subcategory: "EventBridge"
layout: "aws"
page_title: "AWS: aws_events_rules"
description: |-
  Provides details about an AWS EventBridge Rules.
---
<!---
Documentation guidelines:
- Begin data source descriptions with "Provides details about..."
- Use simple language and avoid jargon
- Focus on brevity and clarity
- Use present tense and active voice
- Don't begin argument/attribute descriptions with "An", "The", "Defines", "Indicates", or "Specifies"
- Boolean arguments should begin with "Whether to"
- Use "example" instead of "test" in examples
--->

# Data Source: aws_events_rules

Provides details about an AWS EventBridge Rules.

## Example Usage

### Basic Usage

```terraform
data "aws_events_rules" "example" {
name_prefix = "test"
}
```

## Argument Reference



The following arguments are optional:

* `name_prefix` - (Optional) The prefix matching the rule name.
* `event_bus_name` - (Optional) The name or ARN of the event bus to list the rules for. If you omit this, the default event bus is used.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `rules` - This list of event bus rules.

### `rules` Attribute Reference

* `arn` - The ARN of the rule.
* `description` - The description of the rule.
* `event_bus_name` - The name or ARN of the event bus associated with the rule. If you omit this, the default event bus is used.
* `event_pattern` - The event pattern of the rule. For more information, see [Events and Event Patterns](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-events.html) in the * Amazon EventBridge User Guide *
* `managed_by` - If the rule was created on behalf of your account by an Amazon Web Services service, this field displays the principal name of the service that created the rule.
* `name` - The name of the rule.
* `role_arn` - The Amazon Resource Name (ARN) of the role that is used for target invocation.
* `schedule_expression` - The scheduling expression. For example, “cron(0 20 * * ? *)”, “rate(5 minutes)”. For more information, see [Creating an Amazon EventBridge rule that runs on a schedule](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-create-rule-schedule.html) .
* `state` - The state of the rule. Valid values include DISABLE, ENABLED, ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS.


