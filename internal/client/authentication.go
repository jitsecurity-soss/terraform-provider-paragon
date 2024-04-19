// client.go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

type Client struct {
    baseURL     string
    httpClient  *http.Client
    accessToken string
    username    string
    password    string
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL:    baseURL,
        httpClient: &http.Client{},
    }
}

type AuthResponse struct {
    AccessToken string `json:"accessToken"`
}

func (c *Client) Authenticate(ctx context.Context, username, password string) error {
    url := fmt.Sprintf("%s/auth/login/email", c.baseURL)
    body := map[string]string{
        "username": username,
        "password": password,
    }
    jsonBody, _ := json.Marshal(body)

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
    }

    var authResp AuthResponse
    err = json.NewDecoder(resp.Body).Decode(&authResp)
    if err != nil {
        return err
    }

    c.accessToken = authResp.AccessToken
    c.username = username
    c.password = password
    return nil
}
