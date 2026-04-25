package installer

import (
	"context"
	"github.com/eldius/mineserver-manager/internal/model"
)

type ServerFlavor interface {
	Name() model.MineFlavour
	ListVersions(ctx context.Context) ([]string, error)
	GetVersionInfo(ctx context.Context, version string) (*FlavorVersionInfo, error)
}

type FlavorVersionInfo struct {
	Version     string
	DownloadURL string
	SHA1        string
	JavaVersion int
}
