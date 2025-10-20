package sharedfunc

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/ui/otherUI"
	"DatabaseDB/internal/utils"
	"fmt"
	"image/color"
	"io/ioutil"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func FormPasteDatabase(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, buttonSearch *widget.Button, buttonDelete *widget.Button, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) {

	var ded *dialog.CustomDialog

	createSeparator := func() *canvas.Line {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 1
		return line
	}
	line1 := createSeparator()

	lableName := widget.NewLabel("Name :")
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Name"
	nameContent := container.NewBorder(nil, nil, lableName, nil, nameEntry)

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("No folder selected")

	testConnectionButton := widget.NewButton("Test Connection", func() {

		err := logic.HandleButtonClick(pathEntry.Text, title)
		if err != nil {
			dialog.ShowError(err, mainWindow)
		} else {
			dialog.ShowInformation("Success", "Test connection successful.", mainWindow)
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
					dialog.ShowInformation("Invalid Folder", "The selected folder does not contain a valid LevelDB manifest file.", mainWindow)
				}
			}, mainWindow)
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

				variable.FolderPath = filePath + "/" + title + "-" + nameEntry.Text

				pathEntry.SetText(variable.FolderPath)

			}, mainWindow)
		}

		folderDialog.Show()
	})

	BoxCreateDatabase = widget.NewCheck("Create Database", func(value bool) {

		variable.CreatDatabase = value

	})

	testOpenButton := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, openButton, testConnectionButton),
	)

	buttonOk := widget.NewButton("Add", func() {
		data := map[string]string{
			"Name":     nameEntry.Text,
			"Addres":   pathEntry.Text,
			"Database": title,
		}
		if nameEntry.Text == "" {
			dialog.ShowInformation("Error ", "Please fill in the name field", mainWindow)
			return
		}
		datajson, err := variable.CurrentJson.Load()
		if err != nil {
			fmt.Println("Error opening folder:", err)
		}
		for _, m := range datajson.RecentProjects {
			if nameEntry.Text == m.Name {
				dialog.ShowInformation("Error ", "Your database name is duplicate", mainWindow)
				return
			}
		}

		var addButton bool
		err = logic.HandleButtonClick(pathEntry.Text, title)
		if err == nil {

			err, addButton = variable.CurrentJson.Add(data)
			if err != nil {
				dialog.ShowInformation("error", err.Error(), mainWindow)
				return
			}
		}

		if err != nil {
			dialog.ShowInformation("Error ", string(err.Error()), mainWindow)
		} else {
			if !addButton {

				utils.CheckCondition(rightColumnContentORG)
				utils.CheckCondition(columnEditKey)

				buttonContainer := otherUI.ProjectButton(nameEntry.Text, lastColumnContent, pathEntry.Text, rightColumnContentORG, nameButtonProject, buttonAdd, buttonSearch, buttonDelete, title, columnEditKey, saveKey, mainWindow)
				lastColumnContent.Add(buttonContainer)
				lastColumnContent.Refresh()

				variable.CreatDatabase = false
				ded.Hide()
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

	ded = dialog.NewCustom("Add Key and Value", "Cancel", rightColumnContent, mainWindow)
	ded.Resize(fyne.NewSize(700, 450))
	ded.Show()
}

func FormatFilesDatabase(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error opening folder:", err)
		return false
	}
	var count uint8
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "MANIFEST-") || filepath.Ext(file.Name()) == ".log" {
			count++
		}

		if count == 2 {
			return true
		}
	}
	return false
}
