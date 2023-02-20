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

package exec

import (
	"fmt"

	"github.com/tricorder/src/utils/log"
)

func Run(argv []string) (string, string, error) {
	cmd := NewCommand(argv)
	err := cmd.Start()
	msg := fmt.Sprintf("command=%v stdout=%s stderr=%s error: %v", argv, cmd.Stdout(), cmd.Stderr(), err)
	log.Infof(msg)
	if err != nil {
		return "", "", fmt.Errorf("start failed, message=%s", msg)
	}
	err = cmd.Wait()
	msg = fmt.Sprintf("command=%v stdout=%s stderr=%s error: %v", argv, cmd.Stdout(), cmd.Stderr(), err)
	if err != nil {
		return "", "", fmt.Errorf("wait failed, message=%s", msg)
	}
	return cmd.Stdout(), cmd.Stderr(), nil
}
