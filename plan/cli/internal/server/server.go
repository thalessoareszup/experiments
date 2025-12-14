package server

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Start(addr string, dev bool, webDir string) error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/events", handleSSE)
	mux.HandleFunc("/api/plans", handlePlans)
	mux.HandleFunc("/api/plans/", handlePlanByID)
	mux.HandleFunc("/api/steps/", handleStepByID)

	// Static files or dev proxy
	if dev {
		// Proxy to Vite dev server
		viteURL, _ := url.Parse("http://localhost:5173")
		proxy := httputil.NewSingleHostReverseProxy(viteURL)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
	} else {
		// Serve static files
		mux.HandleFunc("/", staticHandler(webDir))
	}

	// Wrap with CORS middleware
	handler := corsMiddleware(mux)

	return http.ListenAndServe(addr, handler)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func staticHandler(webDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		fullPath := filepath.Join(webDir, path)

		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// SPA fallback - serve index.html for non-file routes
			if !strings.Contains(path, ".") {
				fullPath = filepath.Join(webDir, "index.html")
			}
		}

		// Determine content type
		ext := filepath.Ext(fullPath)
		contentTypes := map[string]string{
			".html": "text/html",
			".js":   "application/javascript",
			".css":  "text/css",
			".json": "application/json",
			".svg":  "image/svg+xml",
			".png":  "image/png",
			".ico":  "image/x-icon",
		}

		if ct, ok := contentTypes[ext]; ok {
			w.Header().Set("Content-Type", ct)
		}

		file, err := os.Open(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		io.Copy(w, file)
	}
}
