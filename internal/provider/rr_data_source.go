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
	_ datasource.DataSourceWithConfigure = &rrDataSource{}
)

func NewRRDataSource() datasource.DataSource {
	return &rrDataSource{}
}

type rrDataSource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *rrDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *rrDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rr"
}

func (d *rrDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rrDataSourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *rrDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state rrModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() && state.Name.IsNull() {
		resp.Diagnostics.AddError("ID or Name is required", "ID or Name must be provided")
		return
	}

	if !state.ID.IsNull() && !state.Name.IsNull() {
		resp.Diagnostics.AddError("ID and Name are mutually exclusive", "Only one of ID or Name can be provided")
		return
	}

	var netdotRR models.RR

	if !state.ID.IsNull() {
		_, err := d.client.GetResourceByID("rr", state.ID.ValueInt64(), &netdotRR)
		if err != nil {
			resp.Diagnostics.AddError("Error reading RR", err.Error())
			return
		}
	} else {
		var rrSearchResults struct {
			RRs []models.RR `xml:"RR"`
		}
		_, err := d.client.Get("/rest/rr?name="+state.Name.ValueString(), &rrSearchResults)
		if err != nil {
			resp.Diagnostics.AddError("Error reading RR", err.Error())
			return
		}

		if len(rrSearchResults.RRs) == 0 {
			resp.Diagnostics.AddError("RR not found", "No RR found with the provided name: "+state.Name.ValueString())
			return
		}

		if len(rrSearchResults.RRs) > 1 {
			resp.Diagnostics.AddError("Multiple RRs found", "Multiple RRs found with the provided name")
			return
		}

		netdotRR = rrSearchResults.RRs[0]
	}

	state = RRToRRModel(netdotRR)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
