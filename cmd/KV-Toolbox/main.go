package main

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/pref"
	"DatabaseDB/internal/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("com.DatabaseDB.KV-Toolbox")

	variable.PrefValue = pref.NewPref(myApp)

	windiw := mainwindow.NewMainWindow("ManageDB")
	windiw.MainWindow(myApp)
}
