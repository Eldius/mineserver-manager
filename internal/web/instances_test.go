package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eldius/mineserver-manager/internal/model"
	"github.com/eldius/mineserver-manager/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	instances map[string]*model.Instance
}

func (f *fakeRepo) SaveInstance(_ context.Context, i *model.Instance) error {
	f.instances[i.ID] = i
	return nil
}

func (f *fakeRepo) GetInstance(_ context.Context, id string) (*model.Instance, error) {
	if i, ok := f.instances[id]; ok {
		return i, nil
	}
	return nil, repository.ErrNotFound
}

func (f *fakeRepo) ListInstances(_ context.Context) ([]model.Instance, error) {
	res := make([]model.Instance, 0, len(f.instances))
	for _, i := range f.instances {
		res = append(res, *i)
	}
	return res, nil
}

func (f *fakeRepo) DeleteInstance(_ context.Context, id string) error {
	delete(f.instances, id)
	return nil
}

func (f *fakeRepo) Close() error { return nil }

func newTestHandler(t *testing.T) (*instanceHandler, *fakeRepo) {
	t.Helper()
	return &instanceHandler{
		repo: &fakeRepo{instances: map[string]*model.Instance{}},
		sb:   newTestSandbox(t),
	}, nil
}

func TestListInstances(t *testing.T) {
	h, _ := newTestHandler(t)
	repo := h.repo.(*fakeRepo)

	inside := filepath.Join(h.sb.InstancesDir(), "my-server")
	require.NoError(t, os.MkdirAll(inside, 0o755))
	outside := t.TempDir()

	repo.instances["1"] = &model.Instance{ID: "1", Name: "my-server", Path: inside}
	repo.instances["2"] = &model.Instance{ID: "2", Name: "elsewhere", Path: outside}

	rec := httptest.NewRecorder()
	h.list(rec, httptest.NewRequest(http.MethodGet, "/api/instances", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"my-server"`)
	assert.Contains(t, body, filepath.Join("instances", "my-server"))
	assert.Contains(t, body, `"manageable":true`)
	assert.Contains(t, body, `"manageable":false`)
}

func TestListInstancesWithoutRepo(t *testing.T) {
	h := &instanceHandler{sb: newTestSandbox(t)}

	rec := httptest.NewRecorder()
	h.list(rec, httptest.NewRequest(http.MethodGet, "/api/instances", nil))

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestStartStopScriptsRunInsideHome(t *testing.T) {
	h, _ := newTestHandler(t)
	repo := h.repo.(*fakeRepo)

	dir := filepath.Join(h.sb.InstancesDir(), "scripted")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	marker := filepath.Join(dir, "started.marker")
	stopped := filepath.Join(dir, "stopped.marker")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "start.sh"),
		[]byte("#!/bin/bash\ntouch \""+marker+"\"\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "stop.sh"),
		[]byte("#!/bin/bash\ntouch \""+stopped+"\"\n"), 0o755))
	repo.instances["1"] = &model.Instance{ID: "1", Name: "scripted", Path: dir}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/instances/1/start", nil)
	req.SetPathValue("id", "1")
	h.start(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.FileExists(t, marker)

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/instances/1/stop", nil)
	req.SetPathValue("id", "1")
	h.stop(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.FileExists(t, stopped)
}

func TestRunScriptRejectsInstanceOutsideHome(t *testing.T) {
	h, _ := newTestHandler(t)
	repo := h.repo.(*fakeRepo)

	outside := t.TempDir()
	repo.instances["1"] = &model.Instance{ID: "1", Name: "elsewhere", Path: outside}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/instances/1/start", nil)
	req.SetPathValue("id", "1")
	h.start(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "outside")
}

func TestRunScriptUnknownInstance(t *testing.T) {
	h, _ := newTestHandler(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/instances/nope/start", nil)
	req.SetPathValue("id", "nope")
	h.start(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestInstallValidatesInput(t *testing.T) {
	h, _ := newTestHandler(t)

	tests := []struct {
		name string
		body string
		want int
	}{
		{"invalid json", `{`, http.StatusBadRequest},
		{"invalid name", `{"name":"../evil"}`, http.StatusBadRequest},
		{"path-ish name", `{"name":"a/b"}`, http.StatusBadRequest},
		{"rcon without password", `{"name":"ok-name","rcon_enabled":true}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/instances", strings.NewReader(tt.body))
			h.install(rec, req)
			assert.Equal(t, tt.want, rec.Code)
		})
	}
}

func TestBackupWritesZipInsideBackupsDir(t *testing.T) {
	h, _ := newTestHandler(t)
	repo := h.repo.(*fakeRepo)

	dir := filepath.Join(h.sb.InstancesDir(), "bkp-me")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "server.properties"), []byte("motd=test"), 0o644))
	repo.instances["1"] = &model.Instance{ID: "1", Name: "bkp-me", Path: dir}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/instances/1/backup", nil)
	req.SetPathValue("id", "1")
	h.backup(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	assert.Contains(t, rec.Body.String(), filepath.Join("backups", "bkp-me_"))

	entries, err := os.ReadDir(h.sb.BackupsDir())
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Contains(t, entries[0].Name(), "_backup.zip")
}
