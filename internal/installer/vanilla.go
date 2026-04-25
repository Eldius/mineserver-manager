package installer

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/model"
	"github.com/eldius/mineserver-manager/internal/mojang"
)

type vanillaFlavor struct {
	client mojang.Client
}

func NewVanillaFlavor(client mojang.Client) ServerFlavor {
	return &vanillaFlavor{
		client: client,
	}
}

func (f *vanillaFlavor) Name() model.MineFlavour {
	return model.MineFlavourVanilla
}

func (f *vanillaFlavor) ListVersions(ctx context.Context) ([]string, error) {
	ver, err := f.client.ListVersions(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing vanilla versions: %w", err)
	}
	var versions []string
	for _, v := range ver.Versions {
		versions = append(versions, v.ID)
	}
	return versions, nil
}

func (f *vanillaFlavor) GetVersionInfo(ctx context.Context, version string) (*FlavorVersionInfo, error) {
	ver, err := f.client.ListVersions(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting versions list: %w", err)
	}
	v, err := ver.GetVersion(version)
	if err != nil {
		return nil, fmt.Errorf("finding version %s: %w", version, err)
	}
	info, err := f.client.GetVersionInfo(ctx, *v)
	if err != nil {
		return nil, fmt.Errorf("getting version info for %s: %w", version, err)
	}

	return &FlavorVersionInfo{
		Version:     info.ID,
		DownloadURL: info.Downloads.Server.URL,
		SHA1:        info.Downloads.Server.SHA1,
		JavaVersion: info.JavaVersion.MajorVersion,
	}, nil
}
