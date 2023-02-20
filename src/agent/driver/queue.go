// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package driver

import (
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/pkg/errors"
)

const DefaultBufferSize = 1000

var (
	ErrEmpty  = errors.New("cannot read data when the queue is empty")
	ErrFull   = errors.New("cannot write data when the queue is full")
	ErrClosed = errors.New("cannot enqueue or dequeue when the queue is closed")
	ErrData   = errors.New("cannot enqueue or dequeue when element types un-matched")
)

type Queue struct {
	// The maximum buffer event size.
	capacity int

	// The size of the elements of this queue in bytes
	sizeByte int

	// Underlying buffer
	buffer *goconcurrentqueue.FIFO
}

// NewQueue returns a sized queue and return an error if has errors during creation.
func NewQueue(bufferSize int) *Queue {
	return &Queue{
		capacity: bufferSize,
		buffer:   goconcurrentqueue.NewFIFO(),
	}
}

type QueueElement struct {
	Name string
	Msg  []byte
}

func (q *Queue) Enqueue(data []byte) error {
	newSize := q.sizeByte + len(data)
	if newSize > q.capacity {
		return ErrFull
	}
	err := q.buffer.Enqueue(data)
	if err != nil {
		return ErrFull
	}
	q.sizeByte = newSize
	return nil
}

func (q *Queue) Dequeue() ([]byte, error) {
	e, err := q.buffer.Dequeue()
	if err != nil {
		return nil, ErrEmpty
	}
	data, ok := e.([]byte)
	if !ok {
		return nil, ErrData
	}
	return data, nil
}

func (q *Queue) Capacity() int {
	return q.capacity
}

func (q *Queue) SizeByte() int {
	return q.sizeByte
}

func (q *Queue) IsFull() bool {
	return q.SizeByte() >= q.Capacity()
}

func (q *Queue) IsEmpty() bool {
	return q.SizeByte() == 0
}

func (q *Queue) HasData() bool {
	return q.SizeByte() >= 0
}
