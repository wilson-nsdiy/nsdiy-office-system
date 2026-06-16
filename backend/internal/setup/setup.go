package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"oa-nsdiy/backend/internal/config"
)

const InstallLockFile = ".installed"

// GetDataDir returns the data directory path.
// Uses the same logic as config.DataDir() for consistency.
func GetDataDir() string {
	return config.DataDir()
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

// CreateInstallLock creates a lock file to prevent re-installation
func CreateInstallLock() error {
	dir := GetDataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	content := fmt.Sprintf("installed_at=%s\n", time.Now().UTC().Format(time.RFC3339))
	return os.WriteFile(GetInstallLockPath(), []byte(content), 0400)
}
