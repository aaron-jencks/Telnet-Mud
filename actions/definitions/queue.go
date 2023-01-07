package definitions

import (
	"sync"
)

type ActionQueue struct {
	Contents       chan Action
	MaxLength      int
	QueueWriteLock sync.Mutex
	QueueReadLock  sync.Mutex
}

func CreateActionQueue(maxCapacity int) *ActionQueue {
	return &ActionQueue{
		Contents:       make(chan Action, maxCapacity),
		MaxLength:      maxCapacity,
		QueueWriteLock: sync.Mutex{},
		QueueReadLock:  sync.Mutex{},
	}
}

func (q *ActionQueue) Enqueue(a Action) {
	q.QueueWriteLock.Lock()
	defer q.QueueWriteLock.Unlock()

	q.Contents <- a
}

func (q *ActionQueue) Dequeue() Action {
	q.QueueReadLock.Lock()
	defer q.QueueReadLock.Unlock()

	return <-q.Contents
}
