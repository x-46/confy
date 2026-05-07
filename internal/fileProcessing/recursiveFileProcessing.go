package fileprocessing

import (
	"os"
	"path/filepath"
)

func RecursiveFileProcessing(rootDir string, fileExtensions []string, callback func(path string) error) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			for _, ext := range fileExtensions {
				if filepath.Ext(path) == ext {
					return callback(path)
				}
			}
		}
		return nil
	})
}
