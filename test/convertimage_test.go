package test

import (
	"fmt"
	"job-recorder-go/internal/convertimage"
	"path/filepath"
	"testing"
)

func TestImg2word(t *testing.T) {
	uploadsPath := filepath.Join("..", "assets", "uploads", "test_image.png")
	ocrPath := filepath.Join("..", "assets", "uploads", "test_image.txt")
	word, err := convertimage.Img2word(&uploadsPath, &ocrPath)
	if err != nil {
		t.Errorf("error when calling img2word: %s", err)
	}
	fmt.Println(word)
	expected := "dialog. ShowInformation ( title: \"Success\", message: \"Image from clipboard uploaded successfully"
	if expected != word {
		t.Fatalf("ocr results not match! expect %s, got %s", expected, word)
	}
}
