package module

import (
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tricorder/src/cli/internal/outputs"
	http_utils "github.com/tricorder/src/utils/http"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Query modules.",
	Long: `Query modules. For example:
1. Query all modules:
$ starship-cli module list
`,
	Run: func(cmd *cobra.Command, args []string) {
		url := http_utils.GetAPIUrl(apiAddress, http_utils.API_ROOT, http_utils.LIST_CODE)
		resp, err := listModules(url)
		if err != nil {
			log.Error(err)
		}

		err = outputs.Output(output, resp)
		if err != nil {
			log.Error(err)
		}
	},
}

func listModules(url string) ([]byte, error) {
	c := http.Client{Timeout: time.Duration(3) * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?fields=%s", url, "id,name,status,create_time,"+
		"ebpf_fmt,ebpf_lang,schema_name,fn,schema_attr"), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return body, nil
}
