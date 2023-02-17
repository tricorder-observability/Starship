// LINT_C_FILE

#include <math.h>
#include <stdio.h>
#include <string.h>

#include "../../../../modules/common/io.h"
#include "cgo/ints.h"
#include "cgo/struct-bindgen.h"

// A simple function to copy entire input buf to output buffer.
void write_ints_to_output() {
  Ints_t *v_ptr = (Ints_t *)get_input_buf();
  printf("A=%d B=%d C=%d D=%d\n", v_ptr->A, v_ptr->B, v_ptr->C, v_ptr->D);
  fflush(stdout);

  v_ptr->A = 1;
  v_ptr->B = 2;
  v_ptr->C = 3;
  v_ptr->D = 4;
  v_ptr->E = 5;

  if (!can_write_to_output_buf(sizeof(*v_ptr))) {
    return;
  }
  copy_to_output(v_ptr, sizeof(*v_ptr));
}

size_t get_ints_t_size() { return sizeof(Packed_ints_t); }

// A simple function to copy and marshal entire input buf to output buffer.
void write_marshaled_struct_to_output() {
  struct event2 *e2 = input_buf.data;
  // simply memcpy will not work, because the struct is not packed as expected.
  // memcpy(&e2, input_buf.data, sizeof(struct event2));
  printf("e2->x = %f, e2->y = %f, e2->z = %d, e2->comm[0] = %d\n", e2->X, e2->Y,
         e2->Z, e2->Comm[0]);
  assert(fabs(e2->X - 1.01) < 0.0001);
  assert(fabs(e2->Y - 2.02) < 0.0001);
  assert(e2->Z == 3);
  assert(e2->Comm[0] == 5);
  e2->A = 1;
  e2->X = 0.03;
  e2->Y = 0.04;
  e2->Z = 5;
  e2->Comm[0] = 1;
  e2->Comm[5] = 10;
  if (can_write_to_output_buf(input_buf.length)) {
    copy_to_output(input_buf.data, input_buf.length);
    printf("copied to output\n");
  }
}

int main() { return 0; }
