package http

import (
	"bytes"
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

// CreateModule creates a new module on the API Server.
// moduleReq is the request data structure, it will be converted to JSON and sent to the API Server.
func (c *Client) CreateModule(moduleReq *CreateModuleReq) (*CreateModuleResp, error) {
	bodyBytes, err := json.Marshal(moduleReq)
	if err != nil {
		return nil, errors.Wrap("creating module", "encode req body", err)
	}
	httpClient := http.Client{Timeout: time.Duration(3) * time.Second}
	req, err := http.NewRequest("POST", c.url+api.CREATE_MODULE_PATH, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, errors.Wrap("creating module", "create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap("creating module", "do request", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap("creating module", "read response body", err)
	}

	var model *CreateModuleResp
	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// DeployModule deploys a module on the API Server.
// moduleId is the ID of the module to be deployed.
func (c *Client) DeployModule(moduleId string) (*DeployModuleResp, error) {
	httpClient := http.Client{Timeout: time.Duration(3) * time.Second}
	url := fmt.Sprintf("%s?id=%s", c.url+api.DEPLOY_MODULE_PATH, moduleId)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, errors.Wrap("deploying module", "create request", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap("deploying module", "do request", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap("deploying module", "read response body", err)
	}

	var model *DeployModuleResp
	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// UndeployModule undeploys a module on the API Server.
// moduleId is the ID of the module to be undeployed.
func (c *Client) UndeployModule(moduleId string) (*UndeployModuleResp, error) {
	httpClient := http.Client{Timeout: time.Duration(3) * time.Second}
	url := fmt.Sprintf("%s?id=%s", c.url+api.UNDEPLOY_MODULE_PATH, moduleId)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, errors.Wrap("undeploying module", "create request", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap("undeploying module", "do request", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap("undeploying module", "read response body", err)
	}

	var model *UndeployModuleResp
	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// DeleteModule deletes a module on the API Server.
// moduleId is the ID of the module to be deleted.
func (c *Client) DeleteModule(moduleId string) (*DeleteModuleResp, error) {
	httpClient := http.Client{Timeout: time.Duration(3) * time.Second}
	url := fmt.Sprintf("%s?id=%s", c.url+api.DELETE_MODULE_PATH, moduleId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap("deleting module", "create request", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap("deleting module", "do request", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap("deleting module", "read response body", err)
	}

	var model *DeleteModuleResp
	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// ListModules lists all modules on the API Server.
// moduleReq is the request data structure, it will be converted to JSON and sent to the API Server.
func (c *Client) ListModules(moduleReq *ListModuleReq) (*ListModuleResp, error) {
	field := "id,name,desire_state,create_time," +
		"ebpf_fmt,ebpf_lang,schema_name,fn,schema_attr"
	if moduleReq != nil && len(moduleReq.Fields) > 0 {
		field = moduleReq.Fields
	}

	httpClient := http.Client{Timeout: time.Duration(3) * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?fields=%s", c.url+api.LIST_MODULE_PATH, field), nil)
	if err != nil {
		return nil, errors.Wrap("listing modules", "create request", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap("listing modules", "do request", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap("listing modules", "read response body", err)
	}

	var model *ListModuleResp
	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
