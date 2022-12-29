package actions

import (
	"sync"
)

type ActionQueue struct {
	Contents       []Action
	MaxLength      int
	QueueWriteLock sync.Mutex
	QueueReadLock  sync.Mutex
	QueueSizeLock  sync.Mutex
}

func (q *ActionQueue) waitForSpace() {
	if q.Length() == q.MaxLength {
		q.QueueSizeLock.Lock()
		defer q.QueueSizeLock.Unlock()
	}
}

func (q *ActionQueue) Enqueue(a Action) {
	q.QueueWriteLock.Lock()
	defer q.QueueWriteLock.Unlock()

	q.waitForSpace()

	q.QueueReadLock.Lock()
	defer q.QueueReadLock.Unlock()

	q.Contents = append(q.Contents, a)

	if q.Length() == q.MaxLength {
		q.QueueSizeLock.Lock()
	}
}

func (q *ActionQueue) Push(a Action) {
	q.QueueWriteLock.Lock()
	defer q.QueueWriteLock.Unlock()

	q.waitForSpace()

	q.QueueReadLock.Lock()
	defer q.QueueReadLock.Unlock()

	q.Contents = append([]Action{a}, q.Contents...)

	if q.Length() == q.MaxLength {
		q.QueueSizeLock.Lock()
	}
}

func (q *ActionQueue) Length() int {
	return len(q.Contents)
}

func (q *ActionQueue) Dequeue() Action {
	q.QueueReadLock.Lock()
	defer q.QueueReadLock.Unlock()

	needsUnlock := q.Length() == q.MaxLength

	if q.Length() > 0 {
		item := q.Contents[0]
		q.Contents = q.Contents[1:]

		if needsUnlock {
			q.QueueSizeLock.Unlock()
		}

		return item
	}
	return Action{}
}
