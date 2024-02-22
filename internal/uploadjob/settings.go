package uploadjob

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

	dailyGoal := widget.NewEntry()
	dailyGoal.SetText(config.String("dailyGoal", "10"))
	tgApiKey := widget.NewEntry()
	tgApiKey.SetText(config.String("tgApi", ""))
	tgReceiverID := widget.NewEntry()
	tgReceiverID.SetText(config.String("tgReceiverID", ""))
	oaiAPIKey := widget.NewEntry()
	oaiAPIKey.SetText(config.String("openaiKey", ""))

	// Save action
	saveAction := func() {
		// Assuming you have the logic to get new values from your UI elements
		newRootFolder := *rootFolder   // Example for folder path, adjust as needed
		newDailyGoal := dailyGoal.Text // Assuming dailyGoal is a binding.Int
		newTgApi := tgApiKey.Text
		newTgReceiverID := tgReceiverID.Text

		// Set new values to config
		config.Set("rootFolder", newRootFolder)
		config.Set("dailyGoal", newDailyGoal)
		config.Set("tgApi", newTgApi)
		config.Set("tgReceiverID", newTgReceiverID)
		config.Set("openaiKey", oaiAPIKey.Text)

		// Call SaveConfig to write changes to file
		SaveConfig()
		// Optionally, show a dialog indicating success
		dialog.ShowInformation("Settings", "Settings saved successfully!", window)
	}
	settingsForm := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Daily Goal", Widget: dailyGoal},
			{Text: "Telegram API key", Widget: tgApiKey},
			{Text: "Telegram receiver ID", Widget: tgReceiverID},
			{Text: "OpenAI API key", Widget: oaiAPIKey},
		},
		OnSubmit: func() { // optional, handle form submission
			log.Println("tg api :", tgApiKey.Text)
			log.Println("tg chatID:", tgReceiverID.Text)
			log.Println("daily goal:", dailyGoal.Text)
			saveAction()
		},
	}
	settingsForm.SubmitText = "Save Settings"
	// Layout your settings UI
	settingsBox := container.NewVBox(
		widget.NewLabel("Telegram Notification Settings"),
		settingsForm,
	)

	// Update the content container to include the new settings form
	content.Objects = []fyne.CanvasObject{settingsBox}
	content.Refresh()

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(settingsForm, layout.NewSpacer(), selectFolderButton),
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
