// client.go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

type CreateEventDestinationRequest struct {
    ProjectID     string                 `json:"projectId"`
    Type          string                 `json:"type"`
    Configuration EventConfiguration     `json:"configuration"`
}

type EventConfiguration struct {
    EmailTo string            `json:"emailTo,omitempty"`
    URL     string            `json:"url,omitempty"`
    Events  []string          `json:"events"`
    Body    WebhookBody        `json:"body,omitempty"`
}

type EventDestination struct {
    ID            string             `json:"id"`
    ProjectID     string             `json:"projectId"`
    Type          string             `json:"type"`
    State         string             `json:"state"`
    Configuration EventConfiguration `json:"configuration"`
    DateDeleted   string             `json:"dateDeleted"`
    DateCreated   string             `json:"dateCreated"`
    DateUpdated   string             `json:"dateUpdated"`
}

func (c *Client) CreateOrUpdateEventDestination(ctx context.Context, projectID, eventID string, req CreateEventDestinationRequest) (*EventDestination, error) {
    var httpMethod string
    var url string

    if eventID == "" {
        // Create a new event destination
        httpMethod = "POST"
        url = fmt.Sprintf("%s/projects/%s/event-destinations", c.baseURL, projectID)
    } else {
        // Update an existing event destination
        httpMethod = "PUT"
        url = fmt.Sprintf("%s/projects/%s/event-destinations/%s", c.baseURL, projectID, eventID)
    }

    req.ProjectID = projectID
    jsonBody, err := json.Marshal(req)
    tflog.Debug(ctx, fmt.Sprintf("before transform : %s", jsonBody))

    if err != nil {
        return nil, err
    }

    httpReq, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+c.accessToken)

    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("failed to create/update event destination with status code: %d", resp.StatusCode)
    }

    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    tflog.Debug(ctx, fmt.Sprintf("body read read read : %s", responseBody))

    var eventDestination EventDestination
    err = json.Unmarshal(responseBody, &eventDestination)
    if err != nil {
        return nil, fmt.Errorf("error decoding response body: %v, json: %s", err, string(responseBody))
    }

    tflog.Debug(ctx, "RETURNING")

    return &eventDestination, nil
}

func (c *Client) GetEventDestination(ctx context.Context, projectID, eventID string) (*EventDestination, error) {
    url := fmt.Sprintf("%s/projects/%s/event-destinations/%s", c.baseURL, projectID, eventID)

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

    if resp.StatusCode == http.StatusNotFound {
        return nil, fmt.Errorf("event destination not found with status code: %d", resp.StatusCode)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get event destination with status code: %d", resp.StatusCode)
    }

    var eventDestination EventDestination
    err = json.NewDecoder(resp.Body).Decode(&eventDestination)
    if err != nil {
        return nil, err
    }

    return &eventDestination, nil
}

func (c *Client) DeleteEventDestination(ctx context.Context, projectID, eventID string) error {
    url := fmt.Sprintf("%s/projects/%s/event-destinations/%s", c.baseURL, projectID, eventID)

    req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bearer "+c.accessToken)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return nil
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to delete event destination with status code: %d", resp.StatusCode)
    }

    return nil
}