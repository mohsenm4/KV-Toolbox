package mainwindow

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/pref"
	"DatabaseDB/internal/utils"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LeftColumn2 struct {
	Container              *fyne.Container
	PreviousClose          *widget.Button
	PreviousProject        *widget.Button
	PreviousRefreshButton  *widget.Button
	ToggleButtonsContainer *fyne.Container
	DarkLight              *fyne.Container
	Pluss                  *widget.Button
	LeveldbButton          *widget.Button
	BottomDatabase         []*widget.Button
}

func (l *MainWindow2) ProjectButton(inputText string, lastColumnContent *fyne.Container, path string) *fyne.Container {
	var refreshButton *widget.Button
	var projectButton *widget.Button
	var closeButton *widget.Button

	projectButton = widget.NewButton(inputText+" - "+l.TypeDB, func() {

		typeDB := strings.Split(projectButton.Text, " - ")
		l.TypeDB = typeDB[len(typeDB)-1]

		if l.LeftColumn.PreviousProject != nil {
			l.LeftColumn.PreviousProject.Importance = widget.MediumImportance
			l.LeftColumn.PreviousClose.Importance = widget.MediumImportance
			l.LeftColumn.PreviousRefreshButton.Importance = widget.MediumImportance
		}
		projectButton.Importance = widget.HighImportance
		closeButton.Importance = widget.HighImportance
		refreshButton.Importance = widget.HighImportance

		l.LeftColumn.PreviousProject = projectButton
		l.LeftColumn.PreviousClose = closeButton
		l.LeftColumn.PreviousRefreshButton = refreshButton

		variable.ItemsAdded = true
		utils.Checkdatabace(path, l.TypeDB)
		l.RightColumn.ButtonAdd.Enable()
		l.RightColumn.SearchButton.Enable()
		l.RightColumn.ButtonDelete.Enable()
		variable.FolderPath = path
		l.RightColumn.LastEnd = nil
		variable.ResultSearch = false
		variable.CurrentPage = 1
		l.RightColumn.LastPage = 0
		variable.PreviousOffsetY = 0
		l.RightColumn.LastStart = nil
		utils.CheckCondition(l.RightColumn.Container)
		utils.CheckCondition(l.EditColumn.Edit2)
		l.UpdatePage()
		l.RightColumn.NameButtonProject.Text = ""
		l.RightColumn.NameButtonProject.Text = inputText + " - " + l.TypeDB

		l.RightColumn.Container.Refresh()
		l.LeftColumn.Container.Refresh()
		l.EditColumn.Container.Refresh()

	})
	buttonContainer := container.NewHBox()

	closeButton = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {

		if l.RightColumn.NameButtonProject.Text == inputText+" - "+l.TypeDB {
			utils.CheckCondition(l.RightColumn.Container)
			utils.CheckCondition(l.EditColumn.Edit2)

			l.RightColumn.ButtonAdd.Disable()
			l.RightColumn.SearchButton.Disable()
			l.RightColumn.ButtonDelete.Disable()

			l.RightColumn.NameButtonProject.Text = ""
			l.RightColumn.NameButtonProject.Refresh()
		}

		for i, r := range l.Pref.ListDB {
			if r.FileAddress == path {
				l.Pref.ListDB = append(l.Pref.ListDB[:i], l.Pref.ListDB[i+1:]...)
				lastColumnContent.Remove(buttonContainer)
				lastColumnContent.Refresh()
			}
		}
	})

	refreshButton = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {

		if l.RightColumn.NameButtonProject.Text == inputText+" - "+l.TypeDB {

			variable.ItemsAdded = true
			utils.Checkdatabace(path, l.TypeDB)
			l.RightColumn.ButtonAdd.Enable()
			l.RightColumn.SearchButton.Enable()
			l.RightColumn.ButtonDelete.Enable()
			variable.FolderPath = path
			l.RightColumn.LastEnd = nil
			variable.ResultSearch = false
			variable.CurrentPage = 1
			l.RightColumn.LastPage = 0
			variable.PreviousOffsetY = 0
			l.RightColumn.LastStart = nil
			utils.CheckCondition(l.RightColumn.Container)
			utils.CheckCondition(l.EditColumn.Edit2)
			l.UpdatePage()

			l.RightColumn.NameButtonProject.Refresh()
		}

	})
	refreshClose := container.NewGridWithColumns(2, refreshButton, closeButton)

	buttonContainer = container.NewBorder(nil, nil, nil, refreshClose, projectButton)
	return buttonContainer
}

func (l *MainWindow2) SetupLastColumn() *fyne.Container {
	lastColumnContent := container.NewVBox()

	dataJson, err := l.Pref.LoadDatabase(pref.KeyListDB)
	if err != nil {
		log.Fatal("Error loading JSON data:", err)
	} else {
		for _, project := range dataJson {
			l.TypeDB = project.Databace
			buttonContainer := l.ProjectButton(project.Name, lastColumnContent, project.FileAddress)
			lastColumnContent.Add(buttonContainer)
		}
	}
	l.Pref.ListDB = dataJson
	return lastColumnContent
}

func (r *MainWindow2) TopLeftColumn2() *fyne.Container {
	r.Objects.Spacer = widget.NewLabel("")
	topLeftColumn := container.NewVBox(
		r.LeftColumn.Pluss,
		r.LeftColumn.ToggleButtonsContainer,
		r.Objects.Spacer,
	)
	return topLeftColumn
}

func (l *MainWindow2) SetupThemeButtons(app fyne.App) *fyne.Container {
	var darkButton *widget.Button
	var lightButton *widget.Button

	darkButton = widget.NewButton("Dark", func() {

		app.Settings().SetTheme(theme.DarkTheme())

		darkButton.Importance = widget.HighImportance
		lightButton.Importance = widget.MediumImportance
		darkButton.Refresh()
		lightButton.Refresh()
	})
	lightButton = widget.NewButton("Light", func() {
		app.Settings().SetTheme(theme.LightTheme())

		lightButton.Importance = widget.HighImportance
		darkButton.Importance = widget.MediumImportance

		darkButton.Refresh()
		lightButton.Refresh()
	})

	setThemeButtonImportance(app, darkButton, lightButton)

	darkLight := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, lightButton, darkButton),
	)
	return darkLight
}

func setThemeButtonImportance(app fyne.App, darkButton, lightButton *widget.Button) {
	t := app.Settings().Theme()
	currentBG := t.Color(theme.ColorNameBackground, app.Settings().ThemeVariant())
	darkBG := theme.DarkTheme().Color(theme.ColorNameBackground, app.Settings().ThemeVariant())
	lightBG := theme.LightTheme().Color(theme.ColorNameBackground, app.Settings().ThemeVariant())

	switch {
	case currentBG == darkBG:
		darkButton.Importance = widget.HighImportance
		lightButton.Importance = widget.MediumImportance
	case currentBG == lightBG:
		lightButton.Importance = widget.HighImportance
		darkButton.Importance = widget.MediumImportance
	default:
		darkButton.Importance = widget.MediumImportance
		lightButton.Importance = widget.MediumImportance
	}

	darkButton.Refresh()
	lightButton.Refresh()
}
