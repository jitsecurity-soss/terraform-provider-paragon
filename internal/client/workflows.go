package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

type Workflow struct {
    ID            string   `json:"id"`
    DateCreated   string   `json:"dateCreated"`
    DateUpdated   string   `json:"dateUpdated"`
    Description   string   `json:"description"`
    ProjectID     string   `json:"projectId"`
    TeamID        string   `json:"teamId"`
    IntegrationID string   `json:"integrationId"`
    WorkflowVersion int    `json:"workflowVersion"`
    Tags          []string `json:"tags"`
}

func (c *Client) GetWorkflows(ctx context.Context, projectID, integrationID string) ([]Workflow, error) {
    url := fmt.Sprintf("%s/projects/%s/workflows?includeDeleted=false&integrationId=%s", c.baseURL, projectID, integrationID)

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+c.accessToken)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get workflows with status code: %d", resp.StatusCode)
    }

    var workflows []Workflow
    err = json.NewDecoder(resp.Body).Decode(&workflows)
    if err != nil {
        return nil, err
    }

    return workflows, nil
}