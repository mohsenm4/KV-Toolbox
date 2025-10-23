package mainwindow

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/utils"
	"fmt"
	"log"

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
	LastLableKeyAndValue   *utils.TappableLabel
	ToggleButtonsContainer *fyne.Container
	DarkLight              *fyne.Container
	Pluss                  *widget.Button
	LeveldbButton          *widget.Button
	BottomDatabase         []*widget.Button
}

func (l *MainWindow2) ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, nameDatabace string) *fyne.Container {
	var refreshButton *widget.Button
	var projectButton *widget.Button
	var closeButton *widget.Button

	projectButton = widget.NewButton(inputText+" - "+nameDatabace, func() {
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
		utils.Checkdatabace(path, nameDatabace)
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
		utils.CheckCondition(l.EditColumn.Container)
		l.UpdatePage()
		l.RightColumn.NameButtonProject.Text = ""
		l.RightColumn.NameButtonProject.Text = inputText + " - " + nameDatabace

		l.RightColumn.Container.Refresh()
		l.LeftColumn.Container.Refresh()
		l.EditColumn.Container.Refresh()

	})
	buttonContainer := container.NewHBox()

	closeButton = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {

		if l.RightColumn.NameButtonProject.Text == inputText+" - "+nameDatabace {
			utils.CheckCondition(l.RightColumn.Container)
			utils.CheckCondition(l.EditColumn.Container)

			l.RightColumn.ButtonAdd.Disable()
			l.RightColumn.SearchButton.Disable()
			l.RightColumn.ButtonDelete.Disable()

			l.RightColumn.NameButtonProject.Text = ""
			l.RightColumn.NameButtonProject.Refresh()
		}

		err := variable.CurrentJson.Remove(inputText)
		if err != nil {
			fmt.Print(err)
		} else {

			lastColumnContent.Remove(buttonContainer)
			lastColumnContent.Refresh()
		}
	})

	refreshButton = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {

		if l.RightColumn.NameButtonProject.Text == inputText+" - "+nameDatabace {

			variable.ItemsAdded = true
			utils.Checkdatabace(path, nameDatabace)
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
			utils.CheckCondition(l.EditColumn.Container)
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

	jsonDataa, err := variable.CurrentJson.Load()
	if err != nil {
		log.Fatal("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {

			buttonContainer := l.ProjectButton(project.Name, lastColumnContent, project.FileAddress, project.Databace)
			lastColumnContent.Add(buttonContainer)
		}
	}

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

	darkButton.Importance = widget.HighImportance

	darkLight := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, lightButton, darkButton),
	)
	return darkLight
}
