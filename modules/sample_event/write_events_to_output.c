#include "../common/cJSON.h"
#include "../common/io.h"

#include "event_bindgen.h"

static_assert(sizeof(struct event_t) == 64, "Size of event2 is not 64");

// A simple function to copy entire input buf to output buffer.
// Return 0 if succeeded.
// Return 1 if failed to malloc output buffer.
int write_events_to_output() {
  struct event_t event;
  event.F = 2 * event.F;
  event.D = 2 * event.D;
  event.I = 2 * event.I;
  event.L = 2 * event.L;

  cJSON *root = cJSON_CreateObject();

  cJSON_AddNumberToObject(root, "F", event.F);
  cJSON_AddNumberToObject(root, "D", event.D);
  cJSON_AddNumberToObject(root, "I", event.I);
  cJSON_AddNumberToObject(root, "L", event.L);
  cJSON_AddStringToObject(root, "Comm", (const char *)event.Comm);

  char *json = NULL;
  json = cJSON_Print(root);
  cJSON_Delete(root);

  int json_size = strlen(json);
  void *buf = malloc_output_buf(json_size);
  if (buf == NULL) {
    return 1;
  }
  copy_to_output(json, json_size);
  // Free allocated memory from JSON_print().
  free(json);
  return 0;
}

// Do nothing
// TODO(yaxiong): Investigate how to remove this and build wasi module without
// main().
int main() { return 0; }
