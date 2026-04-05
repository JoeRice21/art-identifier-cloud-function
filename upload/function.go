package upload

import (
	"encoding/json"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("ImageUploadHandler", ImageUploadHandler)
}

func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := NewTraceContext(r.Context())
	logger := GetLogger(ctx, "controller::ImageUploadHandler")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error("failed to parse form", "error", err)
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		logger.Error("failed to get file", "error", err)
		http.Error(w, "failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	logger.Info("received image", "filename", header.Filename, "size", header.Size)

	signedURL, err := uploadToGCS(ctx, file, header.Filename)
	if err != nil {
		logger.Error("failed to upload to GCS", "error", err)
		http.Error(w, "failed to upload image", http.StatusInternalServerError)
		return
	}
	logger.Info("image uploaded successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"signedURL": signedURL})
}