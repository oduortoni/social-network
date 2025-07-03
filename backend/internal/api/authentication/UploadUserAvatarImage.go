package authentication

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func UploadAvatarImage(imagereader multipart.File, imageheader *multipart.FileHeader) (string, error) {
	const MaxUploadSize = 20 * 1024 * 1024 // 20MB limit

	if imageheader.Size > MaxUploadSize {
		return "maximum size", nil
	}
	// Get the file extension (e.g., .jpg, .png)
	ext := filepath.Ext(imageheader.Filename)
	filename := uuid.New().String()
	filepath := filepath.Join("./UserAvatars", filename+ext)

	// Create the directory if it doesn't exist
	err := os.MkdirAll("./UserAvatars", os.ModePerm)
	if err != nil {
		return "failed to create UserAvatars directory", err
	}

	out, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return "failed to open the file", err
	}
	defer out.Close()

	writtien, errs := io.Copy(out, imagereader)
	if errs != nil {
		return "failed to copy", err
	}
	if writtien > MaxUploadSize {
		return "maximum size", nil
	}

	return filepath, nil
}
