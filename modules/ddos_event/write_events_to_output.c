#include "../common/cJSON.h"
#include "../common/io.h"

#include "event_bindgen.h"

static_assert(sizeof(struct detectionPackets) == 8,
              "Size of detectionPackets is not 8");

// A simple function to copy entire input buf to output buffer.
// Return 0 if succeeded.
// Return 1 if failed to malloc output buffer.
int write_events_to_output() {
  struct detectionPackets *detection_packet = get_input_buf();

  cJSON *root = cJSON_CreateObject();

  cJSON_AddNumberToObject(root, "nb_ddos_packets",
                          detection_packet->nb_ddos_packets);

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
