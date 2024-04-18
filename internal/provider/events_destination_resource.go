// provider/events_destination_resource.go
package provider

import (
    "context"
    "regexp"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/attr"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &eventsDestinationResource{}
    _ resource.ResourceWithConfigure = &eventsDestinationResource{}
)

// NewEventsDestinationResource is a helper function to simplify the provider implementation.
func NewEventsDestinationResource() resource.Resource {
    return &eventsDestinationResource{}
}

// eventsDestinationResource is the resource implementation.
type eventsDestinationResource struct {
    client *client.Client
}

// eventsDestinationResourceModel maps the resource schema data.
type eventsDestinationResourceModel struct {
    ID        types.String   `tfsdk:"id"`
    ProjectID types.String   `tfsdk:"project_id"`
    Events    types.List     `tfsdk:"events"`
    Email     *emailBlock    `tfsdk:"email"`
    Webhook   *webhookBlock  `tfsdk:"webhook"`
}

type emailBlock struct {
    Address types.String `tfsdk:"address"`
}

type webhookBlock struct {
    URL     types.String           `tfsdk:"url"`
    Body    types.Map              `tfsdk:"body"`
    Headers []webhookHeaderBlock   `tfsdk:"headers"`
}

type webhookHeaderBlock struct {
    Key   types.String `tfsdk:"key"`
    Value types.String `tfsdk:"value"`
}

// Configure adds the provider configured client to the resource.
func (r *eventsDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*client.Client)
}

// Metadata returns the resource type name.
func (r *eventsDestinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_events_destination"
}

// Schema defines the schema for the resource.
func (r *eventsDestinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages an events destination.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Identifier of the events destination.",
                Computed:    true,
            },
            "project_id": schema.StringAttribute{
                Description: "Identifier of the project.",
                Required:    true,
            },
            "events": schema.ListAttribute{
                ElementType: types.StringType,
                Description: "List of events to subscribe to.",
                Required:    true,
            },
            "email": schema.SingleNestedAttribute{
                Description: "Email destination configuration.",
                Optional:    true,
                Attributes: map[string]schema.Attribute{
                    "address": schema.StringAttribute{
                        Description: "Email address to send notifications to.",
                        Required:    true,
                        Validators: []validator.String{
                            // Validate string value satisfies the regular expression for alphanumeric characters
                            stringvalidator.RegexMatches(
                                regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
                                "Must be a valid email address",
                            ),
                        },
                    },
                },
            },
            "webhook": schema.SingleNestedAttribute{
                Description: "Webhook destination configuration.",
                Optional:    true,
                Attributes: map[string]schema.Attribute{
                    "url": schema.StringAttribute{
                        Description: "URL to send webhook notifications to.",
                        Required:    true,
                    },
                    "body": schema.MapAttribute{
                        ElementType: types.StringType,
                        Description: "Body to send with the webhook.",
                        Required:    true,
                    },
                    "headers": schema.ListNestedAttribute{
                        Description: "List of headers to include with the webhook.",
                        Optional:    true,
                        NestedObject: schema.NestedAttributeObject{
                            Attributes: map[string]schema.Attribute{
                                "key": schema.StringAttribute{
                                    Description: "Header key.",
                                    Required:    true,
                                },
                                "value": schema.StringAttribute{
                                    Description: "Header value.",
                                    Required:    true,
                                    Sensitive:   true,
                                },
                            },
                        },
                    },
                },
            },
        },
    }
}

