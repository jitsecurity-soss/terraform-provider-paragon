// client.go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "encoding/base64"
    "strings"
    "fmt"
    "net/http"
)

type CLIKeyResponse struct {
    Key string `json:"key"`
}

type CLIKey struct {
    ID           string `json:"id"`
    DateCreated  string `json:"dateCreated"`
    DateUpdated  string `json:"dateUpdated"`
    UserID       string `json:"userId"`
    Name         string `json:"name"`
    Suffix       string `json:"suffix"`
    DateLastUsed string `json:"dateLastUsed"`
}

func (c *Client) CreateCLIKey(ctx context.Context, keyName string) (*CLIKeyResponse, error) {
    url := fmt.Sprintf("%s/auth/login/cli", c.baseURL)

    reqBody := map[string]string{
        "username": c.username,
        "password": c.password,
        "profile":  keyName,
    }
    jsonBody, _ := json.Marshal(reqBody)

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("failed to create CLI key with status code: %d", resp.StatusCode)
    }

    var cliKeyResp CLIKeyResponse
    err = json.NewDecoder(resp.Body).Decode(&cliKeyResp)
    if err != nil {
        return nil, err
    }

    return &cliKeyResp, nil
}

func (c *Client) GetCLIKeys(ctx context.Context, organizationID string) ([]CLIKey, error) {
    url := fmt.Sprintf("%s/organizations/%s/cli-keys", c.baseURL, organizationID)

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
        return nil, fmt.Errorf("failed to get CLI keys with status code: %d", resp.StatusCode)
    }

    var cliKeys []CLIKey
    err = json.NewDecoder(resp.Body).Decode(&cliKeys)
    if err != nil {
        return nil, err
    }

    return cliKeys, nil
}

func (c *Client) GetUserIDFromToken() (string, error) {
    parts := strings.Split(c.accessToken, ".")
    if len(parts) != 3 {
        return "", fmt.Errorf("invalid access token format")
    }

    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return "", err
    }

    var claims map[string]interface{}
    err = json.Unmarshal(payload, &claims)
    if err != nil {
        return "", err
    }

    userID, ok := claims["id"].(string)
    if !ok {
        return "", fmt.Errorf("user ID not found in access token")
    }

    return userID, nil
}

func (c *Client) UpdateCLIKey(ctx context.Context, organizationID, keyID, newName string) (*CLIKey, error) {
    url := fmt.Sprintf("%s/organizations/%s/cli-keys/%s", c.baseURL, organizationID, keyID)

    reqBody := map[string]string{
        "name": newName,
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
        return nil, fmt.Errorf("failed to update CLI key with status code: %d", resp.StatusCode)
    }

    var updatedCLIKey CLIKey
    err = json.NewDecoder(resp.Body).Decode(&updatedCLIKey)
    if err != nil {
        return nil, err
    }

    return &updatedCLIKey, nil
}

func (c *Client) DeleteCLIKey(ctx context.Context, organizationID, keyID string) error {
    url := fmt.Sprintf("%s/organizations/%s/cli-keys/%s", c.baseURL, organizationID, keyID)

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
        return fmt.Errorf("failed to delete CLI key with status code: %d", resp.StatusCode)
    }

    return nil
}
