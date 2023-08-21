package vanilla

import (
	"encoding/json"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	c *http.Client
}

// NewClient creates a new client
func NewClient(t time.Duration) *Client {
	return &Client{c: &http.Client{Timeout: t}}
}

// ListVersions lists all available versions
func (c *Client) ListVersions() (*VersionsResponse, error) {
	res, err := c.c.Get(VersionsURL)
	if err != nil {
		err = fmt.Errorf("getting available versions: %w", err)
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var versions VersionsResponse
	if err = json.NewDecoder(res.Body).Decode(&versions); err != nil {
		err = fmt.Errorf("decoding available versions response: %w", err)
		return nil, err
	}

	return &versions, nil
}

// GetVersionInfo gets a specific version info
func (c *Client) GetVersionInfo(v Version) (*VersionInfoResponse, error) {
	res, err := c.c.Get(v.URL)
	if err != nil {
		err = fmt.Errorf("getting version info: %w", err)
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var version VersionInfoResponse
	if err = json.NewDecoder(res.Body).Decode(&version); err != nil {
		err = fmt.Errorf("decoding version info response: %w", err)
		return nil, err
	}

	return &version, nil
}

// InstallWithConfig installs selected version
func (c *Client) InstallWithConfig(cfgs ...InstallCfg) error {
	cfg := &InstallConfig{
		Start:       GetDefaultScriptParams(),
		SrvProps:    utils.Must(GetDefaultServerProperties()),
		Dest:        "./minecraft",
		VersionName: "latest",
		v:           nil,
	}

	for _, c := range cfgs {
		cfg = c(cfg)
	}

	versions, err := c.ListVersions()
	if err != nil {
		err = fmt.Errorf("getting available versions: %w", err)
		return err
	}

	v, err := versions.GetVersion(cfg.VersionName)
	if err != nil {
		err = fmt.Errorf("getting version from online versions list for name '%s': %w", cfg.VersionName, err)
		return err
	}

	cfg.v, err = c.GetVersionInfo(*v)
	if err != nil {
		err = fmt.Errorf("getting version info for name '%s': %w", cfg.VersionName, err)
		return err
	}

	sf, err := c.DownloadServer(*cfg.v, cfg.Dest)
	if err != nil {
		err = fmt.Errorf("getting version info to install: %w", err)
		return err
	}

	log.Printf("server file: %s", sf)
	return err
}

// Install installs selected version
func (c *Client) Install(v VersionInfoResponse, cfg *InstallConfig) error {
	sf, err := c.DownloadServer(v, cfg.Dest)
	if err != nil {
		err = fmt.Errorf("getting version info to install: %w", err)
		return err
	}

	log.Printf("server file: %s", sf)
	return err
}

// DownloadServer downloads server file
func (c *Client) DownloadServer(v VersionInfoResponse, dest string) (string, error) {
	destFile := filepath.Join(dest, utils.GetFileName(v.Downloads.Server.URL))
	if err := utils.DownloadFile(c.c, v.Downloads.Server.URL, destFile); err != nil {
		err = fmt.Errorf("getting version info: %w", err)
		return "", err
	}

	return destFile, nil
}

// StartScript generates the start script
func (c *Client) StartScript(s StartupParams, dest string) error {
	destFile := filepath.Join(dest, "start.sh")

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server dest file: %w", err)
		return err
	}

	scp, err := s.ToScript()
	if err != nil {
		err = fmt.Errorf("generating start script content: %w", err)
		return err
	}

	if _, err := f.Write([]byte(scp)); err != nil {
		err = fmt.Errorf("writing start script to file: %w", err)
	}
	return nil
}
