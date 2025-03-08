package logic

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/utils"
	"fmt"
	// "DatabaseDB/internal/logic/mainwindowlagic"
)

func HandleButtonClick(test string, nameDatabace string) error {
	err := utils.Checkdatabace(test, nameDatabace)
	if err != nil {
		return err
	}

	if !variable.CreatDatabase {

		nun := variable.NameData.FilterFile(test)
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
