package task

import (
	"time"
)

type Reminder struct {
	TaskID    int
	DueDate   time.Time
	Repeat    RepeatType
	Notified  bool
}

type RepeatType string

const (
	RepeatNone   RepeatType = "none"
	RepeatDaily  RepeatType = "daily"
	RepeatWeekly RepeatType = "weekly"
	RepeatMonthly RepeatType = "monthly"
)

type TaskWithReminder struct {
	Task
	Reminder *Reminder
}

func (s *TaskStore) SetReminder(taskID int, dueDate time.Time, repeat RepeatType) error {
	task, exists := s.Tasks[taskID]
	if !exists {
		return fmt.Errorf("task %d not found", taskID)
	}
	
	task.Reminder = &Reminder{
		TaskID:   taskID,
		DueDate:  dueDate,
		Repeat:   repeat,
		Notified: false,
	}
	
	return nil
}

func (s *TaskStore) GetDueTasks() []TaskWithReminder {
	var due []TaskWithReminder
	now := time.Now()
	
	for _, task := range s.Tasks {
		if task.Reminder != nil && !task.Reminder.Notified {
			if task.Reminder.DueDate.Before(now) {
				due = append(due, TaskWithReminder{*task, task.Reminder})
			}
		}
	}
	
	return due
}

func (s *TaskStore) CheckReminders() []TaskWithReminder {
	var notified []TaskWithReminder
	now := time.Now()
	
	for _, task := range s.Tasks {
		if task.Reminder != nil && !task.Reminder.Notified {
			if task.Reminder.DueDate.Before(now) {
				task.Reminder.Notified = true
				notified = append(notified, TaskWithReminder{*task, task.Reminder})
				
				if task.Reminder.Repeat != RepeatNone {
					s.scheduleNextReminder(task)
				}
			}
		}
	}
	
	return notified
}

func (s *TaskStore) scheduleNextReminder(task *Task) {
	if task.Reminder == nil {
		return
	}
	
	var next time.Time
	switch task.Reminder.Repeat {
	case RepeatDaily:
		next = task.Reminder.DueDate.AddDate(0, 0, 1)
	case RepeatWeekly:
		next = task.Reminder.DueDate.AddDate(0, 0, 7)
	case RepeatMonthly:
		next = task.Reminder.DueDate.AddDate(0, 1, 0)
	default:
		return
	}
	
	task.Reminder = &Reminder{
		TaskID:   task.ID,
		DueDate:  next,
		Repeat:   task.Reminder.Repeat,
		Notified: false,
	}
}