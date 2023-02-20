// Package bcc provides types and APIs for working with BCC-style eBPF C programs.
// Largely wraps BCC's go binding.
package bcc

import (
	"fmt"

	"github.com/iovisor/gobpf/bcc"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/agent/ebpf/common"

	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/pb"
)

// Wraps BCC Module object
type module struct {
	m *bcc.Module
}

func newModule(code string) (*module, error) {
	bccModule := bcc.NewModule(code, []string{} /*cflags*/)
	if bccModule == nil {
		return nil, fmt.Errorf("while creating module, failed to create BCC Module: got nil return value")
	}
	return &module{m: bccModule}, nil
}

// NewPerfBuffer returns a PerfBuffer object with the input name.
func (m *module) NewPerfBuffer(name string) (*PerfBuffer, error) {
	return NewPerfBuffer(m.m, name)
}

// LoadKprobe load the kprobe specified by the input name, and returns the file descriptor pointed to
// the loaded kprobe; returns error if failed.
func (m *module) LoadKprobe(name string) (int, error) {
	return m.m.LoadKprobe(name)
}

// TODO(yzhao): Detect system's CPU count, and set this value as 4 * CPU-count.
// The BPF runtime's default is 2*cpu-count, which is not enough for golang program,
// which tends to call syscalls more concurrently.
//
// A value of -1 maxActive signifies to use the default, which is 2 * NR_CPU (https://stackoverflow.com/a/36517308)
// See also the kernel kprobes documentation
// https://www.kernel.org/doc/Documentation/kprobes.txt
const maxActiveRetProbes = 512

func (m *module) attachKEntryProbe(probeFunc, probeName string) error {
	probe, err := m.m.LoadKprobe(probeName)
	context := fmt.Sprintf("attaching kentryprobe '%s' to '%s'", probeName, probeFunc)
	if err != nil {
		return errors.Wrap(context, "load", err)
	}
	if err := m.m.AttachKprobe(probeFunc, probe, maxActiveRetProbes); err != nil {
		return errors.Wrap(context, "attach", err)
	}
	return nil
}

func (m *module) attachKReturnProbe(probeFunc, probeName string) error {
	probe, err := m.m.LoadKprobe(probeName)
	if err != nil {
		return fmt.Errorf("failed to load %s, error: %v", probeName, err)
	}
	if err := m.m.AttachKretprobe(probeFunc, probe, maxActiveRetProbes); err != nil {
		return fmt.Errorf("failed to attach kretprobe %s, error: %v", probeName, err)
	}
	return nil
}

func (m *module) attachUEntryProbe(binaryPath, probeFunc, probeName string) error {
	probe, err := m.m.LoadUprobe(probeName)
	context := fmt.Sprintf("attaching uentryprobe '%s' to '%s' for '%s'", probeName, probeFunc, binaryPath)
	if err != nil {
		return errors.Wrap(context, "load", err)
	}
	if err := m.m.AttachUprobe(binaryPath, probeFunc, probe, 0); err != nil {
		return errors.Wrap(context, "attach", err)
	}
	return nil
}

func (m *module) attachUReturnProbe(binaryPath, probeFunc, probeName string) error {
	probe, err := m.m.LoadUprobe(probeName)
	if err != nil {
		return fmt.Errorf("failed to load %s for %s, error: %v", probeName, binaryPath, err)
	}
	if err := m.m.AttachUretprobe(binaryPath, probeFunc, probe, 0); err != nil {
		return fmt.Errorf("failed to attach uretprobe %s, error: %v", probeName, err)
	}
	return nil
}

func (m *module) attachTPProbe(probeFunc, probeName string) error {
	probe, err := m.m.LoadTracepoint(probeName)
	context := fmt.Sprintf("attaching tracepoint '%s' to '%s'", probeName, probeFunc)
	if err != nil {
		return errors.Wrap(context, "load", err)
	}
	if err := m.m.AttachTracepoint(probeFunc, probe); err != nil {
		return errors.Wrap(context, "attach", err)
	}
	return nil
}

func (m *module) attachKProbe(probe *ebpfpb.ProbeSpec) error {
	log.Infof("Attaching kprobe %v", probe)
	if probe.Type != ebpfpb.ProbeSpec_KPROBE {
		return fmt.Errorf("must be kprobe, got %v", probe)
	}
	if len(probe.Target) == 0 {
		return fmt.Errorf("while attaching kprobe '%v', target cannot be empty", probe)
	}

	if probe.Entry != "" {
		if err := m.attachKEntryProbe(probe.Target, probe.Entry); err != nil {
			return err
		}
	}
	if probe.Return != "" {
		if err := m.attachKReturnProbe(probe.Target, probe.Return); err != nil {
			return err
		}
	}
	return nil
}

