package service

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

// Test data with actual image signatures
var testImageData = map[ImageFormat][]byte{
	JPEG: {0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46},
	PNG:  {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D},
	GIF:  {0x47, 0x49, 0x46, 0x38, 0x37, 0x61, 0x10, 0x00, 0x10, 0x00},             // GIF87a
	WebP: {0x52, 0x49, 0x46, 0x46, 0x24, 0x08, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}, // RIFF + size + WEBP
	BMP:  {0x42, 0x4D, 0x46, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	TIFF: {0x49, 0x49, 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00}, // Little endian TIFF
}

var gif89aData = []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x10, 0x00, 0x10, 0x00} // GIF89a
var tiffBigEndianData = []byte{0x4D, 0x4D, 0x00, 0x2A, 0x00, 0x08, 0x00, 0x00}      // Big endian TIFF

func TestDetectImageFormat(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expectedFormat ImageFormat
		expectError    bool
	}{
		{
			name:           "Valid JPEG",
			data:           testImageData[JPEG],
			expectedFormat: JPEG,
			expectError:    false,
		},
		{
			name:           "Valid PNG",
			data:           testImageData[PNG],
			expectedFormat: PNG,
			expectError:    false,
		},
		{
			name:           "Valid GIF87a",
			data:           testImageData[GIF],
			expectedFormat: GIF,
			expectError:    false,
		},
		{
			name:           "Valid GIF89a",
			data:           gif89aData,
			expectedFormat: GIF,
			expectError:    false,
		},
		{
			name:           "Valid WebP",
			data:           testImageData[WebP],
			expectedFormat: WebP,
			expectError:    false,
		},
		{
			name:           "Valid BMP",
			data:           testImageData[BMP],
			expectedFormat: BMP,
			expectError:    false,
		},
		{
			name:           "Valid TIFF (little endian)",
			data:           testImageData[TIFF],
			expectedFormat: TIFF,
			expectError:    false,
		},
		{
			name:           "Valid TIFF (big endian)",
			data:           tiffBigEndianData,
			expectedFormat: TIFF,
			expectError:    false,
		},
		{
			name:           "Invalid format",
			data:           []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			expectedFormat: "",
			expectError:    true,
		},
		{
			name:           "Empty data",
			data:           []byte{},
			expectedFormat: "",
			expectError:    true,
		},
		{
			name:           "Too small data",
			data:           []byte{0xFF},
			expectedFormat: "",
			expectError:    true,
		},
		{
			name:           "Incomplete WebP (RIFF only)",
			data:           []byte{0x52, 0x49, 0x46, 0x46, 0x24, 0x08, 0x00, 0x00},
			expectedFormat: "",
			expectError:    true,
		},
		{
			name:           "Fake WebP (RIFF + wrong identifier)",
			data:           []byte{0x52, 0x49, 0x46, 0x46, 0x24, 0x08, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45},
			expectedFormat: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.data)
			format, err := DetectImageFormat(reader)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if format != tt.expectedFormat {
				t.Errorf("Expected format %s, got %s", tt.expectedFormat, format)
			}
		})
	}
}

func TestCheckImageSignature(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
	}{
		{
			name:        "Valid JPEG",
			data:        testImageData[JPEG],
			expectError: false,
		},
		{
			name:        "Valid PNG",
			data:        testImageData[PNG],
			expectError: false,
		},
		{
			name:        "Invalid format",
			data:        []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			expectError: true,
		},
		{
			name:        "Empty data",
			data:        []byte{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.data)
			err := CheckImageSignature(reader)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestIsImageFormat(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expectedFormat ImageFormat
		expectedResult bool
		expectError    bool
	}{
		{
			name:           "JPEG matches JPEG",
			data:           testImageData[JPEG],
			expectedFormat: JPEG,
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "PNG doesn't match JPEG",
			data:           testImageData[PNG],
			expectedFormat: JPEG,
			expectedResult: false,
			expectError:    false,
		},
		{
			name:           "Invalid data",
			data:           []byte{0x00, 0x01, 0x02, 0x03},
			expectedFormat: JPEG,
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.data)
			result, err := IsImageFormat(reader, tt.expectedFormat)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestGetSupportedFormats(t *testing.T) {
	formats := GetSupportedFormats()

	// Check that all expected formats are present
	expectedFormats := []ImageFormat{JPEG, PNG, GIF, WebP, BMP, TIFF}
	formatMap := make(map[ImageFormat]bool)
	for _, format := range formats {
		formatMap[format] = true
	}

	for _, expected := range expectedFormats {
		if !formatMap[expected] {
			t.Errorf("Expected format %s not found in supported formats", expected)
		}
	}

	// Check for duplicates
	if len(formats) != len(formatMap) {
		t.Errorf("Duplicate formats found in supported formats list")
	}

	// Verify minimum expected count
	if len(formats) < 6 {
		t.Errorf("Expected at least 6 supported formats, got %d", len(formats))
	}
}

func TestCalculateMaxBufferSize(t *testing.T) {
	maxSize := calculateMaxBufferSize()

	// WebP requires 12 bytes (RIFF + size + WEBP)
	if maxSize < 12 {
		t.Errorf("Expected max buffer size to be at least 12 bytes, got %d", maxSize)
	}
}

// Test with various reader types
func TestDetectImageFormatWithDifferentReaders(t *testing.T) {
	testData := testImageData[JPEG]

	tests := []struct {
		name   string
		reader io.Reader
	}{
		{
			name:   "bytes.Reader",
			reader: bytes.NewReader(testData),
		},
		{
			name:   "bytes.Buffer",
			reader: bytes.NewBuffer(testData),
		},
		{
			name:   "strings.Reader",
			reader: strings.NewReader(string(testData)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, err := DetectImageFormat(tt.reader)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if format != JPEG {
				t.Errorf("Expected JPEG, got %s", format)
			}
		})
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("Minimum viable JPEG", func(t *testing.T) {
		minJPEG := []byte{0xFF, 0xD8, 0xFF}
		reader := bytes.NewReader(minJPEG)
		format, err := DetectImageFormat(reader)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if format != JPEG {
			t.Errorf("Expected JPEG, got %s", format)
		}
	})

	t.Run("Minimum viable BMP", func(t *testing.T) {
		minBMP := []byte{0x42, 0x4D}
		reader := bytes.NewReader(minBMP)
		format, err := DetectImageFormat(reader)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if format != BMP {
			t.Errorf("Expected BMP, got %s", format)
		}
	})

	t.Run("Reader error handling", func(t *testing.T) {
		// Create a reader that always returns an error
		errorReader := &errorReader{err: errors.New("test error")}
		_, err := DetectImageFormat(errorReader)
		if err == nil {
			t.Error("Expected error from errorReader, but got none")
		}
	})
}

// Mock reader that always returns an error
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

// Benchmark tests
func BenchmarkDetectImageFormat(b *testing.B) {
	testData := testImageData[JPEG]
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(testData)
		_, err := DetectImageFormat(reader)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGetSupportedFormats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetSupportedFormats()
	}
}
