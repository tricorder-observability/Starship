// LINT_C_FILE
#pragma once
#define TASK_COMM_LEN 13

struct event_t {
  float F;
  char C;
  double D;
  int I;
  long long int L;
  short Comm[TASK_COMM_LEN];
};
