package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/config"
	"ngumpul-host/backend/internal/response"
)

const (
	// MaxUploadSizeBytes enforces the hard 10 MB limit as specified in revisi-besar.md
	MaxUploadSizeBytes = 10 << 20
	// MaxDimensionPixels limits image dimensions to prevent decompression bombs
	MaxDimensionPixels = 4096
	// MaxTotalPixels limits overall decoded pixels (16 Megapixels)
	MaxTotalPixels = 16 * 1024 * 1024
)

// Service manages object storage, image validation, upload processing, and asset serving.
type Service struct {
	cfg     *config.Config
	storage Storage
	db      *pgxpool.Pool
}

// NewService creates a storage service using either S3Storage (SeaweedFS) or LocalStorage.
func NewService(cfg *config.Config, db *pgxpool.Pool) (*Service, error) {
	var st Storage
	var err error

	if strings.EqualFold(cfg.StorageDriver, "s3") {
		st, err = NewS3Storage(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize s3 storage: %w", err)
		}
	} else {
		st, err = NewLocalStorage(cfg.StoragePath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize local storage: %w", err)
		}
	}

	return &Service{
		cfg:     cfg,
		storage: st,
		db:      db,
	}, nil
}

// Storage returns the underlying Storage interface.
func (s *Service) Storage() Storage {
	return s.storage
}

// UploadHandler processes image uploads with strict 10 MB limits, MIME sniffing, and dimension bounds.
func (s *Service) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce strict 10 MB maximum upload size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSizeBytes)
	if err := r.ParseMultipartForm(MaxUploadSizeBytes); err != nil {
		response.Error(w, http.StatusBadRequest, "File size exceeds maximum allowed limit of 10 MB")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Missing 'file' field in multipart upload")
		return
	}
	defer file.Close()

	// 2. Read entire file into buffer (capped at 10 MB)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to read uploaded file")
		return
	}

	fileSize := int64(len(fileBytes))
	if fileSize == 0 {
		response.Error(w, http.StatusBadRequest, "Uploaded file is empty")
		return
	}

	// 3. Inspect actual MIME content type (never trust header alone)
	sniffLimit := 512
	if len(fileBytes) < sniffLimit {
		sniffLimit = len(fileBytes)
	}
	detectedMIME := http.DetectContentType(fileBytes[:sniffLimit])

	var ext string
	var width, height int

	switch detectedMIME {
	case "image/jpeg":
		ext = ".jpg"
		cfg, _, err := image.DecodeConfig(bytes.NewReader(fileBytes))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Malformed or unreadable JPEG image")
			return
		}
		width, height = cfg.Width, cfg.Height

	case "image/png":
		ext = ".png"
		cfg, _, err := image.DecodeConfig(bytes.NewReader(fileBytes))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Malformed or unreadable PNG image")
			return
		}
		width, height = cfg.Width, cfg.Height

	case "image/webp":
		ext = ".webp"
		wD, hD, err := parseWebPDimensions(fileBytes)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Malformed or unreadable WebP image")
			return
		}
		width, height = wD, hD

	default:
		// Explicitly reject SVGs and non-supported binary formats per security guidelines
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("Unsupported file format: %s. Only JPEG, PNG, and WebP images are permitted.", detectedMIME))
		return
	}

	// 4. Prevent decompression bomb / resource exhaustion attacks
	if width <= 0 || height <= 0 {
		response.Error(w, http.StatusBadRequest, "Unable to determine valid image dimensions")
		return
	}
	if width > MaxDimensionPixels || height > MaxDimensionPixels {
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("Image dimensions (%dx%d) exceed the maximum allowed limit of %dx%d pixels.", width, height, MaxDimensionPixels, MaxDimensionPixels))
		return
	}
	if int64(width)*int64(height) > MaxTotalPixels {
		response.Error(w, http.StatusBadRequest, "Total image pixel count exceeds safe threshold")
		return
	}

	// 5. Generate safe, unguessable UUID-based object key
	purpose := strings.ToLower(r.FormValue("purpose"))
	var prefix string
	switch purpose {
	case "cover", "project_cover":
		prefix = "covers"
	case "avatar":
		prefix = "avatars"
	default:
		prefix = "uploads"
	}

	objectKey := fmt.Sprintf("%s/%s%s", prefix, uuid.New().String(), ext)

	// 6. Write object to Storage backend (SeaweedFS / S3 or LocalStorage)
	if err := s.storage.Put(r.Context(), objectKey, bytes.NewReader(fileBytes), fileSize, detectedMIME); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to store uploaded media")
		return
	}

	publicURL := s.storage.PublicURL(objectKey)

	// 7. Record metadata in PostgreSQL media_objects table if available
	currentUser := auth.GetUser(r.Context())
	if s.db != nil {
		var ownerID any = nil
		if currentUser != nil {
			ownerID = currentUser.ID
		}
		_, _ = s.db.Exec(context.Background(), `
			INSERT INTO media_objects (owner_id, purpose, object_key, content_type, byte_size, width, height)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (object_key) DO NOTHING
		`, ownerID, strings.ToUpper(purpose), objectKey, detectedMIME, fileSize, width, height)
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"url":           publicURL,
		"object_key":    objectKey,
		"byte_size":     fileSize,
		"content_type":  detectedMIME,
		"width":         width,
		"height":        height,
		"original_name": filepath.Base(handler.Filename),
	})
}

// ServeHandler returns an HTTP handler for streaming stored media files.
func (s *Service) ServeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/uploads/")
		if key == "" || strings.Contains(key, "..") {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		rc, contentType, err := s.storage.Get(r.Context(), key)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer rc.Close()

		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.Header().Set("Cache-Control", "public, max-age=2592000, no-transform")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		_, _ = io.Copy(w, rc)
	})
}

// parseWebPDimensions extracts canvas dimensions from standard lossy, lossless, and extended WebP headers.
func parseWebPDimensions(b []byte) (int, int, error) {
	if len(b) < 30 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WEBP" {
		return 0, 0, fmt.Errorf("invalid webp signature")
	}

	format := string(b[12:16])
	switch format {
	case "VP8 ": // Lossy
		if len(b) < 30 {
			return 0, 0, fmt.Errorf("truncated vp8 payload")
		}
		w := int(b[26]) | (int(b[27]&0x3f) << 8)
		h := int(b[28]) | (int(b[29]&0x3f) << 8)
		return w, h, nil

	case "VP8L": // Lossless
		if len(b) < 25 {
			return 0, 0, fmt.Errorf("truncated vp8l payload")
		}
		val := uint32(b[21]) | (uint32(b[22]) << 8) | (uint32(b[23]) << 16) | (uint32(b[24]) << 24)
		w := int(val&0x3fff) + 1
		h := int((val>>14)&0x3fff) + 1
		return w, h, nil

	case "VP8X": // Extended
		if len(b) < 30 {
			return 0, 0, fmt.Errorf("truncated vp8x payload")
		}
		w := int(b[24]) | (int(b[25]) << 8) | (int(b[26]) << 16) + 1
		h := int(b[27]) | (int(b[28]) << 8) | (int(b[29]) << 16) + 1
		return w, h, nil

	default:
		return 0, 0, fmt.Errorf("unrecognized webp chunk: %s", format)
	}
}
