package client

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/api-server/http/fake"
	pb "github.com/tricorder/src/api-server/pb"
	common "github.com/tricorder/src/pb/module/common"
	"github.com/tricorder/src/pb/module/ebpf"
	"github.com/tricorder/src/pb/module/wasm"
	bazelutils "github.com/tricorder/src/testing/bazel"
	testutils "github.com/tricorder/src/testing/bazel"
	grafanatest "github.com/tricorder/src/testing/grafana"
	pgclienttest "github.com/tricorder/src/testing/pg"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/uuid"
)

var ebpfJson = `
#include <linux/ptrace.h>

BPF_PERF_OUTPUT(events);

// Writes a fixed JSON string to perf buffer.
int sample_json(struct bpf_perf_event_data *ctx) {
  const char word[] = "{\"name\":\"John\", \"age\":30}";
  events.perf_submit(ctx, (void *)word, sizeof(word));
  return 0;
}
`
var wasmJson = `
#include "cJSON.h"
#include "io.h"
#include <assert.h>
#include <stdint.h>
#include <string.h>

struct detectionPackets {
  unsigned long long nb_ddos_packets;
} __attribute__((packed));

static_assert(sizeof(struct detectionPackets) == 8,
			  "Size of detectionPackets is not 8");

// A simple function to copy entire input buf to output buffer.
// Return 0 if succeeded.
// Return 1 if failed to malloc output buffer.
int write_events_to_output() {
  struct detectionPackets *detection_packet = get_input_buf();

  cJSON *root = cJSON_CreateObject();

  cJSON_AddNumberToObject(root, "nb_ddos_packets",
						  detection_packet->nb_ddos_packets);

  char *json = NULL;
  json = cJSON_Print(root);
  cJSON_Delete(root);

  int json_size = strlen(json);
  void *buf = malloc_output_buf(json_size);
  if (buf == NULL) {
	return 1;
  }
  copy_to_output(json, json_size);
  // Free allocated memory from JSON_print().
  free(json);
  return 0;
}

int main() { return 0; }
`

func initWasiSDK() (string, string, string, error) {
	wasiTarBazelFilePath := "external/download_wasi_sdk_from_github_url/file/wasi-sdk.tar.gz"
	wasiSDKTarPath := bazelutils.TestFilePath(wasiTarBazelFilePath)
	wasiSDKPath := bazelutils.CreateTmpDir()
	wasiBuildTmpPath := bazelutils.CreateTmpDir()

	// decompress wasi-sdk.tar.gz toolchain to bazel runtime
	cmd := exec.Command("tar", "-p", "-C", wasiSDKPath, "-zxvf", wasiSDKTarPath, "--no-same-owner")
	_, err := cmd.Output()
	if err != nil {
		return "", "", "", err
	}

	wasiSDKPath += "/wasi-sdk-19.0"

	wasiBazelIncludeFilePath := "modules/common"
	wasmStarshipIncudePath := bazelutils.TestFilePath(wasiBazelIncludeFilePath)
	return wasiSDKPath, wasmStarshipIncudePath, wasiBuildTmpPath, nil
}

func TestListAgent(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testDbFilePath := testutils.GetTmpFile()
	// We'll not cleanup the temp file, as it's troublesome to turn down the http server, and probably not worth it in a
	// test.

	sqliteClient, _ := dao.InitSqlite(testDbFilePath)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	gLock := lock.NewLock()
	waitCond := cond.NewCond()

	wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath, err := initWasiSDK()
	require.Nil(err)

	fakeServer := fake.StartFakeNewServer(sqliteClient, gLock, waitCond, pgClient,
		grafanaURL, wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath)

	// test list agent
	client := NewClient("http://" + fakeServer.String())
	res, err := client.ListAgents(nil)
	require.NoError(err)
	assert.Equal(200, res.Code)
	assert.Equal(0, len(res.Data))

	nodeAgentDao := dao.NodeAgentDao{
		Client: sqliteClient,
	}
	newAgent := dao.NodeAgentGORM{
		AgentID:    "agent_test_id",
		NodeName:   "agent_test_node",
		AgentPodID: "agent_test_pod_id",
	}
	err = nodeAgentDao.SaveAgent(&newAgent)
	require.NoError(err)

	res, err = client.ListAgents(nil)
	require.NoError(err)
	assert.Equal(200, res.Code)
	assert.Equal(1, len(res.Data))
	assert.Equal("agent_test_id", res.Data[0].AgentID)
}

