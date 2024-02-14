package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
	"io"
	"io/ioutil"
	"job-recorder-go/convertimage"
	"log"
	"os"
	"path/filepath"
	"time"
)

func showUploadUI(window fyne.Window, content *fyne.Container, uploadDir *string) {
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
		//dialog.ShowInformation("Success", "Image from clipboard uploaded successfully.", window)
		time.Sleep(2 * time.Second)

		word, err := convertimage.Img2word(&filePath, &ocrPath)
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

func showHistoryUI(content *fyne.Container, rootFolder string) {
	files, err := ioutil.ReadDir(rootFolder)
	if err != nil {
		content.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load history")}
		content.Refresh()
		return
	}

	var fileObjects []fyne.CanvasObject
	for _, file := range files {
		fileObjects = append(fileObjects, widget.NewLabel(file.Name()))
	}

	fileList := container.NewVScroll(container.NewVBox(fileObjects...))
	content.Objects = []fyne.CanvasObject{fileList}
	content.Refresh()
}
func showSettingsUI(window fyne.Window, content *fyne.Container, rootFolder *string) {
	selectFolderButton := widget.NewButton("Select Folder", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			*rootFolder = uri.Path()
			content.Objects = []fyne.CanvasObject{widget.NewLabel("Save path: " + *rootFolder)}
			content.Refresh()
		}, window)
	})

	content.Objects = []fyne.CanvasObject{selectFolderButton}
	content.Refresh()
}

func createHistoryTable(uploadDir string) *widget.Table {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Println("Failed to read upload directory:", err)
		return nil
	}

	// Define your table data
	fileInfos := make([][2]string, len(files))
	for i, fileInfo := range files {
		file, err := fileInfo.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}
		fileInfos[i] = [2]string{file.Name(), file.ModTime().Format(time.RFC1123)}
	}

	table := widget.NewTable(
		func() (int, int) {
			return len(fileInfos), 2 // Rows, Columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // This will create a new cell
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			// Set the cell value. id.Row and id.Col will tell you which cell you're populating
			cell.(*widget.Label).SetText(fileInfos[id.Row][id.Col])
		},
	)

	// Set a minimum width for the table to ensure content is less likely to overlap
	table.SetColumnWidth(0, 300)
	table.SetColumnWidth(1, 300)

	// Customize each column's width (not directly supported, but you can indirectly influence it)
	// For example, you can format your data to ensure it fits within your designated widths
	// This step is more about preparing your data (e.g., truncating file names, adjusting date formats) to fit

	return table
}

func countDaily(today *string) {

}
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Job recorder go!")

	uploadDir := "uploads"
	content := container.NewStack()

	// Setup the counter
	//counterLabel := widget.NewLabel("Images uploaded today: 0")
	//updateCounter(counterLabel, uploadDir) // Initial count update

	// Set up the menu and content area
	menu := setupMenu(myWindow, content, &uploadDir)

	// Create the vertical split container
	split := container.NewHSplit(menu, content)
	split.Offset = 0.2 // Adjust the initial split ratio

	myWindow.SetContent(split)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func setupMenu(window fyne.Window, content *fyne.Container, uploadDir *string) fyne.CanvasObject {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	uploadButton := widget.NewButton("Upload", func() {
		showUploadUI(window, content, uploadDir)
	})

	historyButton := widget.NewButton("History", func() {
		table := createHistoryTable(*uploadDir)
		if table != nil {
			content.Objects = []fyne.CanvasObject{table}
			content.Refresh()
		}
	})

	settingsButton := widget.NewButton("Settings", func() {
		// Implement settings view
	})

	return container.NewVBox(uploadButton, historyButton, settingsButton)
}
