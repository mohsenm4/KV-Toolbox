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
	newWindow := a.NewWindow(title)

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

	lableComment := widget.NewLabel("Comment :")
	commentEntry := widget.NewEntry()
	commentEntry.PlaceHolder = "Comment"
	commentContent := container.NewBorder(nil, nil, lableComment, nil, commentEntry)

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("No folder selected")

	testConnectionButton := widget.NewButton("Test Connection", func() {

		err := logic.HandleButtonClick(pathEntry.Text, title)
		if err != nil {
			dialog.ShowError(err, newWindow)
		} else {
			dialog.ShowInformation("Success", "Test connection successful.", newWindow)
		}
	})
	testConnectionButton.Disable()

	pathEntry.OnChanged = func(text string) {
		if text != "" && !variable.CreatDatabase {
			testConnectionButton.Enable()
		} else if variable.CreatDatabase {
			testConnectionButton.Disable()
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
					dialog.ShowInformation("Invalid Folder", "The selected folder does not contain a valid LevelDB manifest file.", newWindow)
				}
			}, newWindow)
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

			}, newWindow)
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

	buttonCancel := widget.NewButton("Cancel", func() {
		newWindow.Close()
	})

	buttonOk := widget.NewButton("Add", func() {
		data := map[string]string{
			"Name":     nameEntry.Text,
			"Comment":  commentEntry.Text,
			"Addres":   pathEntry.Text,
			"Database": title,
		}
		if nameEntry.Text == "" {
			dialog.ShowInformation("Error ", "Please fill in the name field", newWindow)
			return
		}
		datajson, err := variable.CurrentJson.Load()
		if err != nil {
			fmt.Println("Error opening folder:", err)
		}
		for _, m := range datajson.RecentProjects {
			if nameEntry.Text == m.Name {
				dialog.ShowInformation("Error ", "Your database name is duplicate", newWindow)
				return
			}
		}

		var addButton bool
		err = logic.HandleButtonClick(pathEntry.Text, title)
		if err == nil {

			err, addButton = variable.CurrentJson.Add(data)
			if err != nil {
				dialog.ShowInformation("error", err.Error(), newWindow)

			}
		}

		if err != nil {
			dialog.ShowInformation("Error ", string(err.Error()), newWindow)
		} else {
			if !addButton {

				utils.CheckCondition(rightColumnContentORG)
				utils.CheckCondition(columnEditKey)

				buttonContainer := otherUI.ProjectButton(nameEntry.Text, lastColumnContent, pathEntry.Text, rightColumnContentORG, nameButtonProject, buttonAdd, buttonSearch, buttonDelete, title, columnEditKey, saveKey, mainWindow)
				lastColumnContent.Add(buttonContainer)
				lastColumnContent.Refresh()

				variable.CreatDatabase = false
				newWindow.Close()
			}
		}
	})
	buttonOk.Importance = widget.HighImportance

	rowBotton := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, buttonCancel, buttonOk),
	)

	rightColumnContent := container.NewVBox(
		layout.NewSpacer(),
		nameContent,
		layout.NewSpacer(),
		commentContent,
		layout.NewSpacer(),
		line1,
		layout.NewSpacer(),
		BoxCreateDatabase,
		layout.NewSpacer(),
		pathEntry,
		layout.NewSpacer(),
		testOpenButton,
		layout.NewSpacer(),
		rowBotton,
	)

	newWindow.Resize(fyne.NewSize(700, 400))
	newWindow.CenterOnScreen()
	newWindow.SetContent(rightColumnContent)
	newWindow.Show()
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
