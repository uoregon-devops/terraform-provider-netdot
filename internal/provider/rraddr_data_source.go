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
	_ datasource.DataSourceWithConfigure = &rrAddrDataSource{}
)

func NewRRAddrDataSource() datasource.DataSource {
	return &rrAddrDataSource{}
}

type rrAddrDataSource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *rrAddrDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *rrAddrDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_arecord"
}

func (d *rrAddrDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rrAddrDataSourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *rrAddrDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state rrAddrModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("ID is required", "ID must be provided")
		return
	}

	var netdotRRAddr models.RRAddr

	_, err := d.client.GetResourceByID("rraddr", state.ID.ValueInt64(), &netdotRRAddr)
	if err != nil {
		resp.Diagnostics.AddError("Error reading RR", err.Error())
		return
	}

	state = RRAddrToRRAddrModel(netdotRRAddr)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
