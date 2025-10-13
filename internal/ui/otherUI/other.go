package otherUI

import (
	variable "DatabaseDB"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	// "DatabaseDB/internal/logic/addProjectwindowlogic"

	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

var (
	lastStart             *[]byte
	lastEnd               *[]byte
	lastPage              int
	Orgdata               []dbpak.KVData
	previousClose         *widget.Button
	previousProject       *widget.Button
	previousRefreshButton *widget.Button
	lastLableKeyAndValue  *utils.TappableLabel
)

func SetupLastColumn(rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, buttonSearch *widget.Button, buttonDelete *widget.Button, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) *fyne.Container {
	lastColumnContent := container.NewVBox()

	jsonDataa, err := variable.CurrentJson.Load()
	if err != nil {
		log.Fatal("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {

			buttonContainer := ProjectButton(project.Name, lastColumnContent, project.FileAddress, rightColumnContentORG, nameButtonProject, buttonAdd, buttonSearch, buttonDelete, project.Databace, columnEditKey, saveKey, mainWindow)
			lastColumnContent.Add(buttonContainer)
		}
	}

	return lastColumnContent
}

func SetupThemeButtons(app fyne.App) *fyne.Container {
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

func UpdatePage(rightColumnContent *fyne.Container, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) {

	data, err := logic.FetchPageData(lastStart, lastEnd, lastPage, Orgdata)
	if err != nil {
		return
	}

	if lastPage < variable.CurrentPage {

		if len(rightColumnContent.Objects) >= variable.ItemsPerPage*3 {
			Orgdata = Orgdata[len(data):]
		}

		Orgdata = append(Orgdata, data...)
	} else {

		Orgdata = Orgdata[:len(Orgdata)-len(data)]
		Orgdata = append(data, Orgdata...)

	}

	if len(data) != 0 {
		lastStart = &Orgdata[0].Key
		lastEnd = &Orgdata[len(Orgdata)-1].Key
	}

	var truncatedValue string
	var truncatedKey string

	var arrayContainer []fyne.CanvasObject
	for _, item := range data {

		truncatedKey, truncatedValue = logic.FormatKeyValue(item)

		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, rightColumnContent, columnEditKey, saveKey, mainWindow)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, rightColumnContent, columnEditKey, saveKey, mainWindow)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		arrayContainer = append(arrayContainer, buttonRow)
	}
	if lastPage > variable.CurrentPage {

		rightColumnContent.Objects = append(arrayContainer, rightColumnContent.Objects...)
	} else {

		rightColumnContent.Objects = append(rightColumnContent.Objects, arrayContainer...)

	}

	rightColumnContent.Refresh()
	lastPage = variable.CurrentPage
}

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, buttonSearch *widget.Button, buttonDelete *widget.Button, nameDatabace string, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) *fyne.Container {
	var refreshButton *widget.Button
	var projectButton *widget.Button
	var closeButton *widget.Button

	projectButton = widget.NewButton(inputText+" - "+nameDatabace, func() {
		if previousProject != nil {
			previousProject.Importance = widget.MediumImportance
			previousClose.Importance = widget.MediumImportance
			previousRefreshButton.Importance = widget.MediumImportance

			previousProject.Refresh()
			previousClose.Refresh()
			previousRefreshButton.Refresh()
		}
		projectButton.Importance = widget.HighImportance
		closeButton.Importance = widget.HighImportance
		refreshButton.Importance = widget.HighImportance

		projectButton.Refresh()
		closeButton.Refresh()
		refreshButton.Refresh()

		previousProject = projectButton
		previousClose = closeButton
		previousRefreshButton = refreshButton

		variable.ItemsAdded = true
		utils.Checkdatabace(path, nameDatabace)
		buttonAdd.Enable()
		buttonSearch.Enable()
		buttonDelete.Enable()
		variable.FolderPath = path
		lastEnd = nil
		variable.ResultSearch = false
		variable.CurrentPage = 1
		lastPage = 0
		variable.PreviousOffsetY = 0
		lastStart = nil
		utils.CheckCondition(rightColumnContentORG)
		utils.CheckCondition(columnEditKey)
		UpdatePage(rightColumnContentORG, columnEditKey, saveKey, mainWindow)
		nameButtonProject.Text = ""
		nameButtonProject.Text = inputText + " - " + nameDatabace

		nameButtonProject.Refresh()

	})
	buttonContainer := container.NewHBox()

	closeButton = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {

		if nameButtonProject.Text == inputText+" - "+nameDatabace {
			utils.CheckCondition(rightColumnContentORG)
			utils.CheckCondition(columnEditKey)

			buttonAdd.Disable()
			buttonSearch.Disable()
			buttonDelete.Disable()

			nameButtonProject.Text = ""
			nameButtonProject.Refresh()
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

		if nameButtonProject.Text == inputText+" - "+nameDatabace {

			variable.ItemsAdded = true
			utils.Checkdatabace(path, nameDatabace)
			buttonAdd.Enable()
			buttonSearch.Enable()
			buttonDelete.Enable()
			variable.FolderPath = path
			lastEnd = nil
			variable.ResultSearch = false
			variable.CurrentPage = 1
			lastPage = 0
			variable.PreviousOffsetY = 0
			lastStart = nil
			utils.CheckCondition(rightColumnContentORG)
			utils.CheckCondition(columnEditKey)
			UpdatePage(rightColumnContentORG, columnEditKey, saveKey, mainWindow)

			nameButtonProject.Refresh()
		}

	})
	refreshClose := container.NewGridWithColumns(2, refreshButton, closeButton)

	buttonContainer = container.NewBorder(nil, nil, nil, refreshClose, projectButton)
	return buttonContainer
}

