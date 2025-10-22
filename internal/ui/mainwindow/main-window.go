package mainwindow

import (
	variable "DatabaseDB"
	"image/color"

	Filterbadger "DatabaseDB/internal/filterdatabase/badger"
	FilterLeveldb "DatabaseDB/internal/filterdatabase/leveldb"
	Filterpebbledb "DatabaseDB/internal/filterdatabase/pebble"
	addkeyui "DatabaseDB/internal/ui/addKeyui"
	deletkeyui "DatabaseDB/internal/ui/deletKeyUi"
	"DatabaseDB/internal/ui/otherUI"
	searchkeyui "DatabaseDB/internal/ui/searchKeyui"
	"DatabaseDB/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var saveAndCancle *fyne.Container

var leveldbButton *widget.Button
var BottomDatabase []*widget.Button

func MainWindow(myApp fyne.App) {

	mainWindow := myApp.NewWindow("ManageDB")
	mainWindow.SetMaster()

	iconResource := theme.FyneLogo()
	myApp.SetIcon(iconResource)
	mainWindow.SetIcon(iconResource)

	spacer := widget.NewLabel("")

	// right column show key
	rightColumnAll := container.NewVBox()

	// right column Edit
	var rightColumEdit *fyne.Container

	line := canvas.NewLine(color.Black)
	line.StrokeWidth = 2

	// key top window for colunm keys
	keyRightColunm := widget.NewLabelWithStyle("key", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// value top window for colunm values
	valueRightColunm := widget.NewLabelWithStyle("value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// column key and value
	keyAndRight := container.NewGridWithColumns(6, keyRightColunm, widget.NewLabel(""), valueRightColunm, widget.NewLabel(""))

	// name bottom project in colunm right
	nameButtonProject := widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	saveEditKey := widget.NewButton("Save", nil)
	saveEditKey.Disable()

	cancelEditKey := widget.NewButton("Cancle", func() {
		utils.CheckCondition(rightColumEdit)
	})

	saveAndCancle = container.NewGridWithColumns(2, cancelEditKey, saveEditKey)

	rightColumEdit = container.NewVBox()

	columnEdit := container.NewBorder(widget.NewLabelWithStyle("Edit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), saveAndCancle, nil, nil, rightColumEdit)

	searchButton := widget.NewButton("Search", func() {

		searchkeyui.SearchKeyUi(rightColumnAll, rightColumEdit, saveEditKey, mainWindow)
	})

	buttonAdd := widget.NewButton("Add", func() {
		addkeyui.OpenWindowAddButton(myApp, rightColumnAll, mainWindow)
	})
	buttonAdd.Disable()
	searchButton.Disable()

	buttonDelete := widget.NewButton("Delete", func() {
		deletkeyui.ShowDeleteKeyDialog(rightColumnAll, mainWindow)
	})

	topRightColumn := container.NewVBox(
		nameButtonProject,
		line,
		spacer,
		container.NewGridWithColumns(3, buttonDelete, searchButton, buttonAdd),
		keyAndRight,
	)
	var pluss *widget.Button
	toggleButtonsContainer := container.NewVBox()
	buttonsVisible := false

	buttonDelete.Disable()
	// left column
	leftColumnAll := otherUI.SetupLastColumn(rightColumnAll, nameButtonProject, buttonAdd, searchButton, buttonDelete, rightColumEdit, saveEditKey, mainWindow)
	spacer.Resize(fyne.NewSize(0, 30))

	for _, name := range variable.NameDatabase {

		leveldbButton = widget.NewButton(name, func() {
			toggleButtonsContainer.Objects = nil
			buttonsVisible = false
			switch name {
			case "levelDB":
				variable.NameData = FilterLeveldb.NewFileterLeveldb()
			case "Pebble":
				variable.NameData = Filterpebbledb.NewFileterPebble()
			case "Badger":
				variable.NameData = Filterbadger.NewFileterBadger()
				//case "Redis":
				//	variable.NameData = Filterredis.NewFileterRedis()

			}
			variable.NameData.FormCreate(myApp, name, leftColumnAll, rightColumnAll, nameButtonProject, buttonAdd, searchButton, buttonDelete, rightColumEdit, saveEditKey, mainWindow)
		})
		BottomDatabase = append(BottomDatabase, leveldbButton)
	}

	pluss = widget.NewButton("+", func() {
		if buttonsVisible {

			toggleButtonsContainer.Objects = nil
		} else {

			for _, m := range BottomDatabase {

				toggleButtonsContainer.Add(m)
			}
		}
		buttonsVisible = !buttonsVisible
		toggleButtonsContainer.Refresh()
	})

	mainWindow.SetCloseIntercept(func() {
		dialog.ShowConfirm("close?", "Do you want to go out?", func(confirm bool) {
			if confirm {
				mainWindow.Close()
			}
		}, mainWindow)
	})

	topLeftColumn := container.NewVBox(
		pluss,
		toggleButtonsContainer,
		spacer,
	)

	darkLight := otherUI.SetupThemeButtons(myApp)

	// all window
	containerAll := ColumnContent(rightColumnAll, columnEdit, leftColumnAll, topLeftColumn, darkLight, topRightColumn, rightColumEdit, saveEditKey, mainWindow)
	mainWindow.CenterOnScreen()
	mainWindow.SetContent(containerAll)
	mainWindow.Resize(fyne.NewSize(1200, 700))
	mainWindow.ShowAndRun()
}

func LeftColumn(leftColumnAll *fyne.Container, topLeftColumn *fyne.Container, darkLight *fyne.Container) *fyne.Container {
	lastColumnScrollable := container.NewVScroll(leftColumnAll)

	mainContent := container.NewBorder(topLeftColumn, darkLight, nil, nil, lastColumnScrollable)
	return mainContent
}

func RightColumn(rightColumnAll *fyne.Container, topRightColumn *fyne.Container, rightColumEdit *fyne.Container, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) fyne.CanvasObject {
	rightColumnScrollable := container.NewVScroll(rightColumnAll)

	up := false

	rightColumnScrollable.OnScrolled = func(p fyne.Position) {
		maxScroll := rightColumnAll.MinSize().Height - rightColumnScrollable.Size().Height

		if up && p.Y == 0 && !variable.ResultSearch {
			variable.CurrentPage--
			if variable.CurrentPage < 3 {
				up = false
				variable.CurrentPage = 3
				return
			}
			numberLast := len(rightColumnAll.Objects)
			otherUI.UpdatePage(rightColumnAll, columnEditKey, saveKey, mainWindow)

			rightColumnAll.Objects = rightColumnAll.Objects[:numberLast]

			rightColumnScrollable.Offset.Y = maxScroll / 2
			rightColumnScrollable.Refresh()

		} else if p.Y == maxScroll && !variable.ItemsAdded && !variable.ResultSearch {
			return
		} else if p.Y == maxScroll && variable.ItemsAdded && !variable.ResultSearch {

			variable.CurrentPage++
			numberLast := len(rightColumnAll.Objects)
			otherUI.UpdatePage(rightColumnAll, columnEditKey, saveKey, mainWindow)
			rightColumnScrollable.Offset.Y = maxScroll / 2

			if len(rightColumnAll.Objects) > (variable.ItemsPerPage)*3 {
				rightColumnAll.Objects = rightColumnAll.Objects[len(rightColumnAll.Objects)-numberLast:]
				up = true
			}

		}

	}
	m := container.NewVScroll(columnEditKey)
	rightColumEdit = container.NewBorder(widget.NewLabelWithStyle("Edit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), saveAndCancle, nil, nil, m)

	columns := container.NewHSplit(rightColumnScrollable, rightColumEdit)
	columns.SetOffset(0.80)
	mainContent := container.NewBorder(topRightColumn, nil, nil, nil, columns)

	return mainContent
}

func ColumnContent(rightColumnAll *fyne.Container, columnEdit *fyne.Container, leftColumnAll *fyne.Container, topLeftColumn *fyne.Container, darkLight *fyne.Container, topRightColumn *fyne.Container, rightColumEdit *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) fyne.CanvasObject {

	mainContent := LeftColumn(leftColumnAll, topLeftColumn, darkLight)

	rightColumnScrollable := RightColumn(rightColumnAll, topRightColumn, columnEdit, rightColumEdit, saveKey, mainWindow)

	columns := container.NewHSplit(mainContent, rightColumnScrollable)
	columns.SetOffset(0.10)

	container.NewScroll(columns)
	return columns
}
