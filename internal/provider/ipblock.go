package provider

import (
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ipblockModel struct {
	ID                  types.Int64  `tfsdk:"id"`
	Address             types.String `tfsdk:"address"`
	Prefix              types.Int64  `tfsdk:"prefix"`
	Version             types.Int64  `tfsdk:"version"`
	ASN                 types.Int64  `tfsdk:"asn"`
	ASNID               types.Int64  `tfsdk:"asn_id"`
	Description         types.String `tfsdk:"description"`
	Info                types.String `tfsdk:"info"`
	Interface           types.String `tfsdk:"interface"`
	InterfaceID         types.Int64  `tfsdk:"interface_id"`
	Monitored           types.Bool   `tfsdk:"monitored"`
	RIR                 types.String `tfsdk:"rir"`
	UseNetworkBroadcast types.Bool   `tfsdk:"use_network_broadcast"`
	UsedBy              types.String `tfsdk:"used_by"`
	UsedByID            types.Int64  `tfsdk:"used_by_id"`
	VLAN                types.Int64  `tfsdk:"vlan"`
	VLANID              types.Int64  `tfsdk:"vlan_id"`
	Status              types.String `tfsdk:"status"`
	StatusID            types.Int64  `tfsdk:"status_id"`
	Parent              types.String `tfsdk:"parent"`
	ParentID            types.Int64  `tfsdk:"parent_id"`
	Owner               types.String `tfsdk:"owner"`
	OwnerID             types.Int64  `tfsdk:"owner_id"`
}

func IPBlockToIpblockModel(ipblock models.IpBlock) ipblockModel {
	ipblockModel := ipblockModel{}

	ipblockModel.ID = types.Int64Value(ipblock.ID)
	ipblockModel.Address = types.StringValue(ipblock.Address)
	ipblockModel.Prefix = types.Int64Value(ipblock.Prefix)
	ipblockModel.Version = types.Int64Value(ipblock.Version)
	ipblockModel.Monitored = types.BoolValue(ipblock.Monitored)

	ipblockModel.Description = autoNullString(ipblock.Description)
	ipblockModel.Info = autoNullString(ipblock.Info)
	ipblockModel.RIR = autoNullString(ipblock.RIR)

	ipblockModel.UseNetworkBroadcast = types.BoolValue(ipblock.UseNetworkBroadcast)

	// Xlinks
	ipblockModel.Status = autoNullXlinkString(ipblock.StatusXlink, ipblock.Status)
	ipblockModel.StatusID = autoNullInt64(ipblock.StatusXlink.ID)

	ipblockModel.ASN = autoNullInt64(ipblock.ASN)
	ipblockModel.ASNID = autoNullInt64(ipblock.AsnXlink.ID)

	ipblockModel.Interface = autoNullXlinkString(ipblock.InterfaceXlink, ipblock.Interface)
	ipblockModel.InterfaceID = autoNullInt64(ipblock.InterfaceXlink.ID)

	ipblockModel.Owner = autoNullXlinkString(ipblock.OwnerXlink, ipblock.Owner)
	ipblockModel.OwnerID = autoNullInt64(ipblock.OwnerXlink.ID)

	ipblockModel.Parent = autoNullXlinkString(ipblock.ParentXlink, ipblock.Parent)
	ipblockModel.ParentID = autoNullInt64(ipblock.ParentXlink.ID)

	ipblockModel.UsedBy = autoNullXlinkString(ipblock.UsedByXlink, ipblock.UsedBy)
	ipblockModel.UsedByID = autoNullInt64(ipblock.UsedByXlink.ID)

	ipblockModel.VLAN = autoNullInt64(ipblock.VLAN)
	ipblockModel.VLANID = autoNullInt64(ipblock.VLANXlink.ID)

	return ipblockModel
}

func IPBlockModelToIPBlockQuery(model ipblockModel) netdot.IpBlockQuery {
	queryBuilder := netdot.NewIpBlockQueryBuilder()

	if !model.Address.IsNull() && !model.Address.IsUnknown() {
		queryBuilder.Address(model.Address.ValueString())
	}

	if !model.Prefix.IsNull() && !model.Prefix.IsUnknown() {
		queryBuilder.Prefix(model.Prefix.ValueInt64())
	}

	if !model.Version.IsNull() && !model.Version.IsUnknown() {
		queryBuilder.Version(model.Version.ValueInt64())
	}

	if !model.ASNID.IsNull() && !model.ASNID.IsUnknown() {
		queryBuilder.ASNID(model.ASNID.ValueInt64())
	}

	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		queryBuilder.Description(model.Description.ValueString())
	}

	if !model.Info.IsNull() && !model.Info.IsUnknown() {
		queryBuilder.Info(model.Info.ValueString())
	}

	if !model.InterfaceID.IsNull() && !model.InterfaceID.IsUnknown() {
		queryBuilder.InterfaceID(model.InterfaceID.ValueInt64())
	}

	if !model.Monitored.IsNull() && !model.Monitored.IsUnknown() {
		queryBuilder.Monitored(model.Monitored.ValueBool())
	}

	if !model.RIR.IsNull() && !model.RIR.IsUnknown() {
		queryBuilder.RIR(model.RIR.ValueString())
	}

	if !model.UseNetworkBroadcast.IsNull() && !model.UseNetworkBroadcast.IsUnknown() {
		queryBuilder.UseNetworkBroadcast(model.UseNetworkBroadcast.ValueBool())
	}

	if !model.UsedByID.IsNull() && !model.UsedByID.IsUnknown() {
		queryBuilder.UsedByID(model.UsedByID.ValueInt64())
	}

	if !model.VLANID.IsNull() && !model.VLANID.IsUnknown() {
		queryBuilder.VLANID(model.VLANID.ValueInt64())
	}

	if !model.Status.IsNull() && !model.Status.IsUnknown() {
		queryBuilder.Status(model.Status.ValueString())
	}

	if !model.ParentID.IsNull() && !model.ParentID.IsUnknown() {
		queryBuilder.ParentID(model.ParentID.ValueInt64())
	}

	if !model.OwnerID.IsNull() && !model.OwnerID.IsUnknown() {
		queryBuilder.OwnerID(model.OwnerID.ValueInt64())
	}

	return queryBuilder.Build()
}

