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

package driver

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	linux_headers "github.com/tricorder/src/agent/ebpf/bcc/linux-headers"
	tsdb "github.com/tricorder/src/testing/timescaledb"

	modulepb "github.com/tricorder/src/pb/module"
	commonpb "github.com/tricorder/src/pb/module/common"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
	wasmpb "github.com/tricorder/src/pb/module/wasm"
	testutils "github.com/tricorder/src/testing/bazel"
)

// Tests that module is deployed and data can be polled from perf buffer and write to wasm runtime.
func TestModuleDeploymentAndPoll(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	const sampleJSONBPFCPath = "modules/sample_json/sample_json.bcc.c"
	bccCode, err := testutils.ReadTestFile(sampleJSONBPFCPath)
	require.Nil(err)

	ebpfPB := ebpfpb.Program{
		Lang:           commonpb.Lang_C,
		Code:           bccCode,
		PerfBufferName: "events",
		Probes: []*ebpfpb.ProbeSpec{
			{
				Type:  ebpfpb.ProbeSpec_SAMPLE_PROBE,
				Entry: "sample_json",

				SamplePeriodNanos: 100 * 1000 * 1000,
			},
		},
	}

	wasmRelPath := "modules/sample_json/sample_json.wasm"
	wasmBinaryCode, err := testutils.ReadTestBinFile(wasmRelPath)
	assert.Nil(err)

	wasmPB := wasmpb.Program{
		Code:   wasmBinaryCode,
		FnName: "copy_input_to_output",
		OutputSchema: &commonpb.Schema{
			Name: "data",
			Fields: []*commonpb.DataField{
				{
					Name: "data",
					Type: commonpb.DataField_JSONB,
				},
			},
		},
	}

	modPB := &modulepb.Module{
		Name: "test_module",
		Ebpf: &ebpfPB,
		Wasm: &wasmPB,
	}

	cleaner, pgClient, err := tsdb.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleaner()) }()

	m, err := Deploy(modPB, pgClient)
	require.Nil(err)

	// Starship would create this table in the API server. We have to create table manually here in test.
	err = pgClient.CreateTable(m.outputSchema)
	require.Nil(err)

	time.Sleep(time.Second)
	assert.Nil(m.Poll())

	// TODO: Check database.
	m.Undeploy()
	jsons, err := pgClient.Query(fmt.Sprintf("select data #>> '{}' from %s", m.outputSchema.Name))
	assert.Nil(err)
	require.Greater(len(jsons), 0)
	assert.Equal(`{"age": 30, "name": "John"}`, jsons[0][0])
}

// Tests that the sample event module works as expected.
func TestSampleEventModule(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	assert.Nil(linux_headers.Init())

	const bccPath = "modules/sample_event/sample_event.bcc"
	const wasmRelPath = "modules/sample_event/write_events_to_output.wasm"

	bccCode, err := testutils.ReadTestFile(bccPath)
	require.Nil(err)

	ebpfPB := ebpfpb.Program{
		Code:           bccCode,
		PerfBufferName: "events",
		Probes: []*ebpfpb.ProbeSpec{
			{
				Type:  ebpfpb.ProbeSpec_SAMPLE_PROBE,
				Entry: "sample_event",

				SamplePeriodNanos: 100 * 1000 * 1000,
			},
		},
	}

	wasmBinaryCode, err := testutils.ReadTestBinFile(wasmRelPath)
	assert.Nil(err)

	wasmPB := wasmpb.Program{
		Code:   wasmBinaryCode,
		FnName: "write_events_to_output",
		OutputSchema: &commonpb.Schema{
			Name: "data",
			Fields: []*commonpb.DataField{
				{
					Name: "data",
					Type: commonpb.DataField_JSONB,
				},
			},
		},
	}

	modPB := &modulepb.Module{
		Name: "test_module",
		Ebpf: &ebpfPB,
		Wasm: &wasmPB,
	}

	cleaner, pgClient, err := tsdb.LaunchContainer()
	require.Nil(err)
	defer func() { assert.Nil(cleaner()) }()

	m, err := Deploy(modPB, pgClient)
	require.Nil(err)

	// Starship would create this table in the API server. We have to create table manually here in test.
	err = pgClient.CreateTable(m.outputSchema)
	require.Nil(err)

	time.Sleep(time.Second)
	assert.Nil(m.Poll())

	m.Undeploy()
	jsons, err := pgClient.Query(fmt.Sprintf("select data #>> '{}' from %s", m.outputSchema.Name))
	assert.Nil(err)
	require.Greater(len(jsons), 0)
	assert.Equal(`{"D": 0, "F": 0, "I": 0, "L": 0, "Comm": ""}`, jsons[0][0])
}