// Create creates the resource and sets the initial Terraform state.
func (r *eventsDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan eventsDestinationResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Check if both email and webhook blocks are provided
    if plan.Email != nil && plan.Webhook != nil {
        resp.Diagnostics.AddError(
            "Invalid configuration",
            "Both email and webhook blocks cannot be provided together.",
        )
        return
    }

    // Check if neither email nor webhook block is provided
    if plan.Email == nil && plan.Webhook == nil {
        resp.Diagnostics.AddError(
            "Invalid configuration",
            "Either `email` or `webhook` block must be provided.",
        )
        return
    }

    // Convert events from types.List to []string
    events := make([]string, len(plan.Events.Elements()))
    for i, event := range plan.Events.Elements() {
        events[i] = event.(types.String).ValueString()
    }

    // Create the events destination
    var eventDestination *client.EventDestination
    var err error
    if plan.Email != nil {
        eventDestination, err = r.client.CreateOrUpdateEventDestination(ctx, plan.ProjectID.ValueString(), "", client.CreateEventDestinationRequest{
            Type: "email",
            Configuration: client.EventConfiguration{
                EmailTo: plan.Email.Address.ValueString(),
                Events:  events,
            },
        })
    } else if plan.Webhook != nil {
        headers := make(map[string]string)
        for _, header := range plan.Webhook.Headers {
            headers[header.Key.ValueString()] = header.Value.ValueString()
        }
        body := make(map[string]string)
        for key, value := range plan.Webhook.Body.Elements() {
            body[key] = value.(types.String).ValueString()
        }
        eventDestination, err = r.client.CreateOrUpdateEventDestination(ctx, plan.ProjectID.ValueString(), "", client.CreateEventDestinationRequest{
            Type: "webhook",
            Configuration: client.EventConfiguration{
                URL:     plan.Webhook.URL.ValueString(),
                Body:    body,
                Events:  events,
                Headers: headers,
            },
        })
    }
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating event destination",
            err.Error(),
        )
        return
    }

    // Set state to fully populated data
    plan.ID = types.StringValue(eventDestination.ID)
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read refreshes the Terraform state with the latest data.
// Read refreshes the Terraform state with the latest data.
func (r *eventsDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get the events destination ID from the state
    var state eventsDestinationResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the events destination using the API
    eventDestination, err := r.client.GetEventDestination(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading event destination",
            err.Error(),
        )
        return
    }

    // Update the state with the retrieved data
    state.ID = types.StringValue(eventDestination.ID)
    state.ProjectID = types.StringValue(eventDestination.ProjectID)

    events := make([]attr.Value, len(eventDestination.Configuration.Events))
    for i, event := range eventDestination.Configuration.Events {
        events[i] = types.StringValue(event)
    }
    state.Events = types.ListValueMust(types.StringType, events)

    if eventDestination.Type == "email" {
        state.Email = &emailBlock{
            Address: types.StringValue(eventDestination.Configuration.EmailTo),
        }
        state.Webhook = nil
    } else if eventDestination.Type == "webhook" {
        headers := make([]webhookHeaderBlock, 0, len(eventDestination.Configuration.Headers))
        for key, value := range eventDestination.Configuration.Headers {
            headers = append(headers, webhookHeaderBlock{
                Key:   types.StringValue(key),
                Value: types.StringValue(value),
            })
        }
        body := make(map[string]attr.Value)
        for key, value := range eventDestination.Configuration.Body {
            body[key] = types.StringValue(value)
        }
        state.Email = nil
        state.Webhook = &webhookBlock{
            URL:     types.StringValue(eventDestination.Configuration.URL),
            Body:    types.MapValueMust(types.StringType, body),
            Headers: headers,
        }
    }

    // Set the refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *eventsDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
   // Retrieve values from plan and current state
   var plan, state eventsDestinationResourceModel
   diags := req.Plan.Get(ctx, &plan)
   resp.Diagnostics.Append(diags...)
   diags = req.State.Get(ctx, &state)
   resp.Diagnostics.Append(diags...)
   if resp.Diagnostics.HasError() {
       return
   }

   // Check if both email and webhook blocks are provided
   if plan.Email != nil && plan.Webhook != nil {
       resp.Diagnostics.AddError(
           "Invalid configuration",
           "Both email and webhook blocks cannot be provided together.",
       )
       return
   }

   // Convert events from types.List to []string
   events := make([]string, len(plan.Events.Elements()))
   for i, event := range plan.Events.Elements() {
       events[i] = event.(types.String).ValueString()
   }

   // Update the events destination
   var eventDestination *client.EventDestination
   var err error
   if plan.Email != nil {
       eventDestination, err = r.client.CreateOrUpdateEventDestination(ctx, plan.ProjectID.ValueString(), state.ID.ValueString(), client.CreateEventDestinationRequest{
           Type: "email",
           Configuration: client.EventConfiguration{
               EmailTo: plan.Email.Address.ValueString(),
               Events:  events,
           },
       })
   } else if plan.Webhook != nil {
       headers := make(map[string]string)
       for _, header := range plan.Webhook.Headers {
           headers[header.Key.ValueString()] = header.Value.ValueString()
       }
       body := make(map[string]string)
       for key, value := range plan.Webhook.Body.Elements() {
           body[key] = value.(types.String).ValueString()
       }
       eventDestination, err = r.client.CreateOrUpdateEventDestination(ctx, plan.ProjectID.ValueString(), state.ID.ValueString(), client.CreateEventDestinationRequest{
           Type: "webhook",
           Configuration: client.EventConfiguration{
               URL:     plan.Webhook.URL.ValueString(),
               Body:    body,
               Events:  events,
               Headers: headers,
           },
       })
   }
   if err != nil {
       resp.Diagnostics.AddError(
           "Error updating event destination",
           err.Error(),
       )
       return
   }

   // Set state to fully populated data
   plan.ID = types.StringValue(eventDestination.ID)
   diags = resp.State.Set(ctx, plan)
   resp.Diagnostics.Append(diags...)
   if resp.Diagnostics.HasError() {
       return
   }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *eventsDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
   var state eventsDestinationResourceModel
   diags := req.State.Get(ctx, &state)
   resp.Diagnostics.Append(diags...)
   if resp.Diagnostics.HasError() {
       return
   }

   // Delete the events destination
   err := r.client.DeleteEventDestination(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
   if err != nil {
       resp.Diagnostics.AddError(
           "Error deleting event destination",
           err.Error(),
       )
       return
   }
}