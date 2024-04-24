package provider

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
)

var (
    _ datasource.DataSource              = &workflowDataSource{}
    _ datasource.DataSourceWithConfigure = &workflowDataSource{}
)

func NewWorkflowDataSource() datasource.DataSource {
    return &workflowDataSource{}
}

type workflowDataSource struct {
    client *client.Client
}

type workflowDataSourceModel struct {
    ID              types.String   `tfsdk:"id"`
    ProjectID       types.String   `tfsdk:"project_id"`
    IntegrationID   types.String   `tfsdk:"integration_id"`
    Description     types.String   `tfsdk:"description"`
    DateCreated     types.String   `tfsdk:"date_created"`
    DateUpdated     types.String   `tfsdk:"date_updated"`
    Tags            []types.String `tfsdk:"tags"`
    WorkflowVersion types.Int64    `tfsdk:"workflow_version"`
}

func (d *workflowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*client.Client)
    if !ok {
        return
    }
    d.client = client
}

func (d *workflowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (d *workflowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Fetches a workflow by project ID, integration ID, and description.",
        Attributes: map[string]schema.Attribute{
            "project_id": schema.StringAttribute{
                Description: "The ID of the project.",
                Required:    true,
            },
            "integration_id": schema.StringAttribute{
                Description: "The ID of the integration.",
                Required:    true,
            },
            "description": schema.StringAttribute{
                Description: "The description of the workflow.",
                Required:    true,
            },
            "id": schema.StringAttribute{
                Description: "The ID of the workflow.",
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
    }
}

func (d *workflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var config workflowDataSourceModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := config.ProjectID.ValueString()
    integrationID := config.IntegrationID.ValueString()
    description := config.Description.ValueString()

    workflows, err := d.client.GetWorkflows(ctx, projectID, integrationID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Read Workflows",
            err.Error(),
        )
        return
    }

    var foundWorkflow *client.Workflow
    for _, workflow := range workflows {
        if workflow.Description == description {
            foundWorkflow = &workflow
            break
        }
    }

    if foundWorkflow == nil {
        resp.Diagnostics.AddError(
            "Workflow Not Found",
            fmt.Sprintf("Workflow with description '%s' not found", description),
        )
        return
    }

    state := workflowDataSourceModel{
        ID:              types.StringValue(foundWorkflow.ID),
        ProjectID:       types.StringValue(foundWorkflow.ProjectID),
        IntegrationID:   types.StringValue(foundWorkflow.IntegrationID),
        Description:     types.StringValue(foundWorkflow.Description),
        DateCreated:     types.StringValue(foundWorkflow.DateCreated),
        DateUpdated:     types.StringValue(foundWorkflow.DateUpdated),
        Tags:            client.ConvertStringSliceToTypesStringSlice(foundWorkflow.Tags),
        WorkflowVersion: types.Int64Value(int64(foundWorkflow.WorkflowVersion)),
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
