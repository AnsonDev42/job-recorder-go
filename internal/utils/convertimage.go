package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gookit/config/v2"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func Img2word(imgPath *string, ocrPath *string) (string, error) {
	absImgPath, err := filepath.Abs(*imgPath)
	if err != nil {
		return "", err
	}
	absOutputPath, err := filepath.Abs(*ocrPath)

	preCmd := fmt.Sprintf("shortcuts run \"extract_text_from_image\" -i %s -o %s", absImgPath, absOutputPath)
	cmd := exec.Command("zsh", "-c", preCmd)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// read tmp.txt as result
	results := ""
	file, err := os.Open(absOutputPath)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		results += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return "", err
	}

	return results, nil
}

func SummarizeText(ocrText string) (Job, error) {
	// call openai to summarize the text and return the result as Job struct
	if config.String("openaiKey", "") == "" {
		return Job{}, fmt.Errorf("openai key is empty")
	}
	client := openai.NewClient(config.String("openaiKey"))
	const prompt = "You are a summarizer, your purpose is to summarize the following text, which should be an OCR results " +
		"from a job description, and you should summarize the company name, job title , a brief job description of the job." +
		" Return in json format of these three required fields and only these three, without any comments since I" +
		" process the raw json return. Be aware any company can have developer jobs, so job-title and the company-name might not be the same category and that it okay. If no relevant information in the given texts, still return the key but with value of empty string (\"\")." +
		" Example return: {\"company-name\": \"Microsoft\", \"job-title\":\"Graduate software engineer\",\"job-description\": \"Backend engineer, skill set: python, sql, and golang, estimation salary 35k to 45k GBP\". The following is the text to be summarized: "
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt + ocrText,
				},
			},
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return Job{}, err
	}
	log.Println(resp.Choices[0].Message.Content)
	job, err := processSummarizerResponse(resp.Choices[0].Message.Content)
	if err != nil {
		return Job{}, err

	}
	return job, nil
}

//{
//"company-Name": "American Express",
//"job-title": "Software Engineer",
//"job-description": "Responsible for designing, developing, and maintaining software applications. Required skills include Java, Spring Boot, and Angular. Competitive salary based on experience."
//}

type Job struct {
	CompanyName    string `json:"company-name"`
	JobTitle       string `json:"job-title"`
	JobDescription string `json:"job-description"`
}

func processSummarizerResponse(resp string) (Job, error) {
	// read summary as json and save to struct

	if resp == "" {
		return Job{}, fmt.Errorf("empty response")
	}
	var job Job
	err := json.Unmarshal([]byte(resp), &job)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return Job{}, err
	}

	fmt.Printf("Company Name: %s\n", job.CompanyName)
	fmt.Printf("Job Title: %s\n", job.JobTitle)
	fmt.Printf("Job Description: %s\n", job.JobDescription)
	return job, nil
}
