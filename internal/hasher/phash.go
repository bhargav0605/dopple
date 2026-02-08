package hasher

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
)

var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".tiff": true,
	".webp": true,
	".heic": true,
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return imageExtensions[ext]
}

func perceptualHashImage(path string) (*goimagehash.ImageHash, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	hash, err := goimagehash.DifferenceHash(img)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func hammingDistance(hash1, hash2 *goimagehash.ImageHash) int {
	distance, _ := hash1.Distance(hash2)
	return distance
}
