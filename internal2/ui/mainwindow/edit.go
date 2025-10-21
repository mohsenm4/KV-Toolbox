package mainwindow

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type EditColumn2 struct {
	Container     *fyne.Container
	CancelEditKey *widget.Button
	SaveEditKey   *widget.Button
}

func (e *EditColumn2) SaveAndCancle() *fyne.Container {
	return container.NewGridWithColumns(2, e.CancelEditKey, e.SaveEditKey)
}
