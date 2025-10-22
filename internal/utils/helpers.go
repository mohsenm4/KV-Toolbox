// internal/utils/helpers.go
package utils

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/Databaces/PebbleDB"
	badgerDB "DatabaseDB/internal/Databaces/badger"
	leveldbDB "DatabaseDB/internal/Databaces/leveldb"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var ImageBuffer []byte

// TappableLabel represents a label that responds to click/tap events.
type TappableLabel struct {
	widget.Label
	OnTapped func()
}

// NewTappableLabel creates a new tappable label with a click handler.
func NewTappableLabel(text string, onTapped func()) *TappableLabel {
	label := &TappableLabel{
		Label: widget.Label{
			Text: text,
		},
		OnTapped: onTapped,
	}
	label.ExtendBaseWidget(label)
	return label
}

// Tapped triggers the assigned click handler.
func (t *TappableLabel) Tapped(_ *fyne.PointEvent) {
	t.OnTapped()
}

// TruncateString shortens a string and adds "..." if it exceeds a given length or contains multiple lines.
func TruncateString(input string, length int) string {
	result := input
	if len(result) > length {
		result = result[:length] + "..."
	}
	lines := strings.Split(result, "\n")
	if len(lines) > 1 {
		result = lines[0] + "..."
	}
	return result
}

// ClearContainerIfNotEmpty clears the container if it contains any objects.
func ClearContainerIfNotEmpty(container *fyne.Container) {
	if len(container.Objects) > 0 {
		container.Objects = []fyne.CanvasObject{}
		container.Refresh()
	}
}

// OpenDatabase initializes and opens a selected database client based on its name.
func OpenDatabase(path string, dbName string) error {
	if variable.CurrentDBClient != nil {
		variable.CurrentDBClient.Close()
	}

	switch dbName {
	case "levelDB":
		variable.CurrentDBClient = leveldbDB.NewDataBaseLeveldb(path)
	case "Pebble":
		variable.CurrentDBClient = PebbleDB.NewDataBasePebble(path)
	case "Badger":
		variable.CurrentDBClient = badgerDB.NewDataBaseBadger(path)
	case "Redis":
		// TODO: Add Redis database connection initialization here.
	}

	variable.CurrentDBClient.Open()

	if dbName != "Redis" {
		if _, err := os.Stat(path); os.IsNotExist(err) && !variable.CreatDatabase {
			return err
		}
	}

	return nil
}

// CleanInput trims whitespace from input strings.
func CleanInput(input string) string {
	return strings.TrimSpace(input)
}

// ShowImage loads and displays an image from bytes and allows updating it via file selection.
func ShowImage(key []byte, value []byte, container *fyne.Container, parentWindow fyne.Window) {
	var addImageButton *widget.Button
	var image *canvas.Image

	image = canvas.NewImageFromResource(fyne.NewStaticResource("placeholder.png", value))
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(300, 300))
	container.Add(image)

	addImageButton = widget.NewButton("+", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				fmt.Println("Error opening image or no file selected")
				return
			}

			go func() {
				data, err := ioutil.ReadAll(reader)
				if err != nil {
					fmt.Println("Error reading image:", err)
					return
				}

				image.Resource = fyne.NewStaticResource("image.png", data)
				image.Refresh()
				ImageBuffer = data
			}()
		}, parentWindow)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".gif"}))
		fileDialog.Show()
	})
	container.Add(addImageButton)
}
