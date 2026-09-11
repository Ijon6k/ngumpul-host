package storage

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ngumpul-host/backend/internal/config"
)

func createTestPNG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{R: 121, G: 175, B: 196, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func TestLocalStorage_Operations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "storage-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to init local storage: %v", err)
	}

	ctx := context.Background()
	key := "test/hello.txt"
	content := []byte("hello ngumpul host")

	// 1. Put
	if err := store.Put(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain"); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// 2. Exists
	exists, err := store.Exists(ctx, key)
	if err != nil || !exists {
		t.Fatalf("Expected exists=true, got %v (err: %v)", exists, err)
	}

	// 3. Get
	rc, _, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	readBytes, err := io.ReadAll(rc)
	rc.Close()
	if err != nil || !bytes.Equal(readBytes, content) {
		t.Fatalf("Get content mismatch: got %s", string(readBytes))
	}

	// 4. PublicURL
	url := store.PublicURL(key)
	if url != "/uploads/test/hello.txt" {
		t.Fatalf("Unexpected public URL: %s", url)
	}

	// 5. Delete
	if err := store.Delete(ctx, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	existsAfter, _ := store.Exists(ctx, key)
	if existsAfter {
		t.Fatalf("Expected file to be deleted")
	}
}

func TestUploadHandler_ValidImage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upload-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{StorageDriver: "local", StoragePath: tempDir}
	svc, err := NewService(cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create storage service: %v", err)
	}

	pngData := createTestPNG(120, 90)

	// Create multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("purpose", "cover")
	part, _ := writer.CreateFormFile("file", "cover.png")
	_, _ = part.Write(pngData)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	svc.UploadHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUploadHandler_RejectSVG(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "upload-svg-test-*")
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{StorageDriver: "local", StoragePath: tempDir}
	svc, _ := NewService(cfg, nil)

	svgData := []byte("<svg xmlns='http://www.w3.org/2000/svg'><script>alert(1)</script></svg>")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "malicious.svg")
	_, _ = part.Write(svgData)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	svc.UploadHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected HTTP 400 rejection for SVG, got %d", rec.Code)
	}
}

func TestUploadHandler_RejectOversizedDimensions(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "upload-dim-test-*")
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{StorageDriver: "local", StoragePath: tempDir}
	svc, _ := NewService(cfg, nil)

	// Create header of PNG with 5000x5000 px dimensions
	hugePNG := createTestPNG(10, 10)
	// Mutate IHDR chunk width (bytes 16-19) to 5000 (0x00, 0x00, 0x13, 0x88)
	if len(hugePNG) > 24 {
		hugePNG[16] = 0x00
		hugePNG[17] = 0x00
		hugePNG[18] = 0x13
		hugePNG[19] = 0x88
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "huge.png")
	_, _ = part.Write(hugePNG)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	svc.UploadHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected HTTP 400 for dimensions > 4096, got %d", rec.Code)
	}
}

func TestParseWebPDimensions(t *testing.T) {
	// Extended WebP header mockup: RIFF (4) + size (4) + WEBP (4) + VP8X (4) + size (4) + flags (4) + w (3) + h (3)
	mockVP8X := []byte{
		'R', 'I', 'F', 'F',
		0, 0, 0, 0,
		'W', 'E', 'B', 'P',
		'V', 'P', '8', 'X',
		10, 0, 0, 0,
		0, 0, 0, 0,
		0x7F, 0x03, 0x00, // width = 895 + 1 = 896
		0xDF, 0x01, 0x00, // height = 479 + 1 = 480
	}

	w, h, err := parseWebPDimensions(mockVP8X)
	if err != nil {
		t.Fatalf("parseWebPDimensions failed: %v", err)
	}
	if w != 896 || h != 480 {
		t.Fatalf("Expected 896x480, got %dx%d", w, h)
	}
}
