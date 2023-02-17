/*
 *  Copyright (C) 2023, TriCorder, All Rights Reserved
 *  Copyright (C) 2015, Jhuster, All Rights Reserved
 *  Author: Jhuster(lujun.hust@gmail.com)
 *
 *  https://github.com/Jhuster/TLV
 *
 *  This library is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU Lesser General Public License as published
 *  by the Free Software Foundation; either version 2.1 of the License,
 *  or (at your option) any later version.
 */
#include "tlv_box.h"
#include <stdio.h>
#include <string.h>

#define TEST_TYPE_0 0x00
#define TEST_TYPE_1 0x01
#define TEST_TYPE_2 0x02
#define TEST_TYPE_3 0x03
#define TEST_TYPE_4 0x04
#define TEST_TYPE_5 0x05
#define TEST_TYPE_6 0x06
#define TEST_TYPE_7 0x07
#define TEST_TYPE_8 0x08
#define TEST_TYPE_9 0x09

#define LOG(format, ...) printf(format, ##__VA_ARGS__)

// Create and serialize two TLV objects,
// * the first one encodes some built-int type values;
// * the second one encodes the first TLV object.
//
// Return:
// * 0 on success, non-zero values on failures;
// * also return the two serialized objects with pointers in the argument list.
int test_encoding(tlv_box_t **pbox1, tlv_box_t **pbox2) {
  // Box1
  LOG("Creating tlv box1\n");
  tlv_box_t *box1 = tlv_box_create();

  tlv_box_put_char(box1, TEST_TYPE_1, 'x');
  tlv_box_put_short(box1, TEST_TYPE_2, (short)2);
  tlv_box_put_int(box1, TEST_TYPE_3, (int)3);
  tlv_box_put_long(box1, TEST_TYPE_4, (long)4);
  tlv_box_put_float(box1, TEST_TYPE_5, (float)5.67);
  tlv_box_put_double(box1, TEST_TYPE_6, (double)8.91);
  tlv_box_put_string(box1, TEST_TYPE_7, (char *)"hello world!");
  unsigned char array[6] = {1, 2, 3, 4, 5, 6};
  tlv_box_put_bytes(box1, TEST_TYPE_8, array, 6);

  if (tlv_box_serialize(box1)) {
    LOG("TLV box1 serialization failed !\n");
    return -1;
  }

  LOG("TLV box1 serialization successful, %d bytes\n", tlv_box_get_size(box1));

  // Box2
  LOG("Creating tlv box2\n");
  tlv_box_t *box2 = tlv_box_create();

  tlv_box_put_object(box2, TEST_TYPE_9, box1);

  if (tlv_box_serialize(box2)) {
    LOG("TLV box2 serialization failed!\n");
    return -1;
  }

  LOG("TLV box2 serialization successful, %d bytes\n", tlv_box_get_size(box2));

  *pbox1 = box1;
  *pbox2 = box2;
  return 0;
}

// Decode the seriazlied TLV object `box2` that's created in test_encoding()
int test_decoding(tlv_box_t *box2) {
  tlv_box_t *parsedBox2 =
      tlv_box_parse(tlv_box_get_buffer(box2), tlv_box_get_size(box2));
  LOG("Parse tlv box2 successful, %dbytes\n", tlv_box_get_size(parsedBox2));

  tlv_box_t *parsedBox1;
  if (tlv_box_get_object(parsedBox2, TEST_TYPE_9, &parsedBox1)) {
    LOG("tlv_box_get_object failed!\n");
    return -1;
  }

  LOG("Parse tlv box1 successful, %d bytes\n", tlv_box_get_size(parsedBox1));

  {
    char value;
    if (tlv_box_get_char(parsedBox1, TEST_TYPE_1, &value)) {
      LOG("tlv_box_get_char failed!\n");
      return -1;
    }
    LOG("tlv_box_get_char successful %c\n", value);
  }

  {
    short value;
    if (tlv_box_get_short(parsedBox1, TEST_TYPE_2, &value)) {
      LOG("tlv_box_get_short failed !\n");
      return -1;
    }
    LOG("tlv_box_get_short successful %d\n", value);
  }

  {
    int value;
    if (tlv_box_get_int(parsedBox1, TEST_TYPE_3, &value)) {
      LOG("tlv_box_get_int failed !\n");
      return -1;
    }
    LOG("tlv_box_get_int successful %d\n", value);
  }

  {
    long value;
    if (tlv_box_get_long(parsedBox1, TEST_TYPE_4, &value)) {
      LOG("tlv_box_get_long failed !\n");
      return -1;
    }
    LOG("tlv_box_get_long successful %ld\n", value);
  }

  {
    float value;
    if (tlv_box_get_float(parsedBox1, TEST_TYPE_5, &value)) {
      LOG("tlv_box_get_float failed !\n");
      return -1;
    }
    LOG("tlv_box_get_float successful %f\n", value);
  }

  {
    double value;
    if (tlv_box_get_double(parsedBox1, TEST_TYPE_6, &value)) {
      LOG("tlv_box_get_double failed !\n");
      return -1;
    }
    LOG("tlv_box_get_double successful %f\n", value);
  }

  {
    char value[128];
    int length = 128;
    if (tlv_box_get_string(parsedBox1, TEST_TYPE_7, value, &length)) {
      LOG("tlv_box_get_string failed !\n");
      return -1;
    }
    LOG("tlv_box_get_string successful %s\n", value);
  }

  {
    unsigned char value[128];
    int length = 128;
    if (tlv_box_get_bytes(parsedBox1, TEST_TYPE_8, value, &length)) {
      LOG("tlv_box_get_bytes failed !\n");
      return -1;
    }

    LOG("tlv_box_get_bytes successful:  ");
    int i = 0;
    for (i = 0; i < length; i++) {
      LOG("%d-", value[i]);
    }
    LOG("\n");
  }

  tlv_box_destroy(parsedBox1);
  tlv_box_destroy(parsedBox2);

  return 0;
}

int cleanup(tlv_box_t *box1, tlv_box_t *box2) {
  tlv_box_destroy(box1);
  tlv_box_destroy(box2);
  return 0;
}

int main() {
  tlv_box_t *box1;
  tlv_box_t *box2;

  // Test encoding
  LOG("Test encoding\n");
  if (test_encoding(&box1, &box2)) {
    LOG("Test encoding failed\n");
    return -1;
  }

  LOG("Test encoding successful\n\n");

  // Test decoding
  LOG("Test decoding\n");
  if (test_decoding(box2)) {
    LOG("Test decoding failed\n");
    return -1;
  }

  LOG("Test decoding successful\n\n");

  // Cleanup
  cleanup(box1, box2);
  LOG("Cleanup successful\n");

  return 0;
}
