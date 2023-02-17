// LINT_C_FILE

#include "picohttpparser.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define DEBUG 1

#define MAX_METHOD_LEN 8 // Maximum supported HTTP "method" length
#define MAX_PATH_LEN 120 // Maximum supported HTTP "path" length
#define MAX_HEADERS 32   // Maximum supported HTTP headers
#define MAX_EVENTS                                                             \
  128 // Maximum allowed input events each invoking, enlarge this value by needs

// TODO: define return codes
#define RET_OK 0 // wasmtime favors sandbox return code in [1,126]?
#define RET_BAD_INPUT 2
#define RET_MEM_ERR 3
#define RET_MAX 126

// Parsed information about a HTTP request
// TODO: avoid to use size_t, use a type of specific length, e.g. uint32_t
struct parse_info {
  const char *method;
  size_t method_len;
  const char *path;
  size_t path_len;
  int minor_version;
  struct phr_header headers[MAX_HEADERS];
  size_t num_headers;
};

// Parsed result, which will be passed to the sandbox creator via memory.
// TODO: refine this structure to include more fields that's needed
struct parse_result {
  char method[MAX_METHOD_LEN];
  char path[MAX_PATH_LEN];
};

// Convention between sandbox and its caller (creator):
// * picohttpparser place the parsed result into `result_buf`, and set
// `result_count`
//   to indicate the valid length of the result.
// * agent get the buffer address and valid length, then copy result from the
// buffer.
// TODO: make them a general struct, so we only need to export one function to
// the sandbox creator?
static struct parse_result *result_buf;
static int result_count;

#if DEBUG
static void print_str(const char *title, const char *p, size_t len) {
  if (!p) {
    return;
  }

  printf("%-10s: ", title);
  for (size_t i = 0; i < len; i++) {
    printf("%c", *p++);
  }
  printf("\n");
}

static void print_header(struct phr_header *h) {
  if (!h) {
    return;
  }

  print_str("Header", h->name, h->name_len);
  print_str("", h->value, h->value_len);
}

static void print_parse_result(struct parse_info *r) {
  printf("Parsed request information:\n");

  print_str("Method", r->method, r->method_len);
  print_str("Path", r->path, r->path_len);

  for (size_t i = 0; i < r->num_headers; i++) {
    print_header(&r->headers[i]);
  }
}
#endif

int save_parse_result(struct parse_info *r) {
  if (result_count >= MAX_EVENTS) {
    printf("ERROR: result_buf full, current %d, max %d\r", result_count,
           MAX_EVENTS);
    return -1;
  }

  struct parse_result *p = (struct parse_result *)(result_buf + result_count);
  memset(p, '\0', sizeof(struct parse_result));

  int len = r->method_len > MAX_METHOD_LEN ? MAX_METHOD_LEN : r->method_len;
  memcpy(p->method, r->method, len);

  len = r->path_len > MAX_PATH_LEN ? MAX_PATH_LEN : r->path_len;
  memcpy(p->path, r->path, len);

  // TODO: copy headers if needed
  /* for (int i=0; i<MAX_HEADERS; i++) { */
  /* } */

  result_count++;
  return 0;
}

// Expect <reqN> in the argument list to be a string of standard http header
int main(int argc, char *argv[]) {
  if (argc < 2) {
    printf("Invalid argument list, usage: %s <req1> <req2> ... <reqN>\n",
           argv[0]);
    return RET_BAD_INPUT;
  }

  int num_events = argc;
  if (argc > MAX_EVENTS + 1) {
    printf("Warning: max allowed events per invoking %d, currently given %d, "
           "events exceeds the threshold will be dropped\n",
           MAX_EVENTS, argc - 1);
    num_events = MAX_EVENTS;
  }

  // Allocate buffer for return results
  result_count = 0;
  result_buf =
      (struct parse_result *)malloc(num_events * sizeof(struct parse_result));
  if (!result_buf) {
    printf("Allocate memory for return results failed\n");
    return RET_MEM_ERR;
  }

  printf("Going to process %d events\n", num_events);

  const char *req;
  size_t req_len;
  struct parse_info info;

  for (int i = 1; i < num_events; i++) {
    req = argv[i];
    req_len = strlen(req);
    if (req_len == 0) { // empty request
      printf("Warning: event %d is empty, skip it\n", i);
      continue;
    }

#if DEBUG
    // The "\r\n" characters in req make the print ugly, so just disable it
    /* printf("Raw HTTP request: %s\n", req); */
#endif

    memset(&info, 0, sizeof(struct parse_info));
    info.num_headers = MAX_HEADERS; // reset num headers

    phr_parse_request(req,                 // const char *buf_start,
                      req_len,             // size_t len
                      &info.method,        // const char **method
                      &info.method_len,    // size_t *method_len
                      &info.path,          // const char **path,
                      &info.path_len,      // size_t *path_len
                      &info.minor_version, // int *minor_version
                      info.headers,        // struct phr_header *headers
                      &info.num_headers,   // size_t *num_headers
                      0);                  // size_t last_len
#if DEBUG
    print_parse_result(&info);
#endif
    save_parse_result(&info);
  }

  // TODO: the caller will free the result_buf for us?

  return RET_OK;
}

// Function that will be exported to the sandbox creator
void *get_result_buf() { return (void *)result_buf; }

// Function that will be exported to the sandbox creator
int get_result_count() { return result_count; }

// Function that will be exported to the sandbox creator
int get_result_struct_size() { return sizeof(struct parse_result); }
