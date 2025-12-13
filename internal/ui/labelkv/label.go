package labelkv

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type EditType string

const (
	EditValue EditType = "value"
	EditKey   EditType = "key"
)

type TappableLabel struct {
	widget.Label
	onTapped  func()
	onHovered func()
	hover     bool
	editType  EditType
	keyLabel  *TappableLabel
	key       []byte
}

func (t *TappableLabel) SetKeyLabel(key *TappableLabel) {
	t.keyLabel = key
}

func (t *TappableLabel) GetKeyLabel() *TappableLabel {
	return t.keyLabel
}

func (t *TappableLabel) GetKey() []byte {
	return t.key
}

func (t *TappableLabel) SetKey(key []byte) {
	t.key = key
}

func (t *TappableLabel) SetEditType(editType EditType) {
	t.editType = editType
}

func (t *TappableLabel) GetEditType() EditType {
	return t.editType
}

func NewTappableLabel(text string) *TappableLabel {
	labelee := &TappableLabel{
		Label: widget.Label{
			Text: text,
		},
	}
	labelee.ExtendBaseWidget(labelee)
	return labelee
}

func (t *TappableLabel) SetTopped(f func()) {
	t.onTapped = f
}

func (t *TappableLabel) Tapped(_ *fyne.PointEvent) {
	t.onTapped()
}

func (t *TappableLabel) MouseIn(_ *desktop.MouseEvent) {
	if t.onHovered != nil {
		t.hover = true
		t.onHovered()
	}
}

func (t *TappableLabel) MouseMoved(_ *desktop.MouseEvent) {}
func (t *TappableLabel) MouseOut() {
	t.hover = false
	t.Refresh()
}

func (t *TappableLabel) SetOnHovered(f func()) {
	t.onHovered = f
}

func (t *TappableLabel) SetMouseOut(f func()) {}

func (t *TappableLabel) Refresh() {
	if t.hover {
		t.Label.TextStyle = fyne.TextStyle{Bold: true}
		t.Label.Importance = widget.HighImportance
	} else {
		t.Label.TextStyle = fyne.TextStyle{Bold: false}
		t.Label.Importance = widget.MediumImportance
	}
	t.Label.Refresh()
}
