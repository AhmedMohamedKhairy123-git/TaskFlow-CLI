package task

import (
	"sort"
	"time"
)

type QueryPlan struct {
	Strategy     string
	EstimatedCost int
	IndexUsed    string
	FilterFirst  bool
}

type QueryOptimizer struct {
	stats *QueryStatistics
}

type QueryStatistics struct {
	TotalTasks    int
	PriorityDist  map[Priority]int
	TagDist       map[string]int
	CompleteRatio float64
}

func NewQueryOptimizer(store *TaskStore) *QueryOptimizer {
	return &QueryOptimizer{
		stats: store.collectStatistics(),
	}
}

func (s *TaskStore) collectStatistics() *QueryStatistics {
	stats := &QueryStatistics{
		TotalTasks:   len(s.Tasks),
		PriorityDist: make(map[Priority]int),
		TagDist:      make(map[string]int),
	}
	
	completed := 0
	for _, task := range s.Tasks {
		stats.PriorityDist[task.Priority]++
		
		if task.Completed {
			completed++
		}
		
		for _, tag := range task.Tags {
			stats.TagDist[tag]++
		}
	}
	
	if stats.TotalTasks > 0 {
		stats.CompleteRatio = float64(completed) / float64(stats.TotalTasks)
	}
	
	return stats
}

func (q *QueryOptimizer) OptimizeSearch(criteria SearchCriteria) QueryPlan {
	plan := QueryPlan{}
	
	// Choose best index
	if len(criteria.Priority) > 0 {
		if q.stats.PriorityDist[criteria.Priority[0]] < q.stats.TotalTasks/2 {
			plan.Strategy = "Use priority index"
			plan.IndexUsed = "priority"
			plan.EstimatedCost = q.stats.PriorityDist[criteria.Priority[0]]
			return plan
		}
	}
	
	if len(criteria.Tags) > 0 {
		tag := criteria.Tags[0]
		if q.stats.TagDist[tag] < q.stats.TotalTasks/2 {
			plan.Strategy = "Use tag index"
			plan.IndexUsed = "tag"
			plan.EstimatedCost = q.stats.TagDist[tag]
			return plan
		}
	}
	
	if criteria.Completed != nil {
		if q.stats.CompleteRatio < 0.3 || q.stats.CompleteRatio > 0.7 {
			plan.Strategy = "Use completion index"
			plan.IndexUsed = "complete"
			if *criteria.Completed {
				plan.EstimatedCost = int(float64(q.stats.TotalTasks) * q.stats.CompleteRatio)
			} else {
				plan.EstimatedCost = int(float64(q.stats.TotalTasks) * (1 - q.stats.CompleteRatio))
			}
			return plan
		}
	}
	
	plan.Strategy = "Full table scan"
	plan.IndexUsed = "none"
	plan.EstimatedCost = q.stats.TotalTasks
	return plan
}

func (s *TaskStore) OptimizedSearch(criteria SearchCriteria) SearchResult {
	opt := NewQueryOptimizer(s)
	plan := opt.OptimizeSearch(criteria)
	
	start := time.Now()
	var results []Task
	
	switch plan.IndexUsed {
	case "priority":
		ids := s.priorityIndex.Get(criteria.Priority[0])
		for _, id := range ids {
			if task, exists := s.Tasks[id]; exists {
				if matchCriteria(*task, criteria) {
					results = append(results, *task)
				}
			}
		}
	case "tag":
		ids := s.tagIndex.Get(criteria.Tags[0])
		for _, id := range ids {
			if task, exists := s.Tasks[id]; exists {
				if matchCriteria(*task, criteria) {
					results = append(results, *task)
				}
			}
		}
	default:
		// Full scan
		for _, task := range s.Tasks {
			if matchCriteria(*task, criteria) {
				results = append(results, *task)
			}
		}
	}
	
	return SearchResult{
		Tasks:       results,
		TotalCount:  len(s.Tasks),
		FilteredCount: len(results),
		SearchTime:  time.Since(start),
	}
}