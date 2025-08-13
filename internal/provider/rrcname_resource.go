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
	_ resource.Resource = &rrCnameResource{}
)

func NewRRCnameResource() resource.Resource {
	return &rrCnameResource{}
}

type rrCnameResource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *rrCnameResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (d *rrCnameResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cname"
}

func (d *rrCnameResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rrCnameResourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *rrCnameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rrCnameModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("ID is required", "ID must be provided")
		return
	}

	var netdotRR models.RRCname

	httpStatusCode, err := d.client.GetResourceByID("rrcname", state.ID.ValueInt64(), &netdotRR)
	if err != nil {
		if httpStatusCode != nil && *httpStatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading RR", err.Error())
		return
	}

	newState := RRCnameToRRCnameModel(netdotRR)

	// Set state
	diags := resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *rrCnameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rrCnameModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createQuery := RRCnameModelToRRCnameQuery(plan)

	var newRRCname models.RRCname
	err := r.client.CreateResource("rrcname", createQuery, &newRRCname)
	if err != nil {
		resp.Diagnostics.AddError("Error creating RR", err.Error())
		return
	}

	state := RRCnameToRRCnameModel(newRRCname)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *rrCnameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rrCnameModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var current_state rrCnameModel

	diags = resp.State.Get(ctx, &current_state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateQuery := RRCnameModelToRRCnameQuery(plan)

	if current_state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be greater than 0")
		return
	}

	if current_state.ID.ValueInt64() != plan.ID.ValueInt64() {
		resp.Diagnostics.AddError("ID mismatch", "ID in plan does not match ID in state")
		return
	}

	var updatedRRCname models.RRCname
	err := r.client.UpdateResource("rrcname", current_state.ID.ValueInt64(), updateQuery, &updatedRRCname)
	if err != nil {
		resp.Diagnostics.AddError("Error updating RR", err.Error())
		return
	}

	state := RRCnameToRRCnameModel(updatedRRCname)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *rrCnameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state rrCnameModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be greater than 0")
		return
	}

	qBuilder := netdot.NewRRCnameQueryBuilder()
	qBuilder.SkipDeletingRR(true)
	query := qBuilder.Build()

	// Delete existing order
	err := r.client.DeleteResourceByID("rrcname", state.ID.ValueInt64(), query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting RR",
			"Could not delete RR, unexpected error: "+err.Error(),
		)
		return
	}
}
