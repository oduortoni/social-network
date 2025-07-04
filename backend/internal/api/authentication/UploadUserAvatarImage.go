package authentication

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func DownloadAndSavePicture(profileImage string) (string, error) {
	if profileImage == "" || strings.Contains(profileImage, "default-user") {
		return "No profile Image", nil
	}

	if _, err := os.Stat("./UserAvatars"); os.IsNotExist(err) {
		if err := os.MkdirAll("./UserAvatars", 0755); err != nil {
			return "", err
		}
	}

	// Send GET request to image URL
	resp, err := http.Get(profileImage)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file locally
	savePath := filepath.Join("./UserAvatars", uuid.New().String())
	outFile, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Copy the image bytes to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {

		return "", err
	}
	return savePath, nil
}
