// client.go
package client

import (
    "encoding/json"
    "strings"
)

// tools

type WebhookBody struct {
    DataType string      `json:"dataType"`
    Type     string      `json:"type"`
    Parts    []BodyPart  `json:"parts"`
}

type BodyPart struct {
    DataType string   `json:"dataType,omitempty"`
    Type     string   `json:"type"`
    Value    string   `json:"value,omitempty"`
    Path     []string `json:"path,omitempty"`
    Name     string   `json:"name,omitempty"`
}

func ConvertToAPIFormat(body string) (*WebhookBody, error) {
    var bodyParts []interface{}
    err := json.Unmarshal([]byte(body), &bodyParts)
    if err != nil {
        return nil, err
    }

    var apiBody WebhookBody
    apiBody.DataType = "ANY"
    apiBody.Type = "TOKENIZED"
    apiBody.Parts = make([]BodyPart, 0)

    for _, part := range bodyParts {
        partMap, ok := part.(map[string]interface{})
        if !ok {
            continue
        }

        for _, value := range partMap {
            if strings.HasPrefix(value.(string), "{{") && strings.HasSuffix(value.(string), "}}") {
                path := strings.Trim(value.(string), "{}")
                pathParts := strings.Split(path, ".")
                apiBody.Parts = append(apiBody.Parts, BodyPart{
                    Type: "OBJECT_VALUE",
                    Path: pathParts[1:],
                    Name: pathParts[1],
                })
            } else {
                escapedValue, _ := json.Marshal(value)
                apiBody.Parts = append(apiBody.Parts, BodyPart{
                    DataType: "STRING",
                    Type:     "VALUE",
                    Value:    strings.ReplaceAll(string(escapedValue), `"`, `\"`),
                })
            }
        }
    }

    return &apiBody, nil
}

func ConvertToMultilineFormat(body *WebhookBody) (string, error) {
    var multilineParts []string

    for _, part := range body.Parts {
        if part.Type == "OBJECT_VALUE" {
            multilineParts = append(multilineParts, "{{"+strings.Join(part.Path, ".")+"}}")
        } else if part.Type == "VALUE" {
            unescapedValue := strings.ReplaceAll(part.Value, `\"`, `"`)
            multilineParts = append(multilineParts, unescapedValue)
        }
    }

    multilineBody := strings.Join(multilineParts, "")
    return multilineBody, nil
}
