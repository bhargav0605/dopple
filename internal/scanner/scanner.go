package scanner

import (
	"os"
	"path/filepath"
)

type FileInfo struct {
	Path string
	Size int64
}

func ScanDirectory(rootPath string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Size() > 0 {
			files = append(files, FileInfo{
				Path: path,
				Size: info.Size(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
