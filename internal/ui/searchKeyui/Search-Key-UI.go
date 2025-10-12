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

func SearchKeyUi(rightColumnContent *fyne.Container, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) {
	valueEntry := widget.NewMultiLineEntry()
	valueEntry.SetPlaceHolder("Key for Search")

	buttomSearch := widget.NewButton("Search", nil)
	buttomSearch.Importance = widget.HighImportance

	modalContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomSearch,
	)

	d := dialog.NewCustom("Search in the database", "Close", modalContent, mainWindow)
	d.Resize(fyne.NewSize(600, 300))

	buttomSearch.OnTapped = func() {
		keys, values, err := logic.SearchDatabase(valueEntry.Text)
		if err != nil {
			dialog.ShowInformation("Error", "Such a key is not available in the database", mainWindow)
			return
		}

		utils.CheckCondition(columnEditKey)
		utils.CheckCondition(rightColumnContent)

		var truncatedValue string

		for i := 0; i < len(values); i++ {
			if i > 40 {
				dialog.ShowInformation("Warning",
					"The result of your keys is more than 60. Only the first 60 are shown.\nIf your key is not among these, please search more precisely.",
					mainWindow)
				break
			}

			truncatedKey := utils.TruncateString(string(keys[i]), 20)
			typeValue := mimetype.Detect(values[i])

			if typeValue.Extension() != ".txt" {
				truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
			} else {
				truncatedValue = utils.TruncateString(string(values[i]), 20)
			}

			valueLabel := otherUI.BuidLableKeyAndValue("value", keys[i], values[i], truncatedValue, rightColumnContent, columnEditKey, saveKey, mainWindow)
			keyLabel := otherUI.BuidLableKeyAndValue("key", keys[i], values[i], truncatedKey, rightColumnContent, columnEditKey, saveKey, mainWindow)

			buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
			rightColumnContent.Add(buttonRow)
			rightColumnContent.Refresh()
		}

		d.Hide()
		variable.ResultSearch = true
	}

	d.Show()
}
