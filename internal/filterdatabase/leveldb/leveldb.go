package FilterLeveldb

import (
	"DatabaseDB/internal/filterdatabase"
	sharedfunc "DatabaseDB/internal/filterdatabase/SharedFunc"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

type NameDatabaseLeveldb struct{}

func NewFileterLeveldb() filterdatabase.FilterData {
	return &NameDatabaseLeveldb{}
}

func (l *NameDatabaseLeveldb) FilterFile(path string) bool {
	return sharedfunc.FormatFilesDatabase(path)
}

func (l *NameDatabaseLeveldb) FilterFormat(folderDialog *dialog.FileDialog) {
	folderDialog.SetFilter(storage.NewExtensionFileFilter([]string{".log"}))
}
