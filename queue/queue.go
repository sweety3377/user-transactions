package queue

import (
	"blackwallgroup/internal/model"
	"sync"
)

const queueSize = 1

type Queue struct {
	queue map[string]chan struct{}
	sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{
		queue:   map[string]chan struct{}{},
		RWMutex: sync.RWMutex{},
	}
}

func (q *Queue) Put(req model.CreateTransactionRequest) {
	q.Lock()
	defer q.Unlock()
	q.queue[req.ID] = make(chan struct{}, queueSize)
}

func (q *Queue) Get(req model.CreateTransactionRequest) chan struct{} {
	q.RLocker()
	defer q.RUnlock()
	return q.queue[req.ID]
}

func (q *Queue) Wait(req model.CreateTransactionRequest) {
	q.RLocker()
	defer q.RUnlock()
	<-q.queue[req.ID]
}

func (q *Queue) Release(req model.CreateTransactionRequest) {
	q.Lock()
	defer q.Unlock()
	q.queue[req.ID] <- struct{}{}
}
