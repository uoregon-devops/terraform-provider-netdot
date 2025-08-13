package provider

import (
	"context"
	"fmt"
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigure = &rrCnameDataSource{}
)

func NewRRCnameDataSource() datasource.DataSource {
	return &rrCnameDataSource{}
}

type rrCnameDataSource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *rrCnameDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netdot.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *rrCnameDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cname"
}

func (d *rrCnameDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rrCnameDataSourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *rrCnameDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state rrCnameModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("ID is required", "ID must be provided")
		return
	}

	var netdotRRCname models.RRCname

	_, err := d.client.GetResourceByID("rrcname", state.ID.ValueInt64(), &netdotRRCname)
	if err != nil {
		resp.Diagnostics.AddError("Error reading RR", err.Error())
		return
	}

	state = RRCnameToRRCnameModel(netdotRRCname)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
