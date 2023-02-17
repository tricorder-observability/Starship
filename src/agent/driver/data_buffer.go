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
