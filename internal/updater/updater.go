package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubAPI     = "https://api.github.com/repos/bhargav0605/dopple/releases/latest"
	cacheFile     = ".doppel-version-cache.json"
	cacheDuration = 24 * time.Hour
)

type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type VersionCache struct {
	LastCheck      time.Time `json:"last_check"`
	LatestVersion  string    `json:"latest_version"`
	CurrentVersion string    `json:"current_version"`
}

func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(home, ".cache", "doppel")
	if err := os.MkdirAll(cacheDir, 0750); err != nil {
		return "", err
	}
	return cacheDir, nil
}

func loadCache() (*VersionCache, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}

	cachePath := filepath.Join(cacheDir, cacheFile)
	// #nosec G304 - cachePath is constructed from user home dir and constant filename
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cache VersionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

func saveCache(cache *VersionCache) error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	cachePath := filepath.Join(cacheDir, cacheFile)
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0600)
}

func CheckForUpdates(currentVersion string, forceCheck bool) (string, bool, error) {
	if currentVersion == "dev" {
		return "", false, nil
	}

	if !forceCheck {
		cache, err := loadCache()
		if err == nil && time.Since(cache.LastCheck) < cacheDuration {
			if cache.CurrentVersion == currentVersion {
				updateAvailable := cache.LatestVersion != currentVersion
				return cache.LatestVersion, updateAvailable, nil
			}
		}
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(githubAPI)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", false, err
	}

	latestVersion := release.TagName

	cache := &VersionCache{
		LastCheck:      time.Now(),
		LatestVersion:  latestVersion,
		CurrentVersion: currentVersion,
	}
	_ = saveCache(cache)

	updateAvailable := latestVersion != currentVersion
	return latestVersion, updateAvailable, nil
}

func DownloadAndInstall(currentVersion string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(githubAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release info: %w", err)
	}

	binaryName := getBinaryName()
	var downloadURL string

	for _, asset := range release.Assets {
		if asset.Name == binaryName || strings.Contains(asset.Name, binaryName) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	fmt.Printf("Downloading %s...\n", release.TagName)

	downloadResp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer downloadResp.Body.Close()

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// #nosec G302 - executable binary requires 0755 permissions to be runnable
	newFile, err := os.OpenFile(execPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		_ = os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to create new binary: %w", err)
	}

	_, err = io.Copy(newFile, downloadResp.Body)
	if closeErr := newFile.Close(); closeErr != nil && err == nil {
		err = closeErr
	}

	if err != nil {
		_ = os.Remove(execPath)
		_ = os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to write new binary: %w", err)
	}

	_ = os.Remove(backupPath)

	fmt.Printf("Successfully updated to %s\n", release.TagName)
	return nil
}

func getBinaryName() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	name := fmt.Sprintf("doppel-%s-%s", goos, goarch)
	if goos == "windows" {
		name += ".exe"
	}
	return name
}
