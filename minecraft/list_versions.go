package minecraft

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/minecraft/versions"
	"sort"
	"time"
)

func ListVersions(ctx context.Context) error {
	c := versions.NewClient(versions.WithTimeout(30 * time.Second))

	vers, err := c.ListVersions(ctx)
	if err != nil {
		err = fmt.Errorf("fetching versions from Mojang API: %w", err)
		return err
	}

	sort.Slice(vers.Versions, func(i, j int) bool {
		return vers.Versions[i].ReleaseTime.After(vers.Versions[j].ReleaseTime)
	})

	for _, v := range vers.Versions {
		if vers.Latest.Release == v.ID || vers.Latest.Snapshot == v.ID {
			fmt.Printf("- %s (latest %s)\n", v.ID, v.Type)
		} else {
			fmt.Printf("- %s (%s)\n", v.ID, v.Type)
		}
	}

	return nil
}
