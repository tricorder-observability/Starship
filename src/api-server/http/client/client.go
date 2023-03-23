package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	apiserver "github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/api-server/http/api"
	"github.com/tricorder/src/utils/errors"
)

// Client provides APIs to API Server's HTTP server.
type Client struct {
	// The URL to the API Server.
	url string

	listModuleURL     string
	listAgentURL      string
	createModuleURL   string
	deployModuleURL   string
	undeployModuleURL string
	deleteModuleURL   string
}

// NewClient returns a new Client instance.
func NewClient(url string) *Client {
	return &Client{
		url: url,

		listModuleURL:     api.GetURL(url, api.LIST_MODULE_PATH),
		listAgentURL:      api.GetURL(url, api.LIST_AGENT_PATH),
		createModuleURL:   api.GetURL(url, api.CREATE_MODULE_PATH),
		deployModuleURL:   api.GetURL(url, api.DEPLOY_MODULE_PATH),
		undeployModuleURL: api.GetURL(url, api.UNDEPLOY_MODULE_PATH),
		deleteModuleURL:   api.GetURL(url, api.DELETE_MODULE_PATH),
	}
}

func executeHTTPReq(req *http.Request, resp any) error {
	httpClient := http.Client{Timeout: time.Duration(3) * time.Second}
	httpResp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap("execute http request", "do request", err)
	}

	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap("execute http request", "read response body", err)
	}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return errors.Wrap("listing agents", "decode response body", err)
	}
	return nil
}

// ListAgents returns the list of agents stored on the API Server.
// agentReq is the request data structure, it will be converted to JSON and sent to the API Server.
func (c *Client) ListAgents(agentReq *apiserver.ListAgentReq) (*apiserver.ListAgentResp, error) {
	field := "agent_id,node_name,agent_pod_id,state,create_time,last_update_time"
	if agentReq != nil && len(agentReq.Fields) > 0 {
		field = agentReq.Fields
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?fields=%s", c.listAgentURL, field), nil)
	if err != nil {
		return nil, errors.Wrap("listing agents", "create request", err)
	}

	resp := &apiserver.ListAgentResp{}
	err = executeHTTPReq(req, resp)
	if err != nil {
		return nil, errors.Wrap("listing agents", "execute http request", err)
	}

	return resp, nil
}

// CreateModule creates a new module on the API Server.
// moduleReq is the request data structure, it will be converted to JSON and sent to the API Server.
func (c *Client) CreateModule(moduleReq *apiserver.CreateModuleReq) (*apiserver.CreateModuleResp, error) {
	bodyBytes, err := json.Marshal(moduleReq)
	if err != nil {
		return nil, errors.Wrap("creating module", "encode req body", err)
	}

	req, err := http.NewRequest("POST", c.createModuleURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, errors.Wrap("creating module", "create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp := &apiserver.CreateModuleResp{}
	err = executeHTTPReq(req, resp)
	if err != nil {
		return nil, errors.Wrap("creating module", "execute http request", err)
	}

	return resp, nil
}

// DeployModule deploys a module on the API Server.
// moduleId is the ID of the module to be deployed.
func (c *Client) DeployModule(moduleId string) (*apiserver.DeployModuleResp, error) {
	url := fmt.Sprintf("%s?id=%s", c.deployModuleURL, moduleId)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, errors.Wrap("deploying module", "create request", err)
	}

	resp := &apiserver.DeployModuleResp{}
	err = executeHTTPReq(req, resp)
	if err != nil {
		return nil, errors.Wrap("deploying module", "execute http request", err)
	}

	return resp, nil
}

// UndeployModule undeploys a module on the API Server.
// moduleId is the ID of the module to be undeployed.
func (c *Client) UndeployModule(moduleId string) (*apiserver.UndeployModuleResp, error) {
	url := fmt.Sprintf("%s?id=%s", c.undeployModuleURL, moduleId)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, errors.Wrap("undeploying module", "create request", err)
	}

	resp := &apiserver.UndeployModuleResp{}
	err = executeHTTPReq(req, resp)
	if err != nil {
		return nil, errors.Wrap("undeploying module", "execute http request", err)
	}

	return resp, nil
}

// DeleteModule deletes a module on the API Server.
// moduleId is the ID of the module to be deleted.
func (c *Client) DeleteModule(moduleId string) (*apiserver.DeleteModuleResp, error) {
	url := fmt.Sprintf("%s?id=%s", c.deleteModuleURL, moduleId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap("deleting module", "create request", err)
	}

	resp := &apiserver.DeleteModuleResp{}
	err = executeHTTPReq(req, resp)
	if err != nil {
		return nil, errors.Wrap("deleting module", "execute http request", err)
	}

	return resp, nil
}

// ListModules lists all modules on the API Server.
// moduleReq is the request data structure, it will be converted to JSON and sent to the API Server.
func (c *Client) ListModules(moduleReq *apiserver.ListModuleReq) (*apiserver.ListModuleResp, error) {
	field := "id,name,desire_state,create_time," +
		"ebpf_fmt,ebpf_lang,schema_name,fn,schema_attr"
	if moduleReq != nil && len(moduleReq.Fields) > 0 {
		field = moduleReq.Fields
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?fields=%s", c.listModuleURL, field), nil)
	if err != nil {
		return nil, errors.Wrap("listing modules", "create request", err)
	}

	resp := &apiserver.ListModuleResp{}
	err = executeHTTPReq(req, resp)
	if err != nil {
		return nil, errors.Wrap("listing module", "execute http request", err)
	}

	return resp, nil
}
