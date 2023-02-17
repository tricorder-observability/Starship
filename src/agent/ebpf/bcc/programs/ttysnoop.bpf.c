// Based on https://github.com/iovisor/bcc/blob/tricorder/tools/ttysnoop.py
// Check out the python code in the above link to figure out how to use the
// output events.
// LINT_C_FILE

#include <linux/fs.h>
#include <linux/uio.h>
#include <uapi/linux/ptrace.h>

// NOTE: Added by tricorder to make this self-contained.
#define USER_DATASIZE 1024
#define USER_DATACOUNT 16
#define PTS 13 // /dev/console inode number

#define BUFSIZE USER_DATASIZE
struct data_t {
  int count;
  char buf[BUFSIZE];
};
BPF_ARRAY(data_map, struct data_t, 1);
BPF_PERF_OUTPUT(events);
static int do_tty_write(void *ctx, const char __user *buf, size_t count) {
  int zero = 0, i;
  struct data_t *data;
/* We can't read data to map data before v4.11 */
#if LINUX_VERSION_CODE < KERNEL_VERSION(4, 11, 0)
  struct data_t _data = {};
  data = &_data;
#else
  data = data_map.lookup(&zero);
  if (!data)
    return 0;
#endif
#pragma unroll
  for (i = 0; i < USER_DATACOUNT; i++) {
    // bpf_probe_read_user() can only use a fixed size, so truncate to count
    // in user space:
    if (bpf_probe_read_user(&data->buf, BUFSIZE, (void *)buf))
      return 0;
    if (count > BUFSIZE)
      data->count = BUFSIZE;
    else
      data->count = count;
    events.perf_submit(ctx, data, sizeof(*data));
    if (count < BUFSIZE)
      return 0;
    count -= BUFSIZE;
    buf += BUFSIZE;
  }
  return 0;
};
/**
 * commit 9bb48c82aced (v5.11-rc4) tty: implement write_iter
 * changed arguments of tty_write function
 */
#if LINUX_VERSION_CODE < KERNEL_VERSION(5, 11, 0)
int kprobe__tty_write(struct pt_regs *ctx, struct file *file,
                      const char __user *buf, size_t count) {
  if (file->f_inode->i_ino != PTS)
    return 0;
  return do_tty_write(ctx, buf, count);
}
#else
KFUNC_PROBE(tty_write, struct kiocb *iocb, struct iov_iter *from) {
  const char __user *buf;
  const struct kvec *kvec;
  size_t count;
  if (iocb->ki_filp->f_inode->i_ino != PTS)
    return 0;
/**
 * commit 8cd54c1c8480 iov_iter: separate direction from flavour
 * `type` is represented by iter_type and data_source seperately
 */
#if LINUX_VERSION_CODE < KERNEL_VERSION(5, 14, 0)
  if (from->type != (ITER_IOVEC + WRITE))
    return 0;
#else
  if (from->iter_type != ITER_IOVEC)
    return 0;
  if (from->data_source != WRITE)
    return 0;
#endif
  kvec = from->kvec;
  buf = kvec->iov_base;
  count = kvec->iov_len;
  return do_tty_write(ctx, kvec->iov_base, kvec->iov_len);
}
#endif
