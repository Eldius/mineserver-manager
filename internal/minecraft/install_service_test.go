package minecraft

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/installer"
	"github.com/eldius/mineserver-manager/internal/minecraft/config"
	"github.com/eldius/mineserver-manager/internal/model"
	"github.com/eldius/mineserver-manager/internal/provisioner"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
	"time"
)

type mockDownloader struct {
	mock.Mock
}

func (m *mockDownloader) DownloadServer(ctx context.Context, url, sha1, dest string) (string, error) {
	args := m.Called(ctx, url, sha1, dest)
	return args.String(0), args.Error(1)
}

type mockRuntimeManager struct {
	mock.Mock
}

func (m *mockRuntimeManager) InstallJava(ctx context.Context, dest string, version int, arch, osName string) (string, error) {
	args := m.Called(ctx, dest, version, arch, osName)
	return args.String(0), args.Error(1)
}

type mockProvisioner struct {
	mock.Mock
}

func (m *mockProvisioner) CreateServerProperties(dest string, props *model.ServerProperties) error {
	args := m.Called(dest, props)
	return args.Error(0)
}

func (m *mockProvisioner) CreateStartScript(dest string, opts ...provisioner.StartupOption) error {
	// Variadic arguments with mock.Anything are tricky.
	// We'll use a helper to match them or just accept anything.
	args := m.Called(dest, opts)
	return args.Error(0)
}

func (m *mockProvisioner) CreateStopScript(dest string) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *mockProvisioner) CreateLoggingConfig(dest string, logfileDestDir string) error {
	args := m.Called(dest, logfileDestDir)
	return args.Error(0)
}

func (m *mockProvisioner) CreateEula(dest string, eula *model.Eula) error {
	args := m.Called(dest, eula)
	return args.Error(0)
}

type mockFlavor struct {
	mock.Mock
}

func (m *mockFlavor) Name() model.MineFlavour {
	args := m.Called()
	return args.Get(0).(model.MineFlavour)
}

func (m *mockFlavor) ListVersions(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockFlavor) GetVersionInfo(ctx context.Context, version string) (*installer.FlavorVersionInfo, error) {
	args := m.Called(ctx, version)
	return args.Get(0).(*installer.FlavorVersionInfo), args.Error(1)
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) SaveInstance(ctx context.Context, i *model.Instance) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *mockRepository) GetInstance(ctx context.Context, id string) (*model.Instance, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Instance), args.Error(1)
}

func (m *mockRepository) ListInstances(ctx context.Context) ([]model.Instance, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Instance), args.Error(1)
}

func (m *mockRepository) DeleteInstance(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestInstaller_Install(t *testing.T) {
	defer gock.Off()
	t.Run("should install vanilla server successfully", func(t *testing.T) {
		ctx := context.Background()
		dest, _ := os.MkdirTemp(os.TempDir(), "mine-install-test-*")
		defer func() {
			_ = os.RemoveAll(dest)
		}()

		md := new(mockDownloader)
		mr := new(mockRuntimeManager)
		mp := new(mockProvisioner)
		mf := new(mockFlavor)
		mrepo := new(mockRepository)

		info := &installer.FlavorVersionInfo{
			Version:     "1.20",
			DownloadURL: "https://example.com/server.jar",
			SHA1:        "abc",
			JavaVersion: 17,
		}

		mf.On("GetVersionInfo", mock.Anything, "1.20").Return(info, nil)
		mf.On("Name").Return(model.MineFlavourVanilla)

		md.On("DownloadServer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Sprintf("%s/server.jar", dest), nil)
		mr.On("InstallJava", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Sprintf("%s/java/jdk", dest), nil)
		mp.On("CreateServerProperties", mock.Anything, mock.Anything).Return(nil)
		// Fix variadic mock call
		mp.On("CreateStartScript", mock.Anything, mock.Anything).Return(nil)
		mp.On("CreateStopScript", mock.Anything).Return(nil)
		mp.On("CreateEula", mock.Anything, mock.Anything).Return(nil)
		mrepo.On("SaveInstance", mock.Anything, mock.Anything).Return(nil)

		s := NewInstallService(
			WithTimeout(5*time.Second),
			WithDownloader(md),
			WithRuntimeManager(mr),
			WithProvisioner(mp),
			WithFlavor(mf),
			WithRepository(mrepo),
		)

		err := s.Install(ctx,
			config.WithVersion("1.20"),
			config.ToDestinationFolder(dest),
		)

		assert.Nil(t, err)
		md.AssertExpectations(t)
		mr.AssertExpectations(t)
		mp.AssertExpectations(t)
		mf.AssertExpectations(t)
		mrepo.AssertExpectations(t)
	})
}
