package main

import (
	"DatabaseDB/internal/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("com.DatabaseDB.KV-Toolbox")

	windiw := mainwindow.NewMainWindow("ManageDB")

	windiw.MainWindow(myApp)
}
