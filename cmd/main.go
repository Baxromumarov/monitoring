package main

import (
	"fmt"
	"log"
	"time"

	"github.com/baxromumarov/monitoring"
)

// // SimulatedWork represents a function that does some work
// func SimulatedWork(name string, monitor *monitoring.Monitor) {
// 	fn := monitor.RegisterFunc(name)
// 	defer fn.Finish()
// 	time.Sleep(time.Second * 2)

// 	// Allocate some random memory
// 	data := make([]byte, rand.Intn(1024*1024)) // 0-1MB
// 	_ = data
// }

func main() {
	// Create a new monitor
	monitor := monitoring.NewMonitoring("", 10*time.Minute)
	defer monitor.Stop()

	// Start CPU profiling
	if err := monitor.StartCPUProfile(); err != nil {
		log.Fatal(err)
	}
	defer monitor.StopCPUProfile()

	// Run some simulated work in the background
	// go func() {
	// 	functions := []string{"ProcessData", "CalculateMetrics", "UpdateDatabase", "GenerateReport"}
	// 	for {
	// 		for _, name := range functions {
	// 			SimulatedWork(name, monitor)
	// 		}
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	fmt.Println("Server started at http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop")

	// Keep the main goroutine running
	select {}
}
