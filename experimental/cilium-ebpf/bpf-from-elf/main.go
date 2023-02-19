package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	// Parse CLI arguments
	var elfFile string
	var eventMapName string

	switch len(os.Args) {
	case 1:
		log.Printf("Usage: sudo %s <BPF obj file> [<event map name>]\n", os.Args[0])
		log.Printf("Usage: sudo %s ~/libbpf-bootstrap/examples/c/.output/bootstrap.bpf.o rb\n", os.Args[0])
		return
	case 2:
		elfFile, eventMapName = os.Args[1], ""
	default:
		elfFile, eventMapName = os.Args[1], os.Args[2]
	}

	log.Printf("Input ELF file: %v\n", elfFile)

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	exitCh := make(chan struct{}) // notify other goroutines

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Parse ELF into a spec, which contains BPF map/prog information
	spec, err := ebpf.LoadCollectionSpec(elfFile)
	if err != nil {
		log.Fatalf("Load BPF ELF file failed: %v", err)
	}

	log.Printf("Load CollectionSpec from ELF successful\n\n")
	dumpCollectionSpec(spec)

	// Load BPF map/progs into kernel.
	// Note that BPF map will be created automatically in this phase.
	collection, err := ebpf.LoadCollection(elfFile)
	if err != nil {
		log.Fatalf("Load Collection from ELF failed: %v", err)
	}
	defer collection.Close()

	log.Printf("Load BPF collection (programs/maps) into kernel successful\n")
	dumpCollection(collection)

	// Attach BPF programs to kernel hook points
	if err := attachPrograms(spec, collection, exitCh); err != nil {
		log.Printf("Attach BPF programs to hook points failed: %v\n", err)
	}

	// Poll data from kernel BPF map (ringbuffer, perfbuffer, etc)
	if eventMapName != "" {
		if err := pollData(spec, collection, eventMapName, exitCh); err != nil {
			log.Printf("Poll data from BPF map failed: %v\n", err)
		}
	}

	// Wait
	<-stopper
	log.Printf("Received stop signal, notify goroutines to exit\n")

	close(exitCh)
	time.Sleep(1 * time.Second)
	log.Printf("Agent exited\n")
}

// Attach BPF programs in `collection` into the specified kernel hook.
func attachPrograms(spec *ebpf.CollectionSpec, collection *ebpf.Collection, exitCh chan struct{}) error {
	if spec == nil || collection == nil {
		return fmt.Errorf("Unexpected nil spec or collection: %v %v", spec, collection)
	}

	for progName, prog := range collection.Programs {
		progSpec, ok := spec.Programs[progName]
		if !ok {
			log.Printf("ProgramSpec for %s not found, skip attaching\n", progName)
			continue
		}

		switch progSpec.Type {
		case ebpf.Tracing:
			log.Printf("Attaching BPF program (fentry/fexit) %s\n", progName)
			lk, err := link.AttachTracing(link.TracingOptions{
				Program: prog,
			})
			if err != nil {
				log.Fatalf("Attach fentry/fexit failed: %s", err)
			}

			go func(lk link.Link, progName string) {
				select {
				case <-exitCh:
					log.Printf("Received stop signal, close BPF link for program %s\n", progName)
					lk.Close()
				}
			}(lk, progName)

		case ebpf.Kprobe:
			// Name of the kernel function to trace.
			fn := progSpec.AttachTo

			log.Printf("Attaching BPF program %s to %s", progName, fn)
			kp, err := link.Kprobe(fn, prog, nil)
			if err != nil {
				log.Fatalf("Attach kprobe failed: %s", err)
			}

			go func(lk link.Link, progName string) {
				select {
				case <-exitCh:
					log.Printf("Received stop signal, close BPF link for program %s\n", progName)
					lk.Close()
				}
			}(kp, progName)

		case ebpf.TracePoint:
			log.Printf("Attaching BPF program %v to %s", prog, progName)

			items := strings.Split(progSpec.AttachTo, "/")
			if len(items) != 2 {
				log.Printf("Unknown AttachTo field for program %s, skip attaching\n", progSpec.AttachTo)
				continue
			}

			group, name := items[0], items[1]
			log.Printf("Attaching BPF program %v to %s: group %s, name %s", prog, progName, group, name)

			lk, err := link.Tracepoint(group, name, prog, nil)
			if err != nil {
				log.Fatalf("Attach BPF program %v to %s failed: %v", prog, progName, err)
				continue
			}

			go func(lk link.Link, progName string) {
				select {
				case <-exitCh:
					log.Printf("Received stop signal, close BPF link for program %s\n", progName)
					lk.Close()
				}
			}(lk, progName)

		// We can certainly support more types, add them here
		// case ebpf.TODO

		default:
			log.Printf("Unsupported program type %v\n", progSpec.Type)
		}
	}

	return nil
}

func pollData(spec *ebpf.CollectionSpec, collection *ebpf.Collection, eventMapName string, exitCh chan struct{}) error {
	if spec == nil || collection == nil {
		return fmt.Errorf("Unexpected nil spec or collection: %v %v", spec, collection)
	}

	// Get the BPF map that stores kernel events
	for mapName, m := range collection.Maps {
		mapSpec, ok := spec.Maps[mapName]
		if !ok {
			log.Printf("MapSpec for %s not found, should not happen\n", mapName)
			continue
		}

		if mapName != eventMapName {
			continue
		}

		switch mapSpec.Type {
		case ebpf.RingBuf:
			rd, err := ringbuf.NewReader(m)
			if err != nil {
				log.Fatalf("opening ringbuf reader: %s", err)
			}

			log.Printf("Create ringbuffer for polling data from BPF map %s successful\n", mapName)

			go ringBufReadLoop(rd)

			go func(rd *ringbuf.Reader, mapName string) {
				select {
				case <-exitCh:
					log.Printf("Received stop signal, close polling reader for BPF map %s\n", mapName)
					rd.Close()
				}
			}(rd, mapName)

		// We can certainly support more types, add them here
		// case ebpf.TODO

		default:
			log.Printf("Unsupported BPF map type %v\n", mapSpec.Type)
		}

		// NOTE: we assume only one BPF map is used for polling
		break
	}

	return nil
}

// Read loop for ring buffer
func ringBufReadLoop(rd *ringbuf.Reader) {
	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("received signal, exiting..")
				return
			}

			log.Printf("reading from reader: %s", err)
			continue
		}

		log.Printf("Got an event from BPF map, event size %d\n", len(record.RawSample))
	}
}

// Utilities
func dumpCollectionSpec(spec *ebpf.CollectionSpec) {
	log.Printf("Dumping ELF collection spec\n")
	for mapName, mapSpec := range spec.Maps {
		log.Printf("BPF map: name %s, spec %+v\n", mapName, mapSpec)
	}
	for progName, progSpec := range spec.Programs {
		log.Printf("BPF program: name %s, spec %+v\n", progName, progSpec)
	}
	log.Printf("Dumping ELF collection spec done\n\n")
}

func dumpCollection(spec *ebpf.Collection) {
	log.Printf("Dumping ELF collection\n")
	for mapName, mapObj := range spec.Maps {
		log.Printf("BPF map: name %s, spec %+v\n", mapName, *mapObj)
	}

	// // Program represents BPF program loaded into the kernel.
	for progName, progObj := range spec.Programs {
		log.Printf("BPF program: name %s, spec %+v\n", progName, *progObj)
	}
	log.Printf("Dumping ELF collection done\n\n")
}
