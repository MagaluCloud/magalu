package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkVMImages "magalu.cloud/lib/products/virtual_machine/images"
	"magalu.cloud/sdk"
)

var _ datasource.DataSource = &DataSourceVmImages{}

type DataSourceVmImages struct {
	sdkClient *mgcSdk.Client
	vmImages  sdkVMImages.Service
}

type ImageModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Platform types.String `tfsdk:"platform"`
}

type ImagesModel struct {
	Images []ImageModel `tfsdk:"images"`
}

func NewDataSourceVMIMages() datasource.DataSource {
	return &DataSourceVmImages{}
}

func (r *DataSourceVmImages) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_images"
}

func (r *DataSourceVmImages) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected provider config, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	r.sdkClient = mgcSdk.NewClient(sdk)
	r.vmImages = sdkVMImages.NewService(ctx, r.sdkClient)
}

func (r *DataSourceVmImages) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of available VM Images.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.BoolAttribute{
							Computed:    true,
							Description: "ID of image.",
						},
						"plataform": schema.StringAttribute{
							Computed:    true,
							Description: "The image platform.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The image name.",
						},
					},
				},
			},
		},
	}
	resp.Schema.Description = "Get the available virtual-machine images."
}

const imageActive string = "active"

func (r *DataSourceVmImages) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ImagesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.vmImages.List(sdkVMImages.ListParameters{}, sdkVMImages.ListConfigs{})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, image := range sdkOutput.Images {
		if image.Status != imageActive {
			continue
		}

		platform := ""
		if image.Platform != nil {
			platform = *image.Platform
		}

		data.Images = append(data.Images, ImageModel{
			ID:       types.StringValue(image.Id),
			Name:     types.StringValue(image.Name),
			Platform: types.StringValue(platform),
		})

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
