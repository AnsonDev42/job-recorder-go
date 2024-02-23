package uploadjob

import (
	"fyne.io/fyne/v2/widget"
)

// add a progress bar based on the number of jobs applied and the goal

func CreateProgressBar() *widget.ProgressBar {
	// Create a progress bar
	progress := widget.NewProgressBar()
	progress.SetValue(0)
	return progress
}
