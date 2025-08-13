package provider

import (
	"context"
	"terraform-provider-netdot/internal/netdot"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &netdotProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &netdotProvider{
			version: version,
		}
	}
}

// netdotProvider is the provider implementation.
type netdotProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *netdotProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netdot"
	resp.Version = p.version
}

type netdotProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Schema defines the provider-level schema for configuration data.
func (p *netdotProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *netdotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config netdotProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Bad netdot host",
			"The provider cannot create the netdot API client as there is an unknown configuration value for the netdot API host.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Bad netdot username",
			"The provider cannot create the netdot API client as there is an unknown configuration value for the netdot API username.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Bad netdot password",
			"The provider cannot create the netdot API client as there is an unknown configuration value for the netdot API password.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	netdot_client := netdot.NewClient(config.Host.ValueString(), config.Username.ValueString(), config.Password.ValueString())
	err := netdot_client.Authenticate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create netdot API client",
			"An unexpected error occurred when creating the netdot API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"netdot Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = netdot_client
	resp.ResourceData = netdot_client
}

// DataSources defines the data sources implemented in the provider.
func (p *netdotProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewIpblockDataSource,
		NewRRDataSource,
		NewRRAddrDataSource,
		NewRRCnameDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *netdotProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRResource,
		NewIpblockResource,
		NewRRAddrResource,
		NewRRCnameResource,
	}
}
