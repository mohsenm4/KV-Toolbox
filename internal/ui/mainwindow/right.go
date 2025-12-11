package mainwindow

import (
	variable "DatabaseDB"
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/ui/labelkv"
	"image/color"
	"runtime"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RightColumn struct {
	container            *fyne.Container
	nameButtonProject    *widget.Label
	buttonDelete         *widget.Button
	searchButton         *widget.Button
	buttonAdd            *widget.Button
	keyRightColunm       *widget.Label
	valueRightColunm     *widget.Label
	lastLableKeyAndValue *labelkv.TappableLabel
	lastStart            *[]byte
	lastEnd              *[]byte
	lastPage             int
	orgdata              []dbpak.KVData
}

func NewRightColumn() *RightColumn {
	return &RightColumn{}
}

func (r *MainWindow2) TopRightColumn() *fyne.Container {
	r.Objects.line = canvas.NewLine(color.Black)
	r.Objects.line.StrokeWidth = 2

	container := container.NewVBox(
		r.RightColumn.nameButtonProject,
		r.Objects.line,
		r.Objects.spacer,
		r.RightColumn.Tool(),
		r.RightColumn.KeyAndValue(),
	)
	return container
}

func (r *RightColumn) Tool() *fyne.Container {
	return container.NewGridWithColumns(3, r.buttonDelete, r.searchButton, r.buttonAdd)
}

func (r *RightColumn) KeyAndValue() *fyne.Container {
	return container.NewGridWithColumns(6, r.keyRightColunm, widget.NewLabel(""), r.valueRightColunm, widget.NewLabel(""))
}

func (r *MainWindow2) UpdatePage() {

	data, err := logic.FetchPageData(r.RightColumn.lastStart, r.RightColumn.lastEnd, r.RightColumn.lastPage, r.RightColumn.orgdata)
	if err != nil {
		return
	}

	if r.RightColumn.lastPage < variable.CurrentPage {

		if len(r.RightColumn.orgdata) >= variable.ItemsPerPage*3 {
			r.RightColumn.orgdata = r.RightColumn.orgdata[len(data):]
		}

		r.RightColumn.orgdata = append(r.RightColumn.orgdata, data...)

	} else {

		r.RightColumn.orgdata = r.RightColumn.orgdata[:len(r.RightColumn.orgdata)-len(data)]
		r.RightColumn.orgdata = append(data, r.RightColumn.orgdata...)
	}

	if len(r.RightColumn.orgdata) != 0 {
		r.RightColumn.lastStart = &r.RightColumn.orgdata[0].Key
		r.RightColumn.lastEnd = &r.RightColumn.orgdata[len(r.RightColumn.orgdata)-1].Key
	}

	var truncatedValue string
	var truncatedKey string

	var arrayContainer []fyne.CanvasObject
	for _, item := range data {

		truncatedKey, truncatedValue = logic.FormatKeyValue(item)

		valueLabel := r.NewLabelKV(labelkv.EditValue, item.Key, item.Value, truncatedValue)
		keyLabel := r.NewLabelKV(labelkv.EditKey, item.Key, item.Value, truncatedKey)

		valueLabel.SetKeyLabel(keyLabel)
		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		arrayContainer = append(arrayContainer, buttonRow)
	}

	if r.RightColumn.lastPage > variable.CurrentPage {
		r.RightColumn.container.Objects = applyPageShift(r.RightColumn.container.Objects, arrayContainer, true, variable.ItemsPerPage*3)
	} else {
		r.RightColumn.container.Objects = applyPageShift(r.RightColumn.container.Objects, arrayContainer, false, variable.ItemsPerPage*3)
	}

	arrayContainer = nil
	data = nil
	runtime.GC()
	debug.FreeOSMemory()
	r.RightColumn.container.Refresh()
	r.RightColumn.lastPage = variable.CurrentPage
}

// ApplyPageShift adds or removes items from the objects slice based on the goUp flag.
// If goUp is true, it prepends and if goUp false, it appends the arrayContainer to objects,
func applyPageShift(objects []fyne.CanvasObject, arrayContainer []fyne.CanvasObject, goUp bool, maxItems int) []fyne.CanvasObject {
	isFull := len(objects) == maxItems
	cnt := len(arrayContainer)

	if goUp {
		if isFull {
			copy(objects[cnt:], objects[:len(objects)-cnt])
			copy(objects[:cnt], arrayContainer)
		} else {
			objects = append(arrayContainer, objects...)
		}
	} else {
		if isFull {
			copy(objects[:len(objects)-cnt], objects[cnt:])
			copy(objects[len(objects)-cnt:], arrayContainer)
		} else {
			objects = append(objects, arrayContainer...)
		}
	}

	arrayContainer = nil
	runtime.GC()
	debug.FreeOSMemory()
	return objects
}
