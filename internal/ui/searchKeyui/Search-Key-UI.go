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

		data, err := logic.SearchDatabase(valueEntry.Text)
		if err != nil {
			dialog.ShowInformation("Error", "Such a key is not available in the database", editWindow)
		}

		err = variable.CurrentDBClient.Open()
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), editWindow)

		}
		defer variable.CurrentDBClient.Close()

		utils.CheckCondition(columnEditKey)
		utils.CheckCondition(rightColumnContent)
		var truncatedValue string
		var count int

		for _, item := range data {

			if count > 40 {
				dialog.ShowInformation("Error", "The result of your keys is more than 60 and I will only show the first 60.If your key is not among these, please search more precisely.", mainWindow)
				count = 0
				break
			}
			count++

			value, err := variable.CurrentDBClient.Get(item)
			if err != nil {
				return
			}
			truncatedKey := utils.TruncateString(string(item), 20)

			typeValue := mimetype.Detect([]byte(value))
			if typeValue.Extension() != ".txt" {
				truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
			} else {
				truncatedValue = utils.TruncateString(string(value), 20)

			}

			valueLabel := logic.BuidLableKeyAndValue("value", item, value, truncatedValue, rightColumnContent, columnEditKey, saveKey, mainWindow)
			keyLabel := logic.BuidLableKeyAndValue("key", item, value, truncatedKey, rightColumnContent, columnEditKey, saveKey, mainWindow)

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
