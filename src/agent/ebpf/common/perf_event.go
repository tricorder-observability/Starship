package common

// perf_event.go defines enums for attaching perf events.

// Below is the perf event type enum defined in
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/perf_event.h#L31
//enum perf_type_id {
//	PERF_TYPE_HARDWARE			= 0,
//	PERF_TYPE_SOFTWARE			= 1,
//	PERF_TYPE_TRACEPOINT			= 2,
//	PERF_TYPE_HW_CACHE			= 3,
//	PERF_TYPE_RAW				= 4,
//	PERF_TYPE_BREAKPOINT			= 5,
//
//	PERF_TYPE_MAX,
//};

// Replicated from the above definition
// Used for attaching perf events.
const (
	PerfTypeSoftware = 1
	// TODO: Add additional ones from the above enum, right now do not list them to
	// avoid confusion.
)

// https://elixir.bootlin.com/linux/v4.2/source/include/uapi/linux/perf_event.h#L103
//enum perf_sw_ids {
//	PERF_COUNT_SW_CPU_CLOCK			= 0,
//	PERF_COUNT_SW_TASK_CLOCK		= 1,
//	PERF_COUNT_SW_PAGE_FAULTS		= 2,
//	PERF_COUNT_SW_CONTEXT_SWITCHES		= 3,
//	PERF_COUNT_SW_CPU_MIGRATIONS		= 4,
//	PERF_COUNT_SW_PAGE_FAULTS_MIN		= 5,
//	PERF_COUNT_SW_PAGE_FAULTS_MAJ		= 6,
//	PERF_COUNT_SW_ALIGNMENT_FAULTS		= 7,
//	PERF_COUNT_SW_EMULATION_FAULTS		= 8,
//	PERF_COUNT_SW_DUMMY			= 9,
//
//	PERF_COUNT_SW_MAX,			/* non-ABI */
//};

const (
	PerfCountSWCPUClock = 0
	// TODO: Add additional ones from the above enum, right now do not list them to
	// avoid confusion.
)
