package logic

import (
	variable "DatabaseDB"
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/utils"
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

func SaveValue(key, value []byte) (string, error) {
	err := variable.CurrentDBClient.Open()
	if err != nil {
		return "", fmt.Errorf("error opening database: %w", err)
	}
	defer variable.CurrentDBClient.Close()

	return string(value), variable.CurrentDBClient.Add(key, value)
}

func UpdateKey(oldKey, newKey []byte) (string, error) {
	err := variable.CurrentDBClient.Open()
	if err != nil {
		return "", fmt.Errorf("error opening database: %w", err)
	}
	defer variable.CurrentDBClient.Close()

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

	err := variable.CurrentDBClient.Open()
	if err != nil {
		return data, err
	}
	defer variable.CurrentDBClient.Close()

	if lastEnd == nil && lastStart == nil {
		Orgdata = Orgdata[:0]
	}
	if lastPage < variable.CurrentPage {

		//next page

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		err, data = variable.CurrentDBClient.Read(lastEnd, nil, variable.ItemsPerPage+1)
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
		err, data = variable.CurrentDBClient.Read(nil, lastStart, variable.ItemsPerPage+1)
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
