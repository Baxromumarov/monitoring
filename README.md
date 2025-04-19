# Monitoring Dashboard for Go Functions

A beautiful, user-friendly monitoring dashboard for your Go applications. Instantly visualize function metrics (duration, memory, invocations) in real time, with zero dependencies beyond the Go standard library.


---

## ✨ Features
- **Real-time Dashboard**: See all monitored functions, their total/average duration, memory allocation, and invocation counts.
- **Detailed Function View**: Click any function to see per-call history (timestamp, duration, memory in MB).
- **Live Updating**: Dashboard auto-refreshes every second.
- **Modern UI**: Clean, responsive Go template frontend.
- **Easy Integration**: Just import and wrap your functions.

---

## 🚀 Getting Started

### 1. Clone and Run
```bash
git clone https://github.com/baxromumarov/monitoring.git
cd monitoring
go run cmd/main.go
```

Visit [http://localhost:8080](http://localhost:8080) in your browser.

### 2. Integrate into Your Project

Import the monitor package and wrap your functions:
```go
import "github.com/baxromumarov/monitoring"

func main() {
    monitor := monitoring.NewMonitor(":8080", time.Minute * 10)
    defer monitor.Stop()

    // Wrap your function
    fn := monitor.RegisterFunc("MyFunction")
    defer fn.Finish()
    // ... your code ...
}
```

That's it! All calls will be tracked and visualized.

---

## 📊 What You See
- **Main Page**: Table of all monitored functions with:
  - Function Name (clickable)
  - Total Duration (s)
  - Average Duration (s)
  - Memory Allocation (MB)
  - Number of Invocations
- **Function Detail Page**: Per-call history with timestamp, duration, and memory allocation.

---

## 🛠️ Customization
- Edit `templates/index.html` and `templates/function.html` for UI tweaks.
- Add/remove metrics in `monitor.go` as needed.
- Style with CSS in the templates or add static assets.

---

## 🤝 Contributing
Pull requests and issues are welcome! Please open an issue for bugs or suggestions.

---

Enjoy monitoring! 🚀
