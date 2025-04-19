package monitoring

import (
	"embed"
	"encoding/json"
	"net/http"
	"text/template"
)

//go:embed templates/*.html
var templatesFS embed.FS

var (
	indexTemplate = template.Must(template.New("index.html").Funcs(FuncMap).ParseFS(templatesFS, "templates/index.html"))
	funcTemplate  = template.Must(template.New("function.html").Funcs(FuncMap).ParseFS(templatesFS, "templates/function.html"))
)

func (m *Monitor) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	metrics := m.GetAllFunctionMetrics()
	err := indexTemplate.Execute(w, metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *Monitor) HandleFunction(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Function name is required", http.StatusBadRequest)
		return
	}

	metrics := m.GetFunctionMetrics(name)
	if metrics == nil {
		http.Error(w, "Function not found", http.StatusNotFound)
		return
	}

	err := funcTemplate.Execute(w, metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *Monitor) StartHTTPServer(addr string) error {
	// Static file server for assets (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Main routes
	http.HandleFunc("/", m.HandleIndex)
	http.HandleFunc("/function", m.HandleFunction)

	// API routes
	http.HandleFunc("/api/metrics", m.handleAllMetrics)
	http.HandleFunc("/api/function/", m.handleFunctionMetrics)

	return http.ListenAndServe(addr, nil)
}

func (m *Monitor) handleAllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := m.GetAllFunctionMetrics()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (m *Monitor) handleFunctionMetrics(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/api/function/"):]
	if name == "" {
		http.Error(w, "Function name is required", http.StatusBadRequest)
		return
	}

	metrics := m.GetFunctionMetrics(name)
	if metrics == nil {
		http.Error(w, "Function not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
