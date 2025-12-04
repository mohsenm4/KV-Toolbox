package main

import (
	"DatabaseDB/internal/pref"
	"DatabaseDB/internal/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("com.DatabaseDB.KV-Toolbox")

	window := mainwindow.NewMainWindow("ManageDB")

	window.Pref = pref.NewPref(myApp)

	window.MainWindow(myApp)
}
