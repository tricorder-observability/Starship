// LINT_C_FILE

#include <linux/in6.h>
#include <linux/net.h>
#include <linux/socket.h>
#include <net/inet_sock.h>

// TODO 20221126: hardcode by @ArthurChiao
#define ENABLE_HTTP_TRACING 1
#define ENABLE_CQL_TRACING 0
#define ENABLE_MONGO_TRACING 0
#define ENABLE_PGSQL_TRACING 0
#define ENABLE_MYSQL_TRACING 0
#define ENABLE_MUX_TRACING 0
#define ENABLE_DNS_TRACING 0
#define ENABLE_AMQP_TRACING 0
#define ENABLE_REDIS_TRACING 0
#define ENABLE_NATS_TRACING 0
#define ENABLE_KAFKA_TRACING 0

#define socklen_t size_t

#include <linux/sched.h>

#define __inline inline __attribute__((__always_inline__))

// _VAR suffix indicates that the count of bytes being read is equal to the
// variable's byte size.
#define BPF_PROBE_READ_VAR(value, ptr)                                         \
  bpf_probe_read(&value, sizeof(value), ptr)
#define BPF_PROBE_READ_KERNEL_VAR(value, ptr)                                  \
  bpf_probe_read_kernel(&value, sizeof(value), ptr)

static __inline int32_t read_big_endian_int32(const char *buf) {
  int32_t length;
  BPF_PROBE_READ_VAR(length, buf);
  return bpf_ntohl(length);
}

static __inline int16_t read_big_endian_int16(const char *buf) {
  int16_t val;
  BPF_PROBE_READ_VAR(val, buf);
  return bpf_ntohs(val);
}

// Returns 0 if lhs and rhs compares equal up to n bytes. Otherwise a non-zero
// value is returned. NOTE #1: Cannot use C standard library's strncmp() because
// that cannot be compiled by BCC. NOTE #2: Different from the C standard
// library's strncmp(), this does not distinguish order. NOTE #3: n must be a
// literal so that the BCC runtime can unroll the inner loop. NOTE #4: Loop
// unrolling increases instruction code, be aware when BPF verifier complains
// about.
//          breaching instruction count limit.
// NOTE #5: This function is prefixed with px_ since kernels > 5.17 have a
// builtin version of this function.
static __inline int px_bpf_strncmp(const char *lhs, size_t n, const char *rhs) {
  for (size_t i = 0; i < n; ++i) {
    if (lhs[i] != rhs[i]) {
      return 1;
    }
  }
  return 0;
}

// There is a macro min() defined by a kernel header.
// We prefer being more self-contained, so define this with a different name.
#define DEFINE_MIN_FN_FOR_TYPE(type)                                           \
  static __inline type min_##type(type l, type r) { return l < r ? l : r; }
// Define the function for new types if needed
DEFINE_MIN_FN_FOR_TYPE(uint32_t)
DEFINE_MIN_FN_FOR_TYPE(int64_t)
DEFINE_MIN_FN_FOR_TYPE(uint64_t)
DEFINE_MIN_FN_FOR_TYPE(size_t)
#undef DEFINE_MIN_FN_FOR_TYPE

// This is how Linux converts nanoseconds to clock ticks.
// Used to report PID start times in clock ticks, just like /proc/<pid>/stat
// does.
static __inline uint64_t pl_nsec_to_clock_t(uint64_t x) {
  return div_u64(x, NSEC_PER_SEC / USER_HZ);
}

// TODO 20221126: hardcode by @ArthurChiao just for pass compilation, refer to
// src/stirling/bpf_tools/bcc_wrapper.cc for more info
#define GROUP_LEADER_OFFSET_OVERRIDE 0
#define START_BOOTTIME_OFFSET_OVERRIDE 0
#define START_BOOTTIME_VARNAME                                                 \
  start_boottime // kernel_version >= kLinux5p5VersionCode ? "start_boottime" :
                 // "real_start_time";

// Returns the group_leader offset.
// If GROUP_LEADER_OFFSET_OVERRIDE is defined, it is returned.
// Otherwise, the value is obtained from the definition of header structs.
// The override is important for the case when we don't have an exact header
// match. See user-space TaskStructResolver.
static __inline uint64_t task_struct_group_leader_offset() {
  if (GROUP_LEADER_OFFSET_OVERRIDE != 0) {
    return GROUP_LEADER_OFFSET_OVERRIDE;
  } else {
    return offsetof(struct task_struct, group_leader);
  }
}

// Returns the real_start_time/start_boottime offset.
// If START_BOOTTIME_OFFSET_OVERRIDE is defined, it is returned.
// Otherwise, the value is obtained from the definition of header structs.
// The override is important for the case when we don't have an exact header
// match. See user-space TaskStructResolver.
static __inline uint64_t task_struct_start_boottime_offset() {
  // Find the start_boottime of the current task.
  if (START_BOOTTIME_OFFSET_OVERRIDE != 0) {
    return START_BOOTTIME_OFFSET_OVERRIDE;
  } else {
    return offsetof(struct task_struct, START_BOOTTIME_VARNAME);
  }
}

// Effectively returns:
//   task->group_leader->start_boottime;  // Linux 5.5+
//   task->group_leader->real_start_time; // Linux 5.4 and earlier
static __inline uint64_t read_start_boottime(const struct task_struct *task) {
  uint64_t group_leader_offset = task_struct_group_leader_offset();
  struct task_struct *group_leader_ptr;
  bpf_probe_read(&group_leader_ptr, sizeof(struct task_struct *),
                 (uint8_t *)task + group_leader_offset);

  uint64_t start_boottime_offset = task_struct_start_boottime_offset();
  uint64_t start_boottime = 0;
  bpf_probe_read(&start_boottime, sizeof(uint64_t),
                 (uint8_t *)group_leader_ptr + start_boottime_offset);

  return pl_nsec_to_clock_t(start_boottime);
}

static __inline uint64_t get_tgid_start_time() {
  struct task_struct *task = (struct task_struct *)bpf_get_current_task();
  return read_start_boottime(task);
}

// UPID stands for unique pid.
// Since PIDs can be reused, this attaches the start time of the PID,
// so that the identifier becomes unique.
// Note that this version is node specific; there is also a 'class UPID'
// definition under shared which also includes an Agent ID (ASID),
// to uniquely identify PIDs across a cluster. The ASID is not required here.
struct upid_t {
  // Comes from the process from which this is captured.
  // See https://stackoverflow.com/a/9306150 for details.
  // Use union to give it two names. We use tgid in kernel-space, pid in
  // user-space.
  union {
    uint32_t pid;
    uint32_t tgid;
  };
  uint64_t start_time_ticks;
};

// This file contains definitions that are shared between various kprobes and
// uprobes.

enum message_type_t { kUnknown, kRequest, kResponse };

enum traffic_direction_t {
  kEgress,
  kIngress,
};

// Protocol being used on a connection (HTTP, MySQL, etc.).
// PROTOCOL_LIST: Requires update on new protocols.
// WARNING: Changes to this enum are API-breaking.
// You may add a protocol, but do not change values for existing protocols,
// and do not remove any protocols.
// This is a C-style enum to make it compatible with C BPF code.
// HACK ALERT: This must also match the list in
// //src/shared/protocols/protocols.h
// TODO(oazizi): Find a way to make a common source, while also keeping
// compatibility with BPF.
enum traffic_protocol_t {
  kProtocolUnknown = 0,
  kProtocolHTTP = 1,
  kProtocolHTTP2 = 2,
  kProtocolMySQL = 3,
  kProtocolCQL = 4,
  kProtocolPGSQL = 5,
  kProtocolDNS = 6,
  kProtocolRedis = 7,
  kProtocolNATS = 8,
  kProtocolMongo = 9,
  kProtocolKafka = 10,
  kProtocolMux = 11,
  kProtocolAMQP = 12,
// We use magic enum to iterate through protocols in C++ land,
// and don't want the C-enum-size trick to show up there.
#ifndef __cplusplus
  kNumProtocols
#endif
};

struct protocol_message_t {
  enum traffic_protocol_t protocol;
  enum message_type_t type;
};

// The direction of traffic expected on a probe.
// Values have single bit set, so that they could be used as bit masks.
// WARNING: Do not change the existing mappings (PxL scripts rely on them).
enum endpoint_role_t {
  kRoleClient = 1 << 0,
  kRoleServer = 1 << 1,
  kRoleUnknown = 1 << 2,
};

struct conn_id_t {
  // The unique identifier of the pid/tgid.
  struct upid_t upid;
  // The file descriptor to the opened network connection.
  int32_t fd;
  // Unique id of the conn_id (timestamp).
  uint64_t tsid;
};

// Specifies the corresponding indexes of the entries of a per-cpu array.
enum control_value_index_t {
  // This specify one pid to monitor. This is used during test to eliminate
  // noise.
  // TODO(yzhao): We need a more robust mechanism for production use, which
  // should be able to:
  // * Specify multiple pids up to a certain limit, let's say 1024.
  // * Support efficient lookup inside bpf to minimize overhead.
  kTargetTGIDIndex = 0,
  kStirlingTGIDIndex,
  kNumControlValues,
};
static __inline enum message_type_t infer_http_message(const char *buf,
                                                       size_t count) {
  // Smallest HTTP response is 17 characters:
  // HTTP/1.1 200 OK\r\n
  // Smallest HTTP response is 16 characters:
  // GET x HTTP/1.1\r\n
  if (count < 16) {
    return kUnknown;
  }

  if (buf[0] == 'H' && buf[1] == 'T' && buf[2] == 'T' && buf[3] == 'P') {
    return kResponse;
  }
  if (buf[0] == 'G' && buf[1] == 'E' && buf[2] == 'T') {
    return kRequest;
  }
  if (buf[0] == 'H' && buf[1] == 'E' && buf[2] == 'A' && buf[3] == 'D') {
    return kRequest;
  }
  if (buf[0] == 'P' && buf[1] == 'O' && buf[2] == 'S' && buf[3] == 'T') {
    return kRequest;
  }
  if (buf[0] == 'P' && buf[1] == 'U' && buf[2] == 'T') {
    return kRequest;
  }
  if (buf[0] == 'D' && buf[1] == 'E' && buf[2] == 'L' && buf[3] == 'E' &&
      buf[4] == 'T' && buf[5] == 'E') {
    return kRequest;
  }
  // TODO(oazizi): Should we add CONNECT, OPTIONS, TRACE, PATCH?

  return kUnknown;
}

// Cassandra frame:
//      0         8        16        24        32         40
//      +---------+---------+---------+---------+---------+
//      | version |  flags  |      stream       | opcode  |
//      +---------+---------+---------+---------+---------+
//      |                length                 |
//      +---------+---------+---------+---------+
//      |                                       |
//      .            ...  body ...              .
//      .                                       .
//      .                                       .
//      +----------------------------------------
static __inline enum message_type_t infer_cql_message(const char *buf,
                                                      size_t count) {
  static const uint8_t kError = 0x00;
  static const uint8_t kStartup = 0x01;
  static const uint8_t kReady = 0x02;
  static const uint8_t kAuthenticate = 0x03;
  static const uint8_t kOptions = 0x05;
  static const uint8_t kSupported = 0x06;
  static const uint8_t kQuery = 0x07;
  static const uint8_t kResult = 0x08;
  static const uint8_t kPrepare = 0x09;
  static const uint8_t kExecute = 0x0a;
  static const uint8_t kRegister = 0x0b;
  static const uint8_t kEvent = 0x0c;
  static const uint8_t kBatch = 0x0d;
  static const uint8_t kAuthChallenge = 0x0e;
  static const uint8_t kAuthResponse = 0x0f;
  static const uint8_t kAuthSuccess = 0x10;

  // Cassandra frames have a 9-byte header.
  if (count < 9) {
    return kUnknown;
  }

  // Version contains both version and direction.
  bool request = (buf[0] & 0x80) == 0x00;
  uint8_t version = (buf[0] & 0x7f);
  uint8_t flags = buf[1];
  uint8_t opcode = buf[4];
  int32_t length = read_big_endian_int32(&buf[5]);

  // Cassandra version should 5 or less. Also v2 and lower seem much less
  // popular. For example ScyllaDB only supports v3+.
  if (version < 3 || version > 5) {
    return kUnknown;
  }

  // Only flags 0x1, 0x2, 0x4 and 0x8 are used.
  if ((flags & 0xf0) != 0) {
    return kUnknown;
  }

  // A frame is limited to 256MB in length,
  // but we look for more common frames which should be much smaller in size.
  if (length > 10000) {
    return kUnknown;
  }

  switch (opcode) {
  case kStartup:
  case kOptions:
  case kQuery:
  case kPrepare:
  case kExecute:
  case kRegister:
  case kBatch:
  case kAuthResponse:
    return request ? kRequest : kUnknown;
  case kError:
  case kReady:
  case kAuthenticate:
  case kSupported:
  case kResult:
  case kEvent:
  case kAuthChallenge:
  case kAuthSuccess:
    return !request ? kResponse : kUnknown;
  default:
    return kUnknown;
  }
}

