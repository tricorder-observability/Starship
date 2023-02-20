package load

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "starship-load",
	Short: "CLI tool to load pre-built eBPF+WASM modules into SQLite DB file",
	Long: `CLI tool to load pre-built eBPF+WASM modules into SQLite DB file.
	Requires to specify the paths to BCC source file, WASM compiled binary object file, 
	and module description in JSON format. 
	The content of these files are written to the specified SQLite db file. 
	We use this tool to load pre-built eBPF+WASM modules into Starship's official release.`,
}

// Execute load command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loadCmd)

	rootCmd.SuggestionsMinimumDistance = 1
}
