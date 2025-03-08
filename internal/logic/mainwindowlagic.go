package logic

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/utils"
	"encoding/json"
	"fmt"
	"strings"

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
	err = variable.CurrentDBClient.Open()
	if err != nil {
		return err
	}
	defer variable.CurrentDBClient.Close()

	return nil
}

func SearchDatabase(valueEntry string) ([][]byte, [][]byte, error) {

	var values [][]byte
	err := variable.CurrentDBClient.Open()
	if err != nil {
		return nil, nil, err
	}
	defer variable.CurrentDBClient.Close()

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

	err := variable.CurrentDBClient.Open()
	if err != nil {
		return err
	}
	defer variable.CurrentDBClient.Close()

	key := utils.CleanInput(valueEntry)

	value := QueryKey(valueEntry)
	if value != nil {

		err = variable.CurrentDBClient.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("This key does not exist in the database")
		//dialog.ShowInformation("Error", "This key does not exist in the database", editWindow)
	}
}

func AddKeyLogic(iputKey string, valueFinish []byte) error {

	key := utils.CleanInput(iputKey)

	err := variable.CurrentDBClient.Open()
	if err != nil {
		return err
	}
	defer variable.CurrentDBClient.Close()

	value := QueryKey(iputKey)
	if value != nil {
		//dialog.ShowInformation("Error", "This key has already been added to your database", windowAdd)
		return fmt.Errorf("This key has already been added to your database")
	} else {
		err = variable.CurrentDBClient.Add([]byte(key), valueFinish)
		if err != nil {
			fmt.Print(err.Error())
		}

		return nil
	}

}

func QueryKey(iputKey string) []byte {

	key := utils.CleanInput(iputKey)

	value, err := variable.CurrentDBClient.Get([]byte(key))
	if err != nil {
		fmt.Println("error : delete func logic for get key in databace")
	}
	return value
}

func ProcessValue(value []byte) ([]byte, error) {
	typeValue := mimetype.Detect(value)
	if strings.HasPrefix(typeValue.String(), "application/json") {
		var result json.RawMessage
		err := json.Unmarshal(value, &result)
		if err != nil {
			return nil, err
		}
		return json.MarshalIndent(result, "", "  ")
	}
	return value, nil
}

func SaveValue(key, value []byte, isText bool) error {
	err := variable.CurrentDBClient.Open()
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer variable.CurrentDBClient.Close()

	if isText {
		value = []byte(utils.TruncateString(string(value), 30))
	}

	return variable.CurrentDBClient.Add(key, value)
}

func UpdateKey(oldKey, newKey []byte) error {
	err := variable.CurrentDBClient.Open()
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer variable.CurrentDBClient.Close()

	valueBefore, err := variable.CurrentDBClient.Get(oldKey)
	if err != nil {
		return err
	}

	if err := variable.CurrentDBClient.Delete(oldKey); err != nil {
		return err
	}

	newKey = []byte(utils.CleanInput(string(newKey)))
	return variable.CurrentDBClient.Add(newKey, valueBefore)
}
