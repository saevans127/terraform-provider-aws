// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package events_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"

	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccEventsRulesDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	ruleName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudwatch_event_rule.test"
	dataSource1Name := "data.aws_cloudwatch_event_rules.by_name_prefix"
	dataSource2Name := "data.aws_cloudwatch_event_rules.all"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.EventsEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EventsServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRulesDataSourceConfig_basic(ruleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSource2Name, "rules.0.name", resourceName, names.AttrName),
					acctest.CheckResourceAttrGreaterThanOrEqualValue(dataSource2Name, "rules.#", 1),
					resource.TestCheckResourceAttr(dataSource1Name, "rules.#", "1"),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.arn", resourceName, names.AttrARN),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.description", resourceName, "description"),
					resource.TestCheckResourceAttr(dataSource1Name, "rules.0.event_bus_name", "default"),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.event_pattern", resourceName, "event_pattern"),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.managed_by", resourceName, "managed_by"),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.name", resourceName, names.AttrName),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.state", resourceName, "state"),
				),
			},
		},
	})
}

func TestAccEventsRulesDataSource_custom(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	busName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	ruleName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudwatch_event_rule.test"
	dataSource1Name := "data.aws_cloudwatch_event_rules.by_name_prefix"
	dataSource2Name := "data.aws_cloudwatch_event_rules.all"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.EventsEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EventsServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRulesDataSourceConfig_custom(busName, ruleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSource2Name, "rules.0.name", resourceName, names.AttrName),
					acctest.CheckResourceAttrGreaterThanOrEqualValue(dataSource2Name, "rules.#", 1),
					resource.TestCheckResourceAttr(dataSource1Name, "rules.#", "1"),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.arn", resourceName, names.AttrARN),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.description", resourceName, "description"),
					resource.TestCheckResourceAttr(dataSource1Name, "rules.0.event_bus_name", busName),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.event_pattern", resourceName, "event_pattern"),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.managed_by", resourceName, "managed_by"),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.name", resourceName, names.AttrName),
					resource.TestCheckResourceAttrPair(dataSource1Name, "rules.0.state", resourceName, "state"),
				),
			},
		},
	})
}

func testAccRulesDataSourceConfig_basic(ruleName string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_event_rule" "test" {
  name = %[1]q
  description = "a test rule"
  event_pattern  = <<PATTERN
{
	"source": [
		"aws.ec2"
	]
}
PATTERN
}

data "aws_cloudwatch_event_rules" "by_name_prefix" {
  name_prefix = aws_cloudwatch_event_rule.test.name
}

data "aws_cloudwatch_event_rules" "all" {
  depends_on = [aws_cloudwatch_event_rule.test]
}
`, ruleName)
}

func testAccRulesDataSourceConfig_custom(busName string, ruleName string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_event_bus" "test" {
  name = %[1]q
}
resource "aws_cloudwatch_event_rule" "test" {
  name = %[2]q
  description = "a test rule"
  event_bus_name = aws_cloudwatch_event_bus.test.name
  event_pattern  = <<PATTERN
{
	"source": [
		"aws.ec2"
	]
}
PATTERN
}

data "aws_cloudwatch_event_rules" "by_name_prefix" {
  name_prefix = aws_cloudwatch_event_rule.test.name
  event_bus_name = aws_cloudwatch_event_bus.test.name
}

data "aws_cloudwatch_event_rules" "all" {
  depends_on = [aws_cloudwatch_event_rule.test]
  event_bus_name = aws_cloudwatch_event_bus.test.name
}
`, busName, ruleName)
}
