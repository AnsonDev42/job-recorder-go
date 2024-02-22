package test

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConvertimage(t *testing.T) {
	uploadDir := "testdata"

	// Implement the clipboard reading and image saving logic here
	imgData, _ := os.ReadFile("AMEX.png")
	// Assume imgData is PNG encoded. Save it to the upload directory.
	uploadFileTime := time.Now().Format("2006-01-02-15-04-05.000")
	uploadFileName := uploadFileTime + ".png"
	filePath := filepath.Join(uploadDir, uploadFileName) // Consider generating unique names
	err := os.WriteFile(filePath, imgData, 0644)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
}
