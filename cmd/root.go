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

	"github.com/spf13/cobra"
)

var (
	dryRun     bool
	autoDelete bool
	minSize    int64
	extensions []string
)

var rootCmd = &cobra.Command{
	Use:   "doppel [directory]",
	Short: "Find and remove duplicate files",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show duplicates without deleting")
	rootCmd.Flags().BoolVar(&autoDelete, "auto-delete", false, "Automatically keep first file and delete others")
	rootCmd.Flags().Int64Var(&minSize, "min-size", 0, "Ignore files smaller than this size in bytes")
	rootCmd.Flags().StringSliceVar(&extensions, "extensions", []string{}, "Filter by file extensions (e.g., .jpg,.png)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	directory := args[0]

	files, err := scanner.ScanDirectory(directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	files = filterFiles(files)
	fmt.Printf("Found %d files, hashing...\n", len(files))

	hashed := hasher.HashFiles(files)
	duplicates := detector.FindDuplicates(hashed)

	if len(duplicates) == 0 {
		fmt.Println("No duplicates found!")
		return
	}

	wastedSpace := detector.CalculateWastedSpace(duplicates)
	fmt.Printf("\nFound %d duplicate groups (%.2f MB wasted)\n\n", len(duplicates), float64(wastedSpace)/(1024*1024))

	displayDuplicates(duplicates)
}

func filterFiles(files []scanner.FileInfo) []scanner.FileInfo {
	if minSize == 0 && len(extensions) == 0 {
		return files
	}

	filtered := make([]scanner.FileInfo, 0)
	for _, file := range files {
		if minSize > 0 && file.Size < minSize {
			continue
		}

		if len(extensions) > 0 {
			ext := strings.ToLower(filepath.Ext(file.Path))
			found := false
			for _, allowedExt := range extensions {
				if strings.ToLower(allowedExt) == ext {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, file)
	}

	return filtered
}

func displayDuplicates(groups []detector.DuplicateGroup) {
	for i, group := range groups {
		fmt.Printf("Group %d (%.2f MB, %d files):\n", i+1, float64(group.Size)/(1024*1024), len(group.Files))

		if dryRun {
			for _, file := range group.Files {
				fmt.Printf("  %s\n", file.Path)
			}
			fmt.Println()
			continue
		}

		if autoDelete {
			deleteFiles(group.Files, 0)
			fmt.Println()
			continue
		}

		for j, file := range group.Files {
			fmt.Printf("  [%d] %s\n", j+1, file.Path)
		}

		fmt.Print("\nKeep [1-" + fmt.Sprintf("%d", len(group.Files)) + "/all/skip]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" || input == "skip" {
			fmt.Println("Skipped\n")
			continue
		}

		if input == "all" {
			fmt.Println("Kept all\n")
			continue
		}

		var keepIndex int
		if _, err := fmt.Sscanf(input, "%d", &keepIndex); err != nil || keepIndex < 1 || keepIndex > len(group.Files) {
			fmt.Println("Invalid, skipped\n")
			continue
		}

		deleteFiles(group.Files, keepIndex-1)
		fmt.Println()
	}
}

func deleteFiles(files []scanner.FileInfo, keepIndex int) {
	for i, file := range files {
		if i == keepIndex {
			fmt.Printf("  ✓ %s\n", file.Path)
			continue
		}

		if err := os.Remove(file.Path); err != nil {
			fmt.Printf("  ✗ %s: %v\n", file.Path, err)
		} else {
			fmt.Printf("  ✓ Deleted %s\n", file.Path)
		}
	}
}
