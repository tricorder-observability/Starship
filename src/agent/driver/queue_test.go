package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQueue tests Enqueue and Dequeue methods
func TestQueue(t *testing.T) {
	assert := assert.New(t)

	q := NewQueue(10)
	data, err := q.Dequeue()
	assert.Nil(data)
	assert.NotNil(err)

	err = q.Enqueue([]byte("01234"))
	assert.Nil(err)

	err = q.Enqueue([]byte("56789"))
	assert.Nil(err)

	// Can write empty data
	err = q.Enqueue([]byte{})
	assert.Nil(err)

	err = q.Enqueue([]byte("01234"))
	assert.NotNil(err)

	data, err = q.Dequeue()
	assert.Nil(err)
	assert.Equal(data, []byte("01234"))

	data, err = q.Dequeue()
	assert.Nil(err)
	assert.Equal(data, []byte("56789"))

	data, err = q.Dequeue()
	assert.Nil(err)
	assert.Equal(data, []byte{})
}
