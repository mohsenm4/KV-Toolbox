package variable

import (
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/filterdatabase"
	"sync"
)

var mu sync.RWMutex

var (
	currentDBClient dbpak.DBClient
	currentPage     int
	ItemsPerPage    = 20
	FolderPath      string
	NameData        filterdatabase.FilterData
	itemsAdded      bool
	PreviousOffsetY float32
	resultSearch    bool
	CreatDatabase   bool
)

// CurrentDBClient accessors

func GetCurrentDBClient() dbpak.DBClient {
	mu.RLock()
	defer mu.RUnlock()
	return currentDBClient
}

func SetCurrentDBClient(client dbpak.DBClient) {
	mu.Lock()
	defer mu.Unlock()
	currentDBClient = client
}

func CloseAndSetCurrentDBClient(newClient dbpak.DBClient) {
	mu.Lock()
	defer mu.Unlock()
	if currentDBClient != nil {
		currentDBClient.Close()
	}
	currentDBClient = newClient
}

// CurrentPage accessors

func GetCurrentPage() int {
	mu.RLock()
	defer mu.RUnlock()
	return currentPage
}

func SetCurrentPage(page int) {
	mu.Lock()
	defer mu.Unlock()
	currentPage = page
}

func IncrementCurrentPage() {
	mu.Lock()
	defer mu.Unlock()
	currentPage++
}

func DecrementCurrentPage() int {
	mu.Lock()
	defer mu.Unlock()
	currentPage--
	return currentPage
}

// ItemsAdded accessors

func GetItemsAdded() bool {
	mu.RLock()
	defer mu.RUnlock()
	return itemsAdded
}

func SetItemsAdded(v bool) {
	mu.Lock()
	defer mu.Unlock()
	itemsAdded = v
}

// ResultSearch accessors

func GetResultSearch() bool {
	mu.RLock()
	defer mu.RUnlock()
	return resultSearch
}

func SetResultSearch(v bool) {
	mu.Lock()
	defer mu.Unlock()
	resultSearch = v
}

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
