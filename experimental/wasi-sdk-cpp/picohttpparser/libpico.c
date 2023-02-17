#include "picohttpparser.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define DEBUG 1

#define MAX_REQUEST_LEN 4096 // Maximum supported request length
#define MAX_METHOD_LEN 8     // Maximum supported HTTP "method" length
#define MAX_PATH_LEN 120     // Maximum supported HTTP "path" length
#define MAX_HEADERS 32       // Maximum supported HTTP headers
#define MAX_EVENTS                                                             \
  128 // Maximum allowed input events each invoking, enlarge this value by needs

// TODO: define return codes
#define RET_OK 0 // wasmtime favors sandbox return code in [1,126]?
#define RET_BAD_INPUT 2
#define RET_MEM_ERR 3
#define RET_MAX 126

// http event (e.g. collected by ebpf)
struct http_event {
  char data[MAX_REQUEST_LEN];
};

// Parsed information about a HTTP request
// TODO: avoid to use size_t, use a type of specific length, e.g. uint32_t
struct parsed_req_info {
  const char *method;
  size_t method_len;
  const char *path;
  size_t path_len;
  int minor_version;
  struct phr_header headers[MAX_HEADERS];
  size_t num_headers;
};

// Parsed result about a http_event
// TODO: refine this structure to include more meaningful fields
struct parse_result {
  char method[MAX_METHOD_LEN];
  char path[MAX_PATH_LEN];
};

// Typical workflow of parsing raw events with pico wasm module:
// 1. caller copies raw events into `input_buf`, then calls
// `pico_parse_events()` with `num_events` as parameter;
// 2. pico parsers events, then places the parsed results into `output_buf`, and
// sets `output_item_count`;
// 3. caller gets the output buffer address, output item size and count, then
// copies the results from that buffer.
//
// TODO: make them a general struct, so we only need to export one function to
// the sandbox creator?
static struct http_event *input_buf;
static struct parse_result *output_buf;
static int output_item_count;

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

static void print_parse_result(struct parsed_req_info *r) {
  printf("Parsed request information:\n");

  print_str("Method", r->method, r->method_len);
  print_str("Path", r->path, r->path_len);

  for (size_t i = 0; i < r->num_headers; i++) {
    print_header(&r->headers[i]);
  }
}
#endif

int save_parse_result(struct parsed_req_info *r) {
  if (output_item_count >= MAX_EVENTS) {
    printf("ERROR: output_buf full, current %d, max %d\r", output_item_count,
           MAX_EVENTS);
    return -1;
  }

  struct parse_result *p =
      (struct parse_result *)(output_buf + output_item_count);
  memset(p, '\0', sizeof(struct parse_result));

  int len = r->method_len > MAX_METHOD_LEN ? MAX_METHOD_LEN : r->method_len;
  memcpy(p->method, r->method, len);

  len = r->path_len > MAX_PATH_LEN ? MAX_PATH_LEN : r->path_len;
  memcpy(p->path, r->path, len);

  // TODO: copy headers if needed
  /* for (int i=0; i<MAX_HEADERS; i++) { */
  /* } */

  output_item_count++;
  return 0;
}

// Call this function after input_buf has been filled. Parsed results will be
// stored in output_buf.
//
// Parameters:
// * num_reqs: number of requests in the input_buf
//
// Return:
// * 0 on success
// * error code on failures
int pico_parse_events(int num_reqs) {
  // Reset output stats
  output_item_count = 0;

  if (!input_buf) {
    printf("Input buffer uninitialized\n");
    return RET_MEM_ERR;
  }

  if (!output_buf) {
    printf("Output buffer uninitialized\n");
    return RET_MEM_ERR;
  }

  printf("Received %d http events, going to parse them\n", num_reqs);

  const char *req;
  size_t req_len;
  struct parsed_req_info info;

  for (int i = 0; i < num_reqs; i++) {
    req = ((struct http_event *)(input_buf + i))->data;
    req_len = strlen(req);
    if (req_len == 0) {
      printf("Warning: http event %d is empty, skip it\n", i);
      continue;
    }

#if DEBUG
    // The "\r\n" characters in req make the print ugly, so just disable it
    /* printf("Raw HTTP request: %s\n", req); */
#endif

    memset(&info, 0, sizeof(struct parsed_req_info));
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

  return RET_OK;
}

// Function that will be exported to the sandbox creator
int allocate_input_output_bufs() {
  // Allocate buffer for return results
  input_buf =
      (struct http_event *)malloc(MAX_EVENTS * sizeof(struct http_event));
  if (!input_buf) {
    return 1;
  }

  output_buf =
      (struct parse_result *)malloc(MAX_EVENTS * sizeof(struct parse_result));
  if (!output_buf) {
    free(input_buf);
    return 1;
  }

  printf("Allocate buffer for input/output successful\n");
  return 0;
}

// Function that will be exported to the sandbox creator
int free_input_output_bufs() {
  free(input_buf);
  free(output_buf);
  return 0;
}

// Function that will be exported to the sandbox creator.
// Used to pass input data (http events).
void *get_input_buf() { return (void *)input_buf; }

// Function that will be exported to the sandbox creator
// Used to pass output data (parsed results).
void *get_output_buf() { return (void *)output_buf; }

// Function that will be exported to the sandbox creator
// Used to tell the caller how many parsed results are in the output buffer.
int get_output_item_count() { return output_item_count; }

// Function that will be exported to the sandbox creator
// Used to tell the caller the exact size of one parsed result.
int get_output_item_size() { return sizeof(struct parse_result); }