func TestCreateModule(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testDbFilePath := testutils.GetTmpFile()
	// We'll not cleanup the temp file, as it's troublesome to turn down the http server, and probably not worth it in a
	// test.

	sqliteClient, _ := dao.InitSqlite(testDbFilePath)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	gLock := lock.NewLock()
	waitCond := cond.NewCond()

	wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath, err := initWasiSDK()
	require.Nil(err)

	fakeServer := fake.StartFakeNewServer(sqliteClient, gLock, waitCond, pgClient,
		grafanaURL, wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath)

	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}

	// before create module, list module
	module, err := moduleDao.ListModule([]string{})
	require.NoError(err)
	assert.Equal(0, len(module))

	client := NewClient("http://" + fakeServer.String())

	// test create module
	moduleReq := &http.CreateModuleReq{
		Name: "test_module",
		Wasm: &wasm.Program{
			Code:   []byte("test_code"),
			FnName: "test_fn",
			Fmt:    common.Format_BINARY,
			OutputSchema: &common.Schema{
				Fields: []*common.DataField{
					{
						Name: "test_field",
						Type: common.DataField_JSONB,
					},
				},
			},
		},
		Ebpf: &ebpf.Program{
			Code:           "test_code",
			PerfBufferName: "test_perf_buffer_name",
			Probes: []*ebpf.ProbeSpec{
				{
					Target: "test_target",
					Entry:  "test_entry",
					Return: "test_return",
				},
			},
		},
	}

	res, err := client.CreateModule(moduleReq)
	require.NoError(err)
	assert.Equal(200, res.Code)

	// after create module, list module
	module, err = moduleDao.ListModule([]string{})
	require.NoError(err)
	assert.Equal(1, len(module))
	assert.Equal("test_module", module[0].Name)

	err = moduleDao.DeleteByID(module[0].ID)
	assert.NoError(err)

	// test create module with wasm code,
	// it will compile wasm code to wasm binary and save it to db
	moduleReq = &http.CreateModuleReq{
		Name: "test_module_wasm",
		Wasm: &wasm.Program{
			Code:   []byte(wasmJson),
			FnName: "write_events_to_output",
			Fmt:    common.Format_TEXT,
			OutputSchema: &common.Schema{
				Fields: []*common.DataField{
					{
						Name: "test_field",
						Type: common.DataField_JSONB,
					},
				},
			},
		},
		Ebpf: &ebpf.Program{
			Code:           "test_code",
			PerfBufferName: "test_perf_buffer_name",
			Probes: []*ebpf.ProbeSpec{
				{
					Target: "test_target",
					Entry:  "test_entry",
					Return: "test_return",
				},
			},
		},
	}

	res, err = client.CreateModule(moduleReq)
	require.NoError(err)
	assert.Equal(200, res.Code)

	// after create module, list module
	module, err = moduleDao.ListModule([]string{})
	require.NoError(err)
	assert.Equal(1, len(module))
	assert.Equal("test_module_wasm", module[0].Name)

	// Compare with wasm magic number: \x00\x61\x73\x6d (0x6d736100 in little endian)
	wasmMagic := []byte{0x00, 0x61, 0x73, 0x6d}
	wasmELF := module[0].Wasm
	assert.Equal(wasmELF[:4], wasmMagic)

}

