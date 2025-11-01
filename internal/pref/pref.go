package pref

import (
	"encoding/json"

	"fyne.io/fyne/v2"
)

// keys
const (
	KeyListDB = "ListKey_DB"
	KeyLastDB = "LastDBKey"
	KeyTheme  = "ThemeKey"
)

type Pref struct {
	Preferences fyne.Preferences
	ListDB      []Project
}

type Project struct {
	Name        string `mapstructure:"name"`
	FileAddress string `mapstructure:"fileAddress"`
	Databace    string `mapstructure:"databace"`
}

func NewPref(a fyne.App) *Pref {
	return &Pref{
		Preferences: a.Preferences(),
	}
}

func (p *Pref) LoadDatabase(key string) ([]Project, error) {
	data := p.Preferences.String(key)

	if data == "" {
		return []Project{}, nil
	}

	var items []Project
	err := json.Unmarshal([]byte(data), &items)
	if err != nil {
		return []Project{}, err
	}
	return items, nil
}

func (p *Pref) SaveDatabase(items []Project, key string) error {
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	p.Preferences.SetString(key, string(data))
	return nil
}

func (p *Pref) LoadTheme(key string) string {
	theme := p.Preferences.String(key)
	return theme
}

func (p *Pref) SaveTheme(theme string, key string) {
	p.Preferences.SetString(key, theme)
}
