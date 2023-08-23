package versions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	// ListVersions lists all available versions
	ListVersions() (*VersionsResponse, error)
	// GetVersionInfo gets a specific version info
	GetVersionInfo(v Version) (*VersionInfoResponse, error)
}

type ClientConfig struct {
	Timeout time.Duration
}

type ClientOpt func(config *ClientConfig) *ClientConfig

type apiClient struct {
	cfg ClientConfig
}

// NewClient creates a new client
func NewClient(configs ...ClientOpt) Client {
	cfg := &ClientConfig{
		Timeout: 1 * time.Second,
	}
	for _, c := range configs {
		c(cfg)
	}
	return &apiClient{
		cfg: *cfg,
	}
}

// ListVersions lists all available versions
func (c *apiClient) ListVersions() (*VersionsResponse, error) {
	client := c.httpClient()
	res, err := client.Get(VersionsURL)
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
func (c *apiClient) GetVersionInfo(v Version) (*VersionInfoResponse, error) {
	client := c.httpClient()
	res, err := client.Get(v.URL)
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

func (c *apiClient) httpClient() http.Client {
	return http.Client{Timeout: c.cfg.Timeout}
}

func WithTimeout(d time.Duration) ClientOpt {
	return func(cfg *ClientConfig) *ClientConfig {
		cfg.Timeout = d
		return cfg
	}
}
