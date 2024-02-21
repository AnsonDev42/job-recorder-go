package test

import (
	"fmt"
	"github.com/gookit/config/v2"
	"job-recorder-go/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImg2word(t *testing.T) {
	uploadsPath := filepath.Join("..", "assets", "uploads", "test_image.png")
	ocrPath := filepath.Join("..", "assets", "uploads", "test_image.txt")
	word, err := utils.Img2word(&uploadsPath, &ocrPath)
	if err != nil {
		t.Errorf("error when calling img2word: %s", err)
	}
	fmt.Println(word)
	expected := "dialog. ShowInformation ( title: \"Success\", message: \"Image from clipboard uploaded successfully"
	if expected != word {
		t.Fatalf("ocr results not match! expect %s, got %s", expected, word)
	}
}
func TestImg2wordAMEXDATA(t *testing.T) {
	uploadsPath := filepath.Join("..", "test", "testdata", "AMEX.png")
	ocrPath := filepath.Join("..", "test", "testdata", "AMEX.txt")
	word, err := utils.Img2word(&uploadsPath, &ocrPath)
	if err != nil {
		t.Errorf("error when calling img2word: %s", err)
	}
	fmt.Println(word)
	expectedContainedText := "American Express"
	if !strings.Contains(word, expectedContainedText) {
		t.Fatalf("ocr results not contain %s, got %s", expectedContainedText, word)
	}
}

// struct {

func summarizeText(t *testing.T) utils.Job {
	err := config.LoadFiles("../config/dev-config.json")
	if err != nil {
		t.Errorf("failed to load json")
		t.FailNow()

	}
	key := config.String("openaiKey", "")
	if key == "" {
		t.Errorf("openai key is empty")
		t.FailNow()
	}

	ocrText, _ := os.ReadFile("../test/testdata/AMEX.txt")
	summary, err := utils.SummarizeText(string(ocrText))
	if err != nil {
		t.Errorf("error when calling summarizeText: %s", err)
	}
	return summary
}

func TestProcessSummarizerResponse(t *testing.T) {
	summary := summarizeText(t)

	expected := utils.Job{
		CompanyName:    "American Express",
		JobTitle:       "Software Engineer",
		JobDescription: "Responsible for designing, developing, and maintaining software applications. Required skills include Java, Spring Boot, and Angular. Competitive salary based on experience.",
	}
	//check if the job struct is the same as expected
	if !strings.Contains(summary.JobTitle, "Software") { // might be software developer or engineer
		t.Fatalf("Job title not match! expect %s, got %s", expected.JobTitle, summary.JobTitle)
	}
	if expected.CompanyName != summary.CompanyName {
		t.Fatalf("Company name not match! expect %s, got %s", expected.CompanyName, summary.CompanyName)
	}

}
