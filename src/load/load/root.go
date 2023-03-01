// Copyright (C) 2023 Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
