// LINT_C_FILE
#define TASK_COMM_LEN 13
typedef struct {
  float X;
  char B;
  double Y;
  int Z;
  long long int A;
  short Comm[TASK_COMM_LEN];
  void *Unused_ptr;
} Event_t;
