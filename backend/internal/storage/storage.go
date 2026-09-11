package storage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ngumpul-host/backend/internal/config"
)

// Storage defines the uniform object storage contract for Ngumpul Host.
type Storage interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	PublicURL(key string) string
}

// LocalStorage implements Storage using the local container/host filesystem.
type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create local storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (l *LocalStorage) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	fullPath := filepath.Join(l.basePath, filepath.Clean(key))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, r); err != nil {
		return fmt.Errorf("failed to write local file: %w", err)
	}
	return nil
}

func (l *LocalStorage) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	fullPath := filepath.Join(l.basePath, filepath.Clean(key))
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, "", err
	}
	return f, "", nil
}

func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(l.basePath, filepath.Clean(key))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (l *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := filepath.Join(l.basePath, filepath.Clean(key))
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (l *LocalStorage) PublicURL(key string) string {
	return fmt.Sprintf("/uploads/%s", strings.TrimPrefix(key, "/"))
}

// S3Storage implements Storage using an S3-compatible API (SeaweedFS, MinIO, R2) via Go standard library.
type S3Storage struct {
	endpoint   string
	bucket     string
	region     string
	accessKey  string
	secretKey  string
	publicURL  string
	httpClient *http.Client
}

func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	endpoint := strings.TrimRight(cfg.S3Endpoint, "/")
	if endpoint == "" {
		endpoint = "http://seaweedfs:8333"
	}

	bucket := cfg.S3Bucket
	if bucket == "" {
		bucket = "ngumpul-uploads"
	}

	region := cfg.S3Region
	if region == "" {
		region = "us-east-1"
	}

	publicURL := cfg.S3PublicURL
	if publicURL == "" {
		publicURL = "/uploads"
	}

	s := &S3Storage{
		endpoint:   endpoint,
		bucket:     bucket,
		region:     region,
		accessKey:  cfg.S3AccessKey,
		secretKey:  cfg.S3SecretKey,
		publicURL:  strings.TrimRight(publicURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	return s, nil
}

func (s *S3Storage) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	cleanKey := strings.TrimPrefix(key, "/")
	reqURL := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, cleanKey)

	// Read body into memory if needed for SigV4 payload hash
	var bodyReader io.Reader = r
	var payloadHash string

	if s.accessKey != "" && s.secretKey != "" {
		bodyBytes, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("failed to read payload: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
		hash := sha256.Sum256(bodyBytes)
		payloadHash = hex.EncodeToString(hash[:])
		size = int64(len(bodyBytes))
	} else {
		payloadHash = "UNSIGNED-PAYLOAD"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create s3 request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if size > 0 {
		req.ContentLength = size
	}

	s.signRequest(req, payloadHash)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("s3 put failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("s3 put failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	cleanKey := strings.TrimPrefix(key, "/")
	reqURL := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, cleanKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, "", err
	}

	s.signRequest(req, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855") // empty sha256

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, "", os.ErrNotExist
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, "", fmt.Errorf("s3 get failed with status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	return resp.Body, contentType, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	cleanKey := strings.TrimPrefix(key, "/")
	reqURL := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, cleanKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}

	s.signRequest(req, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("s3 delete failed with status %d", resp.StatusCode)
	}
	return nil
}

func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	cleanKey := strings.TrimPrefix(key, "/")
	reqURL := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, cleanKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, reqURL, nil)
	if err != nil {
		return false, err
	}

	s.signRequest(req, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("s3 head returned status %d", resp.StatusCode)
}

func (s *S3Storage) PublicURL(key string) string {
	return fmt.Sprintf("%s/%s", s.publicURL, strings.TrimPrefix(key, "/"))
}

// signRequest applies standard AWS SigV4 authentication headers if keys are provided.
func (s *S3Storage) signRequest(req *http.Request, payloadHash string) {
	if s.accessKey == "" || s.secretKey == "" {
		return
	}

	now := time.Now().UTC()
	dateStamp := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")

	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)

	u, _ := url.Parse(s.endpoint)
	host := u.Host
	req.Header.Set("Host", host)

	canonicalURI := req.URL.Path
	if canonicalURI == "" {
		canonicalURI = "/"
	}

	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", host, payloadHash, amzDate)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.Method,
		canonicalURI,
		req.URL.RawQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	)

	canonHash := sha256.Sum256([]byte(canonicalRequest))
	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, s.region)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", amzDate, credentialScope, hex.EncodeToString(canonHash[:]))

	signingKey := getSignatureKey(s.secretKey, dateStamp, s.region, "s3")
	signature := hmacSHA256(signingKey, stringToSign)

	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKey,
		credentialScope,
		signedHeaders,
		hex.EncodeToString(signature),
	)

	req.Header.Set("Authorization", authHeader)
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func getSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+key), dateStamp)
	kRegion := hmacSHA256(kDate, regionName)
	kService := hmacSHA256(kRegion, serviceName)
	kSigning := hmacSHA256(kService, "aws4_request")
	return kSigning
}