static __inline enum message_type_t infer_mongo_message(const char *buf,
                                                        size_t count) {
  // Reference:
  // https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/#std-label-wp-request-opcodes.
  // Note: Response side inference for Mongo is not robust, and is not attempted
  // to avoid confusion with other protocols, especially MySQL.
  static const int32_t kOPUpdate = 2001;
  static const int32_t kOPInsert = 2002;
  static const int32_t kReserved = 2003;
  static const int32_t kOPQuery = 2004;
  static const int32_t kOPGetMore = 2005;
  static const int32_t kOPDelete = 2006;
  static const int32_t kOPKillCursors = 2007;
  static const int32_t kOPCompressed = 2012;
  static const int32_t kOPMsg = 2013;

  static const int32_t kMongoHeaderLength = 16;

  if (count < kMongoHeaderLength) {
    return kUnknown;
  }

  int32_t *buf4 = (int32_t *)buf;
  int32_t message_length = buf4[0];

  if (message_length < kMongoHeaderLength) {
    return kUnknown;
  }

  int32_t request_id = buf4[1];

  if (request_id < 0) {
    return kUnknown;
  }

  int32_t response_to = buf4[2];
  int32_t opcode = buf4[3];

  if (opcode == kOPUpdate || opcode == kOPInsert || opcode == kReserved ||
      opcode == kOPQuery || opcode == kOPGetMore || opcode == kOPDelete ||
      opcode == kOPKillCursors || opcode == kOPCompressed || opcode == kOPMsg) {
    if (response_to == 0) {
      return kRequest;
    }
  }

  return kUnknown;
}

// TODO(yzhao): This is for initial development use. Later we need to combine
// with more inference code, as the startup message only appears at the
// beginning of the exchanges between PostgreSQL client and server.
static __inline enum message_type_t infer_pgsql_startup_message(const char *buf,
                                                                size_t count) {
  // Length field: int32, protocol version field: int32, "user" string, 4 bytes.
  const int kMinMsgLen = 4 + 4 + 4;
  if (count < kMinMsgLen) {
    return kUnknown;
  }

  // Assume startup message wont be larger than 10240 (10KiB).
  const int kMaxMsgLen = 10240;
  const int32_t length = read_big_endian_int32(buf);
  if (length < kMinMsgLen) {
    return kUnknown;
  }
  if (length > kMaxMsgLen) {
    return kUnknown;
  }

  const char kPgsqlVer30[] = "\x00\x03\x00\x00";
  if (px_bpf_strncmp((const char *)buf + 4, 4, kPgsqlVer30) != 0) {
    return kUnknown;
  }

  // Next we expect a key like "user", "datestyle" or "extra_float_digits".
  // For inference purposes, we simply look for a short sequence of alphabetic
  // characters.
  for (int i = 0; i < 3; ++i) {
    // Loosely check for an alphabetic character.
    // This is a loose check and still covers some non alphabetic characters
    // (e.g. `\`), but we want to keep the BPF instruction count low.
    if (*((const char *)buf + 8 + i) < 'A') {
      return kUnknown;
    }
  }

  return kRequest;
}

// Regular message format: | byte tag | int32_t len | string payload |
static __inline enum message_type_t infer_pgsql_query_message(const char *buf,
                                                              size_t count) {
  const uint8_t kTagQ = 'Q';
  if (*buf != kTagQ) {
    return kUnknown;
  }
  const int32_t len = read_big_endian_int32(buf + 1);
  // The length field include the field itself of 4 bytes. Also the minimal size
  // command is COPY/MOVE. The minimal length is therefore 8.
  const int32_t kMinPayloadLen = 8;
  // Assume typical query message size is below an artificial limit.
  // 30000 is copied from postgres code base:
  // https://github.com/postgres/postgres/tree/master/src/interfaces/libpq/fe-protocol3.c#L94
  const int32_t kMaxPayloadLen = 30000;
  if (len < kMinPayloadLen || len > kMaxPayloadLen) {
    return kUnknown;
  }
  // If the input includes a whole message (1 byte tag + length), check the last
  // character.
  if ((len + 1 <= (int)count) && (buf[len] != '\0')) {
    return kUnknown;
  }
  return kRequest;
}

// TODO(yzhao): ReadyForQuery message could be nice pattern to check, as it has
// 6 bytes of fixed bit pattern, plus one byte of enum with possible values 'I',
// 'E', 'T'.  But it's usually sent as a suffix of a query response, so it's
// difficult to capture. Research more to see if we can detect this message.

static __inline enum message_type_t infer_pgsql_regular_message(const char *buf,
                                                                size_t count) {
  const int kMinMsgLen = 1 + sizeof(int32_t);
  if (count < kMinMsgLen) {
    return kUnknown;
  }
  return infer_pgsql_query_message(buf, count);
}

static __inline enum message_type_t infer_pgsql_message(const char *buf,
                                                        size_t count) {
  enum message_type_t type = infer_pgsql_startup_message(buf, count);
  if (type != kUnknown) {
    return type;
  }
  return infer_pgsql_regular_message(buf, count);
}

#define PX_AF_UNKNOWN 0xff

const char kControlMapName[] = "control_map";
const char kControlValuesArrayName[] = "control_values";

const int64_t kTraceAllTGIDs = -1;

// Note: A value of 100 results in >4096 BPF instructions, which is too much for
// older kernels.
#define CONN_CLEANUP_ITERS 85
const int kMaxConnMapCleanupItems = CONN_CLEANUP_ITERS;

union sockaddr_t {
  struct sockaddr sa;
  struct sockaddr_in in4;
  struct sockaddr_in6 in6;
};

// This struct contains information collected when a connection is established,
// via an accept() syscall.
struct conn_info_t {
  // Connection identifier (PID, FD, etc.).
  struct conn_id_t conn_id;

  // IP address of the remote endpoint.
  union sockaddr_t addr;

  // The protocol of traffic on the connection (HTTP, MySQL, etc.).
  enum traffic_protocol_t protocol;

  // Classify traffic as requests, responses or mixed.
  enum endpoint_role_t role;

  // Whether the connection uses SSL.
  bool ssl;

  // The number of bytes written/read on this connection.
  int64_t wr_bytes;
  int64_t rd_bytes;

  // The previously reported values of bytes written/read.
  // Used for determining when to send updated conn_stats values.
  int64_t last_reported_bytes;

  // The number of bytes written by application (for uprobe) on this connection.
  int64_t app_wr_bytes;
  // The number of bytes read by application (for uprobe) on this connection.
  int64_t app_rd_bytes;

  // Some stats for protocol inference. Used for threshold-based filtering.
  //
  // How many times the data segments have been classified as the designated
  // protocol.
  int32_t protocol_match_count;
  // How many times traffic inference has been applied on this connection.
  int32_t protocol_total_count;

  // Keep the header of the last packet suspected to be MySQL/Kafka. MySQL/Kafka
  // server does 2 separate read syscalls, first to read the header, and second
  // the body of the packet. Thus, we keep a state. (MySQL): Length(3 bytes) +
  // seq_number(1 byte). (Kafka): Length(4 bytes)
  size_t prev_count;
  char prev_buf[4];
  bool prepend_length_header;
};

// This struct is a subset of conn_info_t. It is used to communicate
// connect/accept events. See conn_info_t for descriptions of the members.
struct conn_event_t {
  union sockaddr_t addr;
  enum endpoint_role_t role;
};

// This struct is a subset of conn_info_t. It is used to communicate close
// events. See conn_info_t for descriptions of the members.
struct close_event_t {
  // The number of bytes written and read at time of close.
  int64_t wr_bytes;
  int64_t rd_bytes;
};

// Data buffer message size. BPF can submit at most this amount of data to a
// perf buffer.
//
// NOTE: This size does not directly affect the size of perf buffer submits, as
// the actual data submitted to perf buffers are determined by attr.msg_size. In
// cases where socket_data_event_t is defined as stack variable, the size can be
// problematic. Currently we only have a few instances in *_test.cc files.
//
// TODO(yzhao): We do not yet have a good sense of the desired size. Things to
// consider:
// * Overhead. This single instance is small. However, we should consider this
// in the context of all possible overhead in BPF program.
// * Complexity. If this buffer is not sufficiently large. We'll need to handle
// chunked message inside user space parsing code. ATM, we saw in one case, when
// gRPC reflection RPC itself is invoked, it can send one FileDescriptorProto
// [1], which often become large. That's the only data point we have right now.
//
// NOTES:
// * Kernel size limit is 32KiB. See https://github.com/iovisor/bcc/issues/2519
// for more details.
//
// [1]
// https://github.com/grpc/grpc-go/blob/master/reflection/serverreflection.go
#define MAX_MSG_SIZE 30720 // 30KiB

// This defines how many chunks a perf_submit can support.
// This applies to messages that are over MAX_MSG_SIZE,
// and effectively makes the maximum message size to be
// CHUNK_LIMIT*MAX_MSG_SIZE.
#define CHUNK_LIMIT 4

// Unique ID to all syscalls and a few other notable functions.
// This applies to events sent to user-space.
enum source_function_t {
  kSourceFunctionUnknown,

  // For syscalls.
  kSyscallAccept,
  kSyscallConnect,
  kSyscallClose,
  kSyscallWrite,
  kSyscallRead,
  kSyscallSend,
  kSyscallRecv,
  kSyscallSendTo,
  kSyscallRecvFrom,
  kSyscallSendMsg,
  kSyscallRecvMsg,
  kSyscallSendMMsg,
  kSyscallRecvMMsg,
  kSyscallWriteV,
  kSyscallReadV,
  kSyscallSendfile,

  // For Go TLS libraries.
  kGoTLSConnWrite,
  kGoTLSConnRead,

  // For SSL libraries.
  kSSLWrite,
  kSSLRead,
};

struct socket_data_event_t {
  // We split attributes into a separate struct, because BPF gets upset if you
  // do lots of size arithmetic. This makes it so that it's attributes followed
  // by message.
  struct attr_t {
    // The timestamp when syscall completed (return probe was triggered).
    uint64_t timestamp_ns;

    // Connection identifier (PID, FD, etc.).
    struct conn_id_t conn_id;

    // The protocol of traffic on the connection (HTTP, MySQL, etc.).
    enum traffic_protocol_t protocol;

    // The server-client role.
    enum endpoint_role_t role;

    // The type of the actual data that the msg field encodes, which is used by
    // the caller to determine how to interpret the data.
    enum traffic_direction_t direction;

    // Whether the traffic was collected from an encrypted channel.
    bool ssl;

    // Represents the syscall or function that produces this event.
    enum source_function_t source_fn;

    // A 0-based position number for this event on the connection, in terms of
    // byte position. The position is for the first byte of this message. Note
    // that write/send have separate sequences than read/recv.
    uint64_t pos;

    // The size of the original message. We use this to truncate msg field to
    // minimize the amount of data being transferred.
    uint32_t msg_size;

    // The amount of data actually being sent to user space. This may be less
    // than msg_size if data had to be truncated, or if the data was stripped
    // because we only want to send metadata (e.g. if the connection data
    // tracking has been disabled).
    uint32_t msg_buf_size;

    // Whether to prepend length header to the buffer for messages first
    // inferred as Kafka. MySQL may also use this in this future. See
    // infer_kafka_message in protocol_inference.h for details.
    bool prepend_length_header;
    uint32_t length_header;
  } attr;
  char msg[MAX_MSG_SIZE];
};

#define CONN_OPEN (1 << 0)
#define CONN_CLOSE (1 << 1)

struct conn_stats_event_t {
  // The timestamp of the stats event.
  uint64_t timestamp_ns;

  struct conn_id_t conn_id;

  // IP address of the remote endpoint.
  union sockaddr_t addr;

  // The server-client role.
  enum endpoint_role_t role;

  // The number of bytes written on this connection.
  int64_t wr_bytes;
  // The number of bytes read on this connection.
  int64_t rd_bytes;

  // Bitmask of flags specifying whether conn open or close have been observed.
  uint32_t conn_events;
};

enum control_event_type_t {
  kConnOpen,
  kConnClose,
};

struct socket_control_event_t {
  enum control_event_type_t type;
  uint64_t timestamp_ns;
  struct conn_id_t conn_id;

  // Represents the syscall or function that produces this event.
  enum source_function_t source_fn;

  union {
    struct conn_event_t open;
    struct close_event_t close;
  };
};

struct connect_args_t {
  const struct sockaddr *addr;
  int32_t fd;
};

