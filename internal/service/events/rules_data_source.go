// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package events

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
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	awstypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwflex "github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// Function annotations are used for datasource registration to the Provider. DO NOT EDIT.
// @FrameworkDataSource("aws_cloudwatch_event_rules", name="Rules")
func newDataSourceRules(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceRules{}, nil
}

const (
	DSNameRules = "Rules Data Source"
)

type dataSourceRules struct {
	framework.DataSourceWithModel[dataSourceRulesModel]
}

func (d *dataSourceRules) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rules": framework.DataSourceComputedListOfObjectAttribute[dataSourceRuleModel](ctx),
			names.AttrNamePrefix: schema.StringAttribute{
				Optional: true,
			},
			"event_bus_name": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// TIP: ==== ASSIGN CRUD METHODS ====
// Data sources only have a read method.
func (d *dataSourceRules) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// TIP: ==== DATA SOURCE READ ====
	// Generally, the Read function should do the following things. Make
	// sure there is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Fetch the config
	// 3. Get information about a resource from AWS
	// 4. Set the ID, arguments, and attributes
	// 5. Set the tags
	// 6. Set the state
	// TIP: -- 1. Get a client connection to the relevant service
	conn := d.Meta().EventsClient(ctx)

	// TIP: -- 2. Fetch the config
	var data dataSourceRulesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 3. Get information about a resource from AWS
	input := eventbridge.ListRulesInput{
		NamePrefix:   fwflex.StringFromFramework(ctx, data.NamePrefix),
		EventBusName: fwflex.StringFromFramework(ctx, data.EventBusName),
	}
	out, err := findRules(ctx, conn, &input)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Events, create.ErrActionReading, DSNameRules, data.NamePrefix.String(), err),
			err.Error(),
		)
		return
	}

	// TIP: -- 4. Set the ID, arguments, and attributes
	// Using a field name prefix allows mapping fields such as `RulesId` to `ID`
	resp.Diagnostics.Append(flex.Flatten(ctx, out, &data.Rules)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 5. Set the tags

	// TIP: -- 6. Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// TIP: ==== DATA STRUCTURES ====
// With Terraform Plugin-Framework configurations are deserialized into
// Go types, providing type safety without the need for type assertions.
// These structs should match the schema definition exactly, and the `tfsdk`
// tag value should match the attribute name.
//
// Nested objects are represented in their own data struct. These will
// also have a corresponding attribute type mapping for use inside flex
// functions.
//
// See more:
// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
type dataSourceRulesModel struct {
	framework.WithRegionModel
	Rules        fwtypes.ListNestedObjectValueOf[dataSourceRuleModel] `tfsdk:"rules"`
	NamePrefix   types.String                                         `tfsdk:"name_prefix"`
	EventBusName types.String                                         `tfsdk:"event_bus_name"`
}
type dataSourceRuleModel struct {
	ARN                types.String `tfsdk:"arn"`
	Description        types.String `tfsdk:"description"`
	EventBusName       types.String `tfsdk:"event_bus_name"`
	EventPattern       types.String `tfsdk:"event_pattern"`
	ManagedBy          types.String `tfsdk:"managed_by"`
	Name               types.String `tfsdk:"name"`
	RoleArn            types.String `tfsdk:"role_arn"`
	ScheduleExpression types.String `tfsdk:"schedule_expression"`
	State              types.String `tfsdk:"state"`
}

func findRules(ctx context.Context, conn *eventbridge.Client, input *eventbridge.ListRulesInput) ([]awstypes.Rule, error) {
	var output []awstypes.Rule

	err := listRulesPages(ctx, conn, input, func(page *eventbridge.ListRulesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		output = append(output, page.Rules...)

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
