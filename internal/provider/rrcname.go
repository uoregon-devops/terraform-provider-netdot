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

type rrCnameModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Cname types.String `tfsdk:"cname"`
	RR    types.String `tfsdk:"rr"`
	RRID  types.Int64  `tfsdk:"rr_id"`
	TTL   types.Int64  `tfsdk:"ttl"`
}

func RRCnameToRRCnameModel(rr models.RRCname) rrCnameModel {
	var finalModel rrCnameModel

	finalModel.ID = types.Int64Value(rr.ID)
	finalModel.Cname = types.StringValue(rr.Cname)
	finalModel.RR = types.StringValue(rr.RR)
	finalModel.RRID = types.Int64Value(rr.RRXlink.ID)
	finalModel.TTL = autoNullInt64(rr.TTL)

	return finalModel
}

func RRCnameModelToRRCnameQuery(model rrCnameModel) netdot.RRCnameQuery {
	rrCnameQuery := netdot.NewRRCnameQueryBuilder()

	if isPopulated(model.TTL) {
		rrCnameQuery.TTL(model.TTL.ValueInt64())
	}

	if isPopulated(model.Cname) {
		rrCnameQuery.Cname(model.Cname.ValueString())
	}

	if isPopulated(model.RRID) {
		rrCnameQuery.RR(model.RRID.ValueInt64())
	}

	return rrCnameQuery.Build()
}

var rrCnameDataSourceSchema = datasourceSchema.Schema{
	Description: "A CNAME points a resource record to another domain name.",
	Attributes: map[string]datasourceSchema.Attribute{
		"id": datasourceSchema.Int64Attribute{
			Description: "ID of the CNAME record.",
			Optional:    true,
		},
		"cname": datasourceSchema.StringAttribute{
			Description: "CNAME target domain.",
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

var rrCnameResourceSchema = resourceSchema.Schema{
	Description: "A CNAME points a resource record to another domain name.",
	Attributes: map[string]resourceSchema.Attribute{
		"id": resourceSchema.Int64Attribute{
			Description: "ID of the CNAME record.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"cname": resourceSchema.StringAttribute{
			Description: "CNAME target domain.",
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
