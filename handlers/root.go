package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ksKao/dota/lib"
)

func RegisterRootRoutes() {
	http.HandleFunc("GET /", rootGetHandler)
	http.HandleFunc("POST /", rootPostHandler)
}

func rootGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "templates/index.html")
}

func rootPostHandler(w http.ResponseWriter, r *http.Request) {
	const maxSize = 500 << 20 // 500 MiB (500 * 1024 * 1024 bytes)
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		http.Error(w, "file too large (max 500 mb)", http.StatusRequestEntityTooLarge)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > maxSize {
		http.Error(w, "file too large (max 500 mb)", http.StatusRequestEntityTooLarge)
		return
	}

	if !strings.HasSuffix(header.Filename, ".dem") {
		http.Error(w, "not a replay file", http.StatusBadRequest)
		return
	}

	start, hooks, err := lib.ParseDemFile(file)
	if err != nil {
		http.Error(w, "fail to parse replay file", http.StatusInternalServerError)
		return
	}

	output := map[string]any{
		"start": start,
		"hooks": hooks,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
