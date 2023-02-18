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

var undeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy module",
	Long: `Undeploy module command. For example:

$ starship-cli module undeploy --id ce8a4fbe_45db_49bb_9568_6688dd84480b
`,
	Run: func(cmd *cobra.Command, args []string) {
		url := http_utils.GetAPIUrl(apiAddress, http_utils.API_ROOT, http_utils.UN_DEPLOY)
		resp, err := undeployModule(url, moduleId)
		if err != nil {
			log.Error(err)
		}

		err = outputs.Output(output, resp)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	undeployCmd.Flags().StringVarP(&moduleId, "id", "i", moduleId, "the id of module.")
	_ = undeployCmd.MarkFlagRequired("id")
}

func undeployModule(url string, moduleId string) ([]byte, error) {
	resp, err := http.Post(fmt.Sprintf("%s?id=%s", url, moduleId), "application/json", nil)
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
