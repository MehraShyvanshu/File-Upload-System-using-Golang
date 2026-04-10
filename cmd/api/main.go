package main

import (
	"log"
	"net/http"
	"project/internal/handlers"
	"project/internal/infrastructure"
	"project/internal/usecase"
)

// corsWrapper wraps a handler function to add CORS headers
type corsWrapper struct {
	handler http.HandlerFunc
}

func (cw *corsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "3600")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	cw.handler(w, r)
}

func main() {

	// initialize repository (SQLite)
	repo := infrastructure.NewSQLiteRepo()

	// initialize usecase
	usecase := usecase.NewFileUsecase(repo)

	// initialize handlers
	fileHandler := handlers.NewFileHandler(usecase)
	authHandler := handlers.NewAuthHandler()

	// Setup routes with CORS middleware
	mux := http.NewServeMux()
	mux.Handle("/", &corsWrapper{handler: fileHandler.Root})
	mux.Handle("/health", &corsWrapper{handler: fileHandler.Health})
	mux.Handle("/login", &corsWrapper{handler: authHandler.Login})
	mux.Handle("/upload", &corsWrapper{handler: authHandler.AuthMiddleware(fileHandler.Upload)})
	mux.Handle("/files", &corsWrapper{handler: authHandler.AuthMiddleware(fileHandler.ListFiles)})

	log.Println("Server running at http://localhost:8080")
	log.Println("Available endpoints:")
	log.Println("  GET  /         - API info")
	log.Println("  GET  /health   - Health check")
	log.Println("  POST /login    - User login (returns JWT token)")
	log.Println("  POST /upload   - Upload file (requires Authorization header)")
	log.Println("  GET  /files    - List all files (requires Authorization header)")
	log.Println("")

	http.ListenAndServe(":8080", mux)
}
