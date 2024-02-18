package uploadjob

import (
	"fyne.io/fyne/v2"
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

	content.Objects = []fyne.CanvasObject{selectFolderButton}
	content.Refresh()
}
