package model

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

type App struct {
	App        fyne.App
	Name       string
	MainWindow MainWindow
}

func NewApp(name string) *App {

	myApp := app.New()

	iconResource := theme.FyneLogo()
	myApp.SetIcon(iconResource)

	return &App{
		Name: name,
		App:  myApp,
	}
}
