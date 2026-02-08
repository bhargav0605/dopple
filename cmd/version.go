package cmd

import (
	"doppel/internal/updater"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var checkUpdate bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&checkUpdate, "check", false, "Check for updates")
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("doppel %s\n", Version)
	if Commit != "unknown" {
		fmt.Printf("Commit: %s\n", Commit)
	}
	if BuildDate != "unknown" {
		fmt.Printf("Built:  %s\n", BuildDate)
	}

	if checkUpdate {
		fmt.Println()
		latest, updateAvailable, err := updater.CheckForUpdates(Version, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check for updates: %v\n", err)
			return
		}

		if updateAvailable {
			fmt.Printf("Latest: %s\n", latest)
			fmt.Println("\033[33mUpdate available\033[0m")
			fmt.Println("Run 'doppel update' to upgrade")
		} else {
			fmt.Println("Up to date")
		}
	}
}
