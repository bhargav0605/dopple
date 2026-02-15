package cmd

import (
	"bufio"
	"doppel/internal/detector"
	"doppel/internal/hasher"
	"doppel/internal/scanner"
	"doppel/internal/updater"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

var (
	dryRun     bool
	autoDelete bool
	minSize    int64
	extensions []string
	showAll    bool
	exact      bool
	threshold  int
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
	rootCmd.Flags().BoolVar(&showAll, "show-all", false, "Show all duplicates first, then delete with single confirmation")
	rootCmd.Flags().BoolVar(&exact, "exact", false, "Use exact byte matching for all files (disable perceptual hashing for images)")
	rootCmd.Flags().IntVar(&threshold, "threshold", 5, "Similarity threshold for images (0-64, lower = more similar)")
	rootCmd.Flags().Int64Var(&minSize, "min-size", 0, "Ignore files smaller than this size in bytes")
	rootCmd.Flags().StringSliceVar(&extensions, "extensions", []string{}, "Filter by file extensions (e.g., .jpg,.png)")
}

func Execute() {
	checkForUpdatesOnStartup()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func checkForUpdatesOnStartup() {
	if Version == "dev" {
		return
	}

	latest, updateAvailable, err := updater.CheckForUpdates(Version, false)
	if err != nil || !updateAvailable {
		return
	}

	fmt.Printf("\033[33mUpdate available: %s -> %s (run 'doppel update')\033[0m\n\n", Version, latest)
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

	hashed := hasher.HashFiles(files, exact)
	duplicates := detector.FindDuplicates(hashed, threshold)

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
	if showAll {
		displayAllThenDelete(groups)
		return
	}

	for i, group := range groups {
		similarityTag := ""
		if group.IsImage {
			similarityTag = fmt.Sprintf(" ~%d%% similar", group.Similarity)
		}
		fmt.Printf("\nGroup %d (%.2f MB, %d files%s):\n", i+1, float64(group.Size)/(1024*1024), len(group.Files), similarityTag)

		displayGroupTable(group.Files)

		if dryRun {
			continue
		}

		if autoDelete {
			deleteFiles(group.Files, 0)
			continue
		}

		fmt.Print("\nKeep [1-" + fmt.Sprintf("%d", len(group.Files)) + "/all/skip]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" || input == "skip" {
			fmt.Println("Skipped")
			continue
		}

		if input == "all" {
			fmt.Println("Kept all")
			continue
		}

		var keepIndex int
		if _, err := fmt.Sscanf(input, "%d", &keepIndex); err != nil || keepIndex < 1 || keepIndex > len(group.Files) {
			fmt.Println("Invalid, skipped")
			continue
		}

		deleteFiles(group.Files, keepIndex-1)
	}
}

func displayGroupTable(files []scanner.FileInfo) {
	table := tablewriter.NewTable(os.Stdout)
	table.Header("#", "Filename", "Location", "Size (MB)")

	for i, file := range files {
		filename := filepath.Base(file.Path)
		location := filepath.Dir(file.Path)
		sizeMB := fmt.Sprintf("%.2f", float64(file.Size)/(1024*1024))
		_ = table.Append(
			fmt.Sprintf("[%d]", i+1),
			filename,
			location,
			sizeMB,
		)
	}

	_ = table.Render()
}

func displayAllThenDelete(groups []detector.DuplicateGroup) {
	fmt.Println("=== All Duplicate Groups ===")

	for i, group := range groups {
		similarityTag := ""
		if group.IsImage {
			similarityTag = fmt.Sprintf(" ~%d%% similar", group.Similarity)
		}
		fmt.Printf("Group %d (%.2f MB, %d files%s):\n", i+1, float64(group.Size)/(1024*1024), len(group.Files), similarityTag)

		table := tablewriter.NewTable(os.Stdout)
		table.Header("Action", "Filename", "Location", "Size (MB)")

		for j, file := range group.Files {
			action := "[DEL]"
			if j == 0 {
				action = "[KEEP]"
			}
			filename := filepath.Base(file.Path)
			location := filepath.Dir(file.Path)
			sizeMB := fmt.Sprintf("%.2f", float64(file.Size)/(1024*1024))
			_ = table.Append(action, filename, location, sizeMB)
		}

		_ = table.Render()
		fmt.Println()
	}

	fmt.Print("\nDelete all duplicates (keep first file in each group)? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "y" && input != "yes" {
		fmt.Println("Cancelled")
		return
	}

	fmt.Println("\nDeleting duplicates...")
	var totalDeleted, totalErrors int

	for i, group := range groups {
		fmt.Printf("\nGroup %d:\n", i+1)
		deleted, errors := deleteFilesCount(group.Files, 0)
		totalDeleted += deleted
		totalErrors += errors
	}

	fmt.Printf("\n✓ Deleted %d files", totalDeleted)
	if totalErrors > 0 {
		fmt.Printf(" (%d errors)", totalErrors)
	}
	fmt.Println()
}

func deleteFilesCount(files []scanner.FileInfo, keepIndex int) (int, int) {
	deleted, errors := 0, 0
	for i, file := range files {
		if i == keepIndex {
			fmt.Printf("  ✓ %s\n", file.Path)
			continue
		}

		if err := os.Remove(file.Path); err != nil {
			fmt.Printf("  ✗ %s: %v\n", file.Path, err)
			errors++
		} else {
			fmt.Printf("  ✓ Deleted %s\n", file.Path)
			deleted++
		}
	}
	return deleted, errors
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
