package module

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/tricorder/src/cli/internal/outputs"
	http_utils "github.com/tricorder/src/utils/http"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete module from staship",
	Long: `Delete module from staship. For example:

$ starship-cli module delete --id 2a339411_7dd8_46ba_9581_e9d41286b564
`,
	Run: func(cmd *cobra.Command, args []string) {
		url := http_utils.GetAPIUrl(apiAddress, http_utils.API_ROOT, http_utils.DELETE_MODULE)
		resp, err := deleteModule(url, moduleId)
		if err != nil {
			log.Error(err)
		}

		err = outputs.Output(output, resp)
		if err != nil {
			log.Error(err)
		}
	},
}

func deleteModule(url string, moduleId string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s?id=%s", url, moduleId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
