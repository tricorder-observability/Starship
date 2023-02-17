// LINT_C_FILE
#include <stdint.h>
typedef struct ints_t {
  // Needs to be capital to enable access front outside of the package
  int32_t A;
  int32_t B;
  int32_t C;
  int32_t D;
  int64_t E;
} Ints_t;

typedef struct __attribute__((__packed__)) {
  char c;
  short s;
  int b;
  long long e;
  char c2;
} Packed_ints_t;
