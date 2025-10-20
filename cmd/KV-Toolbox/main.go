package main

import (
	configApp "DatabaseDB/internal/config"
	"DatabaseDB/internal/model"
	"DatabaseDB/internal/ui/mainwindow"
)

func main() {
	config := configApp.LoadConfig()

	app := model.NewApp("ManageDB")

	app.MainWindow.Config = config

	mainwindow.MainWindow(app.App)
}
