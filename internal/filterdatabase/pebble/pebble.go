package Filterpebbledb

import (
	"DatabaseDB/internal/filterdatabase"
	sharedfunc "DatabaseDB/internal/filterdatabase/SharedFunc"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
