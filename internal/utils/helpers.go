// internal/utils/helpers.go
package utils

import (
	variable "DatabaseDB"
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/Databaces/PebbleDB"
	badgerDB "DatabaseDB/internal/Databaces/badger"
	leveldbb "DatabaseDB/internal/Databaces/leveldb"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var ValueImage []byte

func TruncateString(input string, length int) string {
	nameData := input
	if len(nameData) > length {
		nameData = nameData[:length] + ". . ."
	}
	parts := strings.Split(nameData, "\n")
	if len(parts) > 1 {

		nameData = parts[0] + " . . ."
	}

	return nameData
}

func CheckCondition(rightColumnContent *fyne.Container) {
	if len(rightColumnContent.Objects) > 0 {
		rightColumnContent.Objects = []fyne.CanvasObject{}
		rightColumnContent.Refresh()
	}
}

func Checkdatabace(test string, nameDatabace string) error {
	var newClient dbpak.DBClient

	switch nameDatabace {
	case "levelDB":
		newClient = leveldbb.NewDataBaseLeveldb(test)
	case "Pebble":
		newClient = PebbleDB.NewDataBasePebble(test)
	case "Badger":
		newClient = badgerDB.NewDataBaseBadger(test)
	case "Redis":
		//newClient = Redisdb.NewDataBaseRedis(parts[0], parts[1], parts[2])
	}

	variable.CloseAndSetCurrentDBClient(newClient)

	if err := variable.GetCurrentDBClient().Open(); err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if nameDatabace != "Redis" {
		if _, err := os.Stat(test); os.IsNotExist(err) && !variable.CreatDatabase {
			return err
		}
	}

	return nil
}

func CleanInput(input string) string {
	cleaned := strings.TrimSpace(input)
	return cleaned
}

func ImageShow(key []byte, value []byte, mainContainer *fyne.Container, editWindow fyne.Window) {
	var lableAddpicture *widget.Button
	var image *canvas.Image

	image = canvas.NewImageFromResource(fyne.NewStaticResource("placeholder.png", value))
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(300, 300))
	mainContainer.Add(image)

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
		}, editWindow)

		folderPath.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".gif"}))
		folderPath.Show()
	})
	mainContainer.Add(lableAddpicture)
}
