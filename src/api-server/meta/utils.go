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

package meta

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	log "github.com/sirupsen/logrus"

	"github.com/tricorder/src/utils/pg"
)

func upsert(pgClient *pg.Client, table string, obj runtime.Object, uid types.UID) {
	value, _ := json.Marshal(obj)
	if err := pgClient.JSON().Upsert(table, string(uid), value); err != nil {
		log.Errorf("Upsert error %s", err)
	}
}

func deleteByID(pgClient *pg.Client, table string, uid types.UID) {
	if err := pgClient.JSON().Delete(table, string(uid)); err != nil {
		log.Errorf("DeleteData %s error %s", table, uid)
	}
}
