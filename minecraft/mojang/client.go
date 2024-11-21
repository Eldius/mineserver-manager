package mojang

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/eldius/mineserver-manager/utils"
	"net/http"
	"time"
)

type Client interface {
	// ListVersions lists all available versions
	ListVersions(ctx context.Context) (*VersionsResponse, error)
	// GetVersionInfo gets a specific version info
	GetVersionInfo(ctx context.Context, v Version) (*VersionInfoResponse, error)
	// GetUsersInfo fetch users identification
	GetUsersInfo(users ...string) (UserIDResponse, error)
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
func (c *apiClient) ListVersions(ctx context.Context) (*VersionsResponse, error) {
	client := c.httpClient()
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, VersionsURL, nil)
	if err != nil {
		err = fmt.Errorf("creating mojang query request instance: %w", err)
		return nil, err
	}
	//res, err := mojang.Get(VersionsURL)
	res, err := client.Do(r)
	if err != nil {
		err = fmt.Errorf("getting available mojang: %w", err)
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var versions VersionsResponse
	if err = json.NewDecoder(res.Body).Decode(&versions); err != nil {
		err = fmt.Errorf("decoding available mojang response: %w", err)
		return nil, err
	}

	return &versions, nil
}

// GetVersionInfo gets a specific version info
func (c *apiClient) GetVersionInfo(ctx context.Context, v Version) (*VersionInfoResponse, error) {
	client := c.httpClient()
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, v.URL, nil)
	if err != nil {
		err = fmt.Errorf("getting version info for '%s': %w", v.ID, err)
		return nil, err
	}
	res, err := client.Do(r)
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

func (c *apiClient) GetUsersInfo(users ...string) (UserIDResponse, error) {
	b, err := json.Marshal(users)
	if err != nil {
		err = fmt.Errorf("marshalling users info: %w", err)
		return nil, err
	}
	client := c.httpClient()
	buff := bytes.NewBuffer(b)
	res, err := client.Post(UsersInfoBulkURL, "application/json", buff)
	if err != nil {
		err = fmt.Errorf("getting users info: %w", err)
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var response UserIDResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		err = fmt.Errorf("decoding users info response: %w", err)
		return nil, err
	}

	return response, nil
}

func (c *apiClient) httpClient() http.Client {
	return utils.HTTPClient(c.cfg.Timeout)
}

func WithTimeout(d time.Duration) ClientOpt {
	return func(cfg *ClientConfig) *ClientConfig {
		cfg.Timeout = d
		return cfg
	}
}