func (m *module) attachSyscallProbe(probe *ebpfpb.ProbeSpec) error {
	log.Infof("Attaching syscall probe %v", probe)
	if probe.Type != ebpfpb.ProbeSpec_SYSCALL_PROBE {
		return fmt.Errorf("must be syscall probe, got %v", probe)
	}
	if len(probe.Target) == 0 {
		return fmt.Errorf("while attaching  syscall probe '%v', target cannot be empty", probe)
	}

	syscallName := bcc.GetSyscallFnName(probe.Target)
	if probe.Entry != "" {
		if err := m.attachKEntryProbe(syscallName, probe.Entry); err != nil {
			return err
		}
	}
	if probe.Return != "" {
		if err := m.attachKReturnProbe(syscallName, probe.Return); err != nil {
			return err
		}
	}
	return nil
}

func (m *module) attachProbe(probe *ebpfpb.ProbeSpec) error {
	log.Infof("Attaching probe %v", probe)
	switch probe.Type {
	case ebpfpb.ProbeSpec_KPROBE:
		return m.attachKProbe(probe)
	case ebpfpb.ProbeSpec_SYSCALL_PROBE:
		return m.attachSyscallProbe(probe)
	case ebpfpb.ProbeSpec_UPROBE:
		return m.attachUProbe(probe)
	case ebpfpb.ProbeSpec_TRACEPOINT:
		return m.attachTracepoint(probe)
	case ebpfpb.ProbeSpec_SAMPLE_PROBE:
		return m.attachSampleProbe(probe)
	default:
		return fmt.Errorf("unknown probe type '%d'", probe.Type)
	}
}

// The following values indicate the corresponding argument is ignored by the underlying system
// attachment routines.
const (
	ignoreSampleFreq int = 0
	ignorePID        int = -1
	ignoreCPU        int = -1
	ignoreGroupFD    int = -1
)

// attachSampleProbe attaches a perf event which periodicially got triggered.
func (m *module) attachSampleProbe(probe *ebpfpb.ProbeSpec) error {
	log.Infof("Attaching sample probe %v", probe)

	probeFD, err := m.m.LoadPerfEvent(probe.Entry)
	if err != nil {
		return fmt.Errorf("while attaching sampling perf event, failed to load perf event probe '%s', error: %v",
			pb.FormatOneLine(probe), err)
	}
	log.Printf("SamplePeriodNanos: %d", probe.SamplePeriodNanos)
	// Parameter names:
	// (evType, evConfig int, samplePeriod int, sampleFreq int, pid, cpu, groupFd, fd int)
	err = m.m.AttachPerfEvent(common.PerfTypeSoftware, common.PerfCountSWCPUClock, int(probe.SamplePeriodNanos),
		ignoreSampleFreq, ignorePID, ignoreCPU, ignoreGroupFD, probeFD)
	if err != nil {
		return fmt.Errorf("while attaching sampling perf event, failed to attach perf event, error: %v", err)
	}
	return nil
}

func (m *module) attachUProbe(probe *ebpfpb.ProbeSpec) error {
	log.Infof("Attaching kprobe %v", probe)
	if probe.Type != ebpfpb.ProbeSpec_UPROBE {
		return fmt.Errorf("must be uprobe, got %v", probe)
	}
	context := fmt.Sprintf("while attaching uprobe '%v', target cannot be empty", probe)
	if len(probe.Target) == 0 {
		return errors.Wrap(context, "target", fmt.Errorf("cannot be empty"))
	}

	if len(probe.BinaryPath) == 0 {
		return errors.Wrap(context, "binary path", fmt.Errorf("cannot be empty"))
	}

	if probe.Entry != "" {
		if err := m.attachUEntryProbe(probe.BinaryPath, probe.Target, probe.Entry); err != nil {
			return err
		}
	}
	if probe.Return != "" {
		if err := m.attachUReturnProbe(probe.BinaryPath, probe.Target, probe.Return); err != nil {
			return err
		}
	}
	return nil
}

func (m *module) attachTracepoint(probe *ebpfpb.ProbeSpec) error {
	log.Infof("Attaching kprobe %v", probe)
	if probe.Type != ebpfpb.ProbeSpec_UPROBE {
		return fmt.Errorf("must be tracepoint, got %v", probe)
	}
	if len(probe.Target) == 0 {
		return fmt.Errorf("while attaching tracepoint '%v', target cannot be empty", probe)
	}

	if probe.Entry != "" {
		if err := m.attachTPProbe(probe.Target, probe.Entry); err != nil {
			return err
		}
	}

	return nil
}

func (m *module) Close() {
	m.m.Close()
}
