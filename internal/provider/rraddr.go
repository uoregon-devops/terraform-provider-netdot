package provider

import (
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type rrAddrModel struct {
	ID        types.Int64  `tfsdk:"id"`
	IpBlock   types.String `tfsdk:"ipblock"`
	IpBlockID types.Int64  `tfsdk:"ipblock_id"`
	RR        types.String `tfsdk:"rr"`
	RRID      types.Int64  `tfsdk:"rr_id"`
	TTL       types.Int64  `tfsdk:"ttl"`
}

func RRAddrToRRAddrModel(rr models.RRAddr) rrAddrModel {
	var finalModel rrAddrModel

	finalModel.ID = types.Int64Value(rr.ID)
	finalModel.IpBlock = types.StringValue(rr.IpBlock)
	finalModel.IpBlockID = types.Int64Value(rr.IpBlockXlink.ID)
	finalModel.RR = types.StringValue(rr.RR)
	finalModel.RRID = types.Int64Value(rr.RRXlink.ID)
	finalModel.TTL = autoNullInt64(rr.TTL)

	return finalModel
}

func RRAddrModelToRRAddrQuery(model rrAddrModel) netdot.RRAddrQuery {
	rrAddrQuery := netdot.NewRRAddrQueryBuilder()

	if isPopulated(model.TTL) {
		rrAddrQuery.TTL(model.TTL.ValueInt64())
	}

	if isPopulated(model.IpBlockID) {
		rrAddrQuery.IpBlock(model.IpBlockID.ValueInt64())
	}

	if isPopulated(model.RRID) {
		rrAddrQuery.RR(model.RRID.ValueInt64())
	}

	return rrAddrQuery.Build()
}

var rrAddrDataSourceSchema = datasourceSchema.Schema{
	Attributes: map[string]datasourceSchema.Attribute{
		"id": datasourceSchema.Int64Attribute{
			Description: "The A record's ID.",
			Optional:    true,
		},
		"ipblock": datasourceSchema.StringAttribute{
			Description: "The CIDR of the A record's associated ipblock.",
			Optional:    true,
		},
		"ipblock_id": datasourceSchema.Int64Attribute{
			Description: "The ID of the A record's associated ipblock.",
			Computed:    true,
		},
		"rr": datasourceSchema.StringAttribute{
			Description: "The associated DNS record.",
			Computed:    true,
		},
		"rr_id": datasourceSchema.Int64Attribute{
			Description: "The ID of the associated DNS record.",
			Computed:    true,
		},
		"ttl": datasourceSchema.Int64Attribute{
			Description: "Time to live for the A record in seconds.",
			Computed:    true,
		},
	},
}

var rrAddrResourceSchema = resourceSchema.Schema{
	Attributes: map[string]resourceSchema.Attribute{
		"id": resourceSchema.Int64Attribute{
			Description: "The A record's ID.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"ipblock": resourceSchema.StringAttribute{
			Description: "The CIDR of the A record's associated ipblock.",
			Computed:    true,
		},
		"ipblock_id": resourceSchema.Int64Attribute{
			Description: "The ID of the A record's associated ipblock.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"rr": resourceSchema.StringAttribute{
			Description: "The associated DNS record.",
			Computed:    true,
		},
		"rr_id": resourceSchema.Int64Attribute{
			Description: "The ID of the associated DNS record.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"ttl": resourceSchema.Int64Attribute{
			Description: "Time to live for the A record in seconds.",
			Optional:    true,
			Computed:    true,
			Default:     int64default.StaticInt64(600),
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
	},
}
