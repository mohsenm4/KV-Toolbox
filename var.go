package variable

import (
	dbpak "DatabaseDB/internal/Databaces"
	configApp "DatabaseDB/internal/config"
	"DatabaseDB/internal/filterdatabase"
)

var (
	CurrentDBClient dbpak.DBClient
	CurrentPage     int
	ItemsPerPage    = 16
	FolderPath      string
	CurrentJson     *configApp.Config
	NameData        filterdatabase.FilterData
	ItemsAdded      bool
	PreviousOffsetY float32
	ResultSearch    bool
	CreatDatabase   bool
)

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
