package utils

import (
	"bytes"
	"errors"
	"io"
)

type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	PNG  ImageFormat = "png"
	GIF  ImageFormat = "gif"
	WebP ImageFormat = "webp"
	BMP  ImageFormat = "bmp"
	TIFF ImageFormat = "tiff"
)

type ImageSignature struct {
	Magic     []byte
	MinLength int
	Format    ImageFormat
}

// Define image signatures with proper magic numbers
var imageSignatures = []ImageSignature{
	// JPEG - FF D8 FF
	{Magic: []byte{0xFF, 0xD8, 0xFF}, MinLength: 3, Format: JPEG},

	// PNG - 89 50 4E 47 0D 0A 1A 0A
	{Magic: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, MinLength: 8, Format: PNG},

	// GIF87a - 47 49 46 38 37 61
	{Magic: []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}, MinLength: 6, Format: GIF},

	// GIF89a - 47 49 46 38 39 61
	{Magic: []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}, MinLength: 6, Format: GIF},

	// BMP - 42 4D
	{Magic: []byte{0x42, 0x4D}, MinLength: 2, Format: BMP},

	// TIFF (little endian) - 49 49 2A 00
	{Magic: []byte{0x49, 0x49, 0x2A, 0x00}, MinLength: 4, Format: TIFF},

	// TIFF (big endian) - 4D 4D 00 2A
	{Magic: []byte{0x4D, 0x4D, 0x00, 0x2A}, MinLength: 4, Format: TIFF},
}

var webpSignature = struct {
	RiffMagic []byte
	WebpMagic []byte
	MinLength int
}{
	RiffMagic: []byte{0x52, 0x49, 0x46, 0x46}, // "RIFF"
	WebpMagic: []byte{0x57, 0x45, 0x42, 0x50}, // "WEBP"
	MinLength: 12,
}

func calculateMaxBufferSize() int {
	maxSize := webpSignature.MinLength
	for _, sig := range imageSignatures {
		if sig.MinLength > maxSize {
			maxSize = sig.MinLength
		}
	}
	return maxSize
}

func DetectImageFormat(reader io.Reader) (ImageFormat, error) {
	bufferSize := calculateMaxBufferSize()
	buffer := make([]byte, bufferSize)
	n, err := io.ReadAtLeast(reader, buffer, 2)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return "", errors.New("file too small to determine format")
		}
		return "", err
	}

	// Check webP first
	if n >= webpSignature.MinLength {
		if bytes.Equal(buffer[0:4], webpSignature.RiffMagic) && bytes.Equal(buffer[8:12], webpSignature.WebpMagic) {
			return WebP, nil
		}
	}

	// Check other formats
	for _, sig := range imageSignatures {
		if n >= sig.MinLength && bytes.Equal(buffer[0:sig.MinLength], sig.Magic) {
			return sig.Format, nil
		}
	}
	return "", errors.New("unsupported image format or invalid file signature")
}

// checkImageSignature reads the first few bytes and validates against known image magic numbers.
// func checkImageSignature(reader io.Reader) error {
// 	_, err := DetectImageFormat(reader)
// 	return err
// }

// CheckImageSignature validates if the file is a supported image format (backward compatibility)
func CheckImageSignature(reader io.Reader) error {
	_, err := DetectImageFormat(reader)
	return err
}

// IsImageFormat checks if the file matches a specific image format
func IsImageFormat(reader io.Reader, expectedFormat ImageFormat) (bool, error) {
	detectedFormat, err := DetectImageFormat(reader)
	if err != nil {
		return false, err
	}
	return detectedFormat == expectedFormat, nil
}

// GetSupportedFormats returns a list of all supported image formats
func GetSupportedFormats() []ImageFormat {
	formats := make([]ImageFormat, 0, len(imageSignatures)-1)
	formatSet := make(map[ImageFormat]bool)

	// Add WebP
	formats = append(formats, WebP)
	formatSet[WebP] = true

	// Add other formats (deduplicate)
	for _, sig := range imageSignatures {
		if !formatSet[sig.Format] {
			formats = append(formats, sig.Format)
			formatSet[sig.Format] = true
		}
	}
	return formats
}
