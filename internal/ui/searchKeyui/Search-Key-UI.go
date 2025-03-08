package searchkeyui

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/logic"
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
	editWindow := fyne.CurrentApp().NewWindow("Search in the database")
	editWindow.Resize(fyne.NewSize(600, 300))

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.Resize(fyne.NewSize(500, 500))
	valueEntry.SetPlaceHolder("Key for Search")

	buttomSearch := widget.NewButton("Search", func() {

		keys, values, err := logic.SearchDatabase(valueEntry.Text)
		if err != nil {
			dialog.ShowInformation("Error", "Such a key is not available in the database", editWindow)
		}

		utils.CheckCondition(columnEditKey)
		utils.CheckCondition(rightColumnContent)
		var truncatedValue string

		for i := 0; i < len(values); i++ {
			if i > 40 {
				dialog.ShowInformation("Error", "The result of your keys is more than 60 and I will only show the first 60.If your key is not among these, please search more precisely.", mainWindow)
				break
			}
			truncatedKey := utils.TruncateString(string(keys[i]), 20)

			typeValue := mimetype.Detect(values[i])
			if typeValue.Extension() != ".txt" {
				truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
			} else {
				truncatedValue = utils.TruncateString(string(values[i]), 20)
			}
			valueLabel := logic.BuidLableKeyAndValue("value", keys[i], values[i], truncatedValue, rightColumnContent, columnEditKey, saveKey, mainWindow)
			keyLabel := logic.BuidLableKeyAndValue("key", keys[i], values[i], truncatedKey, rightColumnContent, columnEditKey, saveKey, mainWindow)

			rightColumnContent.Refresh()
			buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
			rightColumnContent.Add(buttonRow)
		}

		editWindow.Close()

		variable.ResultSearch = true

	})
	buttomSearch.Importance = widget.HighImportance
	editContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomSearch,
	)
	editWindow.SetContent(editContent)
	editWindow.Show()
}
