package main

import (
	"DatabaseDB/internal/pref"
	"DatabaseDB/internal/ui/mainwindow"
	"net/http"
	_ "net/http/pprof"

	"fyne.io/fyne/v2/app"
)

func main() {

	go func() {
		http.ListenAndServe("localhost:6060", nil) // pprof server
	}()
	myApp := app.NewWithID("com.DatabaseDB.KV-Toolbox")

	window := mainwindow.NewMainWindow("ManageDB")

	window.Pref = pref.NewPref(myApp)

	window.MainWindow(myApp)
}
