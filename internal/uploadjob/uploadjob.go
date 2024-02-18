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

			// copy the file into the uploadDir
			uploadFileName := filepath.Base(reader.URI().Path())
			uploadFilePath := filepath.Join(*uploadDir, uploadFileName)
			imgData, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			err = os.WriteFile(uploadFilePath, imgData, 0644)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			dialog.ShowInformation("Success", "File uploaded successfully.", window)
		}, window)
		updateCounterCh <- 1
	})

	uploadClipboardButton := widget.NewButton("Upload from Clipboard", func() {
		// Implement the clipboard reading and image saving logic here
		imgData := clipboard.Read(clipboard.FmtImage)
		// Assume imgData is PNG encoded. Save it to the upload directory.
		uploadFileTime := time.Now().Format("2006-01-02-15-04-05.000")
		uploadFileName := uploadFileTime + ".png"
		ocrFileName := uploadFileTime + ".txt"
		filePath := filepath.Join(*uploadDir, uploadFileName) // Consider generating unique names
		ocrPath := filepath.Join(*uploadDir, ocrFileName)
		err := os.WriteFile(filePath, imgData, 0644)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		updateCounterCh <- 1
		//dialog.ShowInformation("Success", "Image from clipboard uploaded successfully.", window)
		time.Sleep(2 * time.Second)

		word, err := utils.Img2word(&filePath, &ocrPath)
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

func CounterUpdator(updateCounterCh chan int, counterLabel *widget.Label, uploadDir string) {
	for range updateCounterCh {
		count, err := UpdateCounterLabel(counterLabel, uploadDir)
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

func UpdateCounterLabel(label *widget.Label, uploadDir string) (count int, err error) {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Println("Failed to read upload directory:", err)
		label.SetText("Images uploaded today: Error")
		return 0, err
	}

	today := time.Now().Format("2006-01-02")
	count = 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), today) && strings.HasSuffix(file.Name(), ".png") {
			count++
		}
	}
	fmt.Println("settings counter!")
	f := fmt.Sprint("You have applied ", count, " jobs, one step closer to your daily goal of ", config.Int("dailyGoal"))
	label.SetText(f)

	return count, nil
}
