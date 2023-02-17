package module

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tricorder/src/cli/internal/outputs"
	http_utils "github.com/tricorder/src/utils/http"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload module WASM code file.",
	Long: `Upload module WASM code file. For example:
$ starship-cli module upload path/to/copy_input_to_output.wasm
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("apiServerAddress: %s, arg[0]:%s \n", apiAddress, args[0])
		url := http_utils.GetAPIUrl(apiAddress, http_utils.API_ROOT, http_utils.UPLOAD)
		resp := uploadWASMFIle(url, args[0])

		err := outputs.Output(output, resp)
		if err != nil {
			log.Error(err)
		}
	},
}

// upload WASM file through http client post
// return WASM file UUID
func uploadWASMFIle(url string, filePath string) []byte {
	request, err := newfileUploadRequest(url, nil, "file", filePath)
	if err != nil {
		log.Error(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		log.Error(err)
		return []byte("")
	} else {
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		return body
	}
}

// Creates a new file upload http request with optional extra params
// paramName, the parameter name of file in form
// path, the path of you need to upload
// params, optional extra params
func newfileUploadRequest(url string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
