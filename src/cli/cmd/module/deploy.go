package module

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	outputs "github.com/tricorder/src/cli/internal/outputs"
	http_utils "github.com/tricorder/src/utils/http"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy Starship module to the platform.",
	Long: "Deploy a previously-created module to the Starship platform.\n" +
		"The deployed module, if succeeded, will be executed by Tricorder.\n" +
		"For example:\n" +
		"    starship-cli module deploy --id=ce8a4fbe_45db_49bb_9568_6688dd84480b",
	Run: func(cmd *cobra.Command, args []string) {
		url := http_utils.GetAPIUrl(apiAddress, http_utils.API_ROOT, http_utils.DEPLOY)
		resp, err := deployModule(url, moduleId)
		if err != nil {
			log.Fatalf("Failed to deploy module, id='%s', error: %v", moduleId, err)
		}
		if len(resp) == 0 {
			log.Fatalf("Failed to deploy module, id='%s', Empty response from API Server", moduleId)
		}
		if err := outputs.Output(output, resp); err != nil {
			log.Fatalf("Failed to write output, error: %v", err)
		}
	},
}

func init() {
	deployCmd.Flags().StringVarP(&moduleId, "id", "i", moduleId, "the id of module.")
	_ = deployCmd.MarkFlagRequired("id")
}

func deployModule(url string, moduleId string) ([]byte, error) {
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
