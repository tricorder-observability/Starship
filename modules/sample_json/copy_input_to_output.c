// LINT_C_FILE

#include <string.h>

#include "../common/io.h"

// A simple function to copy entire input buf to output buffer.
void copy_input_to_output() {
  malloc_output_buf(input_buf.capacity);
  if (can_write_to_output_buf(input_buf.length)) {
    copy_to_output(input_buf.data, input_buf.length);
  }
}

int main() { return 0; }
