package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/vandenbill/social-media-10k-rps/internal/cfg"
	"github.com/vandenbill/social-media-10k-rps/internal/dto"
)

type fileHandler struct {
	cfg *cfg.Cfg
}

func newFileHandler(cfg *cfg.Cfg) *fileHandler {
	return &fileHandler{cfg}
}
func (h *fileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get fromFile form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" {
		http.Error(w, "File must be in JPG or JPEG format", http.StatusBadRequest)
		return
	}

	if handler.Size < 10*1024 || handler.Size > 2*1024*1024 { // 10 KB to 2 MB
		http.Error(w, "File size must be between 10 KB and 2 MB", http.StatusBadRequest)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(h.cfg.S3Region),
		Credentials: credentials.NewStaticCredentials(h.cfg.S3ID, h.cfg.S3SecretKey, ""),
	})
	if err != nil {
		http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
		return
	}

	svc := s3.New(sess)

	fileName := fmt.Sprintf("%s%s", uuid.NewString(), ext)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(h.cfg.S3BucketName),
		Key:    aws.String(fileName),
		ACL:    aws.String("public-read"),
		Body:   file,
	})
	if err != nil {
		http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
		return
	}

	res := dto.ResUpFile{
		ImageUrl: fmt.Sprintf("https://%s.s3.amazonaws.com/%s", h.cfg.S3BucketName, fileName),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
