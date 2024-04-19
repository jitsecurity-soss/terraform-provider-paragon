// client.go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)


type SDKKey struct {
    ID            string      `json:"id"`
    ProjectID     string      `json:"projectId"`
    AuthType      string      `json:"authType"`
    AuthConfig    AuthConfig  `json:"authConfig"`
    DateDeleted   interface{} `json:"dateDeleted"`
    DateCreated   string      `json:"dateCreated"`
    DateUpdated   string      `json:"dateUpdated"`
    Revoked       bool        `json:"revoked"`
    PrivateKey    string      `json:"privateKey,omitempty"`
}

type AuthConfig struct {
    Paragon ParagonConfig `json:"paragon"`
}

type ParagonConfig struct {
    PublicKey string `json:"publicKey"`
    GeneratedDate string      `json:"generatedDate"`
}

func (c *Client) GetSDKKeys(ctx context.Context, projectID string) ([]SDKKey, error) {
    url := fmt.Sprintf("%s/projects/%s/keys", c.baseURL, projectID)

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
        return nil, fmt.Errorf("failed to get SDK keys with status code: %d", resp.StatusCode)
    }

    var sdkKeys []SDKKey
    err = json.NewDecoder(resp.Body).Decode(&sdkKeys)
    if err != nil {
        return nil, err
    }

    return sdkKeys, nil
}

func (c *Client) CreateSDKKey(ctx context.Context, projectID string) (*SDKKey, error) {
    url := fmt.Sprintf("%s/projects/%s/keys", c.baseURL, projectID)

    reqBody := map[string]string{
        "projectId": projectID,
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
        return nil, fmt.Errorf("failed to create SDK key with status code: %d", resp.StatusCode)
    }

    var sdkKey SDKKey
    err = json.NewDecoder(resp.Body).Decode(&sdkKey)
    if err != nil {
        return nil, err
    }

    return &sdkKey, nil
}

func (c *Client) DeleteSDKKey(ctx context.Context, projectID, keyID string) error {
    url := fmt.Sprintf("%s/projects/%s/keys/%s", c.baseURL, projectID, keyID)

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
        return fmt.Errorf("failed to delete SDK key with status code: %d", resp.StatusCode)
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to delete SDK key with status code: %d", resp.StatusCode)
    }

    return nil
}
