package uploadjob

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowSettingsUI(window fyne.Window, content *fyne.Container, rootFolder *string) {
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

	// set the daily goal
	dailyGoal := binding.NewInt()
	dailyGoal.Set(10)
	//dailyGoal :=10
	dW := widget.NewLabelWithData(binding.IntToString(dailyGoal))
	dShowDailyGoal := widget.NewLabel("Current Setting for daily goal is: ")
	dailGoalRow := container.NewHBox(dShowDailyGoal, dW)
	content.Objects = []fyne.CanvasObject{
		container.NewVBox(dailGoalRow, selectFolderButton),
	}
	content.Refresh()
}