func TestListModule(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testDbFilePath := testutils.GetTmpFile()
	// We'll not cleanup the temp file, as it's troublesome to turn down the http server, and probably not worth it in a
	// test.

	sqliteClient, _ := dao.InitSqlite(testDbFilePath)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	gLock := lock.NewLock()
	waitCond := cond.NewCond()

	wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath, err := initWasiSDK()
	require.Nil(err)

	fakeServer := fake.StartFakeNewServer(sqliteClient, gLock, waitCond,
		pgClient, grafanaURL, wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath)

	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}

	// before create module, add module to db
	id := strings.Replace(uuid.New(), "-", "_", -1)
	module := &dao.ModuleGORM{
		ID:                 id,
		DesireState:        int(pb.ModuleState_CREATED_),
		Name:               "TestModule",
		Wasm:               []byte("WasmUid"),
		CreateTime:         time.Now().Format("2006-01-02 15:04:05"),
		EbpfPerfBufferName: "events",
	}

	// save module
	err = moduleDao.SaveModule(module)
	require.NoError(err)

	client := NewClient("http://" + fakeServer.String())

	// test list module
	req := &http.ListModuleReq{}
	res, err := client.ListModules(req)
	require.NoError(err)
	assert.Equal(200, res.Code)
	assert.Equal(1, len(res.Data))
	assert.Equal("TestModule", res.Data[0].Name)
}

func TestDeleteModule(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testDbFilePath := testutils.GetTmpFile()
	// We'll not cleanup the temp file, as it's troublesome to turn down the http server, and probably not worth it in a
	// test.

	sqliteClient, _ := dao.InitSqlite(testDbFilePath)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	gLock := lock.NewLock()
	waitCond := cond.NewCond()

	wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath, err := initWasiSDK()
	require.Nil(err)

	fakeServer := fake.StartFakeNewServer(sqliteClient, gLock, waitCond, pgClient,
		grafanaURL, wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath)

	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}

	// before create module, add module to db
	id := strings.Replace(uuid.New(), "-", "_", -1)
	module := &dao.ModuleGORM{
		ID:                 id,
		DesireState:        int(pb.ModuleState_CREATED_),
		Name:               "TestModule",
		Wasm:               []byte("WasmUid"),
		CreateTime:         time.Now().Format("2006-01-02 15:04:05"),
		EbpfPerfBufferName: "events",
	}

	// save module
	err = moduleDao.SaveModule(module)
	require.NoError(err)

	client := NewClient("http://" + fakeServer.String())

	// test list module
	res, err := client.DeleteModule(id)
	require.NoError(err)
	assert.Equal(200, res.Code)
	assert.Contains(res.Message, "Success")

	// after delete module, list module
	moduleRes, err := moduleDao.ListModule([]string{})
	require.NoError(err)
	assert.Equal(0, len(moduleRes))
}

