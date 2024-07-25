package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func ExploreFolder(path string, callback func(path string, info os.FileInfo) error) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		return callback(path, info)
	})
}
