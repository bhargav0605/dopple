package hasher

import (
	"crypto/sha256"
	"doppel/internal/scanner"
	"encoding/hex"
	"io"
	"os"

	"github.com/corona10/goimagehash"
)

type HashedFile struct {
	FileInfo   scanner.FileInfo
	Hash       string
	PHash      *goimagehash.ImageHash
	IsImage    bool
	Similarity int
}

func HashFiles(files []scanner.FileInfo, exact bool) []HashedFile {
	var images, nonImages []scanner.FileInfo
	for _, file := range files {
		if !exact && isImage(file.Path) {
			images = append(images, file)
		} else {
			nonImages = append(nonImages, file)
		}
	}

	var hashed []HashedFile

	for _, file := range images {
		if phash, err := perceptualHashImage(file.Path); err == nil {
			hashed = append(hashed, HashedFile{
				FileInfo: file,
				PHash:    phash,
				IsImage:  true,
			})
		}
	}

	sizeGroups := groupBySize(nonImages)
	for _, group := range sizeGroups {
		if len(group) < 2 {
			continue
		}

		for _, file := range group {
			if hash, err := hashFile(file.Path); err == nil {
				hashed = append(hashed, HashedFile{
					FileInfo: file,
					Hash:     hash,
					IsImage:  false,
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
