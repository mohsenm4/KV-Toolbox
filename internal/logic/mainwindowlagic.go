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
	keys, err := variable.GetCurrentDBClient().Search([]byte(key))
	if err != nil {
		return nil, nil, err
	}

	if len(keys) == 0 {
		return nil, nil, err
	}

	for _, item := range keys {
		value, err := variable.GetCurrentDBClient().Get(item)
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

	value, err := variable.GetCurrentDBClient().Get([]byte(key))
	if err != nil {
		return fmt.Errorf("failed to get key from database: %w", err)
	}
	if value == nil {
		return fmt.Errorf("this key does not exist in the database")
	}

	return variable.GetCurrentDBClient().Delete([]byte(key))
}

func AddKeyLogic(inputKey string, valueFinish []byte) error {
	key := utils.CleanInput(inputKey)

	value, err := variable.GetCurrentDBClient().Get([]byte(key))
	if err != nil && !errors.Is(err, dberr.ErrKeyNotFound) {
		return fmt.Errorf("failed to get key from database: %w", err)
	}

	if value != nil {
		return fmt.Errorf("key '%s' already exists in the database", key)
	}

	if err := variable.GetCurrentDBClient().Add([]byte(key), valueFinish); err != nil {
		return fmt.Errorf("failed to add key '%s': %w", key, err)
	}

	return nil
}

func QueryKey(inputKey string) ([]byte, error) {

	key := utils.CleanInput(inputKey)

	value, err := variable.GetCurrentDBClient().Get([]byte(key))
	if err != nil {
		return nil, err
	}

	return value, nil
}

func SaveValue(key, value []byte) error {

	_, err := variable.GetCurrentDBClient().Get(key)
	if err != nil {
		return fmt.Errorf("key not found, cannot update value: %w", err)
	}
	return variable.GetCurrentDBClient().Add(key, value)
}

func UpdateKey(oldKey, newKey []byte) (string, error) {

	valueBefore, err := variable.GetCurrentDBClient().Get(oldKey)
	if err != nil {
		return "", fmt.Errorf("failed to get old key: %w", err)
	}

	newKey = []byte(utils.CleanInput(string(newKey)))

	// Add new key first — if this fails, old key is still intact
	if err := variable.GetCurrentDBClient().Add(newKey, valueBefore); err != nil {
		return "", fmt.Errorf("failed to add new key: %w", err)
	}

	// Delete old key only after new key is safely written
	if err := variable.GetCurrentDBClient().Delete(oldKey); err != nil {
		// Rollback: remove the new key we just added
		_ = variable.GetCurrentDBClient().Delete(newKey)
		return "", fmt.Errorf("failed to delete old key: %w", err)
	}

	return string(newKey), nil
}

func FetchPageData(lastStart *[]byte, lastEnd *[]byte, lastPage int, Orgdata []dbpak.KVData) ([]dbpak.KVData, error) {

	var data = make([]dbpak.KVData, 0)
	var err error

	if lastEnd == nil && lastStart == nil {
		Orgdata = nil
	}
	if lastPage < variable.GetCurrentPage() {

		//next page

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		data, err = RangeCursorRead(lastEnd, nil, variable.ItemsPerPage+1)
		if err != nil {
			log.Println(err.Error())
		}

		if len(data) == variable.ItemsPerPage+1 {
			data = data[:variable.ItemsPerPage]
			variable.SetItemsAdded(true)

		} else {
			variable.SetItemsAdded(false)

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
			variable.SetItemsAdded(true)
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

		data, err := variable.GetCurrentDBClient().Read(start, end, 1)
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			return iterms, nil
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
