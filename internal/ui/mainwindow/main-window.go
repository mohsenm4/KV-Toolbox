package mainwindow

import (
	variable "DatabaseDB"
	"fmt"

	dbpak "DatabaseDB/internal/Databaces"
	Filterbadger "DatabaseDB/internal/filterdatabase/badger"
	FilterLeveldb "DatabaseDB/internal/filterdatabase/leveldb"
	Filterpebbledb "DatabaseDB/internal/filterdatabase/pebble"
	"DatabaseDB/internal/pref"
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

	LeftColumn  *LeftColumn2
	RightColumn *RightColumn2
	EditColumn  *EditColumn2
	Objects     *ObjectsMainWindow
	pref        *pref.Pref
}

type ObjectsMainWindow struct {
	Spacer *widget.Label
	Line   *canvas.Line
}

func NewMainWindow(name string) *MainWindow2 {
	mw := &MainWindow2{
		NameWindow: name,
		TypeDB:     "", // default or placeholder DB type
		LeftColumn: &LeftColumn2{
			Container:              container.NewVBox(),
			PreviousClose:          widget.NewButton("", nil),
			PreviousProject:        widget.NewButton("", nil),
			PreviousRefreshButton:  widget.NewButton("", nil),
			ToggleButtonsContainer: container.NewVBox(),
			DarkLight:              container.NewVBox(),
			Pluss:                  widget.NewButton("", nil),
			LeveldbButton:          widget.NewButton("", nil),
			BottomDatabase:         []*widget.Button{},
		},
		RightColumn: &RightColumn2{
			Container:            container.NewVBox(),
			NameButtonProject:    widget.NewLabel(""),
			Spacer:               widget.NewLabel(""),
			ButtonDelete:         widget.NewButton("", nil),
			SearchButton:         widget.NewButton("", nil),
			ButtonAdd:            widget.NewButton("", nil),
			KeyRightColunm:       widget.NewLabel(""),
			ValueRightColunm:     widget.NewLabel(""),
			LastLableKeyAndValue: utils.NewTappableLabel("", nil),
			LastStart:            &[]byte{},
			LastEnd:              &[]byte{},
			LastPage:             0,
			Orgdata:              []dbpak.KVData{},
		},
		EditColumn: &EditColumn2{
			Container:     container.NewVBox(),
			Edit2:         container.NewVBox(),
			CancelEditKey: widget.NewButton("", nil),
			SaveEditKey:   widget.NewButton("", nil),
			ValueEntry:    widget.NewEntry(),
		},
		Objects: &ObjectsMainWindow{
			Line:   canvas.NewLine(theme.PrimaryColor()),
			Spacer: widget.NewLabel(""),
		},
	}

	return mw
}