struct accept_args_t {
  struct sockaddr *addr;
  struct socket *sock_alloc_socket;
};

struct data_args_t {
  // Represents the function from which this argument group originates.
  enum source_function_t source_fn;

  // Did the data event call sock_sendmsg/sock_recvmsg.
  // Used to filter out read/write and readv/writev calls that are not to
  // sockets.
  bool sock_event;

  int32_t fd;

  // For send()/recv()/write()/read().
  const char *buf;

  // For sendmsg()/recvmsg()/writev()/readv().
  const struct iovec *iov;
  size_t iovlen;

  // For sendmmsg()
  unsigned int *msg_len;
};

struct close_args_t {
  int32_t fd;
};

struct sendfile_args_t {
  int32_t out_fd;
  int32_t in_fd;
  size_t count;
};

// MySQL packet:
//      0         8        16        24        32
//      +---------+---------+---------+---------+
//      |        payload_length       | seq_id  |
//      +---------+---------+---------+---------+
//      |                                       |
//      .            ...  body ...              .
//      .                                       .
//      .                                       .
//      +----------------------------------------
// TODO(oazizi/yzhao): This produces too many false positives. Add stronger
// protocol detection.
static __inline enum message_type_t
infer_mysql_message(const char *buf, size_t count,
                    struct conn_info_t *conn_info) {
  static const uint8_t kComQuery = 0x03;
  static const uint8_t kComConnect = 0x0b;
  static const uint8_t kComStmtPrepare = 0x16;
  static const uint8_t kComStmtExecute = 0x17;
  static const uint8_t kComStmtClose = 0x19;

  // Second statement checks whether suspected header matches the length of
  // current packet.
  bool use_prev_buf = (conn_info->prev_count == 4) &&
                      (*((uint32_t *)conn_info->prev_buf) == count);

  if (use_prev_buf) {
    // Check the header_state to find out if the header has been read. MySQL
    // server tends to read in the 4 byte header and the rest of the packet in a
    // separate read.
    count += 4;
  }

  // MySQL packets start with a 3-byte packet length and a 1-byte packet number.
  // The 5th byte on a request contains a command that tells the type.
  if (count < 5) {
    return kUnknown;
  }

  // Convert 3-byte length to uint32_t. But since the 4th byte is supposed to be
  // \x00, directly casting 4-bytes is correct. NOLINTNEXTLINE:
  // readability/casting
  uint32_t len =
      use_prev_buf ? *((uint32_t *)conn_info->prev_buf) : *((uint32_t *)buf);
  len = len & 0x00ffffff;

  uint8_t seq = use_prev_buf ? conn_info->prev_buf[3] : buf[3];
  uint8_t com = use_prev_buf ? buf[0] : buf[4];

  // The packet number of a request should always be 0.
  if (seq != 0) {
    return kUnknown;
  }

  // No such thing as a zero-length request in MySQL protocol.
  if (len == 0) {
    return kUnknown;
  }

  // Assuming that the length of a request is less than 10k characters to avoid
  // false positive flagging as MySQL, which statistically happens frequently
  // for a single-byte check.
  if (len > 10000) {
    return kUnknown;
  }

  // TODO(oazizi): Consider adding more commands (0x00 to 0x1f).
  // Be careful, though: trade-off is higher rates of false positives.
  if (com == kComConnect || com == kComQuery || com == kComStmtPrepare ||
      com == kComStmtExecute || com == kComStmtClose) {
    return kRequest;
  }
  return kUnknown;
}

// Reference: https://kafka.apache.org/protocol.html#protocol_messages
// Request Header v0 => request_api_key request_api_version correlation_id
//     request_api_key => INT16
//     request_api_version => INT16
//     correlation_id => INT32
static __inline enum message_type_t infer_kafka_request(const char *buf) {
  // API is Kafka's terminology for opcode.
  static const int kNumAPIs = 62;
  static const int kMaxAPIVersion = 12;

  const int16_t request_API_key = read_big_endian_int16(buf);
  if (request_API_key < 0 || request_API_key > kNumAPIs) {
    return kUnknown;
  }

  const int16_t request_API_version = read_big_endian_int16(buf + 2);
  if (request_API_version < 0 || request_API_version > kMaxAPIVersion) {
    return kUnknown;
  }

  const int32_t correlation_id = read_big_endian_int32(buf + 4);
  if (correlation_id < 0) {
    return kUnknown;
  }
  return kRequest;
}

static __inline enum message_type_t
infer_kafka_message(const char *buf, size_t count,
                    struct conn_info_t *conn_info) {
  // Second statement checks whether suspected header matches the length of
  // current packet. This shouldn't confuse with MySQL because MySQL uses little
  // endian, and Kafka uses big endian.
  bool use_prev_buf =
      (conn_info->prev_count == 4) &&
      ((size_t)read_big_endian_int32(conn_info->prev_buf) == count);

  if (use_prev_buf) {
    count += 4;
  }

  // length(4 bytes) + api_key(2 bytes) + api_version(2 bytes) +
  // correlation_id(4 bytes)
  static const int kMinRequestLength = 12;
  if (count < kMinRequestLength) {
    return kUnknown;
  }

  const int32_t message_size =
      use_prev_buf ? count : read_big_endian_int32(buf) + 4;

  // Enforcing count to be exactly message_size + 4 to mitigate
  // misclassification. However, this will miss long messages broken into
  // multiple reads.
  if (message_size < 0 || count != (size_t)message_size) {
    return kUnknown;
  }
  const char *request_buf = use_prev_buf ? buf : buf + 4;
  enum message_type_t result = infer_kafka_request(request_buf);

  // Kafka servers read in a 4-byte packet length header first. The first packet
  // in the stream is used to infer protocol, but the header has already been
  // read. One solution is to add another perf_submit of the 4-byte header, but
  // this would impact the instruction limit. Not handling this case causes
  // potential confusion in the parsers. Instead, we set a prepend_length_header
  // field if and only if Kafka has just been inferred for the first time under
  // the scenario described above. Length header is appended to user the buffer
  // in user space.
  if (use_prev_buf && result == kRequest &&
      conn_info->protocol == kProtocolUnknown) {
    conn_info->prepend_length_header = true;
  }
  return result;
}

// Const Reference: https://www.rabbitmq.com/resources/specs/amqp0-9-1.xml
// Frame breakdown Ref: https://www.rabbitmq.com/resources/specs/amqp0-9-1.pdf
static __inline enum message_type_t infer_amqp_message(const char *rbuf,
                                                       size_t count) {
  static const uint16_t kConnectionClass = 10;
  static const uint16_t kBasicClass = 60;

  static const uint16_t kMethodConnectionStart = 10;
  static const uint16_t kMethodConnectionStartOk = 11;
  static const uint16_t kMethodBasicPublish = 40;
  static const uint16_t kMethodBasicDeliver = 60;

  static const uint8_t kFrameMethodType = 1;
  static const uint8_t kMinFrameLength = 8;
  if (count < kMinFrameLength) {
    return kUnknown;
  }

  const uint8_t *buf = (const uint8_t *)rbuf;
  uint8_t frame_type = buf[0];
  // Check only for types Connection Start/Start-OK. Publish/Deliver
  if (frame_type != kFrameMethodType) {
    return kUnknown;
  }

  uint16_t class_id = read_big_endian_int16(rbuf + 7);
  uint16_t method_id = read_big_endian_int16(rbuf + 9);
  // ConnectionStart, ConnectionStartOk, BasicPublish, BasicDeliver are the most
  // likely methods to consider
  if (class_id == kConnectionClass && method_id == kMethodConnectionStart) {
    return kRequest;
  }
  if (class_id == kConnectionClass && method_id == kMethodConnectionStartOk) {
    return kResponse;
  }

  if (class_id == kBasicClass && method_id == kMethodBasicPublish) {
    return kRequest;
  }
  if (class_id == kBasicClass && method_id == kMethodBasicDeliver) {
    return kResponse;
  }

  return kUnknown;
}

static __inline enum message_type_t infer_dns_message(const char *buf,
                                                      size_t count) {
  const int kDNSHeaderSize = 12;

  // Use the maximum *guaranteed* UDP packet size as the max DNS message size.
  // UDP packets can be larger, but this is the typical maximum size for DNS.
  const int kMaxDNSMessageSize = 512;

  // Maximum number of resource records.
  // https://stackoverflow.com/questions/6794926/how-many-a-records-can-fit-in-a-single-dns-response
  const int kMaxNumRR = 25;

  if (count < kDNSHeaderSize || count > kMaxDNSMessageSize) {
    return kUnknown;
  }

  const uint8_t *ubuf = (const uint8_t *)buf;

  uint16_t flags = (ubuf[2] << 8) + ubuf[3];
  uint16_t num_questions = (ubuf[4] << 8) + ubuf[5];
  uint16_t num_answers = (ubuf[6] << 8) + ubuf[7];
  uint16_t num_auth = (ubuf[8] << 8) + ubuf[9];
  uint16_t num_addl = (ubuf[10] << 8) + ubuf[11];

  bool qr = (flags >> 15) & 0x1;
  uint8_t opcode = (flags >> 11) & 0xf;
  uint8_t zero = (flags >> 6) & 0x1;

  if (zero != 0) {
    return kUnknown;
  }

  if (opcode != 0) {
    return kUnknown;
  }

  if (num_questions == 0 || num_questions > 10) {
    return kUnknown;
  }

  uint32_t num_rr = num_questions + num_answers + num_auth + num_addl;
  if (num_rr > kMaxNumRR) {
    return kUnknown;
  }

  return (qr == 0) ? kRequest : kResponse;
}

// Redis request and response messages share the same format.
// See https://redis.io/topics/protocol for the REDIS protocol spec.
//
// TODO(yzhao): Apply simplified parsing to read the content to distinguished
// request & response.
static __inline bool is_redis_message(const char *buf, size_t count) {
  // Redis messages start with an one-byte type marker, and end with \r\n
  // terminal sequence.
  if (count < 3) {
    return false;
  }

  const char first_byte = buf[0];

  if ( // Simple strings start with +
      first_byte != '+' &&
      // Errors start with -
      first_byte != '-' &&
      // Integers start with :
      first_byte != ':' &&
      // Bulk strings start with $
      first_byte != '$' &&
      // Arrays start with *
      first_byte != '*') {
    return false;
  }

  // The last two chars are \r\n, the terminal sequence of all Redis messages.
  if (buf[count - 2] != '\r') {
    return false;
  }
  if (buf[count - 1] != '\n') {
    return false;
  }

  return true;
}

// TODO(ddelnano): Mux protocol traffic is currently misidentified as ssh. Since
// stirling doesn't have ssh support yet, but will need to be addressed. In
// addition, mux seems to send the header and body on its protocol in two
// separate syscalls on the server side.
static __inline enum message_type_t infer_mux_message(const char *buf,
                                                      size_t count) {
  // mux's on the wire format causes false positives for protocol inference
  // In order to address this, we only infer mux messages by the
  // most useful message types and if they are easy to identify
  static const int8_t kTdispatch = 2;
  static const int8_t kRdispatch = -2;
  static const int8_t kTinit = 68;
  static const int8_t kRinit = -68;
  static const int8_t kRerr = -128;
  static const int8_t kRerrOld = 127;
  uint32_t mux_header_size = 8;
  // TODO(ddelnano): Determine why mux-framer text in T/Rinit is
  // 6 bytes after the mux header
  int32_t mux_framer_pos = mux_header_size + 6;

  if (count < mux_header_size) {
    return kUnknown;
  }

  uint32_t length = read_big_endian_int32(buf) + 4;
  enum message_type_t msg_type;

  int32_t type_and_tag = read_big_endian_int32(buf + 4);
  int8_t mux_type = (type_and_tag & 0xff000000) >> 24;
  uint32_t tag = (type_and_tag & 0xffffff);
  switch (mux_type) {
  case kTdispatch:
  case kTinit:
  case kRerrOld:
    msg_type = kRequest;
    break;
  case kRdispatch:
  case kRinit:
  case kRerr:
    msg_type = kResponse;
    break;
  default:
    return kUnknown;
  }

  if (mux_type == kRerr || mux_type == kRerrOld) {
    if (buf[length - 5] != 'c' || buf[length - 4] != 'h' ||
        buf[length - 3] != 'e' || buf[length - 2] != 'c' ||
        buf[length - 1] != 'k')
      return kUnknown;
  }

  if (mux_type == kRinit || mux_type == kTinit) {
    if (buf[mux_framer_pos] != 'm' || buf[mux_framer_pos + 1] != 'u' ||
        buf[mux_framer_pos + 2] != 'x' || buf[mux_framer_pos + 3] != '-' ||
        buf[mux_framer_pos + 4] != 'f' || buf[mux_framer_pos + 5] != 'r' ||
        buf[mux_framer_pos + 6] != 'a' || buf[mux_framer_pos + 7] != 'm' ||
        buf[mux_framer_pos + 8] != 'e' || buf[mux_framer_pos + 9] != 'r')
      return kUnknown;
  }

  if (tag < 1 || tag > ((1 << 23) - 1)) {
    return kUnknown;
  }

  return msg_type;
}

