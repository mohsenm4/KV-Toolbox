package mainwindow

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (mi *MainWindow2) DeleteKeyUi() {

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.Resize(fyne.NewSize(500, 500))
	valueEntry.SetPlaceHolder("Key for Delete")

	buttomDelete := widget.NewButton("Delete", nil)
	buttomDelete.Importance = widget.HighImportance

	editContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomDelete,
		layout.NewSpacer(),
	)

	ded := dialog.NewCustom("Delete in the database", "Close", editContent, mi.Window)
	ded.Resize(fyne.NewSize(600, 300))

	valueEntry.OnChanged = func(s string) {
		if s != "" {
			buttomDelete.Enable()
		} else {
			buttomDelete.Disable()
		}
	}
	buttomDelete.OnTapped = func() {

		message := fmt.Sprintf("Are you sure you want to delete the key: _ %s _?", valueEntry.Text)

		dialog.ShowConfirm("Confirm Delete", message,
			func(response bool) {
				if response {
					err := mi.Logic.DeleteKeyLogic(valueEntry.Text)
					if err != nil {
						dialog.ShowInformation("Error", err.Error(), mi.Window)
					} else {
						ded.Hide()
					}
				}
			}, mi.Window)
	}
	ded.Show()
	mi.Window.Canvas().Focus(valueEntry)
}
