package otherUI

import (

	// "DatabaseDB/internal/logic/addProjectwindowlogic"

	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	lastStart             *[]byte
	lastEnd               *[]byte
	lastPage              int
	Orgdata               []dbpak.KVData
	previousClose         *widget.Button
	previousProject       *widget.Button
	previousRefreshButton *widget.Button
	lastLableKeyAndValue  *utils.TappableLabel
)

func SetupThemeButtons(app fyne.App) *fyne.Container {
	var darkButton *widget.Button
	var lightButton *widget.Button

	darkButton = widget.NewButton("Dark", func() {
		app.Settings().SetTheme(theme.DarkTheme())

		darkButton.Importance = widget.HighImportance
		lightButton.Importance = widget.MediumImportance
		darkButton.Refresh()
		lightButton.Refresh()
	})
	lightButton = widget.NewButton("Light", func() {
		app.Settings().SetTheme(theme.LightTheme())

		lightButton.Importance = widget.HighImportance
		darkButton.Importance = widget.MediumImportance

		darkButton.Refresh()
		lightButton.Refresh()
	})

	darkButton.Importance = widget.HighImportance

	darkLight := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, lightButton, darkButton),
	)
	return darkLight
}
