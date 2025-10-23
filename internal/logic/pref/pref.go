package pref

import (
	"encoding/json"

	"fyne.io/fyne/v2"
)

// keys
var KeyListDB = "ListKey_DB"

type Pref struct {
	Preferences fyne.Preferences
	ListDB      JsonInformation
}

type Project struct {
	Name        string `mapstructure:"name"`
	FileAddress string `mapstructure:"fileAddress"`
	Databace    string `mapstructure:"databace"`
}

type JsonInformation struct {
	RecentProjects []Project `mapstructure:"recentProjects"`
}

func NewPref(a fyne.App) *Pref {
	return &Pref{
		Preferences: a.Preferences(),
	}
}

func (p *Pref) LoadDatabase(key string) (JsonInformation, error) {
	data := p.Preferences.String(key)

	if data == "" {
		return JsonInformation{}, nil
	}

	var items JsonInformation
	err := json.Unmarshal([]byte(data), &items)
	if err != nil {
		return JsonInformation{}, err
	}
	return items, nil
}

func (p *Pref) SaveDatabase(items JsonInformation, key string) error {
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	p.Preferences.SetString(key, string(data))
	return nil
}
