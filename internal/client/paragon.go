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
    "github.com/hashicorp/terraform-plugin-log/tflog"
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

type CreateProjectRequest struct {
    ProjectTitle   string `json:"projectTitle"`
    OrganizationID string `json:"organizationId"`
    Name           string `json:"name"`
}

type CreateProjectResponse struct {
    ID             string   `json:"id"`
    DateCreated    string   `json:"dateCreated"`
    DateUpdated    string   `json:"dateUpdated"`
    Name           string   `json:"name"`
    Website        string   `json:"website"`
    OrganizationID string   `json:"organizationId"`
    Projects       []Project `json:"projects"`
}

type Project struct {
    ID                string      `json:"id"`
    Title             string      `json:"title"`
    OwnerID           string      `json:"ownerId"`
    TeamID            string      `json:"teamId"`
    IsConnectProject  bool        `json:"isConnectProject"`
    IsHidden          bool        `json:"isHidden"`
    DateCreated       string      `json:"dateCreated"`
    DateUpdated       string      `json:"dateUpdated"`
}

func (c *Client) CreateProject(ctx context.Context, organizationID, projectName string) (*Project, error) {
    url := fmt.Sprintf("%s/teams?organizationId=%s", c.baseURL, organizationID)

    reqBody := CreateProjectRequest{
        ProjectTitle:   projectName,
        OrganizationID: organizationID,
        Name:           projectName,
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
        return nil, fmt.Errorf("failed to create project with status code: %d", resp.StatusCode)
    }

    var createProjectResp CreateProjectResponse
    err = json.NewDecoder(resp.Body).Decode(&createProjectResp)
    if err != nil {
        return nil, err
    }

    var connectProject *Project
    for _, project := range createProjectResp.Projects {
        if project.IsConnectProject {
            connectProject = &project
            break
        }
    }

    if connectProject == nil {
        return nil, fmt.Errorf("connect project not found in the response")
    }

    return connectProject, nil
}

func (c *Client) GetProjectByID(ctx context.Context, projectID, teamID string) (*Project, error) {
    url := fmt.Sprintf("%s/projects/%s?teamId=%s", c.baseURL, projectID, teamID)

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
        return nil, fmt.Errorf("failed to get project with status code: %d", resp.StatusCode)
    }

    var project Project
    err = json.NewDecoder(resp.Body).Decode(&project)
    if err != nil {
        return nil, err
    }

    return &project, nil
}

type UpdateProjectTitleRequest struct {
    Title string `json:"title"`
}

func (c *Client) UpdateProjectTitle(ctx context.Context, projectID, teamID, newTitle string) (*Project, error) {
    url := fmt.Sprintf("%s/projects/%s?teamId=%s", c.baseURL, projectID, teamID)

    reqBody := UpdateProjectTitleRequest{
        Title: newTitle,
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
        return nil, fmt.Errorf("failed to update project title with status code: %d", resp.StatusCode)
    }

    var updatedProject Project
    err = json.NewDecoder(resp.Body).Decode(&updatedProject)
    if err != nil {
        return nil, err
    }

    return &updatedProject, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectID, teamID string) error {
    url := fmt.Sprintf("%s/projects/%s?teamId=%s", c.baseURL, projectID, teamID)

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
        return fmt.Errorf("failed to delete project with status code: %d", resp.StatusCode)
    }

    return nil
}


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
    tflog.Debug(ctx, fmt.Sprintf("url to update: %s", url))

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