func BuidLableKeyAndValue(editType string, key []byte, value []byte, nameLabel string, rightColumn *fyne.Container, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) *utils.TappableLabel {
	var label *utils.TappableLabel
	var valueEntry *widget.Entry
	var truncatedText string
	var err error
	var truncatedKey2 string
	var nameLable string

	label = utils.NewTappableLabel(nameLabel, func() {
		saveKey.Disable()
		if lastLableKeyAndValue != nil {
			lastLableKeyAndValue.Importance = widget.MediumImportance
			lastLableKeyAndValue.Refresh()
		}
		label.Importance = widget.HighImportance
		label.Refresh()
		lastLableKeyAndValue = label

		utils.CheckCondition(columnEditKey)
		typeValue := mimetype.Detect([]byte(value))

		labelEdit := widget.NewLabel("")
		columnEditKey.Add(labelEdit)

		if editType == "value" {
			nameLable = fmt.Sprintf("Edit %s - %s", editType, utils.TruncateString(string(value), 20))

			switch {
			case strings.HasPrefix(typeValue.String(), "image/"):
				go utils.ImageShow([]byte(key), []byte(value), columnEditKey, mainWindow)
				truncatedKey2 = fmt.Sprintf("* %s . . .", typeValue.Extension())

			case strings.HasPrefix(typeValue.String(), "text/") || strings.HasPrefix(typeValue.String(), "application/"):
				if strings.HasPrefix(typeValue.String(), "application/json") {
					var result json.RawMessage

					err := json.Unmarshal([]byte(value), &result)
					if err != nil {
						return
					}
					prettyJSON, err := json.MarshalIndent(result, "", "  ")
					if err != nil {
						return
					}
					value = prettyJSON

				}
				valueEntry = configureEntry(columnEditKey, string(value))
				value = []byte(valueEntry.Text)
			}
		} else {
			nameLable = fmt.Sprintf("Edit %s - %s", editType, utils.TruncateString(string(key), 20))

			valueEntry = configureEntry(columnEditKey, string(key))
		}
		labelEdit.SetText(nameLable)

		saveKey.OnTapped = func() {
			if editType == "value" {
				if bytes.Equal(value, []byte(valueEntry.Text)) {
					return
				}
				truncatedKey2, err = logic.SaveValue(key, []byte(valueEntry.Text))
				if err != nil {
					fmt.Println(err.Error())
				}
				value = []byte(truncatedKey2)

			} else {
				if bytes.Equal(key, []byte(valueEntry.Text)) {
					return
				}
				va := logic.QueryKey(valueEntry.Text)
				if va != nil {
					dialog.ShowInformation("Error", "This key has already been added to your database", mainWindow)
					return
				}
				truncatedKey2, err = logic.UpdateKey(key, []byte(valueEntry.Text))
				if err != nil {
					fmt.Println(err.Error())
				}
				key = []byte(truncatedKey2)

			}
			saveKey.Disable()
			truncatedText = utils.TruncateString(truncatedKey2, 20)
			label.SetText(truncatedText)
			labelEdit.SetText(fmt.Sprintf("Edit %s - %s", editType, truncatedText))
			columnEditKey.Refresh()

		}
		columnEditKey.Refresh()

		valueEntry.OnChanged = func(s string) {

			if s == label.Text {
				saveKey.Disable()
			} else {
				saveKey.Enable()
			}
		}
	})
	return label
}

func configureEntry(columnEditKey *fyne.Container, content string) *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Resize(fyne.NewSize(400, 500))
	entry.SetText(content)
	scrollableEntry := container.NewScroll(entry)
	scrollableEntry.SetMinSize(fyne.NewSize(200, 300))
	columnEditKey.Add(scrollableEntry)
	return entry
}
