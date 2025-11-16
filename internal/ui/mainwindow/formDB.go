package mainwindow

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/pref"
	"DatabaseDB/internal/utils"
	"fmt"
	"image/color"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (mi *MainWindow2) FormPasteDatabase(title string) {

	var ded *dialog.CustomDialog

	createSeparator := func() *canvas.Line {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 1
		return line
	}
	line1 := createSeparator()

	lableName := widget.NewLabel("Name :")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name")
	nameContent := container.NewBorder(nil, nil, lableName, nil, nameEntry)

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("No folder selected")

	testConnectionButton := widget.NewButton("Test Connection", func() {

		err := logic.HandleButtonClick(pathEntry.Text, title)
		if err != nil {
			dialog.ShowError(err, mi.Window)
		} else {
			dialog.ShowInformation("Success", "Test connection successful.", mi.Window)
		}
	})
	testConnectionButton.Disable()

	pathEntry.OnChanged = func(text string) {
		if variable.CreatDatabase {
			if !testConnectionButton.Disabled() {
				testConnectionButton.Disable()
			}
			return
		}
		if text != "" {
			if testConnectionButton.Disabled() {
				testConnectionButton.Enable()
			}
		} else {
			if !testConnectionButton.Disabled() {
				testConnectionButton.Disable()
			}
		}
	}
	var BoxCreateDatabase *widget.Check
	openButton := widget.NewButton("Open Folder", func() {
		var folderDialog *dialog.FileDialog
		if !BoxCreateDatabase.Checked {

			folderDialog = dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
				if err != nil {
					fmt.Println("Error opening folder:", err)
					return
				}
				if dir == nil {
					fmt.Print("No folder selected")
					return
				}
				filePath := dir.URI().Path()

				variable.FolderPath = filepath.Dir(filePath)

				if variable.NameData.FilterFile(variable.FolderPath) {
					pathEntry.SetText(variable.FolderPath)
					testConnectionButton.Enable()
				} else {
					dialog.ShowInformation("Invalid Folder", "The selected folder does not contain a valid LevelDB manifest file.", mi.Window)
				}
			}, mi.Window)
			variable.NameData.FilterFormat(folderDialog)
		} else {
			folderDialog = dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
				if err != nil {
					fmt.Println("Error opening folder:", err)
					return
				}
				if lu == nil {
					fmt.Print("No folder selected")
					return
				}
				filePath := lu.Path()

				variable.FolderPath = filePath + "/" + title + "-" + mi.TypeDB

				pathEntry.SetText(variable.FolderPath)

			}, mi.Window)
		}

		folderDialog.Show()
	})

	BoxCreateDatabase = widget.NewCheck("Create Database", func(value bool) {

		if value {
			testConnectionButton.Disable()
		} else if pathEntry.Text != "" {
			testConnectionButton.Enable()
		}
		variable.CreatDatabase = value

	})

	testOpenButton := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, openButton, testConnectionButton),
	)

	buttonOk := widget.NewButton("Add", func() {
		data := pref.Project{
			Name:        nameEntry.Text,
			FileAddress: pathEntry.Text,
			Databace:    title,
		}
		if nameEntry.Text == "" {
			dialog.ShowInformation("Error ", "Please fill in the name field", mi.Window)
			return
		}

		for _, m := range mi.Pref.ListDB {

			if nameEntry.Text == m.Name {
				dialog.ShowInformation("Error ", "Your database name is duplicate", mi.Window)
				return
			}
		}

		var addButton bool
		err := logic.HandleButtonClick(pathEntry.Text, title)
		if err == nil {

			mi.Pref.ListDB = append(mi.Pref.ListDB, data)
			addButton = false
		}

		if err != nil {
			dialog.ShowInformation("Error ", string(err.Error()), mi.Window)
		} else {
			if !addButton {

				utils.CheckCondition(mi.RightColumn.container)
				utils.CheckCondition(mi.EditColumn.edit2)

				buttonContainer := mi.ProjectButton(nameEntry.Text, mi.LeftColumn.container, pathEntry.Text)
				mi.LeftColumn.container.Add(buttonContainer)
				mi.LeftColumn.container.Refresh()

				variable.CreatDatabase = false
				ded.Hide()
				mi.RightColumn.buttonAdd.Disable()
				mi.RightColumn.searchButton.Disable()
				mi.RightColumn.buttonDelete.Disable()
			}
		}
	})
	buttonOk.Importance = widget.HighImportance

	rightColumnContent := container.NewVBox(
		layout.NewSpacer(),
		nameContent,
		layout.NewSpacer(),
		line1,
		layout.NewSpacer(),
		BoxCreateDatabase,
		layout.NewSpacer(),
		pathEntry,
		layout.NewSpacer(),
		testOpenButton,
		layout.NewSpacer(),
		buttonOk,
		layout.NewSpacer(),
	)

	ded = dialog.NewCustom("Add Key and Value", "Cancel", rightColumnContent, mi.Window)
	ded.Resize(fyne.NewSize(700, 450))
	ded.Show()
	mi.Window.Canvas().Focus(nameEntry)

}
