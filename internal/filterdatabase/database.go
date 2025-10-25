package filterdatabase

import (
	"fyne.io/fyne/v2/dialog"
)

type FilterData interface {
	FilterFile(path string) bool
	FilterFormat(folderDialog *dialog.FileDialog)
}
