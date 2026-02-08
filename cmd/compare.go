package cmd

import (
	"bufio"
	"doppel/internal/detector"
	"doppel/internal/hasher"
	"doppel/internal/scanner"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	deleteFrom1 bool
	deleteFrom2 bool
)

var compareCmd = &cobra.Command{
	Use:   "compare [dir1] [dir2]",
	Short: "Compare two directories and find duplicate files between them",
	Args:  cobra.ExactArgs(2),
	Run:   runCompare,
}

func init() {
	rootCmd.AddCommand(compareCmd)
	compareCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show duplicates without deleting")
	compareCmd.Flags().Int64Var(&minSize, "min-size", 0, "Ignore files smaller than this size in bytes")
	compareCmd.Flags().StringSliceVar(&extensions, "extensions", []string{}, "Filter by file extensions (e.g., .jpg,.png)")
	compareCmd.Flags().BoolVar(&deleteFrom1, "delete-from-1", false, "Auto-delete duplicates from directory 1")
	compareCmd.Flags().BoolVar(&deleteFrom2, "delete-from-2", false, "Auto-delete duplicates from directory 2")
}

func runCompare(cmd *cobra.Command, args []string) {
	dir1, dir2 := args[0], args[1]

	fmt.Printf("Scanning directory 1: %s\n", dir1)
	files1, err := scanner.ScanDirectory(dir1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory 1: %v\n", err)
		os.Exit(1)
	}
	files1 = filterFiles(files1)

	fmt.Printf("Scanning directory 2: %s\n", dir2)
	files2, err := scanner.ScanDirectory(dir2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory 2: %v\n", err)
		os.Exit(1)
	}
	files2 = filterFiles(files2)

	fmt.Printf("Found %d files in dir1, %d files in dir2, hashing...\n", len(files1), len(files2))

	hashed1 := hasher.HashFiles(files1, exact)
	hashed2 := hasher.HashFiles(files2, exact)

	duplicates := findCrossDuplicates(hashed1, hashed2, dir1, dir2)

	if len(duplicates) == 0 {
		fmt.Println("No duplicates found between directories!")
		return
	}

	wastedSpace := detector.CalculateWastedSpace(duplicates)
	fmt.Printf("\nFound %d duplicate groups across directories (%.2f MB duplicated)\n\n", len(duplicates), float64(wastedSpace)/(1024*1024))

	displayCompareDuplicates(duplicates, dir1, dir2)
}

func findCrossDuplicates(hashed1, hashed2 []hasher.HashedFile, dir1, dir2 string) []detector.DuplicateGroup {
	hashMap := make(map[string][]scanner.FileInfo)

	for _, h := range hashed1 {
		hashMap[h.Hash] = append(hashMap[h.Hash], h.FileInfo)
	}

	var duplicates []detector.DuplicateGroup
	for _, h := range hashed2 {
		if files, exists := hashMap[h.Hash]; exists {
			allFiles := append(files, h.FileInfo)
			duplicates = append(duplicates, detector.DuplicateGroup{
				Hash:  h.Hash,
				Files: allFiles,
				Size:  h.FileInfo.Size,
			})
			delete(hashMap, h.Hash)
		}
	}

	return duplicates
}

func displayCompareDuplicates(groups []detector.DuplicateGroup, dir1, dir2 string) {
	for i, group := range groups {
		similarityTag := ""
		if group.IsImage {
			similarityTag = fmt.Sprintf(" ~%d%% similar", group.Similarity)
		}
		fmt.Printf("\nDuplicate Group %d (%.2f MB%s):\n", i+1, float64(group.Size)/(1024*1024), similarityTag)

		var dir1Files, dir2Files []scanner.FileInfo
		for _, file := range group.Files {
			if isInDirectory(file.Path, dir1) {
				dir1Files = append(dir1Files, file)
			} else {
				dir2Files = append(dir2Files, file)
			}
		}

		table := tablewriter.NewTable(os.Stdout)
		table.Header("Directory", "Filename", "Location", "Size (MB)")

		for _, file := range dir1Files {
			filename := filepath.Base(file.Path)
			location := filepath.Dir(file.Path)
			sizeMB := fmt.Sprintf("%.2f", float64(file.Size)/(1024*1024))
			table.Append("[1]", filename, location, sizeMB)
		}

		for _, file := range dir2Files {
			filename := filepath.Base(file.Path)
			location := filepath.Dir(file.Path)
			sizeMB := fmt.Sprintf("%.2f", float64(file.Size)/(1024*1024))
			table.Append("[2]", filename, location, sizeMB)
		}

		table.Render()

		if dryRun {
			continue
		}

		if deleteFrom1 {
			deleteCompareFiles(dir1Files, dir1)
			continue
		}

		if deleteFrom2 {
			deleteCompareFiles(dir2Files, dir2)
			continue
		}

		fmt.Print("\nDelete from [1/2/skip]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			deleteCompareFiles(dir1Files, dir1)
		case "2":
			deleteCompareFiles(dir2Files, dir2)
		case "skip", "":
			fmt.Println("Skipped")
		default:
			fmt.Println("Invalid choice, skipped")
		}
	}
}

func deleteCompareFiles(files []scanner.FileInfo, dirName string) {
	for _, file := range files {
		if err := os.Remove(file.Path); err != nil {
			fmt.Printf("  ✗ %s: %v\n", file.Path, err)
		} else {
			fmt.Printf("  ✓ Deleted from %s: %s\n", dirName, file.Path)
		}
	}
}

func isInDirectory(filePath, dirPath string) bool {
	absFile, _ := filepath.Abs(filePath)
	absDir, _ := filepath.Abs(dirPath)
	rel, err := filepath.Rel(absDir, absFile)
	if err != nil {
		return false
	}
	return !filepath.IsAbs(rel) && !startsWith(rel, "..")
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
