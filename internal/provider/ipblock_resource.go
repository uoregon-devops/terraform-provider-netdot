package provider

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"terraform-provider-netdot/internal/netdot"
	"terraform-provider-netdot/internal/netdot/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var ipReservationMutex sync.Mutex

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ipblockResource{}
	_ resource.ResourceWithImportState = &ipblockResource{}
)

// NewIpblockResource is a helper function to simplify the provider implementation.
func NewIpblockResource() resource.Resource {
	return &ipblockResource{}
}

// ipblockResource is the data source implementation.
type ipblockResource struct {
	client *netdot.Client
}

// Configure adds the provider configured client to the data source.
func (d *ipblockResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (d *ipblockResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipblock"
}

func (d *ipblockResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ipblockResourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *ipblockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ipblockModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var netdotIpblock models.IpBlock

	if !state.ID.IsNull() {
		httpStatusCode, err := d.client.GetResourceByID("ipblock", state.ID.ValueInt64(), &netdotIpblock)
		if err != nil {
			if httpStatusCode != nil && *httpStatusCode == http.StatusNotFound {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("Error reading IP block", err.Error())
			return
		}
	}

	state = IPBlockToIpblockModel(netdotIpblock)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ipblockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ipblockModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config ipblockModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryBuilder := IPBlockModelToIPBlockQuery(plan).Builder()

	queryBuilder.SkipReserveFirstN(true)
	queryBuilder.SkipInheritParentOwner(true)

	createQuery := queryBuilder.Build()

	var newIPBlock models.IpBlock

	if config.Address.IsNull() && !config.ParentID.IsNull() {
		ipReservationMutex.Lock()
		defer ipReservationMutex.Unlock()
		existingIpID, addr, err := r.client.GetNextAvailableIP(config.ParentID.ValueInt64(), netdot.IPAllocationStrategyFirstFree)
		if err != nil {
			resp.Diagnostics.AddError("Error getting next available IP", err.Error())
			return
		}
		createQuery.Address = &addr
		if existingIpID == nil {
			err = r.client.CreateResource("ipblock", createQuery, &newIPBlock)
			if err != nil {
				resp.Diagnostics.AddError("Error creating IP block", err.Error())
				return
			}
		} else {
			err = r.client.UpdateResource("ipblock", *existingIpID, createQuery, &newIPBlock)
			if err != nil {
				resp.Diagnostics.AddError("Error creating IP block", err.Error())
				return
			}
		}
	} else {
		err := r.client.CreateResource("ipblock", createQuery, &newIPBlock)
		if err != nil {
			resp.Diagnostics.AddError("Error creating IP block", err.Error())
			return
		}
	}

	state := IPBlockToIpblockModel(newIPBlock)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ipblockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ipblockModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var current_state ipblockModel

	diags = resp.State.Get(ctx, &current_state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateQueryBuilder := IPBlockModelToIPBlockQuery(plan).Builder()
	updateQueryBuilder.SkipReserveFirstN(true)

	updateQuery := updateQueryBuilder.Build()

	var ipblock models.IpBlock
	err := r.client.UpdateResource("ipblock", plan.ID.ValueInt64(), updateQuery, &ipblock)
	if err != nil {
		resp.Diagnostics.AddError("Error updating IP block", err.Error())
		return
	}

	state := IPBlockToIpblockModel(ipblock)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ipblockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ipblockModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteResourceByID("ipblock", state.ID.ValueInt64(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting IP block", err.Error())
		return
	}
}

func (r *ipblockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	myID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), myID)...)
}