// NATS messages are in texts. The role is inferred from the message type.
// See
// https://github.com/nats-io/docs/blob/master/nats_protocol/nats-protocol.md
//
// In case of bpf instruction count limit becomes a problem, we can drop CONNECT
// and INFO message detection, they are only sent once after establishing the
// connection.
static __inline enum message_type_t infer_nats_message(const char *buf,
                                                       size_t count) {
  // NATS messages start with an one-byte type marker, and end with \r\n
  // terminal sequence.
  if (count < 3) {
    return kUnknown;
  }
  // The last two chars are \r\n, the terminal sequence of all NATS messages.
  if (buf[count - 2] != '\r') {
    return kUnknown;
  }
  if (buf[count - 1] != '\n') {
    return kUnknown;
  }
  if (buf[0] == 'C' && buf[1] == 'O' && buf[2] == 'N' && buf[3] == 'N' &&
      buf[4] == 'E' && buf[5] == 'C' && buf[6] == 'T') {
    // kRequest is not precise. Here only means the message is sent by client.
    return kRequest;
  }
  if (buf[0] == 'S' && buf[1] == 'U' && buf[2] == 'B') {
    return kRequest;
  }
  if (buf[0] == 'U' && buf[1] == 'N' && buf[2] == 'S' && buf[3] == 'U' &&
      buf[4] == 'B') {
    return kRequest;
  }
  if (buf[0] == 'P' && buf[1] == 'U' && buf[2] == 'B') {
    return kRequest;
  }
  if (buf[0] == 'I' && buf[1] == 'N' && buf[2] == 'F' && buf[3] == 'O') {
    // kResponse is not precise. Here only means the message is sent by server.
    return kResponse;
  }
  if (buf[0] == 'M' && buf[1] == 'S' && buf[2] == 'G') {
    return kResponse;
  }
  if (buf[0] == '+' && buf[1] == 'O' && buf[2] == 'K') {
    return kResponse;
  }
  if (buf[0] == '-' && buf[1] == 'E' && buf[2] == 'R' && buf[3] == 'R') {
    return kResponse;
  }
  // PING & PONG can be sent by both client and server. Don't use them.
  return kUnknown;
}

static __inline struct protocol_message_t
infer_protocol(const char *buf, size_t count, struct conn_info_t *conn_info) {
  struct protocol_message_t inferred_message;
  inferred_message.protocol = kProtocolUnknown;
  inferred_message.type = kUnknown;

  // The prepend_length_header controls whether a length header is prepended to
  // the buffer in user space.
  conn_info->prepend_length_header = false;

  // TODO(oazizi): Get rid of `inferred_message.type` and convert the functions
  // below to
  //               is_xyz_message().
  //               This is potentially possible because of the fact that we now
  //               infer connection role by considering which side called
  //               accept() vs connect(). Once the clean-up above is done, the
  //               code below can be turned into a chained ternary.
  // PROTOCOL_LIST: Requires update on new protocols.
  if (ENABLE_HTTP_TRACING &&
      (inferred_message.type = infer_http_message(buf, count)) != kUnknown) {
    inferred_message.protocol = kProtocolHTTP;
  } else if (ENABLE_CQL_TRACING && (inferred_message.type = infer_cql_message(
                                        buf, count)) != kUnknown) {
    inferred_message.protocol = kProtocolCQL;
  } else if (ENABLE_MONGO_TRACING &&
             (inferred_message.type = infer_mongo_message(buf, count)) !=
                 kUnknown) {
    inferred_message.protocol = kProtocolMongo;
  } else if (ENABLE_PGSQL_TRACING &&
             (inferred_message.type = infer_pgsql_message(buf, count)) !=
                 kUnknown) {
    inferred_message.protocol = kProtocolPGSQL;
  } else if (ENABLE_MYSQL_TRACING &&
             (inferred_message.type =
                  infer_mysql_message(buf, count, conn_info)) != kUnknown) {
    inferred_message.protocol = kProtocolMySQL;
  } else if (ENABLE_MUX_TRACING && (inferred_message.type = infer_mux_message(
                                        buf, count)) != kUnknown) {
    inferred_message.protocol = kProtocolMux;
  } else if (ENABLE_KAFKA_TRACING &&
             (inferred_message.type =
                  infer_kafka_message(buf, count, conn_info)) != kUnknown) {
    inferred_message.protocol = kProtocolKafka;
  } else if (ENABLE_DNS_TRACING && (inferred_message.type = infer_dns_message(
                                        buf, count)) != kUnknown) {
    inferred_message.protocol = kProtocolDNS;
  } else if (ENABLE_AMQP_TRACING && (inferred_message.type = infer_amqp_message(
                                         buf, count)) != kUnknown) {
    inferred_message.protocol = kProtocolAMQP;
  } else if (ENABLE_REDIS_TRACING && is_redis_message(buf, count)) {
    // For Redis, the message type is left to be kUnknown.
    // The message types are then inferred via traffic direction and
    // client/server role.
    inferred_message.protocol = kProtocolRedis;
  } else if (ENABLE_NATS_TRACING && (inferred_message.type = infer_nats_message(
                                         buf, count)) != kUnknown) {
    inferred_message.protocol = kProtocolNATS;
  }

  conn_info->prev_count = count;
  if (count == 4) {
    conn_info->prev_buf[0] = buf[0];
    conn_info->prev_buf[1] = buf[1];
    conn_info->prev_buf[2] = buf[2];
    conn_info->prev_buf[3] = buf[3];
  }

  return inferred_message;
}

/* #include
 * "src/stirling/source_connectors/socket_tracer/bcc_bpf_intf/socket_trace.h" */
/* #include "src/stirling/upid/upid.h" */

// This keeps instruction count below BPF's limit of 4096 per probe.
#define LOOP_LIMIT 42
#define PROTOCOL_VEC_LIMIT 3

const int32_t kInvalidFD = -1;

// This is the amount of activity required on a connection before a new
// ConnStats event is reported to user-space. It applies to read and write
// traffic combined.
const int kConnStatsDataThreshold = 65536;

// This is the perf buffer for BPF program to export data from kernel to user
// space.
BPF_PERF_OUTPUT(socket_data_events);
BPF_PERF_OUTPUT(socket_control_events);
BPF_PERF_OUTPUT(conn_stats_events);

// This output is used to export notification of processes that have performed
// an mmap.
BPF_PERF_OUTPUT(mmap_events);

// This control_map is a bit-mask that controls which endpoints are traced in a
// connection. The bits are defined in endpoint_role_t enum, kRoleClient or
// kRoleServer. kRoleUnknown is not really used, but is defined for
// completeness. There is a control map element for each protocol.
BPF_PERCPU_ARRAY(control_map, uint64_t, kNumProtocols);

// Map from user-space file descriptors to the connections obtained from
// accept() syscall. Tracks connection from accept() -> close(). Key is {tgid,
// fd}.
BPF_HASH(conn_info_map, uint64_t, struct conn_info_t, 131072);

// Map to indicate which connections (TGID+FD), user-space has disabled.
// This is tracked separately from conn_info_map to avoid any read-write races.
// This particular map is only written from user-space, and only read from BPF.
// The value is a TSID indicating the last TSID to be disabled. Any newer
// TSIDs should still be pushed out to user space. Events on older TSIDs is not
// possible. Key is {tgid, fd}; Value is TSID.
BPF_HASH(conn_disabled_map, uint64_t, uint64_t);

// Map from thread to its ongoing accept() syscall's input argument.
// Tracks accept() call from entry -> exit.
// Key is {tgid, pid}.
BPF_HASH(active_accept_args_map, uint64_t, struct accept_args_t);

// Map from thread to its ongoing connect() syscall's input argument.
// Tracks connect() call from entry -> exit.
// Key is {tgid, pid}.
BPF_HASH(active_connect_args_map, uint64_t, struct connect_args_t);

// Map from thread to its ongoing write() syscall's input argument.
// Tracks write() call from entry -> exit.
// Key is {tgid, pid}.
BPF_HASH(active_write_args_map, uint64_t, struct data_args_t);

// Map from thread to its ongoing read() syscall's input argument.
// Tracks read() call from entry -> exit.
// Key is {tgid, pid}.
BPF_HASH(active_read_args_map, uint64_t, struct data_args_t);

// Map from thread to its ongoing close() syscall's input argument.
// Tracks close() call from entry -> exit.
// Key is {tgid, pid}.
BPF_HASH(active_close_args_map, uint64_t, struct close_args_t);

// Map from thread to its ongoing sendfile syscall's input argument.
// Tracks sendfile() call from entry -> exit.
// Key is {tgid, pid}.
BPF_HASH(active_sendfile_args_map, uint64_t, struct sendfile_args_t);

// BPF programs are limited to a 512-byte stack. We store this value per CPU
// and use it as a heap allocated value.
BPF_PERCPU_ARRAY(socket_data_event_buffer_heap, struct socket_data_event_t, 1);
BPF_PERCPU_ARRAY(conn_stats_event_buffer_heap, struct conn_stats_event_t, 1);

// This array records singular values that are used by probes. We group them
// together to reduce the number of arrays with only 1 element.
BPF_PERCPU_ARRAY(control_values, int64_t, kNumControlValues);

/***********************************************************
 * General helper functions
 ***********************************************************/

static __inline uint64_t gen_tgid_fd(uint32_t tgid, int fd) {
  return ((uint64_t)tgid << 32) | (uint32_t)fd;
}

static __inline void init_conn_id(uint32_t tgid, int32_t fd,
                                  struct conn_id_t *conn_id) {
  conn_id->upid.tgid = tgid;
  conn_id->upid.start_time_ticks = get_tgid_start_time();
  conn_id->fd = fd;
  conn_id->tsid = bpf_ktime_get_ns();
}

static __inline void init_conn_info(uint32_t tgid, int32_t fd,
                                    struct conn_info_t *conn_info) {
  init_conn_id(tgid, fd, &conn_info->conn_id);
  // NOTE: BCC code defaults to 0, because kRoleUnknown is not 0, must
  // explicitly initialize.
  conn_info->role = kRoleUnknown;
  conn_info->addr.sa.sa_family = PX_AF_UNKNOWN;
}

// Be careful calling this function. The automatic creation of BPF map entries
// can result in a BPF map leak if called on unwanted probes. How do we make
// sure we don't leak then? ConnInfoMapManager.ReleaseResources() will clean-up
// the relevant map entries every time a ConnTracker is destroyed.
static __inline struct conn_info_t *get_or_create_conn_info(uint32_t tgid,
                                                            int32_t fd) {
  uint64_t tgid_fd = gen_tgid_fd(tgid, fd);
  struct conn_info_t new_conn_info = {};
  init_conn_info(tgid, fd, &new_conn_info);
  return conn_info_map.lookup_or_init(&tgid_fd, &new_conn_info);
}

static __inline void set_conn_as_ssl(uint32_t tgid, int32_t fd) {
  struct conn_info_t *conn_info = get_or_create_conn_info(tgid, fd);
  if (conn_info == NULL) {
    return;
  }
  conn_info->ssl = true;
}

static __inline struct socket_data_event_t *
fill_socket_data_event(enum source_function_t src_fn,
                       enum traffic_direction_t direction,
                       const struct conn_info_t *conn_info) {
  uint32_t kZero = 0;
  struct socket_data_event_t *event =
      socket_data_event_buffer_heap.lookup(&kZero);
  if (event == NULL) {
    return NULL;
  }
  event->attr.timestamp_ns = bpf_ktime_get_ns();
  event->attr.source_fn = src_fn;
  event->attr.ssl = conn_info->ssl;
  event->attr.direction = direction;
  event->attr.conn_id = conn_info->conn_id;
  event->attr.protocol = conn_info->protocol;
  event->attr.role = conn_info->role;
  event->attr.pos =
      (direction == kEgress) ? conn_info->wr_bytes : conn_info->rd_bytes;
  event->attr.prepend_length_header = conn_info->prepend_length_header;
  BPF_PROBE_READ_VAR(event->attr.length_header, conn_info->prev_buf);
  return event;
}

static __inline struct conn_stats_event_t *
fill_conn_stats_event(const struct conn_info_t *conn_info) {
  uint32_t kZero = 0;
  struct conn_stats_event_t *event =
      conn_stats_event_buffer_heap.lookup(&kZero);
  if (event == NULL) {
    return NULL;
  }

  event->conn_id = conn_info->conn_id;
  event->addr = conn_info->addr;
  event->role = conn_info->role;
  event->wr_bytes = conn_info->wr_bytes;
  event->rd_bytes = conn_info->rd_bytes;
  event->conn_events = 0;
  event->timestamp_ns = bpf_ktime_get_ns();
  return event;
}

