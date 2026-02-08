package cmd

import (
	"bufio"
	"doppel/internal/updater"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	skipConfirm bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update doppel to the latest version",
	Run:   runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation prompt")
}

func runUpdate(cmd *cobra.Command, args []string) {
	if Version == "dev" {
		fmt.Println("Cannot update development build")
		fmt.Println("Please build from source or download a release")
		return
	}

	fmt.Println("Checking for updates...")

	latest, updateAvailable, err := updater.CheckForUpdates(Version, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check for updates: %v\n", err)
		os.Exit(1)
	}

	if !updateAvailable {
		fmt.Printf("Already on latest version (%s)\n", Version)
		return
	}

	fmt.Printf("Update available: %s -> %s\n", Version, latest)

	if !skipConfirm {
		fmt.Print("\nInstall update? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		if input != "y" && input != "yes" {
			fmt.Println("Update cancelled")
			return
		}
	}

	fmt.Println()
	if err := updater.DownloadAndInstall(Version); err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		os.Exit(1)
	}
}
