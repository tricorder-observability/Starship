package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tricorder/src/api-server/http/api"
	"github.com/tricorder/src/utils/errors"
)

// Client provides APIs to API Server's HTTP server.
type Client struct {
	// The URL to the API Server.
	url string
}

// NewClient returns a new Client instance.
func NewClient(url string) *Client {
	return &Client{url: url}
}

// ListAgents returns the list of agents stored on the API Server.
// agentReq is the request data structure, it will be converted to JSON and sent to the API Server.
func (c *Client) ListAgents(agentReq *ListAgentReq) (*ListAgentResp, error) {
	field := "agent_id,node_name,agent_pod_id,state,create_time,last_update_time"
	if agentReq != nil && len(agentReq.Fields) > 0 {
		field = agentReq.Fields
	}

	httpClient := http.Client{Timeout: time.Duration(3) * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?fields=%s", c.url+api.LIST_AGENT_PATH, field), nil)
	if err != nil {
		return nil, errors.Wrap("listing agents", "create request", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap("listing agents", "do request", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap("listing agents", "read response body", err)
	}

	var model *ListAgentResp
	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
