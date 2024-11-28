package task

import (
	"sync"
)

type BulkResult struct {
	SuccessCount int
	FailedCount  int
	Errors       map[int]error
}

func (s *TaskStore) BulkAdd(titles []string) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for i, title := range titles {
		_, err := s.Add(title)
		if err != nil {
			result.Errors[i] = err
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}
	
	return result
}

func (s *TaskStore) BulkComplete(ids []int) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for _, id := range ids {
		if err := s.MarkComplete(id); err != nil {
			result.Errors[id] = err
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}
	
	return result
}

func (s *TaskStore) BulkDelete(ids []int) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for _, id := range ids {
		if err := s.Delete(id); err != nil {
			result.Errors[id] = err
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}
	
	return result
}

func (s *TaskStore) BulkUpdate(updates map[int]Task) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for id, updatedTask := range updates {
		if err := s.Update(id, updatedTask); err != nil {
			result.Errors[id] = err
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}
	
	return result
}

func (s *TaskStore) ParallelBulkComplete(ids []int) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for _, id := range ids {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			
			err := s.MarkComplete(taskID)
			
			mu.Lock()
			if err != nil {
				result.Errors[taskID] = err
				result.FailedCount++
			} else {
				result.SuccessCount++
			}
			mu.Unlock()
		}(id)
	}
	
	wg.Wait()
	return result
}
// Add to task/bulk.go

type BatchProcessor struct {
	batchSize int
	store     *TaskStore
	queue     chan func() error
	wg        sync.WaitGroup
}

func NewBatchProcessor(store *TaskStore, batchSize int, workers int) *BatchProcessor {
	bp := &BatchProcessor{
		batchSize: batchSize,
		store:     store,
		queue:     make(chan func() error, 100),
	}
	
	// Start workers
	for i := 0; i < workers; i++ {
		bp.wg.Add(1)
		go bp.worker()
	}
	
	return bp
}

func (bp *BatchProcessor) worker() {
	defer bp.wg.Done()
	for task := range bp.queue {
		task()
	}
}

func (bp *BatchProcessor) AddTask(task func() error) {
	bp.queue <- task
}

func (bp *BatchProcessor) Stop() {
	close(bp.queue)
	bp.wg.Wait()
}

func (s *TaskStore) BatchAdd(titles []string, batchSize int) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for i := 0; i < len(titles); i += batchSize {
		end := i + batchSize
		if end > len(titles) {
			end = len(titles)
		}
		
		batch := titles[i:end]
		
		// Process batch
		for j, title := range batch {
			_, err := s.Add(title)
			if err != nil {
				result.Errors[i+j] = err
				result.FailedCount++
			} else {
				result.SuccessCount++
			}
		}
	}
	
	return result
}

func (s *TaskStore) BatchProcessWithWorkers(titles []string, workers int) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	ch := make(chan struct {
		index int
		title string
	}, len(titles))
	
	// Queue all items
	for i, title := range titles {
		ch <- struct {
			index int
			title string
		}{i, title}
	}
	close(ch)
	
	// Start workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range ch {
				_, err := s.Add(item.title)
				mu.Lock()
				if err != nil {
					result.Errors[item.index] = err
					result.FailedCount++
				} else {
					result.SuccessCount++
				}
				mu.Unlock()
			}
		}()
	}
	
	wg.Wait()
	return result
}