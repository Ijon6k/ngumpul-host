package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"ngumpul-host/backend/internal/config"
)

type StorageService struct {
	cfg *config.Config
}

func NewStorageService(cfg *config.Config) (*StorageService, error) {
	if err := os.MkdirAll(cfg.StoragePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &StorageService{cfg: cfg}, nil
}

func (s *StorageService) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Max 5 MB upload size
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "File size exceeds maximum allowed (5MB)"})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing file field"})
		return
	}
	defer file.Close()

	// Sniff MIME type
	buff := make([]byte, 512)
	n, err := file.Read(buff)
	if err != nil && err != io.EOF {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Could not read file header"})
		return
	}

	mimeType := http.DetectContentType(buff[:n])
	var ext string
	switch mimeType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	case "image/gif":
		ext = ".gif"
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Unsupported file type: %s. Only JPEG, PNG, WebP, and GIF images are permitted.", mimeType),
		})
		return
	}

	// Rewind file pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to process upload"})
		return
	}

	// Safe unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	targetPath := filepath.Join(s.cfg.StoragePath, filename)

	dst, err := os.Create(targetPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write file"})
		return
	}

	publicURL := fmt.Sprintf("/uploads/%s", filename)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"url":           publicURL,
		"filename":      filename,
		"original_name": filepath.Base(handler.Filename),
		"mime_type":     mimeType,
	})
}

func (s *StorageService) ServeHandler() http.Handler {
	fileServer := http.FileServer(http.Dir(s.cfg.StoragePath))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip /uploads/ prefix
		p := strings.TrimPrefix(r.URL.Path, "/uploads/")
		// Guard against traversal
		if strings.Contains(p, "..") {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		http.StripPrefix("/uploads/", fileServer).ServeHTTP(w, r)
	})
}
