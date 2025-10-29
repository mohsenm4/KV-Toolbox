package mainwindow

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type EditColumn2 struct {
	Edit2         *fyne.Container
	Container     *fyne.Container
	CancelEditKey *widget.Button
	SaveEditKey   *widget.Button
	ValueEntry    *widget.Entry
	FinishValue   string
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

var BaseImage []byte

func hashSlice(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func (m *MainWindow2) ImageShow(key []byte, value []byte) {
	var lableAddpicture *widget.Button
	var image *canvas.Image

	BaseImage = value

	image = canvas.NewImageFromResource(fyne.NewStaticResource("placeholder.png", value))
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(300, 300))
	m.EditColumn.Edit2.Add(image)

	lableAddpicture = widget.NewButton("+", func() {
		folderPath := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil || dir == nil {
				fmt.Println("Error opening folder or no folder selected")
				return
			}
			valueFinish, err := ioutil.ReadAll(dir)
			if err != nil {
				fmt.Print("Error reading file:", err)
				return
			}

			image.Resource = fyne.NewStaticResource("image.png", valueFinish)
			image.Refresh()

			if hashSlice(valueFinish) != hashSlice(BaseImage) {
				m.EditColumn.SaveEditKey.Enable()
			}
			//ValueImage = valueFinish
			m.EditColumn.FinishValue = string(valueFinish)
		}, m.Window)

		folderPath.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".gif"}))
		folderPath.Show()
	})
	m.EditColumn.Edit2.Add(lableAddpicture)
}
