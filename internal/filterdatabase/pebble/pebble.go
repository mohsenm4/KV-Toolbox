package Filterpebbledb

import (
	"DatabaseDB/internal/filterdatabase"
	sharedfunc "DatabaseDB/internal/filterdatabase/SharedFunc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type NameDatabasePebble struct{}

func NewFileterPebble() filterdatabase.FilterData {
	return &NameDatabasePebble{}
}

func (l *NameDatabasePebble) FilterFile(path string) bool {
	return sharedfunc.FormatFilesDatabase(path)
}

func (l *NameDatabasePebble) FilterFormat(folderDialog *dialog.FileDialog) {
	folderDialog.SetFilter(storage.NewExtensionFileFilter([]string{".log"}))
}

func (l *NameDatabasePebble) FormCreate(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) {
	sharedfunc.FormPasteDatabase(a, title, lastColumnContent, rightColumnContentORG, nameButtonProject, buttonAdd, columnEditKey, saveKey, mainWindow)

}
