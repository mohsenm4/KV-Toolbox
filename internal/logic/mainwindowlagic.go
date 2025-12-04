package logic

import (
	variable "DatabaseDB"
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/dberr"
	"DatabaseDB/internal/utils"
	"errors"
	"fmt"
	"log"

	"github.com/gabriel-vasile/mimetype"
	// "DatabaseDB/internal/logic/mainwindowlagic"
)

func HandleButtonClick(path string, nameDatabace string) error {
	err := utils.Checkdatabace(path, nameDatabace)
	if err != nil {
		return err
	}

	if !variable.CreatDatabase {

		nun := variable.NameData.FilterFile(path)
		if !nun {
			return fmt.Errorf("error for no found files database")
		}
	}

	return nil
}

func SearchDatabase(valueEntry string) ([][]byte, [][]byte, error) {

	var values [][]byte

	key := utils.CleanInput(valueEntry)
	err, keys := variable.CurrentDBClient.Search([]byte(key))
	if err != nil {
		return nil, nil, err
	}

	if len(keys) == 0 {
		return nil, nil, err
	}

	for _, item := range keys {
		value, err := variable.CurrentDBClient.Get(item)
		if err != nil {
			return nil, nil, err
		}
		value1 := make([]byte, len(value))
		copy(value1, value)

		values = append(values, value1)
	}

	return keys, values, nil
}

func DeleteKeyLogic(valueEntry string) error {

	key := utils.CleanInput(valueEntry)

	value, err := variable.CurrentDBClient.Get([]byte(key))
	if err != nil {
		fmt.Println("error : delete func logic for get key in databace")
	}
	if value != nil {

		err = variable.CurrentDBClient.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("this key does not exist in the database")
		//dialog.ShowInformation("Error", "This key does not exist in the database", editWindow)
	}
}

func AddKeyLogic(inputKey string, valueFinish []byte) error {
	key := utils.CleanInput(inputKey)

	value, err := variable.CurrentDBClient.Get([]byte(key))
	if err != nil && !errors.Is(err, dberr.ErrKeyNotFound) {
		return fmt.Errorf("failed to get key from database: %w", err)
	}

	if value != nil {
		return fmt.Errorf("key '%s' already exists in the database", key)
	}

	if err := variable.CurrentDBClient.Add([]byte(key), valueFinish); err != nil {
		return fmt.Errorf("failed to add key '%s': %w", key, err)
	}

	return nil
}

func QueryKey(inputKey string) ([]byte, error) {

	key := utils.CleanInput(inputKey)

	value, err := variable.CurrentDBClient.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	return value, nil
}

func SaveValue(key, value []byte) error {

	return variable.CurrentDBClient.Add(key, value)
}

func UpdateKey(oldKey, newKey []byte) (string, error) {

	valueBefore, err := variable.CurrentDBClient.Get(oldKey)
	if err != nil {
		return "", err
	}

	if err := variable.CurrentDBClient.Delete(oldKey); err != nil {
		return "", err
	}

	newKey = []byte(utils.CleanInput(string(newKey)))
	if err := variable.CurrentDBClient.Add(newKey, valueBefore); err != nil {
		return "", err
	}

	return string(newKey), nil
}

func FetchPageData(lastStart *[]byte, lastEnd *[]byte, lastPage int, Orgdata []dbpak.KVData) ([]dbpak.KVData, error) {

	var data = make([]dbpak.KVData, 0)
	var err error

	if lastEnd == nil && lastStart == nil {
		Orgdata = nil
	}
	if lastPage < variable.CurrentPage {

		//next page

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		data, err = RangeCursorRead(lastEnd, nil, variable.ItemsPerPage+1)
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
			return data, err
		}
	} else {

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		data, err = RangeCursorRead(nil, lastStart, variable.ItemsPerPage+1)
		if err != nil {
			log.Println(err.Error())
		}

		if len(data) == variable.ItemsPerPage+1 {
			data = data[1:]
			variable.ItemsAdded = true
		}
		if len(data) == 0 {
			return data, err
		}

	}
	return data, nil
}

func FormatKeyValue(item dbpak.KVData) (string, string) {
	truncatedKey := utils.TruncateString(string(item.Key), 20)

	typeValue := mimetype.Detect(item.Value)
	var truncatedValue string
	if typeValue.Extension() != ".txt" {
		truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
	} else {
		truncatedValue = utils.TruncateString(string(item.Value), 30)
	}

	return truncatedKey, truncatedValue
}

func RangeCursorRead(start, end *[]byte, count int) ([]dbpak.KVData, error) {

	var iterms []dbpak.KVData
	for i := 0; i < count; i++ {

		err, data := variable.CurrentDBClient.Read(start, end, 1)
		if err != nil {
			return nil, err
		}
		item := dbpak.KVData{
			Key:   data[0].Key,
			Value: data[0].Value,
		}
		_, truncatedValue := FormatKeyValue(item)

		if end != nil && start == nil {
			end = &item.Key
		} else {
			start = &item.Key
		}
		iterms = append(iterms, dbpak.KVData{Key: item.Key, Value: []byte(truncatedValue)})
	}

	if end != nil && start == nil {

		for i := 0; i < len(iterms)/2; i++ {
			j := len(iterms) - i - 1
			temp := iterms[i]
			iterms[i] = iterms[j]
			iterms[j] = temp
		}
	}
	return iterms, nil
}
