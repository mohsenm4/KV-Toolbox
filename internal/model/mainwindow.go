package model

import (
	"DatabaseDB/internal/config"

	"fyne.io/fyne/v2"
)

type MainWindow struct {
	Window *fyne.Window
	Config config.Config
}
