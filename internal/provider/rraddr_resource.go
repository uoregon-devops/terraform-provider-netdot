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
	_ resource.Resource = &rrAddrResource{}
)

func NewRRAddrResource() resource.Resource {
	return &rrAddrResource{}
}

type rrAddrResource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *rrAddrResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (d *rrAddrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_arecord"
}

func (d *rrAddrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rrAddrResourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *rrAddrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rrAddrModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("ID is required", "ID must be provided")
		return
	}

	var netdotRR models.RRAddr

	httpStatusCode, err := d.client.GetResourceByID("rraddr", state.ID.ValueInt64(), &netdotRR)
	if err != nil {
		if httpStatusCode != nil && *httpStatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading RR", err.Error())
		return
	}

	newState := RRAddrToRRAddrModel(netdotRR)

	// Set state
	diags := resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *rrAddrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rrAddrModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createQuery := RRAddrModelToRRAddrQuery(plan)

	var newRRAddr models.RRAddr
	err := r.client.CreateResource("rraddr", createQuery, &newRRAddr)
	if err != nil {
		resp.Diagnostics.AddError("Error creating RR", err.Error())
		return
	}

	state := RRAddrToRRAddrModel(newRRAddr)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *rrAddrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rrAddrModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var current_state rrAddrModel

	diags = resp.State.Get(ctx, &current_state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateQuery := RRAddrModelToRRAddrQuery(plan)

	if current_state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be greater than 0")
		return
	}

	if current_state.ID.ValueInt64() != plan.ID.ValueInt64() {
		resp.Diagnostics.AddError("ID mismatch", "ID in plan does not match ID in state")
		return
	}

	var updatedRRAddr models.RRAddr
	err := r.client.UpdateResource("rraddr", current_state.ID.ValueInt64(), updateQuery, &updatedRRAddr)
	if err != nil {
		resp.Diagnostics.AddError("Error updating RR", err.Error())
		return
	}

	state := RRAddrToRRAddrModel(updatedRRAddr)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *rrAddrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state rrAddrModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be greater than 0")
		return
	}

	qBuilder := netdot.NewRRAddrQueryBuilder()
	qBuilder.SkipDeletingRR(true)
	qBuilder.NoChangeStatus(true)
	query := qBuilder.Build()

	// Delete existing order
	err := r.client.DeleteResourceByID("rraddr", state.ID.ValueInt64(), query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting RR",
			"Could not delete RR, unexpected error: "+err.Error(),
		)
		return
	}
}
