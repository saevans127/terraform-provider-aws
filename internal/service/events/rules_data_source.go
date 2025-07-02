// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package events

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

func (d *dataSourceRules) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Meta().EventsClient(ctx)

	var data dataSourceRulesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	resp.Diagnostics.Append(flex.Flatten(ctx, out, &data.Rules)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

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
