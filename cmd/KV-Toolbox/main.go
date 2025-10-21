package main

import (
	variable "DatabaseDB"
	configApp "DatabaseDB/internal/config"
	"DatabaseDB/internal2/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	variable.CurrentJson = configApp.NewConfig()

	windiw := mainwindow.NewMainWindow()
	windiw.MainWindow(myApp)
}
