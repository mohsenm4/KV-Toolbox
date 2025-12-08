package main

import (
	"DatabaseDB/internal/pref"
	"DatabaseDB/internal/ui/mainwindow"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {

	myApp := app.NewWithID("com.DatabaseDB.KV-Toolbox")

	window := mainwindow.NewMainWindow("ManageDB")

	window.Pref = pref.NewPref(myApp)

	window.MainWindow(myApp)

	select {}

	db, err := leveldb.OpenFile("/Users/macbookpro/Desktop/le", nil)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	// value is file in this /Users/macbookpro/Desktop/

	file, err := os.ReadFile("/Users/macbookpro/Desktop/mm.mp3")
	if err != nil {
		panic(err)
	}
	db.Put([]byte("key"), file, nil)
	return
}