/***********************************************************
 * Trace filtering functions
 ***********************************************************/

static __inline bool should_trace_sockaddr_family(sa_family_t sa_family) {
  // PX_AF_UNKNOWN means we never traced the accept/connect, and we don't know
  // the sockaddr family. Trace these because they *may* be a sockaddr of
  // interest.
  return sa_family == PX_AF_UNKNOWN || sa_family == AF_INET ||
         sa_family == AF_INET6;
}

static __inline bool should_trace_conn(struct conn_info_t *conn_info) {
  // While we keep all sa_family types in conn_info_map,
  // we only send connections on INET or UNKNOWN to user-space.
  // Also, it's very important to send the UNKNOWN cases to user-space,
  // otherwise we may have a BPF map leak from the earlier call to
  // get_or_create_conn_info().
  return should_trace_sockaddr_family(conn_info->addr.sa.sa_family);
}

// If this returns false, we still will trace summary stats.
static __inline bool
should_trace_protocol_data(const struct conn_info_t *conn_info) {
  if (conn_info->protocol == kProtocolUnknown) {
    return false;
  }

  uint32_t protocol = conn_info->protocol;
  uint64_t kZero = 0;
  uint64_t control = *control_map.lookup_or_init(&protocol, &kZero);
  return control & conn_info->role;
}

static __inline bool is_stirling_tgid(const uint32_t tgid) {
  int idx = kStirlingTGIDIndex;
  int64_t *stirling_tgid = control_values.lookup(&idx);
  if (stirling_tgid == NULL) {
    return false;
  }
  return *stirling_tgid == tgid;
}

enum target_tgid_match_result_t {
  TARGET_TGID_UNSPECIFIED,
  TARGET_TGID_ALL,
  TARGET_TGID_MATCHED,
  TARGET_TGID_UNMATCHED,
};

static __inline enum target_tgid_match_result_t
match_trace_tgid(const uint32_t tgid) {

  // TODO 20221126, hardcode by @ArthurChiao
  return TARGET_TGID_MATCHED;

  // TODO(yzhao): Use externally-defined macro to replace BPF_MAP. Since this
  // function is called for all PIDs, this optimization is useful.
  int idx = kTargetTGIDIndex;
  int64_t *target_tgid = control_values.lookup(&idx);
  if (target_tgid == NULL) {
    return TARGET_TGID_UNSPECIFIED;
  }
  if (*target_tgid < 0) {
    // Negative value means trace all.
    return TARGET_TGID_ALL;
  }
  if (*target_tgid == tgid) {
    return TARGET_TGID_MATCHED;
  }
  bpf_trace_printk("*target %ld, tgid %ld\n", *target_tgid, tgid);
  return TARGET_TGID_UNMATCHED;
}

static __inline void update_traffic_class(struct conn_info_t *conn_info,
                                          enum traffic_direction_t direction,
                                          const char *buf, size_t count) {
  if (conn_info == NULL) {
    return;
  }
  conn_info->protocol_total_count += 1;

  // Try to infer connection type (protocol) based on data.
  struct protocol_message_t inferred_protocol =
      infer_protocol(buf, count, conn_info);

  // Could not infer the traffic.
  if (inferred_protocol.protocol == kProtocolUnknown ||
      conn_info->protocol == kProtocolMongo) {
    return;
  }

  // Update protocol if not set.
  if (conn_info->protocol == kProtocolUnknown) {
    conn_info->protocol = inferred_protocol.protocol;
  }

  // Update role if not set.
  if (conn_info->role == kRoleUnknown &&
      // As of 2020-01, Redis protocol detection doesn't implement message type
      // detection. There could be more protocols without message type detection
      // in the future.
      inferred_protocol.type != kUnknown) {
    // Classify Role as XOR between direction and req_resp_type:
    //    direction  req_resp_type  => role
    //    ------------------------------------
    //    kEgress    kRequest       => Client
    //    kEgress    KResponse      => Server
    //    kIngress   kRequest       => Server
    //    kIngress   kResponse      => Client
    conn_info->role =
        ((direction == kEgress) ^ (inferred_protocol.type == kResponse))
            ? kRoleClient
            : kRoleServer;
  }
}

/***********************************************************
 * Perf submit functions
 ***********************************************************/

static __inline void read_sockaddr_kernel(struct conn_info_t *conn_info,
                                          const struct socket *socket) {
  // Use BPF_PROBE_READ_KERNEL_VAR since BCC cannot insert them as expected.
  struct sock *sk = NULL;
  BPF_PROBE_READ_KERNEL_VAR(sk, &socket->sk);

  struct sock_common *sk_common = &sk->__sk_common;
  uint16_t family = -1;
  uint16_t port = -1;

  BPF_PROBE_READ_KERNEL_VAR(family, &sk_common->skc_family);
  BPF_PROBE_READ_KERNEL_VAR(port, &sk_common->skc_dport);

  conn_info->addr.sa.sa_family = family;

  if (family == AF_INET) {
    conn_info->addr.in4.sin_port = port;
    BPF_PROBE_READ_KERNEL_VAR(conn_info->addr.in4.sin_addr.s_addr,
                              &sk_common->skc_daddr);
  } else if (family == AF_INET6) {
    conn_info->addr.in6.sin6_port = port;
    BPF_PROBE_READ_KERNEL_VAR(conn_info->addr.in6.sin6_addr,
                              &sk_common->skc_v6_daddr);
  }
}

static __inline void submit_new_conn(struct pt_regs *ctx, uint32_t tgid,
                                     int32_t fd, const struct sockaddr *addr,
                                     const struct socket *socket,
                                     enum endpoint_role_t role,
                                     enum source_function_t source_fn) {
  struct conn_info_t conn_info = {};
  init_conn_info(tgid, fd, &conn_info);
  if (addr != NULL) {
    conn_info.addr = *((union sockaddr_t *)addr);
  } else if (socket != NULL) {
    read_sockaddr_kernel(&conn_info, socket);
  }
  conn_info.role = role;

  uint64_t tgid_fd = gen_tgid_fd(tgid, fd);
  conn_info_map.update(&tgid_fd, &conn_info);

  // While we keep all sa_family types in conn_info_map,
  // we only send connections with supported protocols to user-space.
  // We use the same filter function to avoid sending data of unwanted
  // connections as well.
  if (!should_trace_sockaddr_family(conn_info.addr.sa.sa_family)) {
    return;
  }

  struct socket_control_event_t control_event = {};
  control_event.type = kConnOpen;
  control_event.timestamp_ns = bpf_ktime_get_ns();
  control_event.conn_id = conn_info.conn_id;
  control_event.source_fn = source_fn;
  control_event.open.addr = conn_info.addr;
  control_event.open.role = conn_info.role;

  socket_control_events.perf_submit(ctx, &control_event,
                                    sizeof(struct socket_control_event_t));
}

static __inline void submit_close_event(struct pt_regs *ctx,
                                        struct conn_info_t *conn_info,
                                        enum source_function_t source_fn) {
  struct socket_control_event_t control_event = {};
  control_event.type = kConnClose;
  control_event.timestamp_ns = bpf_ktime_get_ns();
  control_event.conn_id = conn_info->conn_id;
  control_event.source_fn = source_fn;
  control_event.close.rd_bytes = conn_info->rd_bytes;
  control_event.close.wr_bytes = conn_info->wr_bytes;

  socket_control_events.perf_submit(ctx, &control_event,
                                    sizeof(struct socket_control_event_t));
  bpf_trace_printk("Submit one close event, %d\n", sizeof(control_event));
}

// Writes the input buf to event, and submits the event to the corresponding
// perf buffer. Returns the bytes output from the input buf. Note that is not
// the total bytes submitted to the perf buffer, which includes additional
// metadata.
static __inline void perf_submit_buf(struct pt_regs *ctx,
                                     const enum traffic_direction_t direction,
                                     const char *buf, size_t buf_size,
                                     struct conn_info_t *conn_info,
                                     struct socket_data_event_t *event) {
  // Record original size of packet. This may get truncated below before submit.
  event->attr.msg_size = buf_size;

  // This rest of this function has been written carefully to keep the BPF
  // verifier happy in older kernels, so please take care when modifying.
  //
  // Logically, what we'd like is the following:
  //    size_t msg_size = buf_size < sizeof(event->msg) ? buf_size :
  //    sizeof(event->msg); bpf_probe_read(&event->msg, msg_size, buf);
  //    event->attr.msg_size = buf_size;
  //    socket_data_events.perf_submit(ctx, event, size_to_submit);
  //
  // But this does not work in kernel versions 4.14 or older, for various
  // reasons:
  //  1) the verifier does not like a bpf_probe_read with size 0.
  //       - Useful link:
  //       https://www.mail-archive.com/netdev@vger.kernel.org/msg199918.html
  //  2) the verifier does not like a perf_submit that is larger than
  //  sizeof(event).
  //
  // While it is often obvious to us humans that these are not problems,
  // the older verifiers can't prove it to themselves.
  //
  // We often try to provide hints to the verifier using approaches like
  // 'if (msg_size > 0)' around the code, but it turns out that clang is often
  // smarter than the verifier, and optimizes away the structural hints we try
  // to provide the verifier.
  //
  // Solution below involves using a volatile asm statement to prevent clang
  // from optimizing away certain code, so that code can reach the BPF verifier,
  // and convince it that everything is safe.
  //
  // Tested to work on the following kernels:
  //   4.14

  if (buf_size == 0) {
    return;
  }

  // Note that buf_size_minus_1 will be positive due to the if-statement above.
  size_t buf_size_minus_1 = buf_size - 1;

  // Clang is too smart for us, and tries to remove some of the obvious hints we
  // are leaving for the BPF verifier. So we add this NOP volatile statement, so
  // clang can't optimize away some of our if-statements below. By telling clang
  // that buf_size_minus_1 is both an input and output to some black box
  // assembly code, clang has to discard any assumptions on what values this
  // variable can take.
  asm volatile("" : "+r"(buf_size_minus_1) :);

  buf_size = buf_size_minus_1 + 1;

  // 4.14 kernels reject bpf_probe_read with size that they may think is zero.
  // Without the if statement, it somehow can't reason that the bpf_probe_read
  // is non-zero.
  size_t amount_copied = 0;
  if (buf_size_minus_1 < MAX_MSG_SIZE) {
    bpf_probe_read(&event->msg, buf_size, buf);
    amount_copied = buf_size;
  } else if (buf_size_minus_1 < 0x7fffffff) {
    // If-statement condition above is only required to prevent clang from
    // optimizing away the `if (amount_copied > 0)` below.
    bpf_probe_read(&event->msg, MAX_MSG_SIZE, buf);
    amount_copied = MAX_MSG_SIZE;
  }

  // If-statement is redundant, but is required to keep the 4.14 verifier happy.
  if (amount_copied > 0) {
    event->attr.msg_buf_size = amount_copied;
    socket_data_events.perf_submit(ctx, event,
                                   sizeof(event->attr) + amount_copied);
  }
}

static __inline void
perf_submit_wrapper(struct pt_regs *ctx,
                    const enum traffic_direction_t direction, const char *buf,
                    const size_t buf_size, struct conn_info_t *conn_info,
                    struct socket_data_event_t *event) {
  int bytes_sent = 0;
  unsigned int i;

#pragma unroll
  for (i = 0; i < CHUNK_LIMIT; ++i) {
    const int bytes_remaining = buf_size - bytes_sent;
    const size_t current_size =
        (bytes_remaining > MAX_MSG_SIZE && (i != CHUNK_LIMIT - 1))
            ? MAX_MSG_SIZE
            : bytes_remaining;
    perf_submit_buf(ctx, direction, buf + bytes_sent, current_size, conn_info,
                    event);
    bytes_sent += current_size;

    // Move the position for the next event.
    event->attr.pos += current_size;
  }
}

