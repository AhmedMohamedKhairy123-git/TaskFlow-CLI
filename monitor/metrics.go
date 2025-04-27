package monitor

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	mu             sync.RWMutex
	requestCount   int64
	errorCount     int64
	taskCount      int
	avgResponseTime time.Duration
	lastUpdated    time.Time
}

var globalMetrics = &Metrics{
	lastUpdated: time.Now(),
}

func RecordRequest(duration time.Duration) {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	
	globalMetrics.requestCount++
	globalMetrics.avgResponseTime = (globalMetrics.avgResponseTime*time.Duration(globalMetrics.requestCount-1) + duration) / time.Duration(globalMetrics.requestCount)
	globalMetrics.lastUpdated = time.Now()
}

func RecordError() {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	globalMetrics.errorCount++
}

func UpdateTaskCount(count int) {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	globalMetrics.taskCount = count
}

type PerformanceReport struct {
	Requests      int64
	Errors        int64
	ErrorRate     float64
	AvgResponse   time.Duration
	TaskCount     int
	Goroutines    int
	MemoryMB      uint64
	LastUpdated   time.Time
}

func GetPerformanceReport() PerformanceReport {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	report := PerformanceReport{
		Requests:    globalMetrics.requestCount,
		Errors:      globalMetrics.errorCount,
		AvgResponse: globalMetrics.avgResponseTime,
		TaskCount:   globalMetrics.taskCount,
		Goroutines:  runtime.NumGoroutine(),
		MemoryMB:    m.Alloc / 1024 / 1024,
		LastUpdated: globalMetrics.lastUpdated,
	}
	
	if report.Requests > 0 {
		report.ErrorRate = float64(report.Errors) / float64(report.Requests) * 100
	}
	
	return report
}

func (pr PerformanceReport) String() string {
	return fmt.Sprintf(
		"📊 Performance Report:\n"+
		"   Requests: %d\n"+
		"   Errors: %d (%.2f%%)\n"+
		"   Avg Response: %v\n"+
		"   Tasks: %d\n"+
		"   Goroutines: %d\n"+
		"   Memory: %d MB\n"+
		"   Updated: %s",
		pr.Requests, pr.Errors, pr.ErrorRate,
		pr.AvgResponse, pr.TaskCount,
		pr.Goroutines, pr.MemoryMB,
		pr.LastUpdated.Format("15:04:05"),
	)
}

// Middleware for metrics
func MetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			RecordRequest(time.Since(start))
		}()
		next(w, r)
	}
}

// Add monitoring endpoints
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	report := GetPerformanceReport()
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, report.String())
}