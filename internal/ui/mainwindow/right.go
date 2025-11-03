package mainwindow

import (
	variable "DatabaseDB"
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/dberr"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

type RightColumn struct {
	container            *fyne.Container
	nameButtonProject    *widget.Label
	spacer               *widget.Label
	buttonDelete         *widget.Button
	searchButton         *widget.Button
	buttonAdd            *widget.Button
	keyRightColunm       *widget.Label
	valueRightColunm     *widget.Label
	lastLableKeyAndValue *utils.TappableLabel
	lastStart            *[]byte
	lastEnd              *[]byte
	lastPage             int
	orgdata              []dbpak.KVData
}

func NewRightColumn() *RightColumn {
	return &RightColumn{}
}

var Base string
var NameLabel string

func (r *MainWindow2) BuildLabelKeyAndValue(editType string, key []byte, value []byte, nameLabel string) *utils.TappableLabel {
	var label *utils.TappableLabel
	var err error
	// Determine the base value based on the edit type
	label = utils.NewTappableLabel(nameLabel, func() {
		r.EditColumn.SaveEditKey.Disable()
		if r.RightColumn.lastLableKeyAndValue != nil {
			r.RightColumn.lastLableKeyAndValue.Importance = widget.MediumImportance
			r.RightColumn.lastLableKeyAndValue.Refresh()
		}
		label.Importance = widget.HighImportance
		label.Refresh()
		r.RightColumn.lastLableKeyAndValue = label

		utils.CheckCondition(r.EditColumn.Edit2)

		labelEdit := widget.NewLabel("")
		r.EditColumn.Edit2.Add(labelEdit)

		if editType == "value" {
			typeValue := mimetype.Detect([]byte(value))
			Base = string(value)

			switch {
			case strings.HasPrefix(typeValue.String(), "image/"):
				r.ImageShow([]byte(key), []byte(value), typeValue.Extension())
				r.EditColumn.FinishValue = string(value)
				NameLabel = fmt.Sprintf("* %s . . .", typeValue.Extension())

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

				r.EditColumn.ValueEntry = r.ConfigureEntry(string(value))
				value = []byte(r.EditColumn.ValueEntry.Text)
				r.EditColumn.FinishValue = string(value)
				NameLabel = string(value)
			}

		} else {
			Base = string(key)
			NameLabel = string(key)

			r.EditColumn.FinishValue = string(key)
			r.EditColumn.ValueEntry = r.ConfigureEntry(string(key))
		}

		labelEdit.SetText(fmt.Sprintf("Edit %s - %s", editType, utils.TruncateString(NameLabel, 10)))
		r.EditColumn.SaveEditKey.OnTapped = func() {
			if editType == "value" {
				err = logic.SaveValue(key, []byte(r.EditColumn.FinishValue))
				if err != nil {
					fmt.Println(err.Error())
				}
				Base = r.EditColumn.FinishValue
				BaseImage = []byte(r.EditColumn.FinishValue)
				//value = []byte(truncatedKey2)

			} else {
				_, err := logic.QueryKey(r.EditColumn.ValueEntry.Text)
				if !errors.Is(err, dberr.ErrKeyNotFound) {
					dialog.NewConfirm(
						"⚠️ Duplicate Key",
						"This key already exists.\nIf you continue, it might be merged and you could lose one of the values.\nDo you still want to continue?",
						func(confirmed bool) {
							if confirmed {
								r.EditColumn.SaveEditKey.Disable()
								Base, err = logic.UpdateKey(key, []byte(r.EditColumn.ValueEntry.Text))
								if err != nil {
									dialog.ShowInformation("Error", err.Error(), r.Window)
									return
								}
								NameLabel = r.EditColumn.ValueEntry.Text
								dialog.ShowInformation("Success", "The key was added successfully.", r.Window)
								return
							} else {
								dialog.ShowInformation("Cancelled", "Operation cancelled.", r.Window)
								return
							}
						},
						r.Window,
					).Show()
					return
				} else if errors.Is(err, dberr.ErrKeyNotFound) {

					Base, err = logic.UpdateKey([]byte(Base), []byte(r.EditColumn.ValueEntry.Text))
					if err != nil {
						dialog.ShowInformation("Error", err.Error(), r.Window)
						return
					}
					NameLabel = r.EditColumn.ValueEntry.Text
					//r.EditColumn.FinishValue = r.EditColumn.ValueEntry.Text
				} else {
					dialog.ShowInformation("Error", err.Error(), r.Window)
					return
				}
			}

			r.EditColumn.SaveEditKey.Disable()
			value = []byte(r.EditColumn.FinishValue)
			truncatedText := utils.TruncateString(NameLabel, 10)
			label.SetText(truncatedText)
			labelEdit.SetText(fmt.Sprintf("Edit %s - %s", editType, truncatedText))
			r.EditColumn.Edit2.Refresh()
			r.RightColumn.container.Refresh()

		}

		r.EditColumn.ValueEntry.OnChanged = func(s string) {

			if s == Base {
				r.EditColumn.SaveEditKey.Disable()
			} else {
				r.EditColumn.SaveEditKey.Enable()
			}
			r.EditColumn.FinishValue = s
			NameLabel = s
		}
		r.Window.Canvas().Focus(r.EditColumn.ValueEntry)
	})
	return label
}

func (r *MainWindow2) TopRightColumn() *fyne.Container {
	r.Objects.line = canvas.NewLine(color.Black)
	r.Objects.line.StrokeWidth = 2

	container := container.NewVBox(
		r.RightColumn.nameButtonProject,
		r.Objects.line,
		r.RightColumn.spacer,
		r.RightColumn.Tool(),
		r.RightColumn.KeyAndValue(),
	)
	return container
}

func (r *RightColumn) Tool() *fyne.Container {
	return container.NewGridWithColumns(3, r.buttonDelete, r.searchButton, r.buttonAdd)
}

func (r *RightColumn) KeyAndValue() *fyne.Container {
	return container.NewGridWithColumns(6, r.keyRightColunm, widget.NewLabel(""), r.valueRightColunm, widget.NewLabel(""))
}

func (r *MainWindow2) UpdatePage() {

	data, err := logic.FetchPageData(r.RightColumn.lastStart, r.RightColumn.lastEnd, r.RightColumn.lastPage, r.RightColumn.orgdata)
	if err != nil {
		return
	}

	if r.RightColumn.lastPage < variable.CurrentPage {

		if len(r.RightColumn.container.Objects) >= variable.ItemsPerPage*3 {
			r.RightColumn.orgdata = r.RightColumn.orgdata[len(data):]
		}

		r.RightColumn.orgdata = append(r.RightColumn.orgdata, data...)
	} else {

		r.RightColumn.orgdata = r.RightColumn.orgdata[:len(r.RightColumn.orgdata)-len(data)]
		r.RightColumn.orgdata = append(data, r.RightColumn.orgdata...)

	}

	if len(data) != 0 {
		r.RightColumn.lastStart = &r.RightColumn.orgdata[0].Key
		r.RightColumn.lastEnd = &r.RightColumn.orgdata[len(r.RightColumn.orgdata)-1].Key
	}

	var truncatedValue string
	var truncatedKey string

	var arrayContainer []fyne.CanvasObject
	for _, item := range data {

		truncatedKey, truncatedValue = logic.FormatKeyValue(item)

		valueLabel := r.BuildLabelKeyAndValue("value", item.Key, item.Value, truncatedValue)
		keyLabel := r.BuildLabelKeyAndValue("key", item.Key, item.Value, truncatedKey)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		arrayContainer = append(arrayContainer, buttonRow)
	}
	if r.RightColumn.lastPage > variable.CurrentPage {

		r.RightColumn.container.Objects = append(arrayContainer, r.RightColumn.container.Objects...)
	} else {

		r.RightColumn.container.Objects = append(r.RightColumn.container.Objects, arrayContainer...)

	}

	r.RightColumn.container.Refresh()
	r.RightColumn.lastPage = variable.CurrentPage
}
