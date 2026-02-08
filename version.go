package main

var (
	// Version is the current version of doppel
	// This is set during build time using ldflags
	Version = "dev"

	// Commit is the git commit hash
	// This is set during build time using ldflags
	Commit = "unknown"

	// BuildDate is when the binary was built
	// This is set during build time using ldflags
	BuildDate = "unknown"
)
