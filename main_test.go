package main

import (
	"fmt"
	"job-recorder-go/convertimage"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConvertimage(t *testing.T) {
	uploadDir := "uploads"

	// Implement the clipboard reading and image saving logic here
	imgData, _ := os.ReadFile("uploads/test_image.png")
	// Assume imgData is PNG encoded. Save it to the upload directory.
	uploadFileTime := time.Now().Format("2006-01-02-15-04-05.000")
	uploadFileName := uploadFileTime + ".png"
	ocrFileName := uploadFileTime + ".txt"
	filePath := filepath.Join(uploadDir, uploadFileName) // Consider generating unique names
	ocrPath := filepath.Join(uploadDir, ocrFileName)
	err := os.WriteFile(filePath, imgData, 0644)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	//dialog.ShowInformation("Success", "Image from clipboard uploaded successfully.", window)
	time.Sleep(5 * time.Second)
	word, err := convertimage.Img2word(&filePath, &ocrPath)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	//dialog.ShowInformation("OCR Results", word, window)
	fmt.Println(word)
}
