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

type RightColumn2 struct {
	Container            *fyne.Container
	NameButtonProject    *widget.Label
	Line                 *canvas.Line
	Spacer               *widget.Label
	ButtonDelete         *widget.Button
	SearchButton         *widget.Button
	ButtonAdd            *widget.Button
	KeyRightColunm       *widget.Label
	ValueRightColunm     *widget.Label
	LastLableKeyAndValue *utils.TappableLabel
	LastStart            *[]byte
	LastEnd              *[]byte
	LastPage             int
	Orgdata              []dbpak.KVData
}

func NewRightColumn() *RightColumn2 {
	return &RightColumn2{}
}

func (r *MainWindow2) BuildLabelKeyAndValue(editType string, key []byte, value []byte, nameLabel string) *utils.TappableLabel {
	var label *utils.TappableLabel
	var truncatedText string
	var err error
	var truncatedKey2 string
	var nameLable string

	label = utils.NewTappableLabel(nameLabel, func() {
		r.EditColumn.SaveEditKey.Disable()
		if r.RightColumn.LastLableKeyAndValue != nil {
			r.RightColumn.LastLableKeyAndValue.Importance = widget.MediumImportance
			r.RightColumn.LastLableKeyAndValue.Refresh()
		}
		label.Importance = widget.HighImportance
		label.Refresh()
		r.RightColumn.LastLableKeyAndValue = label

		utils.CheckCondition(r.EditColumn.Edit2)
		typeValue := mimetype.Detect([]byte(value))

		labelEdit := widget.NewLabel("")
		r.EditColumn.Edit2.Add(labelEdit)

		if editType == "value" {
			nameLable = fmt.Sprintf("Edit %s - %s", editType, utils.TruncateString(string(value), 20))

			switch {
			case strings.HasPrefix(typeValue.String(), "image/"):
				go utils.ImageShow([]byte(key), []byte(value), r.EditColumn.Edit2, r.Window)
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
				r.EditColumn.ValueEntry = r.ConfigureEntry(string(value))
				value = []byte(r.EditColumn.ValueEntry.Text)
			}
		} else {
			nameLable = fmt.Sprintf("Edit %s - %s", editType, utils.TruncateString(string(key), 20))

			r.EditColumn.ValueEntry = r.ConfigureEntry(string(key))
		}
		labelEdit.SetText(nameLable)

		r.EditColumn.SaveEditKey.OnTapped = func() {
			if editType == "value" {
				truncatedKey2, err = logic.SaveValue(key, []byte(r.EditColumn.ValueEntry.Text))
				if err != nil {
					fmt.Println(err.Error())
				}
				value = []byte(truncatedKey2)

			} else {
				_, err := logic.QueryKey(r.EditColumn.ValueEntry.Text)
				if !errors.Is(err, dberr.ErrKeyNotFound) {
					dialog.NewConfirm(
						"⚠️ Duplicate Key",
						"This key already exists.\nIf you continue, it might be merged and you could lose one of the values.\nDo you still want to continue?",
						func(confirmed bool) {
							if confirmed {
								r.EditColumn.SaveEditKey.Disable()
								truncatedKey2, err = logic.UpdateKey(key, []byte(r.EditColumn.ValueEntry.Text))
								if err != nil {
									dialog.ShowInformation("Error", err.Error(), r.Window)
									return
								}
								key = []byte(truncatedKey2)
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

					truncatedKey2, err = logic.UpdateKey(key, []byte(r.EditColumn.ValueEntry.Text))
					if err != nil {
						dialog.ShowInformation("Error", err.Error(), r.Window)
						return
					}
				} else {
					dialog.ShowInformation("Error", err.Error(), r.Window)
					return
				}
			}

			r.EditColumn.SaveEditKey.Disable()
			truncatedText = utils.TruncateString(truncatedKey2, 20)
			label.SetText(truncatedText)
			labelEdit.SetText(fmt.Sprintf("Edit %s - %s", editType, truncatedText))
			r.EditColumn.Edit2.Refresh()

		}

		r.EditColumn.ValueEntry.OnChanged = func(s string) {

			s = strings.TrimSpace(s)
			if s == label.Text {
				r.EditColumn.SaveEditKey.Disable()
			} else {
				r.EditColumn.SaveEditKey.Enable()
			}
		}
	})
	return label
}

func (r *MainWindow2) TopRightColumn() *fyne.Container {
	r.RightColumn.Line = canvas.NewLine(color.Black)
	r.RightColumn.Line.StrokeWidth = 2

	container := container.NewVBox(
		r.RightColumn.NameButtonProject,
		r.RightColumn.Line,
		r.RightColumn.Spacer,
		r.RightColumn.Tool(),
		r.RightColumn.KeyAndValue(),
	)
	return container
}

func (r *RightColumn2) Tool() *fyne.Container {
	return container.NewGridWithColumns(3, r.ButtonDelete, r.SearchButton, r.ButtonAdd)
}

func (r *RightColumn2) KeyAndValue() *fyne.Container {
	return container.NewGridWithColumns(6, r.KeyRightColunm, widget.NewLabel(""), r.ValueRightColunm, widget.NewLabel(""))
}

func (r *MainWindow2) UpdatePage() {

	data, err := logic.FetchPageData(r.RightColumn.LastStart, r.RightColumn.LastEnd, r.RightColumn.LastPage, r.RightColumn.Orgdata)
	if err != nil {
		return
	}

	if r.RightColumn.LastPage < variable.CurrentPage {

		if len(r.RightColumn.Container.Objects) >= variable.ItemsPerPage*3 {
			r.RightColumn.Orgdata = r.RightColumn.Orgdata[len(data):]
		}

		r.RightColumn.Orgdata = append(r.RightColumn.Orgdata, data...)
	} else {

		r.RightColumn.Orgdata = r.RightColumn.Orgdata[:len(r.RightColumn.Orgdata)-len(data)]
		r.RightColumn.Orgdata = append(data, r.RightColumn.Orgdata...)

	}

	if len(data) != 0 {
		r.RightColumn.LastStart = &r.RightColumn.Orgdata[0].Key
		r.RightColumn.LastEnd = &r.RightColumn.Orgdata[len(r.RightColumn.Orgdata)-1].Key
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
	if r.RightColumn.LastPage > variable.CurrentPage {

		r.RightColumn.Container.Objects = append(arrayContainer, r.RightColumn.Container.Objects...)
	} else {

		r.RightColumn.Container.Objects = append(r.RightColumn.Container.Objects, arrayContainer...)

	}

	r.RightColumn.Container.Refresh()
	r.RightColumn.LastPage = variable.CurrentPage
}
