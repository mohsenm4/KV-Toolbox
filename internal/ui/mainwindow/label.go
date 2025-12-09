package mainwindow

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/dberr"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/ui/labelkv"
	"DatabaseDB/internal/utils"
	"encoding/json"
	"fmt"
	"strings"

	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

var Base string
var NameLabel string

func (r *MainWindow2) NewLabelKV(editType labelkv.EditType, key, value []byte, nameLabel string) *labelkv.TappableLabel {
	keyLabel := labelkv.NewTappableLabel(nameLabel)
	keyLabel.SetTopped(func() {
		r.handleLabelClick(keyLabel, labelkv.EditKey, key, value)
	})
	keyLabel.SetOnHovered(func() {
		keyLabel.Importance = widget.HighImportance
		keyLabel.Refresh()
	})
	return keyLabel
}

func (r *MainWindow2) handleLabelClick(label *labelkv.TappableLabel, editType labelkv.EditType, key, value []byte) {

	// Reset UI
	r.resetLastSelectedLabel(label)
	r.prepareEditArea()

	// Fetch & Process Value
	var (
		finalValue  []byte
		displayText string
		err         error
	)

	if editType == "value" {

		finalValue, displayText, err = r.processValue(key)
		r.EditColumn.finishValue = displayText
	} else {
		finalValue, displayText = r.processKey(key)
	}
	Base = string(finalValue)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = r.AddObjectEdit(editType, key, finalValue)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(r.EditColumn.edit2.Objects) > 0 {
		if lbl, ok := r.EditColumn.edit2.Objects[0].(*widget.Label); ok {
			lbl.SetText(fmt.Sprintf(
				"Edit %s - %s",
				editType,
				utils.TruncateString(displayText, 10),
			))
		}
	}

	// Save button handler
	r.setupSaveButton(label, editType, key, finalValue)

	// Entry change handler
	r.setupEntryOnChange()

	// Focus
	r.focusOnEntry()

}

func (r *MainWindow2) AddObjectEdit(editType labelkv.EditType, key, value []byte) error {

	if editType == labelkv.EditValue {
		typeValue := mimetype.Detect(value)

		switch {
		case strings.HasPrefix(typeValue.String(), "image/"):
			r.ImageShow(key, value, typeValue.Extension())

		case (strings.HasPrefix(typeValue.String(), "application/json")):

			var result json.RawMessage
			err := json.Unmarshal(value, &result)
			if err != nil {
				return err
			}
			prettyJSON, _ := json.MarshalIndent(result, "", "  ")
			// Create Entry
			r.EditColumn.valueEntry = r.ConfigureEntry(string(prettyJSON))
			r.EditColumn.finishValue = string(prettyJSON)

		default:
			// Create Entry
			r.EditColumn.valueEntry = r.ConfigureEntry(string(value))
			r.EditColumn.finishValue = string(value)
		}

	} else {
		// Create Entry
		r.EditColumn.valueEntry = r.ConfigureEntry(string(value))
		r.EditColumn.finishValue = string(value)
	}

	return nil
}

func (r *MainWindow2) resetLastSelectedLabel(current *labelkv.TappableLabel) {
	if r.RightColumn.lastLableKeyAndValue != nil {
		r.RightColumn.lastLableKeyAndValue.Importance = widget.MediumImportance
		r.RightColumn.lastLableKeyAndValue.Refresh()
	}
	current.Importance = widget.HighImportance
	current.Refresh()
	r.RightColumn.lastLableKeyAndValue = current
}

func (r *MainWindow2) prepareEditArea() {
	r.EditColumn.saveEditKey.Disable()
	r.EditColumn.edit2.Objects = nil
	r.EditColumn.edit2.Refresh()

	labelEdit := widget.NewLabel("")
	r.EditColumn.edit2.Add(labelEdit)

	utils.CheckCondition(r.EditColumn.edit2)
}

func (r *MainWindow2) processValue(key []byte) ([]byte, string, error) {

	value, err := variable.CurrentDBClient.Get(key)
	if err != nil && err != dberr.ErrKeyNotFound {
		return nil, "", err
	}

	typeValue := mimetype.Detect(value)

	switch {
	case strings.HasPrefix(typeValue.String(), "image/"):
		return value, fmt.Sprintf("* %s ...", typeValue.Extension()), nil

	case (strings.HasPrefix(typeValue.String(), "application/json")):
		var result json.RawMessage
		err := json.Unmarshal(value, &result)
		if err != nil {
			return value, string(value), nil
		}
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		return prettyJSON, string(prettyJSON), nil

	default:
		return value, string(value), nil
	}
}

func (r *MainWindow2) processKey(key []byte) ([]byte, string) {
	return key, string(key)
}

func (r *MainWindow2) setupSaveButton(label *labelkv.TappableLabel, editType labelkv.EditType, key []byte, value []byte) {

	r.EditColumn.saveEditKey.OnTapped = func() {

		var err error

		if editType == "value" {
			err = logic.SaveValue(key, []byte(r.EditColumn.finishValue))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if strings.HasPrefix(mimetype.Detect([]byte(r.EditColumn.finishValue)).String(), "image/") {
				r.EditColumn.finishValue = fmt.Sprintf("* %s ...", mimetype.Detect([]byte(r.EditColumn.finishValue)).Extension())
			}
		} else {
			// firt Item is new key and second is old key
			_, err := logic.UpdateKey(key, []byte(r.EditColumn.valueEntry.Text))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			for _, v := range r.RightColumn.orgdata {
				if string(v.Key) == string(key) {
					v.Key = []byte(r.EditColumn.valueEntry.Text)
					break
				}
			}
		}

		label.Text = utils.TruncateString(r.EditColumn.finishValue, 10)
		r.EditColumn.edit2.Refresh()
		r.RightColumn.container.Refresh()

		r.EditColumn.saveEditKey.Disable()
	}
}

func (r *MainWindow2) setupEntryOnChange() {
	r.EditColumn.valueEntry.OnChanged = func(s string) {
		if s == Base {
			r.EditColumn.saveEditKey.Disable()
		} else {
			r.EditColumn.saveEditKey.Enable()
		}
		r.EditColumn.finishValue = s
	}
}

func (r *MainWindow2) focusOnEntry() {
	r.Window.Canvas().Focus(r.EditColumn.valueEntry)
}
