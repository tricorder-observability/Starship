package table

import (
	"encoding/json"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/tricorder/src/cli/internal/model"
)

func Output(resp *model.Response) error {
	var stringMapArrays []map[string]string

	bytes, _ := json.Marshal(resp.Data)
	_ = json.Unmarshal(bytes, &stringMapArrays)

	if len(stringMapArrays) < 1 {
		return nil
	}

	var header []string

	for k := range stringMapArrays[0] {
		header = append(header, k)
	}

	var dataT [][]string

	for _, objMap := range stringMapArrays {
		var datum []string
		for _, key := range header {
			datum = append(datum, objMap[key])
		}
		dataT = append(dataT, datum)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(dataT)
	table.Render()

	return nil
}
