package minecraft

import (
	"errors"
	"fmt"
)

var (
	// ErrInstallationFailed is returned when the installation process fails
	ErrInstallationFailed = errors.New("installation failed")
	// ErrInvalidInstanceConfig is returned when the instance configuration is invalid
	ErrInvalidInstanceConfig = errors.New("invalid instance configuration")
	// ErrBackupFileNotFound is returned when a backup file to restore does not exist
	ErrBackupFileNotFound = errors.New("backup file not found")
)

// InstallError provides context for installation failures
type InstallError struct {
	InstanceName string
	Operation    string
	Err          error
}

func (e *InstallError) Error() string {
	return fmt.Sprintf("installation failed for %s during %s: %v", e.InstanceName, e.Operation, e.Err)
}

func (e *InstallError) Unwrap() error {
	return e.Err
}
