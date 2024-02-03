package task

import (
	"strings"
	"time"
)

type SearchCriteria struct {
	Title     string
	Completed *bool
	Priority  []Priority
	Tags      []string
	FromDate  *time.Time
	ToDate    *time.Time
}

type SearchResult struct {
	Tasks       []Task
	TotalCount  int
	FilteredCount int
	SearchTime  time.Duration
}

func (s *TaskStore) Search(criteria SearchCriteria) SearchResult {
	start := time.Now()
	var results []Task
	
	for _, task := range s.Tasks {
		if matchCriteria(*task, criteria) {
			results = append(results, *task)
		}
	}
	
	return SearchResult{
		Tasks:       results,
		TotalCount:  len(s.Tasks),
		FilteredCount: len(results),
		SearchTime:  time.Since(start),
	}
}

func matchCriteria(t Task, c SearchCriteria) bool {
	if c.Title != "" && !strings.Contains(strings.ToLower(t.Title), strings.ToLower(c.Title)) {
		return false
	}
	
	if c.Completed != nil && t.Completed != *c.Completed {
		return false
	}
	
	if len(c.Priority) > 0 {
		match := false
		for _, p := range c.Priority {
			if t.Priority == p {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	
	if len(c.Tags) > 0 {
		tagMatch := false
		for _, searchTag := range c.Tags {
			for _, taskTag := range t.Tags {
				if strings.EqualFold(taskTag, searchTag) {
					tagMatch = true
					break
				}
			}
		}
		if !tagMatch {
			return false
		}
	}
	
	if c.FromDate != nil && t.CreatedAt.Before(*c.FromDate) {
		return false
	}
	
	if c.ToDate != nil && t.CreatedAt.After(*c.ToDate) {
		return false
	}
	
	return true
}