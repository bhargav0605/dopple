package detector

import (
	"doppel/internal/hasher"
	"doppel/internal/scanner"
)

type DuplicateGroup struct {
	Hash       string
	Files      []scanner.FileInfo
	Size       int64
	Similarity int
	IsImage    bool
}

func FindDuplicates(hashed []hasher.HashedFile, threshold int) []DuplicateGroup {
	var images, nonImages []hasher.HashedFile
	for _, h := range hashed {
		if h.IsImage {
			images = append(images, h)
		} else {
			nonImages = append(nonImages, h)
		}
	}

	var duplicates []DuplicateGroup
	duplicates = append(duplicates, findImageDuplicates(images, threshold)...)
	duplicates = append(duplicates, findExactDuplicates(nonImages)...)

	return duplicates
}

func findImageDuplicates(images []hasher.HashedFile, threshold int) []DuplicateGroup {
	var duplicates []DuplicateGroup
	used := make(map[int]bool)

	for i := 0; i < len(images); i++ {
		if used[i] {
			continue
		}

		var group []scanner.FileInfo
		group = append(group, images[i].FileInfo)

		var totalDistance int
		var comparisons int

		for j := i + 1; j < len(images); j++ {
			if used[j] {
				continue
			}

			distance, _ := images[i].PHash.Distance(images[j].PHash)
			if distance <= threshold {
				group = append(group, images[j].FileInfo)
				used[j] = true
				totalDistance += distance
				comparisons++
			}
		}

		if len(group) > 1 {
			avgDistance := 0
			if comparisons > 0 {
				avgDistance = totalDistance / comparisons
			}
			similarity := 100 - (avgDistance * 100 / 64)

			duplicates = append(duplicates, DuplicateGroup{
				Hash:       images[i].PHash.ToString(),
				Files:      group,
				Size:       images[i].FileInfo.Size,
				Similarity: similarity,
				IsImage:    true,
			})
		}
		used[i] = true
	}

	return duplicates
}

func findExactDuplicates(nonImages []hasher.HashedFile) []DuplicateGroup {
	hashGroups := make(map[string][]scanner.FileInfo)

	for _, h := range nonImages {
		hashGroups[h.Hash] = append(hashGroups[h.Hash], h.FileInfo)
	}

	var duplicates []DuplicateGroup
	for hash, files := range hashGroups {
		if len(files) > 1 {
			duplicates = append(duplicates, DuplicateGroup{
				Hash:       hash,
				Files:      files,
				Size:       files[0].Size,
				Similarity: 100,
				IsImage:    false,
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
