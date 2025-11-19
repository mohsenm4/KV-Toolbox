package logic

import (
	variable "DatabaseDB"
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/dberr"
	"DatabaseDB/internal/utils"
	"errors"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
)

type Logic struct{}

func NewLogic() *Logic {
	return &Logic{}
}

func (l *Logic) SearchDatabase(valueEntry string) ([][]byte, [][]byte, error) {

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

func (l *Logic) DeleteKeyLogic(valueEntry string) error {

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

func (l *Logic) AddKeyLogic(inputKey string, valueFinish []byte) error {
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

func (l *Logic) QueryKey(inputKey string) ([]byte, error) {

	key := utils.CleanInput(inputKey)

	value, err := variable.CurrentDBClient.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (l *Logic) SaveValue(key, value []byte) error {

	return variable.CurrentDBClient.Add(key, value)
}

func (l *Logic) UpdateKey(oldKey, newKey []byte) (string, error) {

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

func (l *Logic) FormatKeyValue(item dbpak.KVData) (string, string) {
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
func (l *Logic) GetAllKeys() ([]dbpak.KVData, error) {
	var result []dbpak.KVData
	var startKey, endKey *[]byte

	for {
		// Read the next batch of records (20 items per batch)
		err, batch := variable.CurrentDBClient.Read(startKey, endKey, 200)
		if err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}

		// Append this batch to the final result
		result = append(result, batch...)

		// Set the start key for the next iteration (continue from the last key)
		startKey = &batch[len(batch)-1].Key
	}

	return result, nil
}
