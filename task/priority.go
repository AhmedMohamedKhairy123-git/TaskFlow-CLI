package task

import (
	"errors"
	"sort"
)

type PriorityQueue struct {
	tasks []*Task
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		tasks: make([]*Task, 0),
	}
}

func (pq *PriorityQueue) Push(task *Task) {
	pq.tasks = append(pq.tasks, task)
	pq.sort()
}

func (pq *PriorityQueue) Pop() (*Task, error) {
	if len(pq.tasks) == 0 {
		return nil, errors.New("queue is empty")
	}
	task := pq.tasks[0]
	pq.tasks = pq.tasks[1:]
	return task, nil
}

func (pq *PriorityQueue) Peek() (*Task, error) {
	if len(pq.tasks) == 0 {
		return nil, errors.New("queue is empty")
	}
	return pq.tasks[0], nil
}

func (pq *PriorityQueue) sort() {
	sort.Slice(pq.tasks, func(i, j int) bool {
		if pq.tasks[i].Priority == pq.tasks[j].Priority {
			return pq.tasks[i].CreatedAt.Before(pq.tasks[j].CreatedAt)
		}
		return pq.tasks[i].Priority > pq.tasks[j].Priority
	})
}

func (pq *PriorityQueue) Len() int {
	return len(pq.tasks)
}

// SortByPriority implements sort.Interface for []Task
type SortByPriority []Task

func (s SortByPriority) Len() int      { return len(s) }
func (s SortByPriority) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortByPriority) Less(i, j int) bool {
	if s[i].Priority == s[j].Priority {
		return s[i].CreatedAt.Before(s[j].CreatedAt)
	}
	return s[i].Priority > s[j].Priority
}