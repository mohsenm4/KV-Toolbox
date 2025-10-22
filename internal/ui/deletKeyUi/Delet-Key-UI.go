package deletkeyui

import (
	"DatabaseDB/internal/logic"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowDeleteKeyDialog(rightColumn *fyne.Container, mainWindow fyne.Window) {
	// Input field for key to delete
	keyEntry := widget.NewMultiLineEntry()
	keyEntry.Resize(fyne.NewSize(500, 500))
	keyEntry.SetPlaceHolder("Enter key to delete")

	// Delete button
	deleteButton := widget.NewButton("Delete", nil)
	deleteButton.Importance = widget.HighImportance

	// Layout for dialog content
	dialogContent := container.NewVBox(
		widget.NewLabel("Enter the key you want to delete"),
		keyEntry,
		layout.NewSpacer(),
		deleteButton,
		layout.NewSpacer(),
	)

	// Dialog setup
	deleteDialog := dialog.NewCustom("Delete from Database", "Close", dialogContent, mainWindow)
	deleteDialog.Resize(fyne.NewSize(600, 300))

	// Button click handler
	deleteButton.OnTapped = func() {
		message := fmt.Sprintf("Are you sure you want to delete the key: _%s_?", keyEntry.Text)

		dialog.ShowConfirm("Confirm Delete", message, func(confirmed bool) {
			if confirmed {
				err := logic.DeleteKeyLogic(keyEntry.Text)
				if err != nil {
					dialog.ShowInformation("Error", err.Error(), mainWindow)
				} else {
					deleteDialog.Hide()
				}
			}
		}, mainWindow)
	}

	deleteDialog.Show()
}
