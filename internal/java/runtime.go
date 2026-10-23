package java

import (
	"context"
	"fmt"
	"time"
)

type RuntimeManager interface {
	InstallJava(ctx context.Context, dest string, version int, arch, osName string) (string, error)
}

type microsoftRuntimeManager struct {
	timeout time.Duration
}

func NewRuntimeManager(timeout time.Duration) RuntimeManager {
	return &microsoftRuntimeManager{
		timeout: timeout,
	}
}

func (m *microsoftRuntimeManager) InstallJava(ctx context.Context, dest string, version int, arch, osName string) (string, error) {
	path, err := Install(ctx, dest, version, arch, osName, m.timeout)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrJavaInstallFailed, err)
	}
	return path, nil
}
