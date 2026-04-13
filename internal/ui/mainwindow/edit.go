package mainwindow

import (
	"DatabaseDB/internal/ui/ids"
	"DatabaseDB/internal/ui/labelkv"
	"bytes"
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type EditColumn struct {
	edit2         *fyne.Container
	container     *fyne.Container
	cancelEditKey *widget.Button
	saveEditKey   *widget.Button
	valueEntry    *widget.Entry
	finishValue   string
	editLabel     *widget.Label
}

func (e *MainWindow2) SaveAndCancle() *fyne.Container {
	return container.NewGridWithColumns(2, e.EditColumn.cancelEditKey, e.EditColumn.saveEditKey)
}

func (e *MainWindow2) ConfigureEntry(content string) *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Resize(fyne.NewSize(400, 500))
	entry.SetText(content)
	scrollableEntry := container.NewScroll(entry)
	scrollableEntry.SetMinSize(fyne.NewSize(200, 300))
	e.EditColumn.edit2.Add(scrollableEntry)
	return entry
}

func (m *EditColumn) SetLabelEdit(value string, editType labelkv.EditType) {
	label := fmt.Sprintf("%s %s : %s", ids.EditLabel, editType, value)
	m.editLabel.SetText(label)
	return
}

var BaseImage []byte

func (m *MainWindow2) ImageShow(key []byte, value []byte, types string) {
	var lableAddpicture *widget.Button
	var image *canvas.Image

	BaseImage = value

	image = canvas.NewImageFromResource(fyne.NewStaticResource("placeholder.png", value))
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(300, 300))
	m.EditColumn.edit2.Add(image)

	lableAddpicture = widget.NewButton("+", func() {
		folderPath := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil || dir == nil {
				fmt.Println("Error opening folder or no folder selected")
				return
			}
			valueFinish, err := io.ReadAll(dir)
			if err != nil {
				fmt.Print("Error reading file:", err)
				return
			}

			image.Resource = fyne.NewStaticResource("image.png", valueFinish)
			image.Refresh()

			if !bytes.Equal(valueFinish, BaseImage) {
				m.EditColumn.saveEditKey.Enable()
				m.EditColumn.finishValue = string(valueFinish)
			} else {
				m.EditColumn.saveEditKey.Disable()
			}
			NameLabel = fmt.Sprintf("* %s . . .", types)
			//ValueImage = valueFinish
		}, m.Window)

		folderPath.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".gif"}))
		folderPath.Show()
	})
	m.EditColumn.edit2.Add(lableAddpicture)
}
