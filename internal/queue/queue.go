package queue

import "container/list"

type Queue struct {
	list *list.List
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func (q *Queue) Add(elem any) {
	if elem == nil {
		return
	}
	q.list.PushBack(elem)
}

func (q *Queue) Remove() {
	q.Poll()
}

func (q *Queue) Poll() any {
	elem := q.list.Front()
	if elem == nil {
		return nil
	}
	return q.list.Remove(elem)
}

func (q *Queue) Length() int {
	return q.list.Len()
}
