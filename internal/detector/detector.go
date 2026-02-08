package detector

import (
	"doppel/internal/hasher"
	"doppel/internal/scanner"
)

type DuplicateGroup struct {
	Hash  string
	Files []scanner.FileInfo
	Size  int64
}

func FindDuplicates(hashed []hasher.HashedFile) []DuplicateGroup {
	hashGroups := make(map[string][]scanner.FileInfo)

	for _, h := range hashed {
		hashGroups[h.Hash] = append(hashGroups[h.Hash], h.FileInfo)
	}

	var duplicates []DuplicateGroup
	for hash, files := range hashGroups {
		if len(files) > 1 {
			duplicates = append(duplicates, DuplicateGroup{
				Hash:  hash,
				Files: files,
				Size:  files[0].Size,
			})
		}
	}

	return duplicates
}

func CalculateWastedSpace(groups []DuplicateGroup) int64 {
	var total int64
	for _, group := range groups {
		total += group.Size * int64(len(group.Files)-1)
	}
	return total
}
