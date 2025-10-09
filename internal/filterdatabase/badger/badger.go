package Filterbadger

import (
	"DatabaseDB/internal/filterdatabase"
	sharedfunc "DatabaseDB/internal/filterdatabase/SharedFunc"
	"io/ioutil"
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type NameDatabaseBadger struct{}

func NewFileterBadger() filterdatabase.FilterData {
	return &NameDatabaseBadger{}
}

func (l *NameDatabaseBadger) FilterFile(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Error reading folder:", err)
		return false
	}
	var count uint8
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sst" || filepath.Ext(file.Name()) == ".vlog" {
			count++
		}

		if count == 2 {
			return true
		}
	}
	return false
}

func (l *NameDatabaseBadger) FilterFormat(folderDialog *dialog.FileDialog) {
	folderDialog.SetFilter(storage.NewExtensionFileFilter([]string{".sst", ".vlog"}))
}

func (l *NameDatabaseBadger) FormCreate(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) {
	sharedfunc.FormPasteDatabase(a, title, lastColumnContent, rightColumnContentORG, nameButtonProject, buttonAdd, columnEditKey, saveKey, mainWindow)

}
