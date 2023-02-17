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
