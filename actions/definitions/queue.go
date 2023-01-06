package definitions

import (
	"mud/utils/ui/logger"
	"sync"
)

type ActionQueue struct {
	Contents       []Action
	MaxLength      int
	QueueWriteLock sync.Mutex
	QueueReadLock  sync.Mutex
	QueueSizeLock  sync.Mutex
	QueueEmptyLock sync.Mutex
}

func CreateActionQueue(maxCapacity int) *ActionQueue {
	return &ActionQueue{
		MaxLength:      maxCapacity,
		QueueWriteLock: sync.Mutex{},
		QueueReadLock:  sync.Mutex{},
		QueueSizeLock:  sync.Mutex{},
		QueueEmptyLock: sync.Mutex{},
	}
}

func (q *ActionQueue) waitForSpace() {
	if q.Length() == q.MaxLength {
		logger.Info("Queue is full, waiting for size lock")
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
	logger.Info("Getting Write Lock")
	q.QueueWriteLock.Lock()
	defer q.QueueWriteLock.Unlock()

	logger.Info("Waiting for space")
	q.waitForSpace()

	logger.Info("Getting Read Lock")
	q.QueueReadLock.Lock()
	defer q.QueueReadLock.Unlock()

	logger.Info("Inserting new action")
	q.Contents = append([]Action{a}, q.Contents...)

	if q.Length() == q.MaxLength {
		q.QueueSizeLock.Lock()
	} else if q.Length() == 1 {
		q.QueueEmptyLock.Unlock()
	}
}

func (q *ActionQueue) Length() int {
	return len(q.Contents)
}

func (q *ActionQueue) Dequeue() Action {
	logger.Info("Fetching next action")

	q.QueueReadLock.Lock()
	defer q.QueueReadLock.Unlock()

	needsUnlock := q.Length() == q.MaxLength

	if q.Length() > 0 {
		logger.Info("Dequeueing element")

		item := q.Contents[0]
		q.Contents = q.Contents[1:]

		if needsUnlock {
			logger.Info("Unlocking queue size lock")
			q.QueueSizeLock.Unlock()
		}

		if q.Length() == 0 {
			q.QueueEmptyLock.Lock()
		}

		return item
	}

	logger.Info("Queue was empty")
	return ACTION_NOT_FOUND
}
