package mojang

import (
	"fmt"
	"time"
)

const (
	LatestVersion = "latest"
)

const (
	VersionsURL      = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
	UsersInfoBulkURL = "https://api.minecraftservices.com/minecraft/profile/lookup/bulk/byname"
	//LatestVersion = "latest"
)

// VersionsResponse is a response for mojang query
type VersionsResponse struct {
	Latest   Latest    `json:"latest"`
	Versions []Version `json:"versions"`
}

// Latest is the latest version
type Latest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

// Version represents a version info
type Version struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	Time        time.Time `json:"time"`
	ReleaseTime time.Time `json:"releaseTime"`
}

// VersionInfoResponse is a specific version info
type VersionInfoResponse struct {
	ID                     string      `json:"id"`
	JavaVersion            JavaVersion `json:"javaVersion"`
	Downloads              Downloads   `json:"downloads"`
	ReleaseTime            string      `json:"releaseTime"`
	Time                   string      `json:"time"`
	Type                   string      `json:"type"`
	MinimunLauncherVersion int         `json:"minimumLauncherVersion"`
}

// JavaVersion defines the Java version to run this version
type JavaVersion struct {
	Component    string `json:"component"`
	MajorVersion int    `json:"majorVersion"`
}

// Downloads holds the artifacts for this version
type Downloads struct {
	Client         Artifact `json:"client"`
	ClientMappings Artifact `json:"client_mappings"`
	Server         Artifact `json:"server"`
	ServerMappings Artifact `json:"server_mappings"`
}

// Artifact is the artifact info
type Artifact struct {
	SHA1 string `json:"sha1"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func (r *VersionsResponse) GetLatestRelease() (*Version, error) {
	v, err := r.GetVersion(r.Latest.Release)
	if err != nil {
		err = fmt.Errorf("filter version by id: %w", err)
		return nil, err
	}

	return v, nil
}

func (r *VersionsResponse) GetVersion(v string) (*Version, error) {
	if v == LatestVersion {
		v = r.Latest.Release
	}
	for _, version := range r.Versions {
		if version.ID == v {
			return &version, nil
		}
	}

	return nil, fmt.Errorf("version '%s' not found", v)
}

type UserIDResponse []UserID

type UserID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
