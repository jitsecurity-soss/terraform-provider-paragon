// client.go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

type EnvironmentSecret struct {
    ID         string      `json:"id"`
    Key        string      `json:"key"`
    ProjectID  string      `json:"projectId"`
    Hash       string      `json:"hash"`
    DateCreated string     `json:"dateCreated"`
    DateUpdated string     `json:"dateUpdated"`
    DateDeleted interface{} `json:"dateDeleted"`
}

type CreateEnvironmentSecretRequest struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type UpdateEnvironmentSecretRequest struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

func (c *Client) CreateEnvironmentSecret(ctx context.Context, projectID, key, value string) (*EnvironmentSecret, error) {
    url := fmt.Sprintf("%s/projects/%s/secrets", c.baseURL, projectID)

    reqBody := CreateEnvironmentSecretRequest{
        Key:   key,
        Value: value,
    }
    jsonBody, _ := json.Marshal(reqBody)

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.accessToken)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("failed to create environment secret with status code: %d", resp.StatusCode)
    }

    var secret EnvironmentSecret
    err = json.NewDecoder(resp.Body).Decode(&secret)
    if err != nil {
        return nil, err
    }

    return &secret, nil
}

func (c *Client) GetEnvironmentSecrets(ctx context.Context, projectID string) ([]EnvironmentSecret, error) {
    url := fmt.Sprintf("%s/projects/%s/secrets", c.baseURL, projectID)

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
        return nil, fmt.Errorf("failed to get environment secrets with status code: %d", resp.StatusCode)
    }

    var secrets []EnvironmentSecret
    err = json.NewDecoder(resp.Body).Decode(&secrets)
    if err != nil {
        return nil, err
    }

    return secrets, nil
}

func (c *Client) UpdateEnvironmentSecret(ctx context.Context, projectID, secretID, key, value string) (*EnvironmentSecret, error) {
    url := fmt.Sprintf("%s/projects/%s/secrets/%s", c.baseURL, projectID, secretID)

    reqBody := UpdateEnvironmentSecretRequest{
        Key:   key,
        Value: value,
    }
    jsonBody, _ := json.Marshal(reqBody)

    req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.accessToken)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to update environment secret with status code: %d", resp.StatusCode)
    }

    var updatedSecret EnvironmentSecret
    err = json.NewDecoder(resp.Body).Decode(&updatedSecret)
    if err != nil {
        return nil, err
    }

    return &updatedSecret, nil
}

func (c *Client) DeleteEnvironmentSecret(ctx context.Context, projectID, secretID string) error {
    url := fmt.Sprintf("%s/projects/%s/secrets/%s", c.baseURL, projectID, secretID)

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

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to delete environment secret with status code: %d", resp.StatusCode)
    }

    return nil
}
