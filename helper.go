package monitoring

import (
	"fmt"
	"text/template"
	"time"
)

var (
	FuncMap = template.FuncMap{
		"durationSeconds": func(d time.Duration) string {
			return fmt.Sprintf("%.2f", d.Seconds())
		},
		"formatMB": func(b uint64) string {
			return fmt.Sprintf("%.2f", float64(b)/1024.0/1024.0)
		},
		"formatBytes": func(b uint64) string {
			const unit = 1024
			if b < unit {
				return fmt.Sprintf("%d B", b)
			}
			div, exp := uint64(unit), 0
			for n := b / unit; n >= unit; n /= unit {
				div *= unit
				exp++
			}
			return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
		},
		"formatDuration": func(d time.Duration) string {
			if d < time.Millisecond {
				return fmt.Sprintf("%.2f µs", float64(d.Microseconds()))
			} else if d < time.Second {
				return fmt.Sprintf("%.2f ms", float64(d.Milliseconds()))
			}
			return fmt.Sprintf("%.2f s", d.Seconds())
		},
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"div": func(a, b int64) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
	}
)

