// Copyright (C) 2023  tricorder-observability
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

package pg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	docker "github.com/tricorder/src/testing/docker"
)

const defaultPort = 5432

// NOTE: Cannot use src/testing/pg/fixture.go's LuanchContainer
// Because there will be circular dependency between src/testing/pg and src/utils/pg
func createPGTestFixutre() (*docker.Runner, *Client, error) {
	pgRunner := &docker.Runner{
		ImageName: "postgres",
		RdyMsg:    "database system is ready to accept connections",
		Options:   []string{"--env=POSTGRES_PASSWORD=passwd"},
	}
	err := pgRunner.Launch(10 * time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to start postgres server, error: %v", err)
	}

	pgGatewayIP, err := pgRunner.GetGatewayIP()
	if err != nil {
		return nil, nil, err
	}

	pgPort, err := pgRunner.GetExposedPort(defaultPort)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get postgres container's exposed port for %d, error: %v",
			defaultPort,
			err,
		)
	}
	pgURL := fmt.Sprintf("postgresql://postgres:passwd@%s:%d", strings.TrimSpace(pgGatewayIP), pgPort)

	pgClient := NewClient(pgURL)
	if err := pgClient.Connect(); err != nil {
		return nil, nil, fmt.Errorf("Unable to create client to postgres at %s, error: %v", pgURL, err)
	}
	return pgRunner, pgClient, nil
}

func TestPGClient(t *testing.T) {
	pgRunner, pgClient, err := createPGTestFixutre()
	// Always try to stop the container.
	defer func() {
		if pgClient == nil {
			log.Errorf("pgClient == nil")
			return
		}
		err := pgRunner.Stop()
		if err != nil {
			t.Logf("Failed to stop postgres server, error: %v", err)
		}
	}()
	defer func() {
		if pgClient == nil {
			log.Errorf("pgClient == nil")
			return
		}
		pgClient.Close()
	}()

	if err != nil {
		t.Fatalf("Could not create PGTestFixture, error: %v", err)
	}

	err = pgClient.CreateTable(&Schema{
		Name: "test_table",
		Columns: []Column{
			{
				Name: "id",
				Type: TEXT,
			},
		},
	})
	if err != nil {
		t.Errorf("Unable to create table in database, error: %v", err)
	}

	err = pgClient.CreateHTTPRequestTable()
	if err != nil {
		t.Errorf("Unable to create table in database, error: %v", err)
	}
	requestURL := "http://localhost:8080"
	jsonBody := []byte(`{"client_message": "hello, server!"}`)
	body := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest(http.MethodPost, requestURL, body)
	req.Header.Set("request_id", "7a89d883-9e95-4409-b77d-11b26558a00e")

	err = pgClient.WriteHTTPRequest(req)
	if err != nil {
		t.Errorf("Unable to write to database, error: %v", err)
	}
}

type Object struct {
	metav1.ObjectMeta `       json:"metadata,omitempty"`
	Key               string
}

const tableName = "tests"

// Tests that upsert() can update the value.
func TestPGUpsert(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pgRunner, pgClient, err := createPGTestFixutre()
	require.Nil(err)

	defer func() {
		assert.Nil(pgRunner.Stop())
		pgClient.Close()
	}()

	err = pgClient.CreateTable(GetJSONBTableSchema(tableName))
	assert.Nil(err)

	obj := Object{}
	obj.Name = "@#$%^&*()_+|"
	obj.UID = types.UID("uid1")
	value, _ := json.Marshal(obj)
	err = pgClient.JSON().Upsert(tableName, string(obj.UID), value)
	assert.Nil(err)

	result1 := []*Object{}
	err = pgClient.JSON().List(tableName, &result1)
	assert.Nil(err)
	// Check result upserted
	assert.Equal(1, len(result1))
	assert.Equal(obj.Name, result1[0].Name)
}

func TestPGDelete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pgRunner, pgClient, err := createPGTestFixutre()
	require.Nil(err)

	// Always try to stop the container.
	defer func() {
		pgClient.Close()
		assert.Nil(pgRunner.Stop())
	}()

	err = pgClient.CreateTable(GetJSONBTableSchema(tableName))
	assert.Nil(err)

	obj := Object{}
	obj.Name = "obj1"
	obj.UID = types.UID("uid1")
	value, _ := json.Marshal(obj)
	err = pgClient.JSON().Upsert(tableName, string(obj.UID), value)
	assert.Nil(err)

	result1 := []*Object{}
	err = pgClient.JSON().List(tableName, &result1)
	assert.Nil(err)
	assert.Equal(1, len(result1))

	err = pgClient.JSON().Delete(tableName, string(obj.UID))
	assert.Nil(err)

	result2 := []*Object{}
	err = pgClient.JSON().List(tableName, &result2)
	assert.Nil(err)
	// Check result is deleted
	assert.Equal(0, len(result2))
}

