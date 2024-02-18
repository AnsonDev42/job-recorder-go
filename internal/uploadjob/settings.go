package uploadjob

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gookit/config/v2"
	"log"
	"os"
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

	// set, bind,update json for the daily goal
	dailyGoal := binding.NewInt()
	dailyGoal.Set(config.Int("dailyGoal", 10))
	dW := widget.NewLabelWithData(binding.IntToString(dailyGoal))
	dChangeDailyGoal := widget.NewEntryWithData(binding.IntToString(dailyGoal))
	dChangeDailyGoal.OnChanged = func(string) {
		newGoal, _ := dailyGoal.Get()
		config.Set("dailyGoal", newGoal)
		SaveConfig()
	}

	tgApiKey := widget.NewEntry()
	tgApiKey.SetText(config.String("tgApi", ""))
	tgReceiverID := widget.NewEntry()
	tgReceiverID.SetText(config.String("tgReceiverID", ""))

	// Save button
	saveButton := widget.NewButton("Save telegram settings", func() {
		// Assuming you have the logic to get new values from your UI elements
		newRootFolder := *rootFolder       // Example for folder path, adjust as needed
		newDailyGoal, _ := dailyGoal.Get() // Assuming dailyGoal is a binding.Int
		newTgApi := tgApiKey.Text
		newTgReceiverID := tgReceiverID.Text

		// Set new values to config
		config.Set("rootFolder", newRootFolder)
		config.Set("dailyGoal", newDailyGoal)
		config.Set("tgApi", newTgApi)
		config.Set("tgReceiverID", newTgReceiverID)

		// Call SaveConfig to write changes to file
		SaveConfig()
		// Optionally, show a dialog indicating success
		dialog.ShowInformation("Settings", "Settings saved successfully!", window)
	})

	// Layout your settings UI
	settingsForm := container.NewVBox(
		widget.NewLabel("Telegram Notification Settings"),
		// Include your other settings widgets here
		tgApiKey,
		tgReceiverID,
		saveButton, // Add the save button to the UI
	)

	// Update the content container to include the new settings form
	content.Objects = []fyne.CanvasObject{settingsForm}
	content.Refresh()
	dShowDailyGoal := widget.NewLabel("Current Setting for daily goal is: ")
	dailGoalRow := container.NewHBox(dShowDailyGoal, dW, layout.NewSpacer(), dChangeDailyGoal)

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(dailGoalRow, settingsForm, layout.NewSpacer(), selectFolderButton),
	}
	content.Refresh()
}

func SaveConfig() {
	buf := new(bytes.Buffer)

	_, err := config.DumpTo(buf, config.JSON)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	if _, err := os.Stat(config.String("configPATH")); os.IsNotExist(err) {
		os.MkdirAll(config.String("configPATH"), 0755)
		log.Fatalf("Failed to create the config file in the given path %s", config.String("configPATH"))
	}
	err = os.WriteFile(config.String("configPATH"), buf.Bytes(), 0755)
	if err != nil {
		log.Fatal(err)
	}
}
