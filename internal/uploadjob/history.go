package uploadjob

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"time"
)

func ShowHistoryUI(content *fyne.Container, rootFolder string) {
	files, err := os.ReadDir(rootFolder)
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

func CreateHistoryTable(uploadDir string) *widget.Table {
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
