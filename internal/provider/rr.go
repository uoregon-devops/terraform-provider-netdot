package provider

import (
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type rrModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Active     types.Bool   `tfsdk:"active"`
	AutoUpdate types.Bool   `tfsdk:"auto_update"`
	Expiration types.String `tfsdk:"expiration"`
	Info       types.String `tfsdk:"info"`
	Name       types.String `tfsdk:"name"`
	Zone       types.String `tfsdk:"zone"`
	ZoneID     types.Int64  `tfsdk:"zone_id"`
	FQDN       types.String `tfsdk:"fqdn"`
}

func RRToRRModel(rr models.RR) rrModel {
	var finalModel rrModel

	finalModel.ID = types.Int64Value(rr.ID)
	finalModel.Active = types.BoolValue(rr.Active)
	finalModel.AutoUpdate = types.BoolValue(rr.AutoUpdate)
	finalModel.Expiration = autoNullString(rr.Expiration)
	finalModel.Info = autoNullString(rr.Info)
	finalModel.Name = types.StringValue(rr.Name)
	finalModel.Zone = types.StringValue(rr.Zone)
	finalModel.ZoneID = types.Int64Value(rr.ZoneXlink.ID)
	finalModel.FQDN = types.StringValue(rr.Name + "." + rr.Zone)

	return finalModel
}

func RRModelToRRQuery(model rrModel) netdot.RRQuery {
	rrQuery := netdot.NewRRQueryBuilder()

	if isPopulated(model.Active) {
		rrQuery.Active(model.Active.ValueBool())
	}

	if isPopulated(model.AutoUpdate) {
		rrQuery.AutoUpdate(model.AutoUpdate.ValueBool())
	}

	if isPopulated(model.Expiration) {
		rrQuery.Expiration(model.Expiration.ValueString())
	}

	if isPopulated(model.Info) {
		rrQuery.Info(model.Info.ValueString())
	}

	if isPopulated(model.Name) {
		rrQuery.Name(model.Name.ValueString())
	}

	if isPopulated(model.ZoneID) {
		rrQuery.ZoneID(model.ZoneID.ValueInt64())
	}

	return rrQuery.Build()
}

var rrDataSourceSchema = datasourceSchema.Schema{
	Description: "An RR (resource record) represent a domain name, and is associated with a zone.",
	Attributes: map[string]datasourceSchema.Attribute{
		"id": datasourceSchema.Int64Attribute{
			Description: "ID of the resource record.",
			Optional:    true,
		},
		"name": datasourceSchema.StringAttribute{
			Description: "Name of the resource record.",
			Optional:    true,
		},
		"active": datasourceSchema.BoolAttribute{
			Description: "Active status of the resource record.",
			Computed:    true,
		},
		"auto_update": datasourceSchema.BoolAttribute{
			Description: "Determines whether this record can be updated by an automated process. For example, if a Device Interface changes its name and the change is picked up by an SNMP update.",
			Computed:    true,
		},
		"expiration": datasourceSchema.StringAttribute{
			Description: "Expiration date for this Record (YYYY-MM-DD).",
			Computed:    true,
		},
		"info": datasourceSchema.StringAttribute{
			Description: "Additional information about the resource record.",
			Computed:    true,
		},
		"zone": datasourceSchema.StringAttribute{
			Description: "Zone name for this Record.",
			Computed:    true,
		},
		"zone_id": datasourceSchema.Int64Attribute{
			Description: "Zone ID for this Record.",
			Computed:    true,
		},
		"fqdn": datasourceSchema.StringAttribute{
			Description: "Fully qualified domain name for this Record (name + zone).",
			Computed:    true,
		},
	},
}

var rrResourceSchema = resourceSchema.Schema{
	Description: "An RR (resource record) represent a domain name, and is associated with a zone.",
	Attributes: map[string]resourceSchema.Attribute{
		"id": resourceSchema.Int64Attribute{
			Description: "ID of the resource record.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"name": resourceSchema.StringAttribute{
			Description: "Name of the resource record.",
			Required:    true,
		},
		"active": resourceSchema.BoolAttribute{
			Description: "Active status of the resource record.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(true),
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"auto_update": resourceSchema.BoolAttribute{
			Description: "Determines whether this record can be updated by an automated process. For example, if a Device Interface changes its name and the change is picked up by an SNMP update.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"expiration": resourceSchema.StringAttribute{
			Description: "Expiration date for this Record (YYYY-MM-DD).",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"info": resourceSchema.StringAttribute{
			Description: "Additional information about the resource record.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"zone": resourceSchema.StringAttribute{
			Description: "Zone name for this Record.",
			Required:    true,
		},
		"zone_id": resourceSchema.Int64Attribute{
			Description: "Zone ID for this Record.",
			Computed:    true,
		},
		"fqdn": resourceSchema.StringAttribute{
			Description: "Fully qualified domain name for this Record (name + zone).",
			Computed:    true,
		},
	},
}
