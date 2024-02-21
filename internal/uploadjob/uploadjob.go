package uploadjob

import (
	"context"
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

func uploadJobFromByte(imgData []byte, uploadDir string, updateCounterCh chan int) (string, error) {
	// Assume imgData is PNG encoded. Save it to the upload directory.
	uploadFileTime := time.Now().Format("2006-01-02-15-04-05.000")
	uploadFileName := uploadFileTime + ".png"
	ocrFileName := uploadFileTime + ".txt"
	filePath := filepath.Join(uploadDir, uploadFileName) // Consider generating unique names
	ocrPath := filepath.Join(uploadDir, ocrFileName)
	err := os.WriteFile(filePath, imgData, 0644)
	if err != nil {
		return "", err
	}
	updateCounterCh <- 1
	word, err := utils.Img2word(&filePath, &ocrPath)
	if err != nil {
		return "", err
	}
	return word, nil

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
			_, err = uploadJobFromByte(imgData, *uploadDir, updateCounterCh)
			if err != nil {
				return
			}
		}, window)
	})

	uploadClipboardButton := widget.NewButton("Upload from Clipboard", func() {
		// Implement the clipboard reading and image saving logic here
		imgData := clipboard.Read(clipboard.FmtImage)
		// Assume imgData is PNG encoded. Save it to the upload directory.
		word, err := uploadJobFromByte(imgData, *uploadDir, updateCounterCh)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("OCR Results", word, window)
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
