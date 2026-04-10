package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"project/internal/usecase"
)

type FileHandler struct {
	usecase *usecase.FileUsecase
}

func NewFileHandler(u *usecase.FileUsecase) *FileHandler {
	return &FileHandler{usecase: u}
}

// Health check endpoint
func (h *FileHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// Root endpoint
func (h *FileHandler) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "File Upload API",
		"version":   "1.0.0",
		"endpoints": "GET /health, POST /upload, GET /files",
	})
}

// List all uploaded files
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	files, err := h.usecase.ListFiles()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to list files"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}

// handles upload
func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("DEBUG: Upload handler called")
	log.Printf("DEBUG: Request method: %s, Content-Type: %s", r.Method, r.Header.Get("Content-Type"))

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("DEBUG: ParseMultipartForm error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error parsing form: " + err.Error()})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("DEBUG: FormFile error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error reading file: " + err.Error()})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("DEBUG: ReadAll error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error reading file bytes: " + err.Error()})
		return
	}

	log.Printf("DEBUG: Saving file: %s, size: %d bytes", handler.Filename, len(fileBytes))
	err = h.usecase.SaveFile(handler.Filename, fileBytes)
	if err != nil {
		log.Printf("DEBUG: SaveFile error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save: " + err.Error()})
		return
	}

	log.Println("DEBUG: File uploaded successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
}
