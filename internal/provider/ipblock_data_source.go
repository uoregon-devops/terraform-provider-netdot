package provider

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	// _ datasource.DataSource              = &ipblockDataSource{}
	_ datasource.DataSourceWithConfigure = &ipblockDataSource{}
)

// NewIpblockDataSource is a helper function to simplify the provider implementation.
func NewIpblockDataSource() datasource.DataSource {
	return &ipblockDataSource{}
}

// ipblockDataSource is the data source implementation.
type ipblockDataSource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *ipblockDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ipblockDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipblock"
}

func (d *ipblockDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ipblockDataSourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *ipblockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ipblockModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() && state.Address.IsNull() {
		resp.Diagnostics.AddError("ID or Address is required", "ID or Address must be provided")
		return
	}

	if !state.ID.IsNull() && !state.Address.IsNull() {
		resp.Diagnostics.AddError("ID and Address are mutually exclusive", "Only one of ID or Address can be provided")
		return
	}

	var netdotIpblock models.IpBlock

	if !state.ID.IsNull() {
		_, err := d.client.GetResourceByID("ipblock", state.ID.ValueInt64(), &netdotIpblock)
		if err != nil {
			resp.Diagnostics.AddError("Error reading IP block", err.Error())
			return
		}
	} else {
		var ipblockSearchResults struct {
			IpBlocks []models.IpBlock `xml:"Ipblock"`
		}
		statusCode, err := d.client.Get("/rest/ipblock?address="+state.Address.ValueString(), &ipblockSearchResults)
		if err != nil {
			if statusCode != nil && *statusCode == http.StatusNotFound {
				resp.Diagnostics.AddError("IP block not found", "No IP block found with the provided address: "+state.Address.ValueString())
				return
			}
			resp.Diagnostics.AddError("Error reading IP block", err.Error())
			return
		}

		if len(ipblockSearchResults.IpBlocks) > 1 {
			resp.Diagnostics.AddError("Multiple IP blocks found", "Multiple IP blocks found with the provided address: "+state.Address.ValueString())
			return
		}

		netdotIpblock = ipblockSearchResults.IpBlocks[0]
	}

	state = IPBlockToIpblockModel(netdotIpblock)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