func TestDeployModule(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testDbFilePath := testutils.GetTmpFile()
	// We'll not cleanup the temp file, as it's troublesome to turn down the http server, and probably not worth it in a
	// test.

	sqliteClient, _ := dao.InitSqlite(testDbFilePath)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	gLock := lock.NewLock()
	waitCond := cond.NewCond()

	wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath, err := initWasiSDK()
	require.Nil(err)

	fakeServer := fake.StartFakeNewServer(sqliteClient, gLock, waitCond, pgClient,
		grafanaURL, wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath)

	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}

	// before create module, add module to db
	moduleID := strings.Replace(uuid.New(), "-", "_", -1)
	module := &dao.ModuleGORM{
		ID:                 moduleID,
		Ebpf:               ebpfJson,
		Wasm:               []byte("moduleString"),
		CreateTime:         time.Date(2022, 12, 31, 14, 30, 0, 0, time.Local).Format("2006-01-02 15:04:05"),
		DesireState:        int(pb.ModuleState_CREATED_),
		Name:               "test-module-foo",
		EbpfFmt:            0,
		EbpfLang:           0,
		EbpfPerfBufferName: "events",

		SchemaName: "out_put_name",
		SchemaAttr: "[{\"name\":\"data\",\"type\":5}]",
		Fn:         "copy_input_to_output",
		WasmFmt:    0,
		WasmLang:   0,
	}

	// save module
	err = moduleDao.SaveModule(module)
	require.NoError(err)

	moduleRes, err := moduleDao.ListModule([]string{})
	require.NoError(err)
	assert.Equal(1, len(moduleRes))
	assert.Equal("test-module-foo", moduleRes[0].Name)
	assert.Equal(int(pb.ModuleState_CREATED_), moduleRes[0].DesireState)

	client := NewClient("http://" + fakeServer.String())

	// test deploy module
	res, err := client.DeployModule(moduleID)
	require.NoError(err)
	assert.Equal(200, res.Code)
	assert.Contains(res.Message, "prepare to deploy module")

	// check module state
	moduleRes, err = moduleDao.ListModule([]string{})
	require.NoError(err)
	assert.Equal(1, len(moduleRes))
	assert.Equal("test-module-foo", moduleRes[0].Name)
	assert.Equal(int(pb.ModuleState_DEPLOYED), moduleRes[0].DesireState)
}

func TestUndeployModule(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testDbFilePath := testutils.GetTmpFile()
	// We'll not cleanup the temp file, as it's troublesome to turn down the http server, and probably not worth it in a
	// test.

	sqliteClient, _ := dao.InitSqlite(testDbFilePath)

	cleanerFn, grafanaURL, err := grafanatest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleanerFn()) }()

	pgClientCleanerFn, pgClient, err := pgclienttest.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(pgClientCleanerFn()) }()

	gLock := lock.NewLock()
	waitCond := cond.NewCond()

	wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath, err := initWasiSDK()
	require.Nil(err)

	fakeServer := fake.StartFakeNewServer(sqliteClient, gLock, waitCond, pgClient,
		grafanaURL, wasiSDK, wasiStarshipIncludePath, wasiBuildTmpPath)

	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}

	// before create module, add module to db
	moduleID := strings.Replace(uuid.New(), "-", "_", -1)
	module := &dao.ModuleGORM{
		ID:                 moduleID,
		Ebpf:               ebpfJson,
		Wasm:               []byte("moduleString"),
		CreateTime:         time.Date(2022, 12, 31, 14, 30, 0, 0, time.Local).Format("2006-01-02 15:04:05"),
		DesireState:        int(pb.ModuleState_CREATED_),
		Name:               "test-module-foo",
		EbpfFmt:            0,
		EbpfLang:           0,
		EbpfPerfBufferName: "events",

		SchemaName: "out_put_name",
		SchemaAttr: "[{\"name\":\"data\",\"type\":5}]",
		Fn:         "copy_input_to_output",
		WasmFmt:    0,
		WasmLang:   0,
	}

	// save module
	err = moduleDao.SaveModule(module)
	require.NoError(err)

	moduleRes, err := moduleDao.ListModule([]string{})
	require.NoError(err)
	assert.Equal(1, len(moduleRes))
	assert.Equal("test-module-foo", moduleRes[0].Name)
	assert.Equal(int(pb.ModuleState_CREATED_), moduleRes[0].DesireState)

	client := NewClient("http://" + fakeServer.String())

	// test deploy module
	res, err := client.UndeployModule(moduleID)
	require.NoError(err)
	assert.Equal(200, res.Code)
	assert.Contains(res.Message, "success")

	// check module state
	moduleRes, err = moduleDao.ListModule([]string{})
	require.NoError(err)
	assert.Equal(1, len(moduleRes))
	assert.Equal("test-module-foo", moduleRes[0].Name)
	assert.Equal(int(pb.ModuleState_UNDEPLOYED), moduleRes[0].DesireState)
}
