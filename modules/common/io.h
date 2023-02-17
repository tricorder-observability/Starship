// LINT_C_FILE

#pragma once

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// This file is meant to be included by a WASM C/C++ program,
// and provides APIs for userspace driver to allocate memory for the WASM
// program.
//
// The WASM program is expected to export all of the non-static functions
// defined here.
//
// The WASM program's main data process API needs to call get_input_buf()
// to get the linear memory as input, and get_output_buf() to get output memory.

// Describes a linear memory buffer. Used for sharing data with userspace code.
struct buffer_t {
  // Points to the memory buffer.
  void *data;

  // Actual amount of data written to the above pointed address.
  // TODO(yzhao): The plan is to let userspace write length, and let WASM code
  // to update This field is not used right now.
  uint32_t length;

  // The total capacity of this memory buffer.
  uint32_t capacity;
};

// Used for input and output.
static struct buffer_t input_buf = {};
static struct buffer_t output_buf = {};

static void malloc_buf(uint32_t capacity, struct buffer_t *buf) {
  buf->data = malloc(capacity);
  buf->length = 0;
  buf->capacity = capacity;
}

static void free_buf(struct buffer_t *buf) {
  free(buf->data);
  buf->data = NULL;
  buf->length = 0;
  buf->capacity = 0;
}

void *malloc_input_buf(uint32_t capacity) {
  malloc_buf(capacity, &input_buf);
  return input_buf.data;
}
void *malloc_output_buf(uint32_t capacity) {
  malloc_buf(capacity, &output_buf);
  return output_buf.data;
}

void free_input_buf() { free_buf(&input_buf); }
void free_output_buf() { free_buf(&input_buf); }

void *get_input_buf() { return input_buf.data; }
void *get_output_buf() { return output_buf.data; }

uint32_t get_input_buf_len() { return input_buf.length; }
uint32_t get_input_buf_cap() { return input_buf.capacity; }
uint32_t get_output_buf_len() { return output_buf.length; }
uint32_t get_output_buf_cap() { return output_buf.capacity; }

void set_input_buf_len(uint32_t len) { input_buf.length = len; }
void set_output_buf_len(uint32_t len) { output_buf.length = len; }

// Returns non-0 if output buffer's free space is larger than the input length.
// The return value is used as bool, but WASM does not allow bool type.
int can_write_to_output_buf(uint32_t len) {
  return output_buf.length + len <= output_buf.capacity;
}

// Blindly copy data to the output buffer.
// Caller must check the capacity.
void copy_to_output(const void *data, uint32_t len) {
  memcpy(output_buf.data + output_buf.length, data, len);
  output_buf.length += len;
}

// This does not work, as the memory layout of struct buffer_t in WASM is
// different than C. struct buffer_t* get_input_buf() {
//    return &input_buf;
// }
// struct buffer_t* get_output_buf() {
//    return &output_buf;
// }

// This does not work, only return int32 or int64 (need to run on this option to
// use 64 bit memory space)
// struct buffer_t get_input_buf() {
//   return input_buf;
// }
//
// The following procedure also does not work, it does not work because there is
// no WAT instruction to return memory address:
// 1. Build a wat function, and insert it into the runtime extension
// 2. Allocate memory in the wasmtime runtime API from outside of the wasm
// runtime
// 3. Let the wasm data processing code call this wat function to get the
// allocated linear memory.
