package uploadjob

import (
	"context"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gookit/config/v2"
	"github.com/nikoksr/notify"
	"golang.design/x/clipboard"
	"io"
	"job-recorder-go/internal/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func uploadJobFromByte(imgData []byte, uploadDir string, updateCounterCh chan int) (string, string, error) {
	// Assume imgData is PNG encoded. Save it to the upload directory.
	uploadFileTime := time.Now().Format("2006-01-02-15-04-05.000")
	uploadFileName := uploadFileTime + ".png"
	ocrFileName := uploadFileTime + ".txt"
	filePath := filepath.Join(uploadDir, uploadFileName) // Consider generating unique names
	ocrPath := filepath.Join(uploadDir, ocrFileName)
	err := os.WriteFile(filePath, imgData, 0644)
	if err != nil {
		return "", uploadFileTime, err
	}
	updateCounterCh <- 1
	word, err := utils.Img2word(&filePath, &ocrPath)
	if err != nil {
		return "", uploadFileTime, err
	}
	return word, uploadFileTime, nil

}
func saveJobSummary(job utils.Job, uploadDir string, uploadFileTime string) error {
	// Save the job summary to a file
	summaryFileName := uploadFileTime + ".json"
	// check if uploadDir/summary exists, if not create it
	summaryDir := filepath.Join(uploadDir, "summary")
	if _, err := os.Stat(summaryDir); os.IsNotExist(err) {
		err := os.Mkdir(summaryDir, 0755)
		if err != nil {
			return err
		}
	}
	summaryPath := filepath.Join(uploadDir, "summary", summaryFileName)
	jsonBytes, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return err
	}
	err = os.WriteFile(summaryPath, jsonBytes, 0644)
	return err
}
func ShowUploadUI(window fyne.Window, content *fyne.Container, uploadDir *string, updateCounterCh chan int) {
	uploadFileButton := widget.NewButton("Upload File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return // Handle error or cancellation
			}
			defer func(reader fyne.URIReadCloser) {
				err := reader.Close()
				if err != nil {

				}
			}(reader)
			imgData, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			ocrResult, ftime, err := uploadJobFromByte(imgData, *uploadDir, updateCounterCh)
			if err != nil {
				return
			}
			job, err := utils.SummarizeText(ocrResult)
			go func() {
				err := saveJobSummary(job, *uploadDir, ftime)
				if err != nil {
					dialog.ShowError(err, window)
				}
			}()
			dialog.ShowInformation("Summary result:", fmt.Sprint(job), window)
		}, window)
	})

	uploadClipboardButton := widget.NewButton("Upload from Clipboard", func() {
		// Implement the clipboard reading and image saving logic here
		imgData := clipboard.Read(clipboard.FmtImage)
		// Assume imgData is PNG encoded. Save it to the upload directory.
		ocrResult, ftime, err := uploadJobFromByte(imgData, *uploadDir, updateCounterCh)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		job, err := utils.SummarizeText(ocrResult)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		go func() {
			err := saveJobSummary(job, *uploadDir, ftime)
			if err != nil {
				dialog.ShowError(err, window)
			}
		}()
		dialog.ShowInformation("Summary result:", fmt.Sprint(job), window)

		//fmt.Println(word)
		return
	})

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(uploadFileButton, uploadClipboardButton),
	}
	content.Refresh()
}

func CounterUpdator(updateCounterCh chan int, counterLabel *widget.Label) {
	for range updateCounterCh {
		count, err := UpdateCounterLabel(counterLabel)
		if err != nil {
			fmt.Println("failed to update counter")
		}
		err = notify.Send(
			context.Background(),
			"Another hardworking day!",
			fmt.Sprint("You have applied", count, "jobs, one step closer to your daily goal of ", config.Int("dailyGoal")),
		)
		if err != nil {
			fmt.Errorf("failed to send the notification")
		}
	}
}
func CountTodayJobs() (int, error) {
	uploadsDir := config.String("rootFolder")
	files, err := os.ReadDir(uploadsDir)
	if err != nil {
		log.Println("Failed to read upload directory:", err)
		return 0, err
	}

	today := time.Now().Format("2006-01-02")
	count := 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), today) && strings.HasSuffix(file.Name(), ".png") {
			count++
		}
	}
	return count, nil
}
func UpdateCounterLabel(label *widget.Label) (count int, err error) {
	count, err = CountTodayJobs()
	if err != nil {
		return 0, err
	}
	log.Println("settings counter!")
	f := fmt.Sprint("You have applied ", count, " jobs, one step closer to your daily goal of ", config.Int("dailyGoal"))
	label.SetText(f)

	return count, nil
}

func SummarizeTodaysWork() (string, error) {
	// Get today's date as a string prefix
	today := time.Now().Format("2006-01-02")

	// Directory where the files are stored
	uploadsDir := config.String("rootFolder")

	// Open the directory
	files, err := os.ReadDir(uploadsDir)
	if err != nil {
		log.Fatal("Error reading directory:", err)
		return "", err
	}
	count, err := CountTodayJobs() // create a summary string: currently just concatenating all the OCR
	if err != nil {
		return "summarizer error", err
	}
	var summary string
	summary += fmt.Sprintf("Today's work summary: (%s,%s) \n", count, config.Int("dailyGoal", -1))
	summary += "---------------------\n"
	// Iterate over the files in the directory
	for _, file := range files {
		// Check if the file name starts with today's date and has a .txt extension
		if strings.HasPrefix(file.Name(), today) && strings.HasSuffix(file.Name(), ".txt") {
			// Read the file
			content, err := os.ReadFile(fmt.Sprintf("%s/%s", uploadsDir, file.Name()))
			if err != nil {
				fmt.Println("Error reading file:", file.Name(), err)
				continue // Skip to the next file upon error
			}
			// Concatenate the content to the summary
			summary += string(content) + "\n" // Adding a newline for separation
		}
	}

	return summary, nil
}

func SendSummary() {
	summary, err := SummarizeTodaysWork()
	err = notify.Send(
		context.Background(),
		"Today's work summary",
		summary,
	)
	if err != nil {
		fmt.Errorf("failed to send the notification")
	}
}
