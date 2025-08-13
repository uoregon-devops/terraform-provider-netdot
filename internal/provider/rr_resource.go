package provider

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &rrResource{}
)

func NewRResource() resource.Resource {
	return &rrResource{}
}

type rrResource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *rrResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (d *rrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rr"
}

func (d *rrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rrResourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *rrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rrModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("ID is required", "ID must be provided")
		return
	}

	var netdotRR models.RR

	httpStatusCode, err := d.client.GetResourceByID("rr", state.ID.ValueInt64(), &netdotRR)
	if err != nil {
		if httpStatusCode != nil && *httpStatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading RR", err.Error())
		return
	}

	newState := RRToRRModel(netdotRR)

	// Set state
	diags := resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *rrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rrModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var zoneSearchResults struct {
		RRs []struct {
			ID int64 `xml:"id,attr"`
		} `xml:"Zone"`
	}

	httpCode, err := r.client.Get("/rest/zone?name="+plan.Zone.ValueString(), &zoneSearchResults)
	if err != nil {
		if httpCode != nil && *httpCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Zone not found", "No zone found named \""+plan.Zone.ValueString()+"\"")
			return
		}
		resp.Diagnostics.AddError("Error reading RR", err.Error())
		return
	}

	plan.ZoneID = types.Int64Value(zoneSearchResults.RRs[0].ID)
	createQuery := RRModelToRRQuery(plan)

	var newRR models.RR
	err = r.client.CreateResource("rr", createQuery, &newRR)
	if err != nil {
		resp.Diagnostics.AddError("Error creating RR", err.Error())
		return
	}

	state := RRToRRModel(newRR)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *rrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rrModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var current_state rrModel

	diags = resp.State.Get(ctx, &current_state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var zoneSearchResults struct {
		RRs []struct {
			ID int64 `xml:"id,attr"`
		} `xml:"Zone"`
	}

	httpCode, err := r.client.Get("/rest/zone?name="+plan.Zone.ValueString(), &zoneSearchResults)
	if err != nil {
		if httpCode != nil && *httpCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Zone not found", "No zone found named \""+plan.Zone.ValueString()+"\"")
			return
		}
		resp.Diagnostics.AddError("Error reading RR", err.Error())
		return
	}

	plan.ZoneID = types.Int64Value(zoneSearchResults.RRs[0].ID)

	updateQuery := RRModelToRRQuery(plan)

	if current_state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be greater than 0")
		return
	}

	if current_state.ID.ValueInt64() != plan.ID.ValueInt64() {
		resp.Diagnostics.AddError("ID mismatch", "ID in plan does not match ID in state")
		return
	}

	var updatedRR models.RR
	err = r.client.UpdateResource("rr", current_state.ID.ValueInt64(), updateQuery, &updatedRR)
	if err != nil {
		resp.Diagnostics.AddError("Error updating RR", err.Error())
		return
	}

	state := RRToRRModel(updatedRR)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *rrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state rrModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be greater than 0")
		return
	}

	// Delete existing order
	err := r.client.DeleteResourceByID("rr", state.ID.ValueInt64(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting RR",
			"Could not delete RR, unexpected error: "+err.Error(),
		)
		return
	}
}
