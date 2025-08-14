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

type rrNsModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Ns   types.String `tfsdk:"name_server"`
	RR   types.String `tfsdk:"rr"`
	RRID types.Int64  `tfsdk:"rr_id"`
	TTL  types.Int64  `tfsdk:"ttl"`
}

func RRNsToRRNsModel(rr models.RRNs) rrNsModel {
	var finalModel rrNsModel

	finalModel.ID = types.Int64Value(rr.ID)
	finalModel.Ns = types.StringValue(rr.NsDName)
	finalModel.RR = types.StringValue(rr.RR)
	finalModel.RRID = types.Int64Value(rr.RRXlink.ID)
	finalModel.TTL = autoNullInt64(rr.TTL)

	return finalModel
}

func RRNsModelToRRNsQuery(model rrNsModel) netdot.RRNsQuery {
	rrNsQuery := netdot.NewRRNsQueryBuilder()

	if isPopulated(model.TTL) {
		rrNsQuery.TTL(model.TTL.ValueInt64())
	}

	if isPopulated(model.Ns) {
		rrNsQuery.Ns(model.Ns.ValueString())
	}

	if isPopulated(model.RRID) {
		rrNsQuery.RR(model.RRID.ValueInt64())
	}

	return rrNsQuery.Build()
}

var rrNsDataSourceSchema = datasourceSchema.Schema{
	Description: "An NS record designates an authoritative name server for a domain.",
	Attributes: map[string]datasourceSchema.Attribute{
		"id": datasourceSchema.Int64Attribute{
			Description: "ID of the NS record.",
			Optional:    true,
		},
		"name_server": datasourceSchema.StringAttribute{
			Description: "Target name server.",
			Optional:    true,
		},
		"rr": datasourceSchema.StringAttribute{
			Description: "Associated resource record name.",
			Computed:    true,
		},
		"rr_id": datasourceSchema.Int64Attribute{
			Description: "ID of the associated resource record.",
			Computed:    true,
		},
		"ttl": datasourceSchema.Int64Attribute{
			Description: "Time to live for the CNAME record in seconds.",
			Computed:    true,
		},
	},
}

var rrNsResourceSchema = resourceSchema.Schema{
	Description: "An NS record designates an authoritative name server for a domain.",
	Attributes: map[string]resourceSchema.Attribute{
		"id": resourceSchema.Int64Attribute{
			Description: "ID of the NS record.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"name_server": resourceSchema.StringAttribute{
			Description: "Target name server.",
			Required:    true,
		},
		"rr": resourceSchema.StringAttribute{
			Description: "Associated resource record name.",
			Computed:    true,
		},
		"rr_id": resourceSchema.Int64Attribute{
			Description: "ID of the associated resource record.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"ttl": resourceSchema.Int64Attribute{
			Description: "Time to live for the CNAME record in seconds.",
			Optional:    true,
			Computed:    true,
			Default:     int64default.StaticInt64(600),
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
	},
}