static __inline void perf_submit_iovecs(
    struct pt_regs *ctx, const enum traffic_direction_t direction,
    const struct iovec *iov, const size_t iovlen, const size_t total_size,
    struct conn_info_t *conn_info, struct socket_data_event_t *event) {
  // NOTE: The syscalls for scatter buffers, {send,recv}msg()/{write,read}v(),
  // access buffers in array order. That means they read or fill iov[0], then
  // iov[1], and so on. They return the total size of the written or read data.
  // Therefore, when loop through the buffers, both the number of buffers and
  // the total size need to be checked. More details can be found on their man
  // pages.
  int bytes_sent = 0;
#pragma unroll
  for (int i = 0; i < LOOP_LIMIT && i < iovlen && bytes_sent < total_size;
       ++i) {
    struct iovec iov_cpy;
    BPF_PROBE_READ_VAR(iov_cpy, &iov[i]);

    const int bytes_remaining = total_size - bytes_sent;
    const size_t iov_size = min_size_t(iov_cpy.iov_len, bytes_remaining);

    // TODO(oazizi/yzhao): Should switch this to go through perf_submit_wrapper.
    //                     We don't have the BPF instruction count to do so
    //                     right now.
    perf_submit_buf(ctx, direction, iov_cpy.iov_base, iov_size, conn_info,
                    event);
    bytes_sent += iov_size;

    // Move the position for the next event.
    event->attr.pos += iov_size;
  }

  // TODO(oazizi): If there is data left after the loop limit, we should still
  // report the remainder
  //               with a data-less event.
}

/***********************************************************
 * Map cleanup functions
 ***********************************************************/

int conn_cleanup_uprobe(struct pt_regs *ctx) {
  int n = (int)PT_REGS_PARM1(ctx);
  struct conn_id_t *conn_id_list = (struct conn_id_t *)PT_REGS_PARM2(ctx);

#pragma unroll
  for (int i = 0; i < CONN_CLEANUP_ITERS; ++i) {
    struct conn_id_t conn_id = conn_id_list[i];

    // Moving this break above or into the for loop causes us to breach the BPF
    // instruction count. Has not been investigated. Just keep it here for now.
    if (i >= n) {
      break;
    }

    uint64_t tgid_fd = gen_tgid_fd(conn_id.upid.tgid, conn_id.fd);

    // Before deleting, make sure we have the correct generation by checking the
    // TSID. We don't want to accidentally delete a newer generation that has
    // since come into existence.

    struct conn_info_t *conn_info = conn_info_map.lookup(&tgid_fd);
    if (conn_info != NULL && conn_info->conn_id.tsid == conn_id.tsid) {
      conn_info_map.delete(&tgid_fd);
    }

    uint64_t *tsid = conn_disabled_map.lookup(&tgid_fd);
    if (tsid != NULL && *tsid == conn_id.tsid) {
      conn_disabled_map.delete(&tgid_fd);
    }
  }

  return 0;
}

/***********************************************************
 * BPF syscall processing functions
 ***********************************************************/

// Table of what events to send to user-space:
//
// SockAddr   | Protocol   ||  Connect/Accept   |   Data      | Close
// -----------|------------||-------------------|-------------|-------
// INET       | Unknown    ||  Yes              |   Summary   | Yes
// INET       | Known      ||  N/A              |   Full      | Yes
// Other      | Unknown    ||  No               |   No        | No
// Other      | Known      ||  N/A              |   No        | No
// Unknown    | Unknown    ||  No*              |   Summary   | Yes
// Unknown    | Known      ||  N/A              |   Full      | Yes
//
// *: Only applicable to accept() syscalls where addr is nullptr. We won't know
// the remote addr.
//    Since no useful information is traced, just skip it. Will be treated as a
//    case where we missed the accept.

static __inline void
process_syscall_connect(struct pt_regs *ctx, uint64_t id,
                        const struct connect_args_t *args) {
  uint32_t tgid = id >> 32;
  int ret_val = PT_REGS_RC(ctx);

  if (match_trace_tgid(tgid) == TARGET_TGID_UNMATCHED) {
    return;
  }

  if (args->fd < 0) {
    return;
  }

  // We allow EINPROGRESS to go through, which indicates that a NON_BLOCK socket
  // is undergoing handshake.
  //
  // In case connect() eventually fails, any write or read on the fd would fail
  // nonetheless, and we won't see spurious events.
  //
  // In case a separate connect() is called concurrently in another thread, and
  // succeeds immediately, any write or read on the fd would be attributed to
  // the new connection.
  if (ret_val < 0 && ret_val != -EINPROGRESS) {
    return;
  }

  submit_new_conn(ctx, tgid, args->fd, args->addr, /*socket*/ NULL, kRoleClient,
                  kSyscallConnect);
}

static __inline void process_syscall_accept(struct pt_regs *ctx, uint64_t id,
                                            const struct accept_args_t *args) {
  uint32_t tgid = id >> 32;
  int ret_fd = PT_REGS_RC(ctx);

  if (match_trace_tgid(tgid) == TARGET_TGID_UNMATCHED) {
    return;
  }

  if (ret_fd < 0) {
    return;
  }

  submit_new_conn(ctx, tgid, ret_fd, args->addr, args->sock_alloc_socket,
                  kRoleServer, kSyscallAccept);
}

// TODO(oazizi): This is badly broken (but better than before).
//               Suppose a server with a UDP socket has the following sequence:
//                 recvmsg(/*sockfd*/ 5, /*msgaddr*/ A);
//                 recvmsg(/*sockfd*/ 5, /*msgaddr*/ B);
//                 sendmsg(/*sockfd*/ 5, /*msgaddr*/ B);
//                 sendmsg(/*sockfd*/ 5, /*msgaddr*/ A);
//
// This function will produce incorrect results, because it will never register
// B. Everything will be attributed to the first address recorded on the socket.
//
// Note that even if we record address changes, the sequence above will
// not be correct for the last sendmsg in the sequence above.
//
// Problem is our ConnTracker model is not suitable for UDP, where there is no
// connection. For a TCP server, accept() sets the remote address, and all
// messages on that socket are to/from that remote address. For a UDP server,
// there is no such thing. Every datagram has an address specified with it. If
// we try to record and submit the "connection", then it may not be the right
// remote endpoint for all messages on that socket.
//
// In this example, process_implicit_conn() will get triggered on the first
// recvmsg, and then everything on sockfd=5 will assume to be on that
// address...which is clearly wrong.
static __inline void process_implicit_conn(struct pt_regs *ctx, uint64_t id,
                                           const struct connect_args_t *args,
                                           enum source_function_t source_fn) {
  uint32_t tgid = id >> 32;

  if (match_trace_tgid(tgid) == TARGET_TGID_UNMATCHED) {
    return;
  }

  if (args->fd < 0) {
    return;
  }

  uint64_t tgid_fd = gen_tgid_fd(tgid, args->fd);

  struct conn_info_t *conn_info = conn_info_map.lookup(&tgid_fd);
  if (conn_info != NULL) {
    return;
  }

  submit_new_conn(ctx, tgid, args->fd, args->addr, /*socket*/ NULL,
                  kRoleUnknown, source_fn);
}

static __inline bool should_send_data(uint32_t tgid,
                                      uint64_t conn_disabled_tsid,
                                      bool force_trace_tgid,
                                      struct conn_info_t *conn_info) {
  // Never trace stirling.
  if (is_stirling_tgid(tgid)) {
    return false;
  }

  // Never trace any connections that user-space has asked us to disable.
  if (conn_info->conn_id.tsid <= conn_disabled_tsid) {
    return false;
  }

  // Only trace data for protocols of interest, or if forced on.
  return (force_trace_tgid || should_trace_protocol_data(conn_info));
}

static __inline void update_conn_stats(struct pt_regs *ctx,
                                       struct conn_info_t *conn_info,
                                       enum traffic_direction_t direction,
                                       ssize_t bytes_count) {
  // Update state of the connection.
  switch (direction) {
  case kEgress:
    conn_info->wr_bytes += bytes_count;
    break;
  case kIngress:
    conn_info->rd_bytes += bytes_count;
    break;
  }

  // Only send event if there's been enough of a change.
  // TODO(oazizi): Add elapsed time since last send as a triggering condition
  // too.
  uint64_t total_bytes = conn_info->wr_bytes + conn_info->rd_bytes;
  bool meets_activity_threshold =
      total_bytes >= conn_info->last_reported_bytes + kConnStatsDataThreshold;
  if (meets_activity_threshold) {
    struct conn_stats_event_t *event = fill_conn_stats_event(conn_info);
    if (event != NULL) {
      conn_stats_events.perf_submit(ctx, event,
                                    sizeof(struct conn_stats_event_t));
    }

    conn_info->last_reported_bytes = conn_info->rd_bytes + conn_info->wr_bytes;
  }
}

static __inline void process_data(const bool vecs, struct pt_regs *ctx,
                                  uint64_t id,
                                  const enum traffic_direction_t direction,
                                  const struct data_args_t *args,
                                  ssize_t bytes_count, bool ssl) {
  uint32_t tgid = id >> 32;

  if (!vecs && args->buf == NULL) {
    return;
  }

  if (vecs && (args->iov == NULL || args->iovlen <= 0)) {
    return;
  }

  if (args->fd < 0) {
    return;
  }

  if (bytes_count <= 0) {
    // This read()/write() call failed, or processed nothing.
    return;
  }

  enum target_tgid_match_result_t match_result = match_trace_tgid(tgid);
  if (match_result == TARGET_TGID_UNMATCHED) {
    return;
  }
  bool force_trace_tgid = (match_result == TARGET_TGID_MATCHED);

  struct conn_info_t *conn_info = get_or_create_conn_info(tgid, args->fd);
  if (conn_info == NULL) {
    return;
  }

  if (!should_trace_conn(conn_info)) {
    return;
  }

  uint64_t tgid_fd = gen_tgid_fd(tgid, args->fd);
  uint64_t *conn_disabled_tsid_ptr = conn_disabled_map.lookup(&tgid_fd);
  uint64_t conn_disabled_tsid =
      (conn_disabled_tsid_ptr == NULL) ? 0 : *conn_disabled_tsid_ptr;

  // Only process plaintext data.
  if (conn_info->ssl == ssl) {
    // TODO(yzhao): Split the interface such that the singular buf case and
    // multiple bufs in msghdr are handled separately without mixed interface.
    // The plan is to factor out helper functions for lower-level
    // functionalities, and call them separately for each case.
    if (!vecs) {
      update_traffic_class(conn_info, direction, args->buf, bytes_count);
    } else {
      struct iovec iov_cpy;
      size_t buf_size = 0;
      // With vectorized buffers, there can be empty elements sent.
      // For protocol inference, it requires a non empty buffer to get the real
      // data

#pragma unroll
      for (size_t i = 0; i < PROTOCOL_VEC_LIMIT && i < args->iovlen; i++) {
        BPF_PROBE_READ_VAR(iov_cpy, &args->iov[i]);
        buf_size = min_size_t(iov_cpy.iov_len, bytes_count);
        if (buf_size != 0) {
          update_traffic_class(conn_info, direction, iov_cpy.iov_base,
                               buf_size);
          break;
        }
      }
    }

    if (should_send_data(tgid, conn_disabled_tsid, force_trace_tgid,
                         conn_info)) {
      struct socket_data_event_t *event =
          fill_socket_data_event(args->source_fn, direction, conn_info);
      if (event == NULL) {
        // event == NULL not expected to ever happen.
        return;
      }

      // TODO(yzhao): Same TODO for split the interface.
      if (!vecs) {
        perf_submit_wrapper(ctx, direction, args->buf, bytes_count, conn_info,
                            event);
      } else {
        // TODO(yzhao): iov[0] is copied twice, once in calling
        // update_traffic_class(), and here. This happens to the write probes as
        // well, but the calls are placed in the entry and return probes
        // respectively. Consider remove one copy.
        perf_submit_iovecs(ctx, direction, args->iov, args->iovlen, bytes_count,
                           conn_info, event);
      }
    }
  }

  // TODO(oazizi): For conn stats, we should be using the encrypted traffic to
  // do the accounting,
  //               but that will break things with how we track data positions.
  //               For now, keep using plaintext data. In the future, this if
  //               statement should be:
  //                     if (!ssl) { ... }
  if (conn_info->ssl == ssl) {
    update_conn_stats(ctx, conn_info, direction, bytes_count);
  }

  return;
}

// These wrappers around process_data are carefully written so that they call
// process_data(), with a constant for `vecs`. Normally this would be done with
// meta-programming--for example through a template parameter--but C does not
// support that. By using a hard-coded constant for vecs (true or false), we
// enable clang to clone process_data and optimize away some code paths
// depending on whether we are using the iovecs or buf-based version. This is
// important for reducing the number of BPF instructions, since each syscall
// only needs one particular version.
// TODO(oazizi): Split process_data() into two versions, so we don't have to
// count on
//               Clang function cloning, which is not directly controllable.

static __inline void
process_syscall_data(struct pt_regs *ctx, uint64_t id,
                     const enum traffic_direction_t direction,
                     const struct data_args_t *args, ssize_t bytes_count) {
  process_data(/* vecs */ false, ctx, id, direction, args, bytes_count,
               /* ssl */ false);
}

