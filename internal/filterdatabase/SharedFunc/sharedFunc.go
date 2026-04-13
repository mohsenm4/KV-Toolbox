package sharedfunc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FormatFilesDatabase(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error opening folder:", err)
		return false
	}
	var count uint8
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "MANIFEST-") || filepath.Ext(file.Name()) == ".log" {
			count++
		}

		if count == 2 {
			return true
		}
	}
	return false
}
