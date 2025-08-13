package provider

import (
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type hostModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Active     types.Bool   `tfsdk:"active"`
	AutoUpdate types.Bool   `tfsdk:"auto_update"`
	Expiration types.String `tfsdk:"expiration"`
	Info       types.String `tfsdk:"info"`
	Name       types.String `tfsdk:"name"`
	Zone       types.String `tfsdk:"zone"`
	ZoneID     types.Int64  `tfsdk:"zone_id"`
	FQDN       types.String `tfsdk:"fqdn"`
	// nested objects
	IpBlocks ipblockModel `tfsdk:"ip_addresses"`
}

func RRToHostModel(rr models.RR) hostModel {
	var finalModel hostModel

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

func HostModelToRRQuery(model hostModel) netdot.RRQuery {
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

var hostResourceSchema = rrResourceSchema

func appendSchemaAttributesAsSetNetstedAttribute(schema *resourceSchema.Schema, attributeName string, attributes map[string]resourceSchema.Attribute) {
	schema.Attributes[attributeName] = resourceSchema.SetNestedAttribute{
		NestedObject: resourceSchema.NestedAttributeObject{
			Attributes: attributes,
		},
		Optional: true,
	}
}

func init() {
	appendSchemaAttributesAsSetNetstedAttribute(&hostResourceSchema, "ip_addresses", ipblockResourceSchema.Attributes)
}
