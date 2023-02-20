// Copyright (C) 2023  Tricorder Observability
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

package yaml

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/tricorder/src/cli/internal/model"
)

func Output(data *model.Response) error {
	bytes, e := yaml.Marshal(data)
	if e != nil {
		return e
	}
	_, e = fmt.Printf("%v", string(bytes))
	return e
}
