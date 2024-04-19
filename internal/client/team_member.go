// client.go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

type TeamMember struct {
    ID             string      `json:"id"`
    Name           string      `json:"name"`
    Email          string      `json:"email"`
    UserID         string      `json:"userId"`
    Role           string      `json:"role"`
    OrganizationID string      `json:"organizationId"`
}

type TeamInvite struct {
    ID           string `json:"id"`
    DateCreated  string `json:"dateCreated"`
    DateUpdated  string `json:"dateUpdated"`
    Status       string `json:"status"`
    Role         string `json:"role"`
    Email        string `json:"email"`
    Team         Team   `json:"team"`
}

type InviteTeamMemberRequest struct {
    Role   string   `json:"role"`
    Emails []string `json:"emails"`
}

func (c *Client) GetTeamMembers(ctx context.Context, teamID string) ([]TeamMember, error) {
    url := fmt.Sprintf("%s/teams/%s/members", c.baseURL, teamID)

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
        return nil, fmt.Errorf("failed to get team members with status code: %d", resp.StatusCode)
    }

    var members []TeamMember
    err = json.NewDecoder(resp.Body).Decode(&members)
    if err != nil {
        return nil, err
    }

    return members, nil
}

func (c *Client) GetTeamInvites(ctx context.Context, teamID string) ([]TeamInvite, error) {
    url := fmt.Sprintf("%s/teams/%s/invite", c.baseURL, teamID)

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
        return nil, fmt.Errorf("failed to get team invites with status code: %d", resp.StatusCode)
    }

    var invites []TeamInvite
    err = json.NewDecoder(resp.Body).Decode(&invites)
    if err != nil {
        return nil, err
    }

    return invites, nil
}

func (c *Client) InviteTeamMember(ctx context.Context, teamID, role, email string) ([]TeamInvite, error) {
    url := fmt.Sprintf("%s/teams/%s/invite", c.baseURL, teamID)

    reqBody := InviteTeamMemberRequest{
        Role:   role,
        Emails: []string{email},
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
        return nil, fmt.Errorf("failed to invite team member with status code: %d", resp.StatusCode)
    }

    var invites []TeamInvite
    err = json.NewDecoder(resp.Body).Decode(&invites)
    if err != nil {
        return nil, err
    }

    return invites, nil
}

func (c *Client) UpdateTeamMemberRole(ctx context.Context, teamID, memberID, role string) (*TeamMember, error) {
    url := fmt.Sprintf("%s/teams/%s/members/%s", c.baseURL, teamID, memberID)

    reqBody := map[string]string{
        "role": role,
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
        return nil, fmt.Errorf("failed to update team member role with status code: %d", resp.StatusCode)
    }

    var updatedMember TeamMember
    err = json.NewDecoder(resp.Body).Decode(&updatedMember)
    if err != nil {
        return nil, err
    }

    return &updatedMember, nil
}

func (c *Client) DeleteTeamMember(ctx context.Context, teamID, memberID string) error {
    url := fmt.Sprintf("%s/teams/%s/members/%s", c.baseURL, teamID, memberID)
    tflog.Debug(ctx, fmt.Sprintf("url to delete: %s", url))

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
        return fmt.Errorf("status code: %d", http.StatusNotFound)
    }

    // Weirdly enough - if team member is not found, the status code is 403 with this body:
    // {
    //     "message": "Unable to find team member.",
    //     "code": "13200",
    //     "status": 403,
    //     "meta": {
    //         "teamMemberId": "<member_id>"
    //     }
    // }
    if resp.StatusCode == http.StatusForbidden {
        var errResp struct {
            Message string `json:"message"`
            Code    string `json:"code"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
            return fmt.Errorf("failed to decode error response: %v", err)
        }

        if errResp.Message == "Unable to find team member." || errResp.Code == "13200" {
            return fmt.Errorf("status code: %d", http.StatusNotFound)
        }
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to delete team member with status code: %d", resp.StatusCode)
    }

    return nil
}

func (c *Client) DeleteTeamInvite(ctx context.Context, teamID, inviteID string) error {
    url := fmt.Sprintf("%s/teams/%s/invite/%s", c.baseURL, teamID, inviteID)

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
        return fmt.Errorf("status code: %d", http.StatusNotFound)
    }

    // Weirdly enough - if invite is not found, the status code is 403 with this body:
//     {
//     "message": "Unable to find invite.",
//     "code": "13101",
//     "status": 403,
//     "meta": {
//         "inviteId": "<invite_id>>"
//     }
// }
    if resp.StatusCode == http.StatusForbidden {
        var errResp struct {
            Message string `json:"message"`
            Code    string `json:"code"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
            return fmt.Errorf("failed to decode error response: %v", err)
        }

        if errResp.Message == "Unable to find invite." || errResp.Code == "13101" {
            return fmt.Errorf("status code: %d", http.StatusNotFound)
        }
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to delete team invite with status code: %d", resp.StatusCode)
    }

    return nil
}
