package variable

import (
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/filterdatabase"
)

var (
	CurrentDBClient dbpak.DBClient
	CurrentPage     int
	ItemsPerPage    = 20
	FolderPath      string
	NameData        filterdatabase.FilterData
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
