package mainwindow

import (
	variable "DatabaseDB"

	dbpak "DatabaseDB/internal/Databaces"
	Filterbadger "DatabaseDB/internal/filterdatabase/badger"
	FilterLeveldb "DatabaseDB/internal/filterdatabase/leveldb"
	Filterpebbledb "DatabaseDB/internal/filterdatabase/pebble"
	"DatabaseDB/internal/pref"
	"DatabaseDB/internal/ui/ids"
	"DatabaseDB/internal/ui/labelkv"
	"DatabaseDB/internal/ui/them"
	"DatabaseDB/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/widget"
)

type MainWindow2 struct {
	Window     fyne.Window
	NameWindow string
	TypeDB     string
	//DBService *service.DBService
	//Storage   *service.StorageService

	LeftColumn  *LeftColumn
	RightColumn *RightColumn
	EditColumn  *EditColumn
	Objects     *ObjectsMainWindow
	Pref        *pref.Pref
}

type ObjectsMainWindow struct {
	spacer *widget.Label
	line   *canvas.Line
}

func NewMainWindow(name string) *MainWindow2 {
	leftColumn := &LeftColumn{
		container:              container.NewVBox(),
		previousClose:          widget.NewButtonWithIcon("", theme.CancelIcon(), nil),
		previousProject:        widget.NewButton("", nil), // dinamic name of project
		previousRefreshButton:  widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), nil),
		toggleButtonsContainer: container.NewVBox(),
		darkLight:              container.NewVBox(),
		pluss:                  widget.NewButton(ids.Plass, nil),
		leveldbButton:          widget.NewButton("", nil), // dinamic name of database
		bottomDatabase:         []*widget.Button{},
	}

	rightColumn := &RightColumn{
		container:            container.NewVBox(),
		nameButtonProject:    widget.NewLabel(""), // dinamic name of project
		buttonDelete:         widget.NewButton(ids.DeleteButtonMain, nil),
		searchButton:         widget.NewButton(ids.SearchButtonMain, nil),
		buttonAdd:            widget.NewButton(ids.AddButtonMain, nil),
		keyRightColunm:       widget.NewLabelWithStyle(ids.KeyRightColunm, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		valueRightColunm:     widget.NewLabelWithStyle(ids.ValueRightColunm, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		lastLableKeyAndValue: labelkv.NewTappableLabel(""), // dinamic last label key and value
		lastStart:            &[]byte{},
		lastEnd:              &[]byte{},
		lastPage:             0,
		orgdata:              []dbpak.KVData{},
	}

	editColumn := &EditColumn{
		container:     container.NewVBox(),
		edit2:         container.NewVBox(),
		cancelEditKey: widget.NewButton(ids.CancelButtonEdit, nil),
		saveEditKey:   widget.NewButton(ids.SaveButtonEdit, nil),
		valueEntry:    widget.NewEntry(),
		editLabel:     widget.NewLabelWithStyle(ids.EditLabel, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	}

	object := &ObjectsMainWindow{
		line:   canvas.NewLine(theme.PrimaryColor()),
		spacer: widget.NewLabel(""),
	}

	mw := &MainWindow2{
		NameWindow:  name,
		TypeDB:      "", // default or placeholder DB type
		LeftColumn:  leftColumn,
		RightColumn: rightColumn,
		EditColumn:  editColumn,
		Objects:     object,
	}

	return mw
}

func (m *MainWindow2) MainWindow(myApp fyne.App) {

	m.Window = myApp.NewWindow(m.NameWindow)
	m.Window.SetMaster()

	m.Objects.spacer = widget.NewLabel("")

	// name bottom project in colunm right
	m.RightColumn.nameButtonProject = widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	m.EditColumn.saveEditKey.Disable()

	m.EditColumn.cancelEditKey.OnTapped = func() {
		utils.CheckCondition(m.EditColumn.edit2)
	}

	m.RightColumn.searchButton.OnTapped = func() {
		m.SearchKeyUi()
	}

	m.RightColumn.buttonAdd.OnTapped = func() {
		m.OpenAddDialog()
	}
	m.RightColumn.buttonAdd.Disable()
	m.RightColumn.searchButton.Disable()

	m.RightColumn.buttonDelete.OnTapped = func() {
		m.DeleteKeyUi()
	}

	buttonsVisible := false

	m.RightColumn.buttonDelete.Disable()
	// left column
	m.LeftColumn.container = m.SetupLastColumn()
	m.Objects.spacer.Resize(fyne.NewSize(0, 30))

	for _, name := range variable.NameDatabase {
		dbName := name // کپی محلی برای capture صحیح در closure

		m.LeftColumn.leveldbButton = widget.NewButton(dbName, func() {
			m.LeftColumn.toggleButtonsContainer.Objects = nil
			buttonsVisible = false
			m.TypeDB = dbName
			switch dbName {
			case "levelDB":
				variable.NameData = FilterLeveldb.NewFileterLeveldb()
			case "Pebble":
				variable.NameData = Filterpebbledb.NewFileterPebble()
			case "Badger":
				variable.NameData = Filterbadger.NewFileterBadger()
				//case "Redis":
				//	variable.NameData = Filterredis.NewFileterRedis()
			}

			m.FormPasteDatabase(dbName)
		})
		m.LeftColumn.bottomDatabase = append(m.LeftColumn.bottomDatabase, m.LeftColumn.leveldbButton)
	}

	m.LeftColumn.pluss.OnTapped = func() {
		if buttonsVisible {
			m.LeftColumn.toggleButtonsContainer.Objects = nil
		} else {
			for _, m2 := range m.LeftColumn.bottomDatabase {
				m.LeftColumn.toggleButtonsContainer.Add(m2)
			}
		}
		buttonsVisible = !buttonsVisible
		m.LeftColumn.toggleButtonsContainer.Refresh()
	}

	m.Window.SetCloseIntercept(func() {
		dialog.ShowConfirm("close?", "Do you want to go out?", func(confirm bool) {
			if confirm {
				m.Pref.SaveDatabase(m.Pref.ListDB, pref.KeyListDB)

				keyTheme := them.GetThemeKey(myApp)
				m.Pref.SaveTheme(keyTheme, pref.KeyTheme)
				m.Window.Close()
			}
		}, m.Window)
	})

	m.LeftColumn.darkLight = m.SetupThemeButtons(myApp)

	// all window
	containerAll := m.ColumnContent()
	m.Window.CenterOnScreen()
	m.Window.SetContent(containerAll)
	m.Window.Resize(fyne.NewSize(1200, 700))
	m.Window.ShowAndRun()
}

func (m *MainWindow2) LeftColumn2() fyne.CanvasObject {
	lastColumnScrollable := container.NewVScroll(m.LeftColumn.container)

	mainContent := container.NewBorder(m.TopLeftColumn2(), m.LeftColumn.darkLight, nil, nil, lastColumnScrollable)
	return mainContent
}

func (mi *MainWindow2) RightColumn2() fyne.CanvasObject {
	if mi.RightColumn.container == nil {
		mi.RightColumn.container = container.NewVBox()
	}
	rightColumnScrollable := container.NewVScroll(mi.RightColumn.container)

	up := false

	rightColumnScrollable.OnScrolled = func(p fyne.Position) {
		maxScroll := mi.RightColumn.container.MinSize().Height - rightColumnScrollable.Size().Height

		if up && p.Y == 0 && !variable.GetResultSearch() {
			newPage := variable.DecrementCurrentPage()
			if newPage < 3 {
				up = false
				variable.SetCurrentPage(3)
				return
			}
			numberLast := len(mi.RightColumn.container.Objects)
			mi.UpdatePage()

			mi.RightColumn.container.Objects = mi.RightColumn.container.Objects[:numberLast]

			rightColumnScrollable.Offset.Y = maxScroll / 2
			rightColumnScrollable.Refresh()

		} else if p.Y == maxScroll && !variable.GetItemsAdded() && !variable.GetResultSearch() {
			return
		} else if p.Y == maxScroll && variable.GetItemsAdded() && !variable.GetResultSearch() {

			variable.IncrementCurrentPage()
			mi.UpdatePage()
			rightColumnScrollable.Offset.Y = maxScroll / 2

			if len(mi.RightColumn.container.Objects) == (variable.ItemsPerPage)*3 {
				up = true
			}

		}

	}

	m := container.NewVScroll(mi.EditColumn.edit2)
	mi.EditColumn.container = container.NewBorder(
		mi.EditColumn.editLabel,
		mi.SaveAndCancle(),
		nil, nil, m,
	)
	mi.EditColumn.container.Refresh()

	columns := container.NewHSplit(rightColumnScrollable, mi.EditColumn.container)
	columns.SetOffset(0.80)
	mainContent := container.NewBorder(mi.TopRightColumn(), nil, nil, nil, columns)

	return mainContent
}

func (m *MainWindow2) ColumnContent() fyne.CanvasObject {

	mainContent := m.LeftColumn2()

	rightColumnScrollable := m.RightColumn2()

	columns := container.NewHSplit(mainContent, rightColumnScrollable)
	columns.SetOffset(0.10)

	return columns
}
