package mainwindow

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

func (r *MainWindow2) SearchKeyUi() {
	valueEntry := widget.NewMultiLineEntry()
	valueEntry.SetPlaceHolder("Key for Search")

	buttomSearch := widget.NewButton("Search", nil)
	buttomSearch.Importance = widget.HighImportance

	modalContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomSearch,
		layout.NewSpacer(),
	)

	d := dialog.NewCustom("Search in the database", "Close", modalContent, r.Window)
	d.Resize(fyne.NewSize(600, 300))

	buttomSearch.OnTapped = func() {
		keys, values, err := logic.SearchDatabase(valueEntry.Text)
		if err != nil {
			dialog.ShowInformation(
				"Database Error",
				"An error occurred while searching the database.\nPlease check your input or try again.",
				r.Window,
			)
			return
		} else if len(keys) == 0 && len(values) == 0 {
			dialog.ShowInformation(
				"No Results Found",
				"No data was found for the entered key.\nPlease make sure the key exists in the database.",
				r.Window,
			)
			return
		}

		utils.CheckCondition(r.EditColumn.Container)
		utils.CheckCondition(r.RightColumn.Container)

		var truncatedValue string

		for i := 0; i < len(values); i++ {
			if i > 40 {
				dialog.ShowInformation("Warning",
					"The result of your keys is more than 60. Only the first 60 are shown.\nIf your key is not among these, please search more precisely.",
					r.Window)
				break
			}

			truncatedKey := utils.TruncateString(string(keys[i]), 20)
			typeValue := mimetype.Detect(values[i])

			if typeValue.Extension() != ".txt" {
				truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
			} else {
				truncatedValue = utils.TruncateString(string(values[i]), 20)
			}

			valueLabel := otherUI.BuidLableKeyAndValue("value", keys[i], values[i], truncatedValue, r.RightColumn.Container, r.EditColumn.Container, r.EditColumn.SaveEditKey, r.Window)
			keyLabel := otherUI.BuidLableKeyAndValue("key", keys[i], values[i], truncatedKey, r.RightColumn.Container, r.EditColumn.Container, r.EditColumn.SaveEditKey, r.Window)

			buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
			r.RightColumn.Container.Add(buttonRow)
			r.RightColumn.Container.Refresh()
		}

		d.Hide()
		variable.ResultSearch = true
	}

	d.Show()
}
