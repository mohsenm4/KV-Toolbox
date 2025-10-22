package searchkeyui

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/ui/otherUI"
	"DatabaseDB/internal/utils"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

func ShowSearchKeyDialog(rightColumn *fyne.Container, editColumn *fyne.Container, saveButton *widget.Button, mainWindow fyne.Window) {
	keyEntry := widget.NewMultiLineEntry()
	keyEntry.SetPlaceHolder("Enter key to search")

	searchButton := widget.NewButton("Search", nil)
	searchButton.Importance = widget.HighImportance

	dialogContent := container.NewVBox(
		widget.NewLabel("Enter the key you want to search for"),
		keyEntry,
		layout.NewSpacer(),
		searchButton,
		layout.NewSpacer(),
	)

	searchDialog := dialog.NewCustom("Search in Database", "Close", dialogContent, mainWindow)
	searchDialog.Resize(fyne.NewSize(600, 300))

	searchButton.OnTapped = func() {
		keys, values, err := logic.SearchDatabase(keyEntry.Text)
		if err != nil {
			dialog.ShowInformation(
				"Database Error",
				"An error occurred while searching the database.\nPlease check your input or try again.",
				mainWindow,
			)
			return
		} else if len(keys) == 0 && len(values) == 0 {
			dialog.ShowInformation(
				"No Results Found",
				"No data was found for the entered key.\nPlease make sure the key exists in the database.",
				mainWindow,
			)
			return
		}

		utils.ClearContainerIfNotEmpty(editColumn)
		utils.ClearContainerIfNotEmpty(rightColumn)

		var truncatedValue string

		for i := 0; i < len(values); i++ {
			if i > 40 {
				dialog.ShowInformation("Warning",
					"The number of results exceeds 60. Only the first 60 are shown.\nIf your key is not among these, please search more precisely.",
					mainWindow)
				break
			}

			truncatedKey := utils.TruncateString(string(keys[i]), 20)
			detectedType := mimetype.Detect(values[i])

			if detectedType.Extension() != ".txt" {
				truncatedValue = fmt.Sprintf("* %s . . .", detectedType.Extension())
			} else {
				truncatedValue = utils.TruncateString(string(values[i]), 20)
			}

			valueLabel := otherUI.BuidLableKeyAndValue("value", keys[i], values[i], truncatedValue, rightColumn, editColumn, saveButton, mainWindow)
			keyLabel := otherUI.BuidLableKeyAndValue("key", keys[i], values[i], truncatedKey, rightColumn, editColumn, saveButton, mainWindow)

			row := container.NewGridWithColumns(2, keyLabel, valueLabel)
			rightColumn.Add(row)
			rightColumn.Refresh()
		}

		searchDialog.Hide()
		variable.ResultSearch = true
	}

	searchDialog.Show()
}
