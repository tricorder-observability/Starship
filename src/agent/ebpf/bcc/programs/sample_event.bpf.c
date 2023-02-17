#include <linux/ptrace.h>

#include "event.h"

BPF_PERF_OUTPUT(events);

int syscall__probe_entry_read(struct pt_regs *ctx, int fd, char *buf,
                              size_t count) {
  struct event_t event = {};
  const char word[] = "hello world!";
  events.perf_submit(ctx, &event, sizeof(event));
  return 0;
}
