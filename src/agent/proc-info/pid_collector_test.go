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

package proc_info

import (
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/tricorder/src/api-server/pb"
	testuitls "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/file"
)

func TestGKEFormat(t *testing.T) {
	assert := assert.New(t)
	tmpDir := testuitls.CreateTmpDir()
	destination := path.Join(
		tmpDir,
		"fs/cgroup/cpu,cpuacct/kubepods/pod8dbc5577-d0e2-4706-8787-57d52c03ddf2/"+
			"14011c7d92a9e513dfd69211da0413dbf319a5e45a02b354ba6e98e10272542d/cgroup.procs",
	)
	assert.Nil(file.Create(destination))

	ci := &pb.ContainerInfo{
		Id:       "docker://14011c7d92a9e513dfd69211da0413dbf319a5e45a02b354ba6e98e10272542d",
		PodUid:   "8dbc5577-d0e2-4706-8787-57d52c03ddf2",
		QosClass: "Burstable",
	}
	_, err := grabProcessInfo(tmpDir, ci)
	assert.Nil(err)
}

func TestGKEFormat2(t *testing.T) {
	assert := assert.New(t)
	tmpDir := testuitls.CreateTmpDir()
	destination := path.Join(
		tmpDir,
		"fs/cgroup/cpu,cpuacct/kubepods/burstable/podc458de04-9784-4f7a-990e-cefe26b511f0/"+
			"01aa0bfe91e8a58da5f1f4db469fa999fe9263c702111e611445cde2b9cb0c1a/cgroup.procs",
	)
	assert.Nil(file.Create(destination))

	ci := &pb.ContainerInfo{
		Id:       "docker://01aa0bfe91e8a58da5f1f4db469fa999fe9263c702111e611445cde2b9cb0c1a",
		PodUid:   "c458de04-9784-4f7a-990e-cefe26b511f0",
		QosClass: "Burstable",
	}
	_, err := grabProcessInfo(tmpDir, ci)
	assert.Nil(err)
}

func TestStandardFormatDocker(t *testing.T) {
	assert := assert.New(t)
	tmpDir := testuitls.CreateTmpDir()
	destination := path.Join(
		tmpDir,
		"fs/cgroup/cpu,cpuacct/kubepods.slice/kubepods-pod8dbc5577_d0e2_4706_8787_57d52c03ddf2.slice/"+
			"docker-14011c7d92a9e513dfd69211da0413dbf319a5e45a02b354ba6e98e10272542d.scope/cgroup.procs",
	)
	assert.Nil(file.Create(destination))

	ci := &pb.ContainerInfo{
		Id:       "docker://14011c7d92a9e513dfd69211da0413dbf319a5e45a02b354ba6e98e10272542d",
		PodUid:   "8dbc5577_d0e2_4706_8787_57d52c03ddf2",
		QosClass: "Burstable",
	}
	_, err := grabProcessInfo(tmpDir, ci)
	assert.Nil(err)
}

func TestStandardFormatCRIO(t *testing.T) {
	assert := assert.New(t)
	tmpDir := testuitls.CreateTmpDir()
	destination := path.Join(
		tmpDir,
		"fs/cgroup/cpu,cpuacct/kubepods.slice/kubepods-pod8dbc5577_d0e2_4706_8787_57d52c03ddf2.slice/"+
			"crio-14011c7d92a9e513dfd69211da0413dbf319a5e45a02b354ba6e98e10272542d.scope/cgroup.procs",
	)
	assert.Nil(file.Create(destination))

	ci := &pb.ContainerInfo{
		Id:       "crio://14011c7d92a9e513dfd69211da0413dbf319a5e45a02b354ba6e98e10272542d",
		PodUid:   "8dbc5577_d0e2_4706_8787_57d52c03ddf2",
		QosClass: "Burstable",
	}
	_, err := grabProcessInfo(tmpDir, ci)
	assert.Nil(err)
}

func TestOpenShiftFormat(t *testing.T) {
	assert := assert.New(t)
	tmpDir := testuitls.CreateTmpDir()
	destination := path.Join(
		tmpDir,
		"fs/cgroup/cpu,cpuacct/kubepods.slice/kubepods-burstable.slice/"+
			"kubepods-burstable-pod9b7969b2_aad0_47d4_b11c_4acfd1ce018e.slice/"+
			"crio-9b9ccc15d288aa0f7d3bf7b583993921bf261edfeff3467765ab81e687c6a889.scope/cgroup.procs",
	)
	assert.Nil(file.Create(destination))

	ci := &pb.ContainerInfo{
		Id:       "crio://9b9ccc15d288aa0f7d3bf7b583993921bf261edfeff3467765ab81e687c6a889",
		PodUid:   "9b7969b2-aad0-47d4-b11c-4acfd1ce018e",
		QosClass: "Burstable",
	}
	_, err := grabProcessInfo(tmpDir, ci)
	assert.Nil(err)
}

func TestBareMetalFormat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	basePath := testuitls.CreateTmpDir()
	destination := path.Join(
		basePath,
		"fs/cgroup/cpu,cpuacct/system.slice/containerd.service/"+
			"kubepods-besteffort-pod1544eb37_e4f7_49eb_8cc4_3d01c41be77b.slice:cri-containerd:"+
			"8618d3540ce713dd59ed0549719643a71dd482c40c21685773e7ac1291b004f5/cgroup.procs",
	)
	assert.Nil(file.Create(destination))
	assert.Nil(file.Write(destination, " "+strconv.Itoa(os.Getpid())+" \n  \n   "))

	ci := &pb.ContainerInfo{
		Id:       "cri-containerd://8618d3540ce713dd59ed0549719643a71dd482c40c21685773e7ac1291b004f5",
		PodUid:   "1544eb37-e4f7-49eb-8cc4-3d01c41be77b",
		QosClass: "BestEffort",
	}

	procInfo, err := grabProcessInfo(basePath, ci)
	assert.Nil(err)
	require.Equal(1, len(procInfo.ProcList))
	assert.Equal(int32(os.Getpid()), procInfo.ProcList[0].Id)
	assert.Greater(procInfo.ProcList[0].CreateTime, int64(0))
}
