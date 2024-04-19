// client.go
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)


type Organization struct {
    ID                    string `json:"id"`
    DateCreated           string `json:"dateCreated"`
    DateUpdated           string `json:"dateUpdated"`
    Name                  string `json:"name"`
    Website               string `json:"website"`
    Type                  string `json:"type"`
    Purpose               string `json:"purpose"`
    Referral              string `json:"referral"`
    Size                  string `json:"size"`
    Role                  string `json:"role"`
    CompletedQualification bool   `json:"completedQualification"`
    FeatureFlagMeta       map[string]interface{} `json:"featureFlagMeta"`
}

func (c *Client) GetOrganizations(ctx context.Context) ([]Organization, error) {
    url := fmt.Sprintf("%s/organizations", c.baseURL)

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
        return nil, fmt.Errorf("failed to get organizations with status code: %d", resp.StatusCode)
    }

    var organizations []Organization
    err = json.NewDecoder(resp.Body).Decode(&organizations)
    if err != nil {
        return nil, err
    }

    return organizations, nil
}
