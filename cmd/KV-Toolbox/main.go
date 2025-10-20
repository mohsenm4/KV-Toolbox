package main

import (
	configApp "DatabaseDB/internal/config"
	"DatabaseDB/internal/model"
	"DatabaseDB/internal/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	config := configApp.LoadConfig()

	app := model.NewApp("ManageDB")

	app.MainWindow.Config = config

	mainwindow.MainWindow(myApp)
}
