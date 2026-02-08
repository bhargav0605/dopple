package hasher

import (
	"crypto/sha256"
	"doppel/internal/scanner"
	"encoding/hex"
	"io"
	"os"
)

type HashedFile struct {
	FileInfo scanner.FileInfo
	Hash     string
}

func HashFiles(files []scanner.FileInfo) []HashedFile {
	sizeGroups := groupBySize(files)
	var hashed []HashedFile

	for _, group := range sizeGroups {
		if len(group) < 2 {
			continue
		}

		for _, file := range group {
			if hash, err := hashFile(file.Path); err == nil {
				hashed = append(hashed, HashedFile{
					FileInfo: file,
					Hash:     hash,
				})
			}
		}
	}

	return hashed
}

func groupBySize(files []scanner.FileInfo) map[int64][]scanner.FileInfo {
	groups := make(map[int64][]scanner.FileInfo)

	for _, file := range files {
		groups[file.Size] = append(groups[file.Size], file)
	}

	return groups
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