var ipblockResourceSchema = resourceSchema.Schema{
	Attributes: map[string]resourceSchema.Attribute{
		"id": resourceSchema.Int64Attribute{
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"address": resourceSchema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"prefix": resourceSchema.Int64Attribute{
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(32),
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"version": resourceSchema.Int64Attribute{
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(4),
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"status": resourceSchema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Static | Reserved | Available | Subnet | Container",
			Default:     stringdefault.StaticString("Static"),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"status_id": resourceSchema.Int64Attribute{
			Computed: true,
		},
		"asn": resourceSchema.Int64Attribute{
			Computed: true,
		},
		"asn_id": resourceSchema.Int64Attribute{
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"description": resourceSchema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"info": resourceSchema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"interface": resourceSchema.StringAttribute{
			Computed: true,
		},
		"interface_id": resourceSchema.Int64Attribute{
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"monitored": resourceSchema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"owner": resourceSchema.StringAttribute{
			Computed: true,
		},
		"owner_id": resourceSchema.Int64Attribute{
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"parent": resourceSchema.StringAttribute{
			Computed: true,
		},
		"parent_id": resourceSchema.Int64Attribute{
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"rir": resourceSchema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"use_network_broadcast": resourceSchema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"used_by": resourceSchema.StringAttribute{
			Computed: true,
		},
		"used_by_id": resourceSchema.Int64Attribute{
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"vlan": resourceSchema.Int64Attribute{
			Computed: true,
		},
		"vlan_id": resourceSchema.Int64Attribute{
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
	},
}

var ipblockDataSourceSchema = datasourceSchema.Schema{
	Attributes: map[string]datasourceSchema.Attribute{
		"id": datasourceSchema.Int64Attribute{
			Optional: true,
		},
		"address": datasourceSchema.StringAttribute{
			Optional: true,
		},
		"prefix": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"version": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"status": datasourceSchema.StringAttribute{
			Computed: true,
		},
		"status_id": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"asn": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"asn_id": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"description": datasourceSchema.StringAttribute{
			Computed: true,
		},
		"info": datasourceSchema.StringAttribute{
			Computed: true,
		},
		"interface": datasourceSchema.StringAttribute{
			Computed: true,
		},
		"interface_id": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"monitored": datasourceSchema.BoolAttribute{
			Computed: true,
		},
		"owner": datasourceSchema.StringAttribute{
			Computed: true,
		},
		"owner_id": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"parent": datasourceSchema.StringAttribute{
			Computed: true,
		},
		"parent_id": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"rir": datasourceSchema.StringAttribute{
			Computed: true,
		},
		"use_network_broadcast": datasourceSchema.BoolAttribute{
			Computed: true,
		},
		"used_by": datasourceSchema.StringAttribute{
			Computed: true,
		},
		"used_by_id": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"vlan": datasourceSchema.Int64Attribute{
			Computed: true,
		},
		"vlan_id": datasourceSchema.Int64Attribute{
			Computed: true,
		},
	},
}
