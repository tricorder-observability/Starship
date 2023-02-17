package http

import (
	"fmt"
	"strings"
)

// API path component
const (
	API_ROOT      = "api"
	LIST_CODE     = "listCode"
	ADD_CODE      = "addCode"
	UPLOAD        = "uploadFile"
	DEPLOY        = "deploy"
	UN_DEPLOY     = "undeploy"
	DELETE_MODULE = "deleteCode"
)

// GetAPIUrl returns a http URL that corresponds to the requested path.
func GetAPIUrl(addr, root, path string) string {
	return fmt.Sprintf("http://%s", strings.Join([]string{addr, root, path}, "/"))
}
