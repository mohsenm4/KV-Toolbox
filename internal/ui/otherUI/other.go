package otherUI

import (
	variable "DatabaseDB"
	"fmt"
	"log"

	// "DatabaseDB/internal/logic/addProjectwindowlogic"

	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

var (
	lastStart             *[]byte
	lastEnd               *[]byte
	Orgdata               []dbpak.KVData
	lastPage              int
	previousClose         *widget.Button
	previousProject       *widget.Button
	previousRefreshButton *widget.Button
	lastLableKeyAndValue  *utils.TappableLabel
)

func SetupLastColumn(rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) *fyne.Container {
	lastColumnContent := container.NewVBox()

	jsonDataa, err := variable.CurrentJson.Load()
	if err != nil {
		log.Fatal("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {

			buttonContainer := ProjectButton(project.Name, lastColumnContent, project.FileAddress, rightColumnContentORG, nameButtonProject, buttonAdd, project.Databace, columnEditKey, saveKey, mainWindow)
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

	lightButton.Importance = widget.HighImportance

	darkLight := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, lightButton, darkButton),
	)
	return darkLight
}

func UpdatePage(rightColumnContent *fyne.Container, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) {

	var data = make([]dbpak.KVData, 0)
	var err error
	err = variable.CurrentDBClient.Open()
	if err != nil {
		return
	}
	defer variable.CurrentDBClient.Close()

	if lastEnd == nil && lastStart == nil {
		Orgdata = Orgdata[:0]
	}
	if lastPage < variable.CurrentPage {
		//next page

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		err, data = variable.CurrentDBClient.Read(lastEnd, nil, variable.ItemsPerPage+1)
		if err != nil {
			log.Println(err.Error())
		}

		if len(data) == variable.ItemsPerPage+1 {
			data = data[:variable.ItemsPerPage]
			variable.ItemsAdded = true

		} else {
			variable.ItemsAdded = false

		}
		if len(data) == 0 {
			return
		}
		if len(rightColumnContent.Objects) >= variable.ItemsPerPage*3 {
			Orgdata = Orgdata[len(data):]
		}

		Orgdata = append(Orgdata, data...)
	} else {

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		err, data = variable.CurrentDBClient.Read(nil, lastStart, variable.ItemsPerPage+1)
		if err != nil {
			log.Println(err.Error())
		}

		if len(data) == variable.ItemsPerPage+1 {
			data = data[1:]
			variable.ItemsAdded = true
		}
		if len(data) == 0 {
			return
		}
		Orgdata = Orgdata[:len(Orgdata)-len(data)]
		Orgdata = append(data, Orgdata...)

	}

	lastStart = &Orgdata[0].Key
	lastEnd = &Orgdata[len(Orgdata)-1].Key

	var truncatedValue string
	var arrayContainer []fyne.CanvasObject
	for _, item := range data {

		truncatedKey := utils.TruncateString(string(item.Key), 20)

		typeValue := mimetype.Detect(item.Value)
		if typeValue.Extension() != ".txt" {

			truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
		} else {
			truncatedValue = utils.TruncateString(string(item.Value), 30)

		}

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

	data = data[:0]
	rightColumnContent.Refresh()
	lastPage = variable.CurrentPage
}

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, nameDatabace string, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) *fyne.Container {
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

	label = utils.NewTappableLabel(nameLabel, func() {
		if lastLableKeyAndValue != nil {
			lastLableKeyAndValue.Importance = widget.MediumImportance
			lastLableKeyAndValue.Refresh()
		}
		label.Importance = widget.HighImportance
		label.Refresh()
		lastLableKeyAndValue = label

		utils.CheckCondition(columnEditKey)
		columnEditKey.Add(widget.NewLabel(fmt.Sprintf("Edit %s - %s", editType, nameLabel)))

		if editType == "value" {
			processedValue, err := logic.ProcessValue(value)
			if err == nil {
				value = processedValue
			}

			valueEntry = configureEntry(columnEditKey, string(value))
		} else {
			valueEntry = configureEntry(columnEditKey, string(key))
		}

		saveKey.OnTapped = func() {
			if editType == "value" {
				err := logic.SaveValue(key, []byte(valueEntry.Text), true)
				if err != nil {
					fmt.Println(err.Error())
				}
			} else {
				err := logic.UpdateKey(key, []byte(valueEntry.Text))
				if err != nil {
					fmt.Println(err.Error())
				}
			}

			truncatedText = utils.TruncateString(valueEntry.Text, 20)
			label.SetText(truncatedText)
			label.Refresh()
		}

		columnEditKey.Refresh()
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
