package provider

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &rrNsResource{}
)

func NewRRNsResource() resource.Resource {
	return &rrNsResource{}
}

type rrNsResource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *rrNsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (d *rrNsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ns"
}

func (d *rrNsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rrNsResourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *rrNsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rrNsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("ID is required", "ID must be provided")
		return
	}

	var netdotRR models.RRNs

	httpStatusCode, err := d.client.GetResourceByID("rrns", state.ID.ValueInt64(), &netdotRR)
	if err != nil {
		if httpStatusCode != nil && *httpStatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading RRNS", err.Error())
		return
	}

	newState := RRNsToRRNsModel(netdotRR)

	// Set state
	diags := resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *rrNsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rrNsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createQuery := RRNsModelToRRNsQuery(plan)

	var newRRNs models.RRNs
	err := r.client.CreateResource("rrns", createQuery, &newRRNs)
	if err != nil {
		resp.Diagnostics.AddError("Error creating RRNS", err.Error())
		return
	}

	state := RRNsToRRNsModel(newRRNs)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *rrNsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rrNsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var current_state rrNsModel

	diags = resp.State.Get(ctx, &current_state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateQuery := RRNsModelToRRNsQuery(plan)

	if current_state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be greater than 0")
		return
	}

	if current_state.ID.ValueInt64() != plan.ID.ValueInt64() {
		resp.Diagnostics.AddError("ID mismatch", "ID in plan does not match ID in state")
		return
	}

	var updatedRRNs models.RRNs
	err := r.client.UpdateResource("rrns", current_state.ID.ValueInt64(), updateQuery, &updatedRRNs)
	if err != nil {
		resp.Diagnostics.AddError("Error updating RRNS", err.Error())
		return
	}

	state := RRNsToRRNsModel(updatedRRNs)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *rrNsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state rrNsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be greater than 0")
		return
	}

	qBuilder := netdot.NewRRNsQueryBuilder()
	query := qBuilder.Build()

	// Delete existing order
	err := r.client.DeleteResourceByID("rrns", state.ID.ValueInt64(), query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting RRNS",
			"Could not delete RRNS, unexpected error: "+err.Error(),
		)
		return
	}
}
