package local

import (
	"container/list"
	"errors"
	"sync"

	"github.com/davidzech/raidbutts/pkg/job"
)

type Queue struct {
	m    sync.Mutex
	list *list.List
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func (q *Queue) Publish(j *job.Job) error {
	q.m.Lock()
	defer q.m.Unlock()
	q.list.PushBack(j)
	return nil
}

func (q *Queue) Consume() (*job.Job, error) {
	q.m.Lock()
	defer q.m.Unlock()
	if q.list.Len() == 0 {
		return nil, errors.New("empty queue")
	}
	job := q.list.Remove(q.list.Back()).(*job.Job)
	return job, nil
}

func (q *Queue) Len() int {
	return q.list.Len()
}
