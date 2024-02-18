package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gookit/config/v2"
	"golang.design/x/clipboard"
	"job-recorder-go/internal/uploadjob"
	"job-recorder-go/internal/utils"
)

func countDaily(today *string) {

}
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Job recorder go!")

	uploadDir := "assets/uploads"
	content := container.NewStack()

	// Setup the counter
	counterLabel := widget.NewLabel("Images uploaded today: 0")
	counterLabel.TextStyle = fyne.TextStyle{Bold: true}
	counterLabel.Resize(fyne.NewSize(200, 100))
	uploadjob.UpdateCounterLabel(counterLabel, uploadDir) // Initial count update
	// Load App config
	configPATH := "config/dev-config.json"
	err := config.LoadFiles(configPATH)
	if err != nil {
		_ = fmt.Errorf("failed to load json")
	}
	config.Set("configPATH", configPATH)

	// Setup telegram notification
	err = utils.SetupTelegramBot(config.String("tgApi"), config.String("tgReceiverID"))
	if err != nil {
		panic("error setting up telegram notification!")
	}
	// Setup Counter updator
	updateCounterCh := make(chan int)
	go uploadjob.CounterUpdator(updateCounterCh, counterLabel, uploadDir)

	// Set up the menu and content area
	menu := setupMenu(myWindow, content, &uploadDir, updateCounterCh)
	menuContentSplit := container.NewHSplit(menu, content)
	menuContentSplit.Offset = 0.2 // Adjust the initial split ratio
	mainContent := container.NewVSplit(container.New(layout.NewCenterLayout(), counterLabel), menuContentSplit)
	mainContent.Offset = 0.5

	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("Job-Recorder",
			fyne.NewMenuItem("Show", func() {
				myWindow.Show()
			}))
		fyne.NewMenuItem("Quit", func() {
			myApp.Quit()
		})
		fyne.NewMenuItem("Copy from clipboard", func() {
			//todo: atomic upload from the uploadjob UI
			//uploadjob.UploadFromClipboard(&uploadDir, updateCounterCh)
			dialog.ShowInformation("Not implemented yet", "This feature is not implemented yet", myWindow)
		})

		desk.SetSystemTrayMenu(m)
	}

	myWindow.SetContent(widget.NewLabel("Fyne System Tray"))
	myWindow.SetCloseIntercept(func() {
		myWindow.Hide()
	})
	myWindow.SetContent(mainContent)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()

}

func setupMenu(window fyne.Window, content *fyne.Container, uploadDir *string, updateCounterCh chan int) fyne.CanvasObject {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	uploadButton := widget.NewButton("Upload", func() {
		uploadjob.ShowUploadUI(window, content, uploadDir, updateCounterCh)
	})

	historyButton := widget.NewButton("History", func() {
		table := uploadjob.CreateHistoryTable(*uploadDir)
		if table != nil {
			content.Objects = []fyne.CanvasObject{table}
			content.Refresh()
		}
	})

	settingsButton := widget.NewButton("Settings", func() {
		// Implement settings view
		uploadjob.ShowSettingsUI(window, content, uploadDir)
	})
	return container.NewVBox(uploadButton, historyButton, settingsButton)
}
