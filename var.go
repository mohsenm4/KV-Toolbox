package variable

import (
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/filterdatabase"
)

var (
	CurrentDBClient dbpak.DBClient
	FolderPath      string
	NameData        filterdatabase.FilterData
	ItemsAdded      bool
	ResultSearch    bool
	CreatDatabase   bool
)
var CH = make(chan []dbpak.KVData, 50)

var (
	NameDatabase = []string{
		"levelDB",
		"Pebble",
		"Badger",
		//"Redis",
	}
)

/*
export GOOS=darwin
export GOARCH=arm64
export CC=clang
export CGO_ENABLED=0


-------------------
export GOOS=windows
export GOARCH=amd64 || export GOARCH=386
export CC=x86_64-w64-mingw32-gcc || export CC=i686-w64-mingw32-gcc
export CGO_ENABLED='1'


*/
