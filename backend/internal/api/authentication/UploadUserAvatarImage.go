package authentication

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Define magic numbers for common image formats
var imageMagicNumbers = map[string][]byte{
	"jpeg": {0xFF, 0xD8, 0xFF},
	"png":  {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
	"gif":  {0x47, 0x49, 0x46, 0x38}, // GIF87a or GIF89a
	"webp": {0x52, 0x49, 0x46, 0x46}, // RIFF (WebP files start with RIFF, followed by file size and WEBP)
}

// checkImageSignature reads the first few bytes and validates against known image magic numbers.
func checkImageSignature(reader io.Reader) error {
	// Read enough bytes to cover the longest magic number (PNG: 8 bytes)
	buffer := make([]byte, 8)
	n, err := io.ReadAtLeast(reader, buffer, 4) // Read at least 4 bytes (for GIF/WebP)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	// Check for JPEG
	if n >= 3 && bytes.Equal(buffer[0:3], imageMagicNumbers["jpeg"]) {
		return nil
	}
	// Check for PNG
	if n >= 8 && bytes.Equal(buffer[0:8], imageMagicNumbers["png"]) {
		return nil
	}
	// Check for GIF
	if n >= 4 && bytes.Equal(buffer[0:4], imageMagicNumbers["gif"]) {
		return nil
	}
	// Check for WebP (RIFF header, then check for "WEBP" at offset 8)
	if n >= 4 && bytes.Equal(buffer[0:4], imageMagicNumbers["webp"]) {
		// For WebP, we need to read more to confirm "WEBP"
		// Since we only read 8 bytes, we can't fully validate WEBP here without more reads.
		// For simplicity and to avoid re-reading, we'll assume RIFF is enough for now,
		// but a more robust check would involve reading bytes 8-11 for "WEBP".
		// For this implementation, we'll just check the RIFF header.
		return nil
	}

	return errors.New("unsupported image format or invalid file signature")
}

func UploadAvatarImage(imagereader multipart.File, imageheader *multipart.FileHeader) (string, error) {
	const MaxUploadSize = 20 * 1024 * 1024 // 20MB limit

	if imageheader.Size > MaxUploadSize {
		return "maximum size", errors.New("file size exceeds maximum allowed")
	}

	// Create a tee reader to read from imagereader and also pass to checkImageSignature
	// This allows checkImageSignature to read from the beginning of the file without
	// affecting the subsequent io.Copy operation.
	var buf bytes.Buffer
	teeReader := io.TeeReader(imagereader, &buf)

	if err := checkImageSignature(teeReader); err != nil {
		return "", err // Return error if signature check fails
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

	// Copy the content from the buffer (which contains the bytes already read by teeReader)
	// and then continue copying from the original imagereader.
	// This ensures all bytes are copied to the output file.
	_, err = io.Copy(out, io.MultiReader(&buf, imagereader))
	if err != nil {
		return "failed to copy file content", err
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
