package mainwindow

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type EditColumn2 struct {
	Edit2         *fyne.Container
	Container     *fyne.Container
	CancelEditKey *widget.Button
	SaveEditKey   *widget.Button
	ValueEntry    *widget.Entry
}

func (e *MainWindow2) SaveAndCancle() *fyne.Container {
	return container.NewGridWithColumns(2, e.EditColumn.CancelEditKey, e.EditColumn.SaveEditKey)
}

func (e *MainWindow2) ConfigureEntry(content string) *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Resize(fyne.NewSize(400, 500))
	entry.SetText(content)
	scrollableEntry := container.NewScroll(entry)
	scrollableEntry.SetMinSize(fyne.NewSize(200, 300))
	e.EditColumn.Edit2.Add(scrollableEntry)
	return entry
}
