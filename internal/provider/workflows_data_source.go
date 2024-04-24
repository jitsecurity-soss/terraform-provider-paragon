package provider

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"

)

var (
    _ datasource.DataSource              = &workflowsDataSource{}
    _ datasource.DataSourceWithConfigure = &workflowsDataSource{}
)

func NewWorkflowsDataSource() datasource.DataSource {
    return &workflowsDataSource{}
}

type workflowsDataSource struct {
    client *client.Client
}

type workflowsDataSourceModel struct {
    ProjectID     types.String                  `tfsdk:"project_id"`
    IntegrationID types.String                 `tfsdk:"integration_id"`
    Workflows     []workflowDataSourceModel     `tfsdk:"workflows"`
}


func (d *workflowsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*client.Client)
    if !ok {
        return
    }
    d.client = client
}

func (d *workflowsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflows"
}

func (d *workflowsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Fetches a list of workflows by project ID and integration ID.",
        Attributes: map[string]schema.Attribute{
            "project_id": schema.StringAttribute{
                Description: "The ID of the project.",
                Required:    true,
            },
            "integration_id": schema.StringAttribute{
                Description: "The ID of the integration.",
                Required:    true,
            },
            "workflows": schema.ListNestedAttribute{
                Description: "The list of workflows.",
                Computed:    true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.StringAttribute{
                            Description: "The ID of the workflow.",
                            Computed:    true,
                        },
                        "description": schema.StringAttribute{
                            Description: "The description of the workflow.",
                            Computed:    true,
                        },
                        "project_id": schema.StringAttribute{
                            Description: "The ID of the project.",
                            Computed:    true,
                        },
                        "integration_id": schema.StringAttribute{
                            Description: "The ID of the integration.",
                            Computed:    true,
                        },
                        "date_created": schema.StringAttribute{
                            Description: "The creation date of the workflow.",
                            Computed:    true,
                        },
                        "date_updated": schema.StringAttribute{
                            Description: "The last update date of the workflow.",
                            Computed:    true,
                        },
                        "tags": schema.ListAttribute{
                            Description: "The tags associated with the workflow.",
                            Computed:    true,
                            ElementType: types.StringType,
                        },
                        "workflow_version": schema.Int64Attribute{
                            Description: "The version of the workflow.",
                            Computed:    true,
                        },
                    },
                },
            },
        },
    }
}

func (d *workflowsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var config workflowsDataSourceModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := config.ProjectID.ValueString()
    integrationID := config.IntegrationID.ValueString()

    workflows, err := d.client.GetWorkflows(ctx, projectID, integrationID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Read Workflows",
            err.Error(),
        )
        return
    }

    var workflowModels []workflowDataSourceModel
    for _, workflow := range workflows {
        workflowModel := workflowDataSourceModel{
            ID:              types.StringValue(workflow.ID),
            Description:     types.StringValue(workflow.Description),
            ProjectID:       types.StringValue(projectID),
            IntegrationID:   types.StringValue(integrationID),
            DateCreated:     types.StringValue(workflow.DateCreated),
            DateUpdated:     types.StringValue(workflow.DateUpdated),
            Tags:            client.ConvertStringSliceToTypesStringSlice(workflow.Tags),
            WorkflowVersion: types.Int64Value(int64(workflow.WorkflowVersion)),
        }
        workflowModels = append(workflowModels, workflowModel)
    }

    state := workflowsDataSourceModel{
        ProjectID:     types.StringValue(projectID),
        IntegrationID: types.StringValue(integrationID),
        Workflows:     workflowModels,
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
