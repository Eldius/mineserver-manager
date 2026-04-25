package installer

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/java"
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
	path, err := java.Install(ctx, dest, version, arch, osName, m.timeout)
	if err != nil {
		return "", fmt.Errorf("installing java: %w", err)
	}
	return path, nil
}
