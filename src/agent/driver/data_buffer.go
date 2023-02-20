// Copyright (C) 2023  tricorder-observability
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

type DataBuffer struct {
	q *Queue
}

func NewDefaultDataBuffer() *DataBuffer {
	return &DataBuffer{q: NewQueue(DefaultBufferSize)}
}

func (d *DataBuffer) Produce(pollData map[string][]byte) error {
	for _, data := range pollData {
		err := d.q.Enqueue(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DataBuffer) Consume() []byte {
	if !d.q.HasData() {
		return nil
	}
	e, err := d.q.Dequeue()
	if err != nil {
		return nil
	}
	return e
}
