package uploadjob

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"strconv"
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
	todayJobs, err := GetTodaysJobFromFile()
	if err != nil {
		return nil
	}
	// Define your table data
	fileInfos := make([][4]string, len(todayJobs))
	for i, job := range todayJobs {
		fileInfos[i] = [4]string{job.JobTitle, job.CompanyName, job.JobDescription}
	}

	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(fileInfos), 3 // Rows, Columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // This will create a new cell
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			// Set the cell value. id.Row and id.Col will tell you which cell you're populating
			cell.(*widget.Label).SetText(fileInfos[id.Row][id.Col])
		},
	)
	// set custom header: show column names and row numbers
	table.CreateHeader = func() fyne.CanvasObject {
		return container.NewHBox(
			widget.NewLabel(""),
			widget.NewLabel(""),
			widget.NewLabel(""),
		)
	}
	headerLabels := []string{"Position", "Company", "Summary", "Fourth Column"} // Labels for headers
	table.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		if id.Col >= 0 && id.Col < len(headerLabels) {
			header := template.(*fyne.Container).Objects[id.Col].(*widget.Label) // Access the specific header label by index
			header.SetText(headerLabels[id.Col])                                 // Set the text for the header based on the column
		}
		if id.Row >= 0 {
			header := template.(*fyne.Container).Objects[0].(*widget.Label) // Access the specific header label by index
			header.SetText(strconv.Itoa(id.Row))                            // Set the text for the header based on the column
		}
	}

	// Set a minimum width for the table to ensure content is less likely to overlap
	table.SetColumnWidth(0, 100)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 300)

	// Customize each column's width (not directly supported, but you can indirectly influence it)
	// For example, you can format your data to ensure it fits within your designated widths
	// This step is more about preparing your data (e.g., truncating file names, adjusting date formats) to fit

	return table
}
