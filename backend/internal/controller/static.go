package controller

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Static(staticDir string, api http.Handler) http.Handler {
	if strings.TrimSpace(staticDir) == "" {
		return api
	}

	if _, err := os.Stat(filepath.Join(staticDir, "index.html")); err != nil {
		return api
	}

	fileServer := http.FileServer(http.Dir(staticDir))
	indexPath := filepath.Join(staticDir, "index.html")
	staticFS := os.DirFS(staticDir)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if strings.HasPrefix(requestPath, "/api/") || strings.HasPrefix(requestPath, "/telegram/") || requestPath == "/healthz" {
			api.ServeHTTP(w, r)

			return
		}

		if requestPath == "/" {
			http.ServeFile(w, r, indexPath)

			return
		}

		cleanPath := filepath.Clean(strings.TrimPrefix(requestPath, "/"))
		if !fs.ValidPath(cleanPath) {
			http.ServeFile(w, r, indexPath)

			return
		}

		info, err := fs.Stat(staticFS, cleanPath)
		if err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)

			return
		}

		http.ServeFile(w, r, indexPath)
	})
}
