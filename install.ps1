# PowerShell installation script for Windows

$ErrorActionPreference = "Stop"

$INSTALL_DIR = "$env:LOCALAPPDATA\doppel"
$BINARY_NAME = "doppel.exe"

Write-Host "Installing doppel..." -ForegroundColor Green

# Create install directory
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null

# Check if binary exists locally
if (Test-Path ".\$BINARY_NAME") {
    Write-Host "Installing from local binary..." -ForegroundColor Yellow
    Copy-Item ".\$BINARY_NAME" -Destination "$INSTALL_DIR\$BINARY_NAME" -Force
} else {
    Write-Host "Binary not found. Please build first with: go build -o doppel.exe" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "✓ doppel installed to $INSTALL_DIR\$BINARY_NAME" -ForegroundColor Green
Write-Host ""

# Check if install directory is in PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    Write-Host "Adding $INSTALL_DIR to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$INSTALL_DIR",
        "User"
    )
    Write-Host "✓ PATH updated. Please restart your terminal." -ForegroundColor Green
} else {
    Write-Host "✓ Installation complete! Run 'doppel' from anywhere." -ForegroundColor Green
}
