package selection

import "errors"

/*
[BERITH]
Circular queue structure for binarySearch by random value
*/
type Queue struct {
	storage []Range
	size    int
	front   int
	rear    int
}

func (q *Queue) setQueueAsCandidates(candidateCount int) *Queue {
	return &Queue{
		storage: make([]Range, candidateCount + 1),
		size:    candidateCount + 1,
		front:   0,
		rear:    0,
	}
}

func (q *Queue) enqueue(r Range) error {
	if (q.rear + 1) % q.size == q.front {
		return errors.New("Queue is full")
	}

	q.storage[q.rear] = r
	q.rear = (q.rear + 1) % q.size
	return nil
}

func (q *Queue) dequeue() (Range, error) {
	if q.front == q.rear {
		return Range{}, errors.New("Queue is empty")
	}

	result := q.storage[q.front]
	q.front = (q.front + 1) % q.size
	return result, nil
}