static __inline void
process_syscall_data_vecs(struct pt_regs *ctx, uint64_t id,
                          const enum traffic_direction_t direction,
                          const struct data_args_t *args, ssize_t bytes_count) {
  process_data(/* vecs */ true, ctx, id, direction, args, bytes_count,
               /* ssl */ false);
}

static __inline void
process_syscall_sendfile(struct pt_regs *ctx, uint64_t id,
                         const struct sendfile_args_t *args,
                         ssize_t bytes_count) {
  uint32_t tgid = id >> 32;

  if (args->out_fd < 0) {
    return;
  }

  if (bytes_count <= 0) {
    // This sendfile call failed, or processed nothing.
    return;
  }

  enum target_tgid_match_result_t match_result = match_trace_tgid(tgid);
  if (match_result == TARGET_TGID_UNMATCHED) {
    return;
  }
  bool force_trace_tgid = (match_result == TARGET_TGID_MATCHED);

  struct conn_info_t *conn_info = get_or_create_conn_info(tgid, args->out_fd);
  if (conn_info == NULL) {
    return;
  }

  if (!should_trace_conn(conn_info)) {
    return;
  }

  uint64_t tgid_fd = gen_tgid_fd(tgid, args->out_fd);
  uint64_t *conn_disabled_tsid_ptr = conn_disabled_map.lookup(&tgid_fd);
  uint64_t conn_disabled_tsid =
      (conn_disabled_tsid_ptr == NULL) ? 0 : *conn_disabled_tsid_ptr;

  if (should_send_data(tgid, conn_disabled_tsid, force_trace_tgid, conn_info)) {
    struct socket_data_event_t *event =
        fill_socket_data_event(kSyscallSendfile, kEgress, conn_info);
    if (event == NULL) {
      // event == NULL not expected to ever happen.
      return;
    }

    event->attr.pos = conn_info->wr_bytes;
    event->attr.msg_size = bytes_count;
    event->attr.msg_buf_size = 0;
    socket_data_events.perf_submit(ctx, event, sizeof(event->attr));
  }

  update_conn_stats(ctx, conn_info, kEgress, bytes_count);

  return;
}

static __inline void
process_syscall_close(struct pt_regs *ctx, uint64_t id,
                      const struct close_args_t *close_args) {
  bpf_trace_printk("Enter process_syscall_close, %d\n",
                   sizeof(struct socket_control_event_t));
  uint32_t tgid = id >> 32;
  int ret_val = PT_REGS_RC(ctx);

  if (close_args->fd < 0) {
    bpf_trace_printk("Exit before submit: close_args->fd < 0, %d\n", 0);
    return;
  }

  if (ret_val < 0) {
    bpf_trace_printk("Exit before submit: ret_val < 0, %d\n", 0);
    // This close() call failed.
    return;
  }

  if (match_trace_tgid(tgid) == TARGET_TGID_UNMATCHED) {
    bpf_trace_printk("Exit before submit: TARGET_TGID_UNMATCHED, %d\n",
                     TARGET_TGID_UNMATCHED);
    return;
  }

  uint64_t tgid_fd = gen_tgid_fd(tgid, close_args->fd);
  struct conn_info_t *conn_info = conn_info_map.lookup(&tgid_fd);
  if (conn_info == NULL) {
    bpf_trace_printk("Exit before submit: conn_info === NULL, %d\n", 0);
    return;
  }

  // Only submit event to user-space if there was a corresponding open or data
  // event reported. This is to avoid polluting the perf buffer.
  if (should_trace_sockaddr_family(conn_info->addr.sa.sa_family) ||
      conn_info->wr_bytes != 0 || conn_info->rd_bytes != 0) {
    bpf_trace_printk("Going to submit, %d\n",
                     sizeof(struct socket_control_event_t));
    submit_close_event(ctx, conn_info, kSyscallClose);

    // Report final conn stats event for this connection.
    struct conn_stats_event_t *event = fill_conn_stats_event(conn_info);
    if (event != NULL) {
      event->conn_events = event->conn_events | CONN_CLOSE;
      conn_stats_events.perf_submit(ctx, event,
                                    sizeof(struct conn_stats_event_t));
    }
  } else {
    bpf_trace_printk("Not going to submit, %d\n",
                     sizeof(struct socket_control_event_t));
  }

  conn_info_map.delete(&tgid_fd);
}

/***********************************************************
 * BPF syscall probe function entry-points
 ***********************************************************/

// The following functions are the tracing function entry points.
// There is an entry probe and a return probe for each syscall.
// Information from both the entry and return probes are required
// before a syscall can be processed.
//
// General structure:
//    Entry probe: responsible for recording arguments.
//    Return probe: responsible for retrieving recorded arguments,
//                  extracting the return value,
//                  and processing the syscall with the combined context.
//
// Syscall signatures are listed. Look for detailed synopses in man pages.

// int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen);
int syscall__probe_entry_connect(struct pt_regs *ctx, int sockfd,
                                 const struct sockaddr *addr,
                                 socklen_t addrlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct connect_args_t connect_args = {};
  connect_args.fd = sockfd;
  connect_args.addr = addr;
  active_connect_args_map.update(&id, &connect_args);

  return 0;
}

int syscall__probe_ret_connect(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Unstash arguments, and process syscall.
  const struct connect_args_t *connect_args =
      active_connect_args_map.lookup(&id);
  if (connect_args != NULL) {
    process_syscall_connect(ctx, id, connect_args);
  }

  active_connect_args_map.delete(&id);
  return 0;
}

// int accept(int sockfd, struct sockaddr *addr, socklen_t *addrlen);
int syscall__probe_entry_accept(struct pt_regs *ctx, int sockfd,
                                struct sockaddr *addr, socklen_t *addrlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct accept_args_t accept_args = {};
  accept_args.addr = addr;
  active_accept_args_map.update(&id, &accept_args);

  return 0;
}

int syscall__probe_ret_accept(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Unstash arguments, and process syscall.
  struct accept_args_t *accept_args = active_accept_args_map.lookup(&id);
  if (accept_args != NULL) {
    process_syscall_accept(ctx, id, accept_args);
  }

  active_accept_args_map.delete(&id);
  return 0;
}

// int accept4(int sockfd, struct sockaddr *addr, socklen_t *addrlen, int
// flags);
int syscall__probe_entry_accept4(struct pt_regs *ctx, int sockfd,
                                 struct sockaddr *addr, socklen_t *addrlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct accept_args_t accept_args = {};
  accept_args.addr = addr;
  active_accept_args_map.update(&id, &accept_args);

  return 0;
}

int syscall__probe_ret_accept4(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Unstash arguments, and process syscall.
  struct accept_args_t *accept_args = active_accept_args_map.lookup(&id);
  if (accept_args != NULL) {
    process_syscall_accept(ctx, id, accept_args);
  }

  active_accept_args_map.delete(&id);
  return 0;
}

// ssize_t write(int fd, const void *buf, size_t count);
int syscall__probe_entry_write(struct pt_regs *ctx, int fd, char *buf,
                               size_t count) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct data_args_t write_args = {};
  write_args.source_fn = kSyscallWrite;
  write_args.fd = fd;
  write_args.buf = buf;
  active_write_args_map.update(&id, &write_args);

  return 0;
}

int syscall__probe_ret_write(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  struct data_args_t *write_args = active_write_args_map.lookup(&id);
  if (write_args != NULL && write_args->sock_event) {
    process_syscall_data(ctx, id, kEgress, write_args, bytes_count);
  }

  active_write_args_map.delete(&id);
  return 0;
}

// ssize_t send(int sockfd, const void *buf, size_t len, int flags);
int syscall__probe_entry_send(struct pt_regs *ctx, int sockfd, char *buf,
                              size_t len) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct data_args_t write_args = {};
  write_args.source_fn = kSyscallSend;
  write_args.fd = sockfd;
  write_args.buf = buf;
  active_write_args_map.update(&id, &write_args);

  return 0;
}

int syscall__probe_ret_send(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  struct data_args_t *write_args = active_write_args_map.lookup(&id);
  if (write_args != NULL) {
    process_syscall_data(ctx, id, kEgress, write_args, bytes_count);
  }

  active_write_args_map.delete(&id);
  return 0;
}

// ssize_t read(int fd, void *buf, size_t count);
int syscall__probe_entry_read(struct pt_regs *ctx, int fd, char *buf,
                              size_t count) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct data_args_t read_args = {};
  read_args.source_fn = kSyscallRead;
  read_args.fd = fd;
  read_args.buf = buf;
  active_read_args_map.update(&id, &read_args);

  return 0;
}

int syscall__probe_ret_read(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  struct data_args_t *read_args = active_read_args_map.lookup(&id);
  if (read_args != NULL && read_args->sock_event) {
    process_syscall_data(ctx, id, kIngress, read_args, bytes_count);
  }

  active_read_args_map.delete(&id);
  return 0;
}

// ssize_t recv(int sockfd, void *buf, size_t len, int flags);
int syscall__probe_entry_recv(struct pt_regs *ctx, int sockfd, char *buf,
                              size_t len) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct data_args_t read_args = {};
  read_args.source_fn = kSyscallRecv;
  read_args.fd = sockfd;
  read_args.buf = buf;
  active_read_args_map.update(&id, &read_args);

  return 0;
}

int syscall__probe_ret_recv(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  struct data_args_t *read_args = active_read_args_map.lookup(&id);
  if (read_args != NULL) {
    process_syscall_data(ctx, id, kIngress, read_args, bytes_count);
  }

  active_read_args_map.delete(&id);
  return 0;
}

// ssize_t sendto(int sockfd, const void *buf, size_t len, int flags,
//                const struct sockaddr *dest_addr, socklen_t addrlen);
int syscall__probe_entry_sendto(struct pt_regs *ctx, int sockfd, char *buf,
                                size_t len, int flags,
                                const struct sockaddr *dest_addr,
                                socklen_t addrlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  if (dest_addr != NULL) {
    struct connect_args_t connect_args = {};
    connect_args.fd = sockfd;
    connect_args.addr = dest_addr;
    active_connect_args_map.update(&id, &connect_args);
  }

  // Stash arguments.
  struct data_args_t write_args = {};
  write_args.source_fn = kSyscallSendTo;
  write_args.fd = sockfd;
  write_args.buf = buf;
  active_write_args_map.update(&id, &write_args);

  return 0;
}

int syscall__probe_ret_sendto(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Potential issue: If sentto() addr is provided by a TCP connection, the
  // syscall may ignore it, but we would still trace it. In practice, TCP
  // connections should not be using sendto() with an addr argument.
  //
  // From the man page:
  //   If sendto() is used on a connection-mode (SOCK_STREAM, SOCK_SEQPACKET)
  //   socket, the arguments dest_addr and addrlen are ignored (and the error
  //   EISCONN may be returned when they  are not NULL and 0)
  //
  //   EISCONN
  //   The connection-mode socket was connected already but a recipient was
  //   specified. (Now either this error is returned, or the recipient
  //   specification is ignored.)

  // Unstash arguments, and process syscall.
  const struct connect_args_t *connect_args =
      active_connect_args_map.lookup(&id);
  if (connect_args != NULL && bytes_count > 0) {
    process_implicit_conn(ctx, id, connect_args, kSyscallSendTo);
  }
  active_connect_args_map.delete(&id);

  // Unstash arguments, and process syscall.
  struct data_args_t *write_args = active_write_args_map.lookup(&id);
  if (write_args != NULL) {
    process_syscall_data(ctx, id, kEgress, write_args, bytes_count);
  }

  active_write_args_map.delete(&id);

  return 0;
}

// ssize_t recvfrom(int sockfd, void *buf, size_t len, int flags,
//                  struct sockaddr *src_addr, socklen_t *addrlen);
int syscall__probe_entry_recvfrom(struct pt_regs *ctx, int sockfd, char *buf,
                                  size_t len, int flags,
                                  struct sockaddr *src_addr,
                                  socklen_t *addrlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  if (src_addr != NULL) {
    struct connect_args_t connect_args = {};
    connect_args.fd = sockfd;
    connect_args.addr = src_addr;
    active_connect_args_map.update(&id, &connect_args);
  }

  // Stash arguments.
  struct data_args_t read_args = {};
  read_args.source_fn = kSyscallRecvFrom;
  read_args.fd = sockfd;
  read_args.buf = buf;
  active_read_args_map.update(&id, &read_args);

  return 0;
}

int syscall__probe_ret_recvfrom(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  const struct connect_args_t *connect_args =
      active_connect_args_map.lookup(&id);
  if (connect_args != NULL && bytes_count > 0) {
    process_implicit_conn(ctx, id, connect_args, kSyscallRecvFrom);
  }
  active_connect_args_map.delete(&id);

  // Unstash arguments, and process syscall.
  struct data_args_t *read_args = active_read_args_map.lookup(&id);
  if (read_args != NULL) {
    process_syscall_data(ctx, id, kIngress, read_args, bytes_count);
  }
  active_read_args_map.delete(&id);

  return 0;
}

