package mainwindow

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

func (r *MainWindow2) SearchKeyUi() {
	valueEntry := widget.NewMultiLineEntry()
	valueEntry.SetPlaceHolder("Key for Search")

	buttomSearch := widget.NewButton("Search", nil)
	buttomSearch.Importance = widget.HighImportance
	buttomSearch.Disable()

	modalContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomSearch,
		layout.NewSpacer(),
	)

	d := dialog.NewCustom("Search in the database", "Close", modalContent, r.Window)
	d.Resize(fyne.NewSize(600, 300))

	valueEntry.OnChanged = func(s string) {
		if s != "" {
			buttomSearch.Enable()
		} else {
			buttomSearch.Disable()
		}
	}

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

		utils.CheckCondition(r.EditColumn.edit2)
		utils.CheckCondition(r.RightColumn.container)

		newList := widget.NewList(
			func() int {
				return len(values)
			},
			func() fyne.CanvasObject {
				keyLabel := widget.NewLabel("key")
				valueLabel := widget.NewLabel("value")
				buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
				return buttonRow
			},
			func(i widget.ListItemID, obj fyne.CanvasObject) {

				item := values[i]

				typeValue := mimetype.Detect(item)
				var truncatedValue string
				if typeValue.Extension() != ".txt" {
					truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
				} else {
					truncatedValue = utils.TruncateString(string(item), 20)
				}
				truncatedKey := utils.TruncateString(string(keys[i]), 20)

				keyLabel := r.BuildLabelKeyAndValue("key", keys[i], item, truncatedKey)
				valueLabel := r.BuildLabelKeyAndValue("value", keys[i], item, truncatedValue)

				row := obj.(*fyne.Container)
				row.Objects[0] = keyLabel
				row.Objects[1] = valueLabel

			},
		)
		r.RightColumn.list = newList
		r.RightColumn.container.Objects = nil
		r.RightColumn.container.Add(newList)
		r.RightColumn.container.Refresh()
		d.Hide()
		variable.ResultSearch = true
	}

	d.Show()
	r.Window.Canvas().Focus(valueEntry)
}
