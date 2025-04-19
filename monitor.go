package monitoring

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

type CallMetrics struct {
	Timestamp time.Time `json:"timestamp"`
	Duration  float64   `json:"duration"`
	MemAlloc  uint64    `json:"mem_alloc"`
}

// FunctionMetrics holds per-function metrics
type FunctionMetrics struct {
	Name        string        `json:"name"`
	Duration    float64       `json:"duration"`
	AvgDuration float64       `json:"avg_duration"`
	MemAlloc    uint64        `json:"mem_alloc"`
	Invocations int           `json:"invocations"`
	Calls       []CallMetrics `json:"calls"` // 🆕 per-call details

}

// Monitor tracks system and function metrics
type Monitor struct {
	mu         sync.Mutex
	functions  map[string]*FunctionMetrics
	startTime  time.Time
	cpuProfile *os.File
}

type FuncMonitor struct {
	name      string
	startTime time.Time
	startMem  uint64
	monitor   *Monitor
}

func NewMonitor(addr string) *Monitor {
	if addr == "" {
		// If no address is provided, start the HTTP server on a default port
		addr = ":8080"
	}

	m := &Monitor{
		functions: make(map[string]*FunctionMetrics),
		startTime: time.Now(),
	}

	go func() {
		fmt.Println("Starting HTTP server on", addr)
		if err := m.StartHTTPServer(addr); err != nil {
			panic(err)
		}
	}()

	return m

}

func (m *Monitor) RegisterFunc(name string) *FuncMonitor {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return &FuncMonitor{
		name:      name,
		startTime: time.Now(),
		startMem:  mem.Alloc,
		monitor:   m,
	}
}

func (f *FuncMonitor) Finish() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	duration := time.Since(f.startTime)
	memUsed := mem.Alloc - f.startMem

	f.monitor.mu.Lock()
	defer f.monitor.mu.Unlock()

	if metrics, exists := f.monitor.functions[f.name]; exists {
		metrics.Duration += duration.Seconds()
		metrics.MemAlloc += memUsed
		metrics.Invocations++
		metrics.Calls = append(metrics.Calls, CallMetrics{
			Timestamp: time.Now(),
			Duration:  duration.Seconds(),
			MemAlloc:  memUsed,
		})
	} else {
		f.monitor.functions[f.name] = &FunctionMetrics{
			Name:        f.name,
			Duration:    duration.Seconds(),
			MemAlloc:    memUsed,
			Invocations: 1,
			Calls: []CallMetrics{
				{
					Timestamp: time.Now(),
					Duration:  duration.Seconds()	,
					MemAlloc:  memUsed,
				},
			},
		}
	}
}

func (m *Monitor) GetMetrics() map[string]*FunctionMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics := make(map[string]*FunctionMetrics)
	for name, metric := range m.functions {
		metric.AvgDuration = metric.Duration / float64(metric.Invocations)
		metrics[name] = metric
	}
	return metrics
}

func (m *Monitor) GetFunctionMetrics(name string) *FunctionMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metric, exists := m.functions[name]; exists {
		metric.AvgDuration = metric.Duration / float64(metric.Invocations)
		return metric
	}
	return nil
}

func (m *Monitor) GetSystemMetrics() map[string]interface{} {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return map[string]interface{}{
		"uptime":      time.Since(m.startTime).Seconds(),
		"goroutines":  runtime.NumGoroutine(),
		"allocations": mem.Alloc,
		"totalAlloc":  mem.TotalAlloc,
		"sys":         mem.Sys,
		"numGC":       mem.NumGC,
	}
}

func (m *Monitor) StartCPUProfile() error {
	if m.cpuProfile != nil {
		return nil
	}

	f, err := os.Create("cpu.prof")
	if err != nil {
		return err
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return err
	}

	m.cpuProfile = f
	return nil
}

func (m *Monitor) StopCPUProfile() {
	if m.cpuProfile == nil {
		return
	}

	pprof.StopCPUProfile()
	m.cpuProfile.Close()
	m.cpuProfile = nil
}

func (m *Monitor) WriteHeapProfile() error {
	f, err := os.Create("heap.prof")
	if err != nil {
		return err
	}
	defer f.Close()

	return pprof.WriteHeapProfile(f)
}

func (m *Monitor) WriteBlockProfile() error {
	f, err := os.Create("block.prof")
	if err != nil {
		return err
	}
	defer f.Close()

	return pprof.Lookup("block").WriteTo(f, 0)
}

func (m *Monitor) WriteMutexProfile() error {
	f, err := os.Create("mutex.prof")
	if err != nil {
		return err
	}
	defer f.Close()

	return pprof.Lookup("mutex").WriteTo(f, 0)
}

func (m *Monitor) WriteProfile() error {
	if err := m.WriteHeapProfile(); err != nil {
		return err
	}
	if err := m.WriteBlockProfile(); err != nil {
		return err
	}
	if err := m.WriteMutexProfile(); err != nil {
		return err
	}

	return nil
}

func (m *Monitor) Stop() {
	m.StopCPUProfile()
}

func (m *Monitor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.functions = make(map[string]*FunctionMetrics)
	m.startTime = time.Now()

	if m.cpuProfile != nil {
		m.cpuProfile.Close()
		m.cpuProfile = nil
	}
}

func (m *Monitor) GetAllFunctionMetrics() []*FunctionMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics := make([]*FunctionMetrics, 0, len(m.functions))
	for _, metric := range m.functions {
		metric.AvgDuration = metric.Duration / float64(metric.Invocations)
		metrics = append(metrics, metric)
	}
	return metrics
}

func (m *Monitor) PrintAllFunctionMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, metric := range m.functions {
		metric.AvgDuration = metric.Duration / float64(metric.Invocations)
		fmt.Printf("Function: %s, Total Duration: %.2f s, Avg Duration: %.2f s, Mem Alloc: %d MB, Invocations: %d\n",
			metric.Name, metric.Duration, metric.AvgDuration, metric.MemAlloc/(1024*1024), metric.Invocations)

		for _, call := range metric.Calls {
			fmt.Printf("  Call at %s, Duration: %.2f s, Mem Alloc: %d MB\n",
				call.Timestamp.Format(time.DateTime), call.Duration, call.MemAlloc/(1024*1024))
		}
		fmt.Println()
	}
}
