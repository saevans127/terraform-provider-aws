// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package events_test

// **PLEASE DELETE THIS AND ALL TIP COMMENTS BEFORE SUBMITTING A PR FOR REVIEW!**
//
// TIP: ==== INTRODUCTION ====
// Thank you for trying the skaff tool!
//
// You have opted to include these helpful comments. They all include "TIP:"
// to help you find and remove them when you're done with them.
//
// While some aspects of this file are customized to your input, the
// scaffold tool does *not* look at the AWS API and ensure it has correct
// function, structure, and variable names. It makes guesses based on
// commonalities. You will need to make significant adjustments.
//
// In other words, as generated, this is a rough outline of the work you will
// need to do. If something doesn't make sense for your situation, get rid of
// it.

import (
	// TIP: ==== IMPORTS ====
	// This is a common set of imports but not customized to your code since
	// your code hasn't been written yet. Make sure you, your IDE, or
	// goimports -w <file> fixes these imports.
	//
	// The provider linter wants your imports to be in two groups: first,
	// standard library (i.e., "fmt" or "strings"), second, everything else.
	//
	// Also, AWS Go SDK v2 may handle nested structures differently than v1,
	// using the services/eventbridge/types package. If so, you'll
	// need to import types and reference the nested types, e.g., as
	// types.<Type Name>.
	"fmt"
	//"strings"
	"testing"

	//"github.com/YakDriver/regexache"
	//"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	//"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	//"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	//"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	//"github.com/hashicorp/terraform-provider-aws/internal/conns"
	//"github.com/hashicorp/terraform-provider-aws/internal/create"

	// TIP: You will often need to import the package that this test file lives
	// in. Since it is in the "test" context, it must import the package to use
	// any normal context constants, variables, or functions.
	//tfevents "github.com/hashicorp/terraform-provider-aws/internal/service/events"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// TIP: File Structure. The basic outline for all test files should be as
// follows. Improve this data source's maintainability by following this
// outline.
//
// 1. Package declaration (add "_test" since this is a test file)
// 2. Imports
// 3. Unit tests
// 4. Basic test
// 5. Disappears test
// 6. All the other tests
// 7. Helper functions (exists, destroy, check, etc.)
// 8. Functions that return Terraform configurations

// TIP: ==== ACCEPTANCE TESTS ====
// This is an example of a basic acceptance test. This should test as much of
// standard functionality of the data source as possible, and test importing, if
// applicable. We prefix its name with "TestAcc", the service, and the
// data source name.
//
// Acceptance test access AWS and cost money to run.
func TestAccEventsRulesDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	//busName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
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