// ssize_t sendmsg(int sockfd, const struct msghdr *msg, int flags);
int syscall__probe_entry_sendmsg(struct pt_regs *ctx, int sockfd,
                                 const struct user_msghdr *msghdr) {
  uint64_t id = bpf_get_current_pid_tgid();

  if (msghdr != NULL) {
    // Stash arguments.
    if (msghdr->msg_name != NULL) {
      struct connect_args_t connect_args = {};
      connect_args.fd = sockfd;
      connect_args.addr = msghdr->msg_name;
      active_connect_args_map.update(&id, &connect_args);
    }

    // Stash arguments.
    struct data_args_t write_args = {};
    write_args.source_fn = kSyscallSendMsg;
    write_args.fd = sockfd;
    write_args.iov = msghdr->msg_iov;
    write_args.iovlen = msghdr->msg_iovlen;
    active_write_args_map.update(&id, &write_args);
  }

  return 0;
}

int syscall__probe_ret_sendmsg(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  const struct connect_args_t *connect_args =
      active_connect_args_map.lookup(&id);
  if (connect_args != NULL && bytes_count > 0) {
    process_implicit_conn(ctx, id, connect_args, kSyscallSendMsg);
  }
  active_connect_args_map.delete(&id);

  // Unstash arguments, and process syscall.
  struct data_args_t *write_args = active_write_args_map.lookup(&id);
  if (write_args != NULL) {
    process_syscall_data_vecs(ctx, id, kEgress, write_args, bytes_count);
  }

  active_write_args_map.delete(&id);
  return 0;
}

int syscall__probe_entry_sendmmsg(struct pt_regs *ctx, int sockfd,
                                  struct mmsghdr *msgvec, unsigned int vlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // TODO(oazizi): Right now, we only trace the first message in a sendmmsg()
  // call.
  if (msgvec != NULL && vlen >= 1) {
    // Stash arguments.
    if (msgvec[0].msg_hdr.msg_name != NULL) {
      struct connect_args_t connect_args = {};
      connect_args.fd = sockfd;
      connect_args.addr = msgvec[0].msg_hdr.msg_name;
      active_connect_args_map.update(&id, &connect_args);
    }

    // Stash arguments.
    struct data_args_t write_args = {};
    write_args.source_fn = kSyscallSendMMsg;
    write_args.fd = sockfd;
    write_args.iov = msgvec[0].msg_hdr.msg_iov;
    write_args.iovlen = msgvec[0].msg_hdr.msg_iovlen;
    write_args.msg_len = &msgvec[0].msg_len;
    active_write_args_map.update(&id, &write_args);
  }

  return 0;
}

int syscall__probe_ret_sendmmsg(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  int num_msgs = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  const struct connect_args_t *connect_args =
      active_connect_args_map.lookup(&id);
  if (connect_args != NULL && num_msgs > 0) {
    process_implicit_conn(ctx, id, connect_args, kSyscallSendMMsg);
  }
  active_connect_args_map.delete(&id);

  // Unstash arguments, and process syscall.
  struct data_args_t *write_args = active_write_args_map.lookup(&id);
  if (write_args != NULL && num_msgs > 0) {
    // msg_len is defined as unsigned int, so we have to use the same here.
    // This is different than most other syscalls that use ssize_t.
    unsigned int bytes_count = 0;
    BPF_PROBE_READ_VAR(bytes_count, write_args->msg_len);
    process_syscall_data_vecs(ctx, id, kEgress, write_args, bytes_count);
  }
  active_write_args_map.delete(&id);

  return 0;
}

// ssize_t recvmsg(int sockfd, struct msghdr *msg, int flags);
int syscall__probe_entry_recvmsg(struct pt_regs *ctx, int sockfd,
                                 struct user_msghdr *msghdr) {
  uint64_t id = bpf_get_current_pid_tgid();

  if (msghdr != NULL) {
    // Stash arguments.
    if (msghdr->msg_name != NULL) {
      struct connect_args_t connect_args = {};
      connect_args.fd = sockfd;
      connect_args.addr = msghdr->msg_name;
      active_connect_args_map.update(&id, &connect_args);
    }

    // Stash arguments.
    struct data_args_t read_args = {};
    read_args.source_fn = kSyscallRecvMsg;
    read_args.fd = sockfd;
    read_args.iov = msghdr->msg_iov;
    read_args.iovlen = msghdr->msg_iovlen;
    active_read_args_map.update(&id, &read_args);
  }

  return 0;
}

int syscall__probe_ret_recvmsg(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  const struct connect_args_t *connect_args =
      active_connect_args_map.lookup(&id);
  if (connect_args != NULL && bytes_count > 0) {
    process_implicit_conn(ctx, id, connect_args, kSyscallRecvMsg);
  }
  active_connect_args_map.delete(&id);

  // Unstash arguments, and process syscall.
  struct data_args_t *read_args = active_read_args_map.lookup(&id);
  if (read_args != NULL) {
    process_syscall_data_vecs(ctx, id, kIngress, read_args, bytes_count);
  }

  active_read_args_map.delete(&id);
  return 0;
}

// int recvmmsg(int sockfd, struct mmsghdr *msgvec, unsigned int vlen,
//              int flags, struct timespec *timeout);
int syscall__probe_entry_recvmmsg(struct pt_regs *ctx, int sockfd,
                                  struct mmsghdr *msgvec, unsigned int vlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // TODO(oazizi): Right now, we only trace the first message in a recvmmsg()
  // call.
  if (msgvec != NULL && vlen >= 1) {
    // Stash arguments.
    if (msgvec[0].msg_hdr.msg_name != NULL) {
      struct connect_args_t connect_args = {};
      connect_args.fd = sockfd;
      connect_args.addr = msgvec[0].msg_hdr.msg_name;
      active_connect_args_map.update(&id, &connect_args);
    }

    // Stash arguments.
    struct data_args_t read_args = {};
    read_args.source_fn = kSyscallRecvMMsg;
    read_args.fd = sockfd;
    read_args.iov = msgvec[0].msg_hdr.msg_iov;
    read_args.iovlen = msgvec[0].msg_hdr.msg_iovlen;
    read_args.msg_len = &msgvec[0].msg_len;
    active_read_args_map.update(&id, &read_args);
  }

  return 0;
}

int syscall__probe_ret_recvmmsg(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  int num_msgs = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  const struct connect_args_t *connect_args =
      active_connect_args_map.lookup(&id);
  if (connect_args != NULL && num_msgs > 0) {
    process_implicit_conn(ctx, id, connect_args, kSyscallRecvMMsg);
  }
  active_connect_args_map.delete(&id);

  // Unstash arguments, and process syscall.
  struct data_args_t *read_args = active_read_args_map.lookup(&id);
  if (read_args != NULL && num_msgs > 0) {
    // msg_len is defined as unsigned int, so we have to use the same here.
    // This is different than most other syscalls that use ssize_t.
    unsigned int bytes_count = 0;
    BPF_PROBE_READ_VAR(bytes_count, read_args->msg_len);
    process_syscall_data_vecs(ctx, id, kIngress, read_args, bytes_count);
  }
  active_read_args_map.delete(&id);

  return 0;
}

// ssize_t writev(int fd, const struct iovec *iov, int iovcnt);
int syscall__probe_entry_writev(struct pt_regs *ctx, int fd,
                                const struct iovec *iov, int iovlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct data_args_t write_args = {};
  write_args.source_fn = kSyscallWriteV;
  write_args.fd = fd;
  write_args.iov = iov;
  write_args.iovlen = iovlen;
  active_write_args_map.update(&id, &write_args);

  return 0;
}

int syscall__probe_ret_writev(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  struct data_args_t *write_args = active_write_args_map.lookup(&id);
  if (write_args != NULL && write_args->sock_event) {
    process_syscall_data_vecs(ctx, id, kEgress, write_args, bytes_count);
  }

  active_write_args_map.delete(&id);
  return 0;
}

// ssize_t readv(int fd, const struct iovec *iov, int iovcnt);
int syscall__probe_entry_readv(struct pt_regs *ctx, int fd, struct iovec *iov,
                               int iovlen) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct data_args_t read_args = {};
  read_args.source_fn = kSyscallReadV;
  read_args.fd = fd;
  read_args.iov = iov;
  read_args.iovlen = iovlen;
  active_read_args_map.update(&id, &read_args);

  return 0;
}

int syscall__probe_ret_readv(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  struct data_args_t *read_args = active_read_args_map.lookup(&id);
  if (read_args != NULL && read_args->sock_event) {
    process_syscall_data_vecs(ctx, id, kIngress, read_args, bytes_count);
  }

  active_read_args_map.delete(&id);
  return 0;
}

// int close(int fd);
int syscall__probe_entry_close(struct pt_regs *ctx, int fd) {
  uint64_t id = bpf_get_current_pid_tgid();
  bpf_trace_printk("Enter syscall__probe_entry_close, pid %lld\n", id);

  // Stash arguments.
  struct close_args_t close_args;
  close_args.fd = fd;
  active_close_args_map.update(&id, &close_args);

  /* bpf_trace_printk("Exit syscall__probe_entry_close, pid %lld\n", id); */
  return 0;
}

int syscall__probe_ret_close(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  bpf_trace_printk("Enter syscall__probe_ret_close, pid %lld\n", id);

  // Unstash arguments, and process syscall.
  const struct close_args_t *close_args = active_close_args_map.lookup(&id);
  if (close_args != NULL) {
    bpf_trace_printk("Close args != NULL, %d\n", 0);
    process_syscall_close(ctx, id, close_args);
  } else {
    bpf_trace_printk("Close args == NULL, %d\n", 0);
  }

  active_close_args_map.delete(&id);
  return 0;
}

// int close(int fd);
int syscall__probe_entry_sendfile(struct pt_regs *ctx, int out_fd, int in_fd,
                                  off_t *offset, size_t count) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Stash arguments.
  struct sendfile_args_t args;
  args.out_fd = out_fd;
  args.in_fd = in_fd;
  args.count = count;
  active_sendfile_args_map.update(&id, &args);

  return 0;
}

int syscall__probe_ret_sendfile(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  ssize_t bytes_count = PT_REGS_RC(ctx);

  // Unstash arguments, and process syscall.
  const struct sendfile_args_t *args = active_sendfile_args_map.lookup(&id);
  if (args != NULL) {
    process_syscall_sendfile(ctx, id, args, bytes_count);
  }

  active_sendfile_args_map.delete(&id);
  return 0;
}

// void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t
// offset)
int syscall__probe_entry_mmap(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  struct upid_t upid = {};
  upid.tgid = id >> 32;
  upid.start_time_ticks = get_tgid_start_time();

  mmap_events.perf_submit(ctx, &upid, sizeof(upid));

  return 0;
}

// Trace kernel function:
// struct socket *sock_alloc(void)
// which is called inside accept4() syscall to allocate socket data structure.
// Only need a return probe, as the function does not accept any arguments.
int probe_ret_sock_alloc(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();

  // Only trace sock_alloc() called by accept()/accept4().
  struct accept_args_t *accept_args = active_accept_args_map.lookup(&id);
  if (accept_args == NULL) {
    return 0;
  }

  if (accept_args->sock_alloc_socket == NULL) {
    accept_args->sock_alloc_socket = (struct socket *)PT_REGS_RC(ctx);
  }

  return 0;
}

// Trace kernel function:
// int security_socket_sendmsg(struct socket *sock, struct msghdr *msg, int
// size) which is called by write/writev/send/sendmsg.
int probe_entry_security_socket_sendmsg(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  struct data_args_t *write_args = active_write_args_map.lookup(&id);
  if (write_args != NULL) {
    write_args->sock_event = true;
  }
  return 0;
}

// Trace kernel function:
// int security_socket_recvmsg(struct socket *sock, struct msghdr *msg, int
// size)
int probe_entry_security_socket_recvmsg(struct pt_regs *ctx) {
  uint64_t id = bpf_get_current_pid_tgid();
  struct data_args_t *read_args = active_read_args_map.lookup(&id);
  if (read_args != NULL) {
    read_args->sock_event = true;
  }
  return 0;
}

/*
// OpenSSL tracing probes.
#include "src/stirling/source_connectors/socket_tracer/bcc_bpf/openssl_trace.c"

// Go HTTP2 tracing probes.
#include "src/stirling/source_connectors/socket_tracer/bcc_bpf/go_http2_trace.c"

// GoTLS tracing probes.
#include "src/stirling/source_connectors/socket_tracer/bcc_bpf/go_tls_trace.c"

// gRPC-c tracing probes.
#include "src/stirling/source_connectors/socket_tracer/bcc_bpf/grpc_c_trace.c"
*/
