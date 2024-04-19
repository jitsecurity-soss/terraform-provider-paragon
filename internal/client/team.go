// client.go
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)


type Team struct {
    ID             string                 `json:"id"`
    DateCreated    string                 `json:"dateCreated"`
    DateUpdated    string                 `json:"dateUpdated"`
    Name           string                 `json:"name"`
    Website        string                 `json:"website"`
    OrganizationID string                 `json:"organizationId"`
    Organization   Organization           `json:"organization"`
}

func (c *Client) GetTeams(ctx context.Context) ([]Team, error) {
    url := fmt.Sprintf("%s/teams", c.baseURL)

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
        return nil, fmt.Errorf("failed to get teams with status code: %d", resp.StatusCode)
    }

    var teams []Team
    err = json.NewDecoder(resp.Body).Decode(&teams)
    if err != nil {
        return nil, err
    }

    return teams, nil
}

func (c *Client) GetTeamByID(ctx context.Context, teamID string) (*Team, error) {
    url := fmt.Sprintf("%s/teams/%s", c.baseURL, teamID)

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
        return nil, fmt.Errorf("failed to get team with status code: %d", resp.StatusCode)
    }

    var team Team
    err = json.NewDecoder(resp.Body).Decode(&team)
    if err != nil {
        return nil, err
    }

    return &team, nil
}
