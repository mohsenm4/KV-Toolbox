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
	"github.com/gabriel-vasile/mimetype"
)

const (
	ThemeDark   = "dark"
	ThemeLight  = "light"
	ThemeCustom = "custom"
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
	All         []dbpak.KVData
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
		pluss:                  widget.NewButton("+", nil),
		leveldbButton:          widget.NewButton("", nil), // dinamic name of database
		bottomDatabase:         []*widget.Button{},
	}

	rightColumn := &RightColumn{
		container:            container.NewVBox(),
		nameButtonProject:    widget.NewLabel(""), // dinamic name of project
		buttonDelete:         widget.NewButton("Delete", nil),
		searchButton:         widget.NewButton("Search", nil),
		buttonAdd:            widget.NewButton("Add", nil),
		keyRightColunm:       widget.NewLabelWithStyle("key", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		valueRightColunm:     widget.NewLabelWithStyle("value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		lastLableKeyAndValue: utils.NewTappableLabel("", nil), // dinamic last label key and value
		list:                 widget.NewList(nil, nil, nil),
	}

	editColumn := &EditColumn{
		container:     container.NewVBox(),
		edit2:         container.NewVBox(),
		cancelEditKey: widget.NewButton("Cancel", nil),
		saveEditKey:   widget.NewButton("Save", nil),
		valueEntry:    widget.NewEntry(),
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

	mytheme := m.Pref.LoadTheme(pref.KeyTheme)

	if mytheme == ThemeDark {
		myApp.Settings().SetTheme(theme.DarkTheme())
	} else if mytheme == ThemeLight {
		myApp.Settings().SetTheme(theme.LightTheme())
	}

	m.Objects.spacer = widget.NewLabel("")

	// key top window for colunm keys
	m.RightColumn.keyRightColunm = widget.NewLabelWithStyle("key", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// value top window for colunm values
	m.RightColumn.valueRightColunm = widget.NewLabelWithStyle("value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// name bottom project in colunm right
	m.RightColumn.nameButtonProject = widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	m.EditColumn.saveEditKey = widget.NewButton("Save", func() {})
	m.EditColumn.saveEditKey.Disable()

	m.EditColumn.cancelEditKey = widget.NewButton("Cancle", func() {
		utils.CheckCondition(m.EditColumn.edit2)
	})

	m.RightColumn.searchButton = widget.NewButton("Search", func() {
		m.SearchKeyUi()
	})

	m.EditColumn.container = container.NewBorder(widget.NewLabelWithStyle("Edit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), m.SaveAndCancle(), nil, nil, m.EditColumn.edit2)

	m.RightColumn.buttonAdd = widget.NewButton("Add", func() {
		m.OpenAddDialog()
	})
	m.RightColumn.buttonAdd.Disable()
	m.RightColumn.searchButton.Disable()

	m.RightColumn.buttonDelete = widget.NewButton("Delete", func() {
		m.DeleteKeyUi()
	})

	buttonsVisible := false

	m.RightColumn.buttonDelete.Disable()
	// left column
	m.LeftColumn.container = m.SetupLastColumn()
	m.Objects.spacer.Resize(fyne.NewSize(0, 30))

	for _, name := range variable.NameDatabase {

		m.LeftColumn.leveldbButton = widget.NewButton(name, func() {
			m.LeftColumn.toggleButtonsContainer.Objects = nil
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
		m.LeftColumn.bottomDatabase = append(m.LeftColumn.bottomDatabase, m.LeftColumn.leveldbButton)
	}

	m.LeftColumn.pluss = widget.NewButton("+", func() {
		if buttonsVisible {

			m.LeftColumn.toggleButtonsContainer.Objects = nil
		} else {

			for _, m2 := range m.LeftColumn.bottomDatabase {

				m.LeftColumn.toggleButtonsContainer.Add(m2)
			}
		}
		buttonsVisible = !buttonsVisible
		m.LeftColumn.toggleButtonsContainer.Refresh()
	})

	m.Window.SetCloseIntercept(func() {
		dialog.ShowConfirm("close?", "Do you want to go out?", func(confirm bool) {
			if confirm {
				m.Pref.SaveDatabase(m.Pref.ListDB, pref.KeyListDB)

				keyTheme := getThemeKey(myApp)
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
	// Create the "Edit" section (bottom-right panel)
	m := container.NewVScroll(mi.EditColumn.edit2)
	mi.EditColumn.container = container.NewBorder(
		widget.NewLabelWithStyle("Edit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		mi.SaveAndCancle(),
		nil, nil, m,
	)

	// Create an empty container for the right list (initially blank)
	mi.RightColumn.container = container.NewMax()

	// In the future, you can initialize it with default data, e.g.:
	// mi.UpdateRightList(defaultItems)

	// Combine the list and edit sections into a horizontal split
	columns := container.NewHSplit(mi.RightColumn.container, mi.EditColumn.container)
	columns.SetOffset(0.80)

	// Wrap everything into the main content container
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

func getThemeKey(app fyne.App) string {
	t := app.Settings().Theme()
	currentBG := t.Color(theme.ColorNameBackground, app.Settings().ThemeVariant())
	darkBG := theme.DarkTheme().Color(theme.ColorNameBackground, app.Settings().ThemeVariant())
	lightBG := theme.LightTheme().Color(theme.ColorNameBackground, app.Settings().ThemeVariant())

	switch {
	case currentBG == darkBG:
		return ThemeDark
	case currentBG == lightBG:
		return ThemeLight
	default:
		return ThemeCustom
	}
}

func (mi *MainWindow2) UpdateRightList(all []dbpak.KVData) {

	newList := widget.NewList(
		func() int {
			return len(all)
		},
		func() fyne.CanvasObject {
			keyLabel := widget.NewLabel("key")
			valueLabel := widget.NewLabel("value")
			buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
			return buttonRow
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {

			item := all[i]

			typeValue := mimetype.Detect(item.Value)
			var truncatedValue string
			if typeValue.Extension() != ".txt" {
				truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
			} else {
				truncatedValue = utils.TruncateString(string(item.Value), 20)
			}
			truncatedKey := utils.TruncateString(string(item.Key), 20)

			keyLabel := mi.BuildLabelKeyAndValue("key", item.Key, item.Value, truncatedKey)
			valueLabel := mi.BuildLabelKeyAndValue("value", item.Key, item.Value, truncatedValue)

			row := obj.(*fyne.Container)
			row.Objects[0] = keyLabel
			row.Objects[1] = valueLabel

		},
	)

	mi.RightColumn.container.Objects = nil
	mi.RightColumn.container.Add(newList)
	mi.RightColumn.container.Refresh()

	mi.RightColumn.list = newList
}
