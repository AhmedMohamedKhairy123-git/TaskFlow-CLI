package task

import (
	"sort"
	"time"
)

type Analytics struct {
	TotalTasks      int
	CompletedTasks  int
	PendingTasks    int
	CompletionRate  float64
	ByPriority      map[Priority]int
	ByTag           map[string]int
	AveragePerDay   float64
	OldestTask      *Task
	NewestTask      *Task
	LongestTitle    string
}

func (s *TaskStore) GetAnalytics() Analytics {
	analytics := Analytics{
		ByPriority: make(map[Priority]int),
		ByTag:      make(map[string]int),
	}
	
	analytics.TotalTasks = len(s.Tasks)
	
	var oldestTime time.Time
	var newestTime time.Time
	
	for _, task := range s.Tasks {
		if task.Completed {
			analytics.CompletedTasks++
		}
		
		analytics.ByPriority[task.Priority]++
		
		for _, tag := range task.Tags {
			analytics.ByTag[tag]++
		}
		
		if len(task.Title) > len(analytics.LongestTitle) {
			analytics.LongestTitle = task.Title
		}
		
		if oldestTime.IsZero() || task.CreatedAt.Before(oldestTime) {
			oldestTime = task.CreatedAt
			analytics.OldestTask = task
		}
		
		if task.CreatedAt.After(newestTime) {
			newestTime = task.CreatedAt
			analytics.NewestTask = task
		}
	}
	
	analytics.PendingTasks = analytics.TotalTasks - analytics.CompletedTasks
	
	if analytics.TotalTasks > 0 {
		analytics.CompletionRate = float64(analytics.CompletedTasks) / float64(analytics.TotalTasks) * 100
	}
	
	if analytics.TotalTasks > 0 {
		firstDay := oldestTime.Truncate(24 * time.Hour)
		days := int(time.Since(firstDay).Hours() / 24)
		if days > 0 {
			analytics.AveragePerDay = float64(analytics.TotalTasks) / float64(days)
		}
	}
	
	return analytics
}

func (s *TaskStore) GetTopTags(limit int) []struct {
	Tag   string
	Count int
} {
	var tags []struct {
		Tag   string
		Count int
	}
	
	counts := make(map[string]int)
	for _, task := range s.Tasks {
		for _, tag := range task.Tags {
			counts[tag]++
		}
	}
	
	for tag, count := range counts {
		tags = append(tags, struct {
			Tag   string
			Count int
		}{tag, count})
	}
	
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Count > tags[j].Count
	})
	
	if len(tags) > limit {
		tags = tags[:limit]
	}
	
	return tags
}