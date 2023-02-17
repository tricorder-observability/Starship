package model

// Response represents the api server response model.
// In order to parse the API server's response and facilitate later formatting of the output
type Response struct {
	Data    []map[string]interface{} `json:"data"`
	Code    string                   `json:"code"`
	Message string                   `json:"message"`
}