func (m *MainWindow2) MainWindow(myApp fyne.App) {

	m.Window = myApp.NewWindow(m.NameWindow)
	m.Window.SetMaster()
	m.pref = pref.NewPref(myApp)

	m.Objects.Spacer = widget.NewLabel("")

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

	m.EditColumn.SaveEditKey = widget.NewButton("Save", func() {})
	m.EditColumn.SaveEditKey.Disable()

	m.EditColumn.CancelEditKey = widget.NewButton("Cancle", func() {
		utils.CheckCondition(m.EditColumn.Edit2)
	})

	m.RightColumn.SearchButton = widget.NewButton("Search", func() {
		m.SearchKeyUi()
	})

	m.EditColumn.Container = container.NewBorder(widget.NewLabelWithStyle("Edit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), m.SaveAndCancle(), nil, nil, m.EditColumn.Edit2)

	m.RightColumn.ButtonAdd = widget.NewButton("Add", func() {
		m.OpenAddDialog()
	})
	m.RightColumn.ButtonAdd.Disable()
	m.RightColumn.SearchButton.Disable()

	m.RightColumn.ButtonDelete = widget.NewButton("Delete", func() {
		m.DeleteKeyUi()
	})

	buttonsVisible := false

	m.RightColumn.ButtonDelete.Disable()
	// left column
	m.LeftColumn.Container = m.SetupLastColumn()
	m.Objects.Spacer.Resize(fyne.NewSize(0, 30))

	for _, name := range variable.NameDatabase {

		m.LeftColumn.LeveldbButton = widget.NewButton(name, func() {
			m.LeftColumn.ToggleButtonsContainer.Objects = nil
			buttonsVisible = false
			m.TypeDB = name
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

			m.FormPasteDatabase(name)
		})
		m.LeftColumn.BottomDatabase = append(m.LeftColumn.BottomDatabase, m.LeftColumn.LeveldbButton)
	}

	m.LeftColumn.Pluss = widget.NewButton("+", func() {
		if buttonsVisible {

			m.LeftColumn.ToggleButtonsContainer.Objects = nil
		} else {

			for _, m2 := range m.LeftColumn.BottomDatabase {

				m.LeftColumn.ToggleButtonsContainer.Add(m2)
			}
		}
		buttonsVisible = !buttonsVisible
		m.LeftColumn.ToggleButtonsContainer.Refresh()
	})

	m.Window.SetCloseIntercept(func() {
		dialog.ShowConfirm("close?", "Do you want to go out?", func(confirm bool) {
			if confirm {
				m.pref.SaveDatabase(m.pref.ListDB, pref.KeyListDB)

				m.Window.Close()
			}
		}, m.Window)
	})

	m.LeftColumn.DarkLight = m.SetupThemeButtons(myApp)

	// all window
	containerAll := m.ColumnContent()
	m.Window.CenterOnScreen()
	m.Window.SetContent(containerAll)
	m.Window.Resize(fyne.NewSize(1200, 700))
	m.Window.ShowAndRun()
}

func (m *MainWindow2) LeftColumn2() fyne.CanvasObject {
	lastColumnScrollable := container.NewVScroll(m.LeftColumn.Container)

	mainContent := container.NewBorder(m.TopLeftColumn2(), m.LeftColumn.DarkLight, nil, nil, lastColumnScrollable)
	return mainContent
}

func (mi *MainWindow2) RightColumn2() fyne.CanvasObject {
	if mi.RightColumn.Container == nil {
		mi.RightColumn.Container = container.NewVBox()
	}
	if mi.TopRightColumn() == nil {
		fmt.Println("")
	}
	rightColumnScrollable := container.NewVScroll(mi.RightColumn.Container)

	up := false

	rightColumnScrollable.OnScrolled = func(p fyne.Position) {
		maxScroll := mi.RightColumn.Container.MinSize().Height - rightColumnScrollable.Size().Height

		if up && p.Y == 0 && !variable.ResultSearch {
			variable.CurrentPage--
			if variable.CurrentPage < 3 {
				up = false
				variable.CurrentPage = 3
				return
			}
			numberLast := len(mi.RightColumn.Container.Objects)
			mi.UpdatePage()

			mi.RightColumn.Container.Objects = mi.RightColumn.Container.Objects[:numberLast]

			rightColumnScrollable.Offset.Y = maxScroll / 2
			rightColumnScrollable.Refresh()

		} else if p.Y == maxScroll && !variable.ItemsAdded && !variable.ResultSearch {
			return
		} else if p.Y == maxScroll && variable.ItemsAdded && !variable.ResultSearch {

			variable.CurrentPage++
			numberLast := len(mi.RightColumn.Container.Objects)
			mi.UpdatePage()
			rightColumnScrollable.Offset.Y = maxScroll / 2

			if len(mi.RightColumn.Container.Objects) > (variable.ItemsPerPage)*3 {
				mi.RightColumn.Container.Objects = mi.RightColumn.Container.Objects[len(mi.RightColumn.Container.Objects)-numberLast:]
				up = true
			}

		}

	}

	m := container.NewVScroll(mi.EditColumn.Edit2)
	mi.EditColumn.Container = container.NewBorder(
		widget.NewLabelWithStyle("Edit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		mi.SaveAndCancle(),
		nil, nil, m,
	)
	mi.EditColumn.Container.Refresh()

	columns := container.NewHSplit(rightColumnScrollable, mi.EditColumn.Container)
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