func TestPGClean(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pgRunner, pgClient, err := createPGTestFixutre()
	require.Nil(err)

	// Always try to stop the container.
	defer func() {
		pgClient.Close()
		assert.Nil(pgRunner.Stop())
	}()

	err = pgClient.CreateTable(GetJSONBTableSchema(tableName))
	assert.Nil(err)

	obj := Object{}
	obj.Name = "obj1"
	obj.UID = types.UID("uid1")
	value, _ := json.Marshal(obj)
	err = pgClient.JSON().Upsert(tableName, string(obj.UID), value)
	assert.Nil(err)

	if err = pgClient.Clean(tableName); err != nil {
		t.Error(err)
	}

	result := []*Object{}
	err = pgClient.JSON().List(tableName, &result)
	assert.Nil(err)
	// Check result clean
	assert.Equal(0, len(result))
}

// Tests that TestJSONGet can get a JSON object from the data base.
func TestJSONGet(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	pgRunner, pgClient, err := createPGTestFixutre()
	require.Nil(err)
	// Always try to stop the container.
	defer func() {
		pgClient.Close()
		assert.Nil(pgRunner.Stop())
	}()

	tableName := "test1"
	err = pgClient.CreateTable(GetJSONBTableSchema(tableName))
	assert.Nil(err)
	// Insert the first object
	obj1 := Object{}
	obj1.Name = "obj1"
	obj1.UID = types.UID("uid1")
	value1, _ := json.Marshal(obj1)
	err = pgClient.JSON().Upsert(tableName, string(obj1.UID), value1)
	assert.Nil(err)

	target := Object{}
	err = pgClient.JSON().Get(tableName, &target, fmt.Sprintf("WHERE data#>>'{metadata,uid}'='%s'", string(obj1.UID)))
	assert.Nil(err)
	assert.Equal(obj1.Name, target.Name)
}

func TestPGListObjects(t *testing.T) {
	assert := assert.New(t)
	pgRunner, pgClient, err := createPGTestFixutre()
	assert.Nil(err)
	// Always try to stop the container.
	defer func() {
		pgClient.Close()
		assert.Nil(pgRunner.Stop())
	}()

	tableName := "test1"
	err = pgClient.CreateTable(GetJSONBTableSchema(tableName))
	assert.Nil(err)
	// insert the first object
	obj1 := Object{}
	obj1.Name = "obj1"
	obj1.UID = types.UID("uid1")
	value1, _ := json.Marshal(obj1)
	err = pgClient.JSON().Upsert(tableName, string(obj1.UID), value1)
	assert.Nil(err)

	result1 := []*Object{}
	err = pgClient.JSON().List(tableName, &result1)
	assert.Nil(err)
	assert.Equal(1, len(result1))
	assert.Equal("obj1", result1[0].Name)

	// insert another object
	obj2 := Object{}
	obj2.Name = "obj2"
	obj2.UID = types.UID("uid2")
	value2, _ := json.Marshal(obj2)
	err = pgClient.JSON().Upsert(tableName, string(obj2.UID), value2)
	assert.Nil(err)

	result2 := []*Object{}
	err = pgClient.JSON().List(tableName, &result2)
	assert.Nil(err)
	assert.Equal(2, len(result2))
}

// Tests that WriteRecord can write a text record into the data base.
func TestWriteRecord(t *testing.T) {
	log.SetReportCaller(true)
	assert := assert.New(t)

	pgRunner, pgClient, err := createPGTestFixutre()
	assert.Nil(err)

	defer func() {
		assert.Nil(pgRunner.Stop())
		pgClient.Close()
	}()

	schema := &Schema{
		Name: "test_table",
		Columns: []Column{
			{
				Name: "id",
				Type: TEXT,
			},
		},
	}
	assert.Nil(pgClient.CreateTable(schema))
	assert.Nil(pgClient.WriteRecord([]interface{}{"1234"}, schema))
	records, err := pgClient.Query("select * from test_table")
	assert.Nil(err)
	assert.Equal([][]interface{}{{"1234"}}, records)
}

// Tests that WriteRecord return error when input value count and schema column count are not equal.
func TestWriteRecordFailUnequalCount(t *testing.T) {
	assert := assert.New(t)
	pgClient := NewClient("wont-connect")
	schema := &Schema{
		Name: "test_table",
		Columns: []Column{
			{
				Name: "test",
				Type: 100,
			},
		},
	}
	err := pgClient.WriteRecord([]interface{}{}, schema)
	assert.ErrorContains(err, "field count differs from the schema's column count")
}
