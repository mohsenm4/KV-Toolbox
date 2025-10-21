package mainwindow

import (
	"DatabaseDB/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RightColumn2 struct {
	Container            *fyne.Container
	NameButtonProject    *widget.Label
	Line                 *canvas.Line
	Spacer               *widget.Label
	ButtonDelete         *widget.Button
	SearchButton         *widget.Button
	ButtonAdd            *widget.Button
	KeyRightColunm       *widget.Label
	ValueRightColunm     *widget.Label
	ParentWindow         fyne.Window
	LastLableKeyAndValue *utils.TappableLabel
}

func NewRightColumn(parentWindow *fyne.Window) *RightColumn2 {
	return &RightColumn2{
		ParentWindow: *parentWindow,
	}
}

func (r *RightColumn2) TopRightColumn() *fyne.Container {
	r.Line.StrokeWidth = 2
	container := container.NewVBox(
		r.NameButtonProject,
		r.Line,
		r.Spacer,
		r.Tool(),
		r.KeyAndValue(),
	)
	return container
}

func (r *RightColumn2) Tool() *fyne.Container {
	return container.NewGridWithColumns(3, r.ButtonDelete, r.SearchButton, r.ButtonAdd)
}

func (r *RightColumn2) KeyAndValue() *fyne.Container {
	return container.NewGridWithColumns(6, r.KeyRightColunm, widget.NewLabel(""), r.ValueRightColunm, widget.NewLabel(""))
}
