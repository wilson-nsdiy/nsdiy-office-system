package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const InstallLockFile = ".installed"

// GetDataDir returns the directory where the SQLite database file is stored.
// It checks the current working directory for the .db file.
func GetDataDir() string {
	// Default to current directory
	return "."
}

// GetInstallLockPath returns the full path to .installed lock file
func GetInstallLockPath() string {
	return filepath.Join(GetDataDir(), InstallLockFile)
}

// NeedsSetup checks if the system needs initial setup.
// Returns true if the .installed lock file does not exist.
func NeedsSetup() bool {
	_, err := os.Stat(GetInstallLockPath())
	return os.IsNotExist(err)
}

// createInstallLock creates a lock file to prevent re-installation
func CreateInstallLock() error {
	content := fmt.Sprintf("installed_at=%s\n", time.Now().UTC().Format(time.RFC3339))
	return os.WriteFile(GetInstallLockPath(), []byte(content), 0400)
}
