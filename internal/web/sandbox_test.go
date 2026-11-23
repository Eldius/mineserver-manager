package web

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSandbox(t *testing.T) *Sandbox {
	t.Helper()
	sb, err := NewSandbox(filepath.Join(t.TempDir(), "home"))
	require.NoError(t, err)
	return sb
}

func TestSandboxResolveKeepsPathsInsideHome(t *testing.T) {
	sb := newTestSandbox(t)

	tests := []struct {
		name string
		rel  string
		want string
	}{
		{"plain relative", "instances/my-server", "instances/my-server"},
		{"dot segments are cleaned", "instances/../backups/b.zip", "backups/b.zip"},
		{"parent traversal is tamed", "../../etc/passwd", "etc/passwd"},
		{"absolute-looking path is home-relative", "/etc/passwd", "etc/passwd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sb.Resolve(tt.rel)
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(sb.Home(), tt.want), got)
			assert.True(t, sb.Contains(got))
		})
	}
}

func TestSandboxResolveRejectsSymlinkEscape(t *testing.T) {
	sb := newTestSandbox(t)

	outside := t.TempDir()
	require.NoError(t, os.Symlink(outside, filepath.Join(sb.Home(), "evil-link")))

	_, err := sb.Resolve("evil-link/secrets.txt")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPathOutsideHome)
}

func TestSandboxResolveFollowsSymlinkToInside(t *testing.T) {
	sb := newTestSandbox(t)

	require.NoError(t, os.Symlink("instances", filepath.Join(sb.Home(), "link")))
	got, err := sb.Resolve("link/my-server")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(sb.InstancesDir(), "my-server"), got)
}

func TestSandboxContainsAndRel(t *testing.T) {
	sb := newTestSandbox(t)

	assert.True(t, sb.Contains(sb.Home()))
	assert.True(t, sb.Contains(filepath.Join(sb.Home(), "instances", "a")))
	assert.False(t, sb.Contains(filepath.Dir(sb.Home())))
	assert.False(t, sb.Contains(sb.Home()+"-suffixed"))

	assert.Equal(t, filepath.Join("instances", "a"), sb.Rel(filepath.Join(sb.Home(), "instances", "a")))
	outside := filepath.Join(t.TempDir(), "x")
	assert.Equal(t, outside, sb.Rel(outside))
}

func TestNewSandboxCreatesWorkingDirs(t *testing.T) {
	sb := newTestSandbox(t)

	for _, d := range []string{sb.InstancesDir(), sb.BackupsDir(), sb.TempDir()} {
		st, err := os.Stat(d)
		require.NoError(t, err)
		assert.True(t, st.IsDir())
	}
	assert.NotEqual(t, os.TempDir(), sb.TempDir())
}

func TestResolveExistingWalksUpToExistingAncestor(t *testing.T) {
	sb := newTestSandbox(t)

	deep := filepath.Join(sb.Home(), "a", "b", "c")
	got, err := resolveExisting(deep)
	require.NoError(t, err)
	assert.Equal(t, deep, got)
}

func TestErrPathOutsideHomeIsWrapped(t *testing.T) {
	sb := newTestSandbox(t)
	outside := t.TempDir()
	require.NoError(t, os.Symlink(outside, filepath.Join(sb.Home(), "evil")))

	err := func() error {
		_, err := sb.Resolve("evil")
		return err
	}()
	assert.True(t, errors.Is(err, ErrPathOutsideHome))
}
