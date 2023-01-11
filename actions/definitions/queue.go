package definitions

type ActionQueue struct {
	Contents  chan Action
	MaxLength int
}

func CreateActionQueue(maxCapacity int) *ActionQueue {
	return &ActionQueue{
		Contents:  make(chan Action, maxCapacity),
		MaxLength: maxCapacity,
	}
}

func (q *ActionQueue) Enqueue(a Action) {
	q.Contents <- a
}

func (q *ActionQueue) Dequeue() Action {
	return <-q.Contents
}
