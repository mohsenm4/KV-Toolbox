package mainwindow

import (
	variable "DatabaseDB"

	Filterbadger "DatabaseDB/internal/filterdatabase/badger"
	FilterLeveldb "DatabaseDB/internal/filterdatabase/leveldb"
	Filterpebbledb "DatabaseDB/internal/filterdatabase/pebble"
	addkeyui "DatabaseDB/internal/ui/addKeyui"
	deletkeyui "DatabaseDB/internal/ui/deletKeyUi"
	"DatabaseDB/internal/ui/otherUI"
	"DatabaseDB/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type MainWindow2 struct {
	Window     fyne.Window
	NameWindow string
	//DBService *service.DBService
	//Storage   *service.StorageService

	LeftColumn  *LeftColumn2
	RightColumn *RightColumn2
	EditColumn  *EditColumn2
	Objects     ObjectsMainWindow
}

type ObjectsMainWindow struct {
	Spacer *widget.Label
}

var saveAndCancle *fyne.Container

var leveldbButton *widget.Button
var BottomDatabase []*widget.Button

func NewMainWindow() *MainWindow2 {
	return &MainWindow2{}
}

func (m *MainWindow2) MainWindow(myApp fyne.App) {

	m.Window = myApp.NewWindow(m.NameWindow)
	m.Window.SetMaster()

	m.Objects.Spacer = widget.NewLabel("")

	// right column show key
	m.RightColumn.Container = container.NewVBox()

	// key top window for colunm keys
	m.RightColumn.KeyRightColunm = widget.NewLabelWithStyle("key", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// value top window for colunm values
	m.RightColumn.ValueRightColunm = widget.NewLabelWithStyle("value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// name bottom project in colunm right
	m.RightColumn.NameButtonProject = widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	m.EditColumn.SaveEditKey = widget.NewButton("Save", nil)
	m.EditColumn.SaveEditKey.Disable()

	m.EditColumn.CancelEditKey = widget.NewButton("Cancle", func() {
		utils.CheckCondition(m.EditColumn.Container)
	})

	m.EditColumn.Container = container.NewBorder(widget.NewLabelWithStyle("Edit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), m.EditColumn.SaveAndCancle(), nil, nil, m.EditColumn.Container)

	m.RightColumn.SearchButton = widget.NewButton("Search", func() {
		m.SearchKeyUi()
	})

	buttonAdd := widget.NewButton("Add", func() {
		addkeyui.OpenWindowAddButton(m.RightColumn.Container, m.Window)
	})
	buttonAdd.Disable()
	m.RightColumn.SearchButton.Disable()

	buttonDelete := widget.NewButton("Delete", func() {
		deletkeyui.DeleteKeyUi(m.RightColumn.Container, m.Window)
	})

	var pluss *widget.Button
	toggleButtonsContainer := container.NewVBox()
	buttonsVisible := false

	buttonDelete.Disable()
	// left column
	leftColumnAll := otherUI.SetupLastColumn(m.RightColumn.Container, m.RightColumn.NameButtonProject, buttonAdd, m.RightColumn.SearchButton, buttonDelete, m.EditColumn.Container, m.EditColumn.SaveEditKey, m.Window)
	m.Objects.Spacer.Resize(fyne.NewSize(0, 30))

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
			variable.NameData.FormCreate(myApp, name, leftColumnAll, m.RightColumn.Container, m.RightColumn.NameButtonProject, buttonAdd, m.RightColumn.SearchButton, buttonDelete, m.EditColumn.Container, m.EditColumn.SaveEditKey, m.Window)
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

	m.Window.SetCloseIntercept(func() {
		dialog.ShowConfirm("close?", "Do you want to go out?", func(confirm bool) {
			if confirm {
				m.Window.Close()
			}
		}, m.Window)
	})

	topLeftColumn := container.NewVBox(
		pluss,
		toggleButtonsContainer,
		m.Objects.Spacer,
	)

	darkLight := otherUI.SetupThemeButtons(myApp)

	// all window
	containerAll := ColumnContent(m.RightColumn.Container, m.EditColumn.Container, leftColumnAll, topLeftColumn, darkLight, m.RightColumn.TopRightColumn(), m.EditColumn.Container, m.EditColumn.SaveEditKey, m.Window)
	m.Window.CenterOnScreen()
	m.Window.SetContent(containerAll)
	m.Window.Resize(fyne.NewSize(1200, 700))
	m.Window.ShowAndRun()
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
