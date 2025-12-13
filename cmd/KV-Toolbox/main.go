package main

import (
	"DatabaseDB/internal/pref"
	"DatabaseDB/internal/ui/mainwindow"
	"DatabaseDB/internal/ui/them"

	"fyne.io/fyne/v2/app"
)

func main() {

	myApp := app.NewWithID("com.DatabaseDB.KV-Toolbox")

	window := mainwindow.NewMainWindow("ManageDB")

	window.Pref = pref.NewPref(myApp)

	mytheme := window.Pref.LoadTheme(pref.KeyTheme)

	them.SetThemeByKey(myApp, mytheme)
	window.MainWindow(myApp)

}
