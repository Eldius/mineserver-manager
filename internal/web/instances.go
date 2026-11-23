package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	cfg "github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/minecraft"
	mcconfig "github.com/eldius/mineserver-manager/internal/minecraft/config"
	"github.com/eldius/mineserver-manager/internal/model"
	"github.com/eldius/mineserver-manager/internal/repository"
	"github.com/eldius/mineserver-manager/internal/utils"
)

const (
	rconDefaultPort   = 25575
	scriptRunTimeout  = 30 * time.Second
	instanceNameRegex = `^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`
)

type instanceHandler struct {
	repo repository.Repository // may be nil when the DB could not be opened
	sb   *Sandbox
}

type instanceResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	InstallDate time.Time `json:"install_date"`
	Running     bool      `json:"running"`
	Manageable  bool      `json:"manageable"`
}

// toResponse maps a stored instance to its API representation. Path is
// rendered relative to the home folder when inside it, and Manageable
// reports whether the API is allowed to touch it (instances installed by
// the CLI may live outside the API home folder).
func (h *instanceHandler) toResponse(i model.Instance) instanceResponse {
	res := instanceResponse{
		ID:          i.ID,
		Name:        i.Name,
		Path:        i.Path,
		InstallDate: i.InstallDate,
	}
	abs, err := utils.AbsolutePath(i.Path)
	if err != nil {
		return res
	}
	resolved, err := resolveExisting(abs)
	if err == nil {
		abs = resolved
	}
	res.Manageable = h.sb.Contains(abs)
	res.Path = h.sb.Rel(abs)
	if res.Manageable {
		res.Running = instanceRunning(abs)
	}
	return res
}

func (h *instanceHandler) requireRepo(w http.ResponseWriter) bool {
	if h.repo == nil {
		writeError(w, http.StatusServiceUnavailable, "instance repository unavailable")
		return false
	}
	return true
}

// findInstance looks an instance up by ID and confines its path to the
// home folder, writing the proper HTTP error when not possible.
func (h *instanceHandler) findInstance(w http.ResponseWriter, r *http.Request) (*model.Instance, string, bool) {
	inst, err := h.repo.GetInstance(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "instance not found")
		} else {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("getting instance: %v", err))
		}
		return nil, "", false
	}
	abs, err := utils.AbsolutePath(inst.Path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("resolving instance path: %v", err))
		return nil, "", false
	}
	resolved, err := resolveExisting(abs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("resolving instance path: %v", err))
		return nil, "", false
	}
	if !h.sb.Contains(resolved) {
		writeError(w, http.StatusForbidden, fmt.Sprintf("%v: instance lives outside the api home folder", ErrPathOutsideHome))
		return nil, "", false
	}
	return inst, resolved, true
}

// list answers GET /instances with every stored instance.
func (h *instanceHandler) list(w http.ResponseWriter, r *http.Request) {
	if !h.requireRepo(w) {
		return
	}
	instances, err := h.repo.ListInstances(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("listing instances: %v", err))
		return
	}
	res := make([]instanceResponse, 0, len(instances))
	for _, i := range instances {
		res = append(res, h.toResponse(i))
	}
	writeJSON(w, http.StatusOK, map[string]any{"instances": res})
}

type installRequest struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	MemoryLimit  string `json:"memory_limit"`
	Motd         string `json:"motd"`
	LevelName    string `json:"level_name"`
	Seed         string `json:"seed"`
	RconEnabled  bool   `json:"rcon_enabled"`
	RconPassword string `json:"rcon_password"`
}

// install answers POST /instances, installing a new server into
// <home>/instances/<name>. It runs synchronously and may take a few
// minutes (server jar + JRE downloads).
func (h *instanceHandler) install(w http.ResponseWriter, r *http.Request) {
	var req installRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if ok, _ := regexp.MatchString(instanceNameRegex, req.Name); !ok {
		writeError(w, http.StatusBadRequest, "invalid instance name: use letters, digits, '.', '_' or '-' (max 64 chars)")
		return
	}
	if req.Version == "" {
		req.Version = "latest"
	}
	if req.MemoryLimit == "" {
		req.MemoryLimit = "1g"
	}
	if req.RconEnabled && req.RconPassword == "" {
		writeError(w, http.StatusBadRequest, "rcon_password is required when rcon is enabled")
		return
	}

	dest, err := h.sb.Resolve(filepath.Join("instances", req.Name))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := os.Stat(dest); err == nil {
		writeError(w, http.StatusConflict, fmt.Sprintf("instance %q already exists", req.Name))
		return
	} else if !os.IsNotExist(err) {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("checking destination folder: %v", err))
		return
	}

	svcOpts := []minecraft.InstallServiceOpt{
		minecraft.WithTimeout(cfg.GetMinecraftApiTimeout()),
		minecraft.WithDownloadTimeout(cfg.GetMinecraftDownloadTimeout()),
	}
	if h.repo != nil {
		svcOpts = append(svcOpts, minecraft.WithRepository(h.repo))
	}
	svc := minecraft.NewInstallService(svcOpts...)

	instanceOpts := []mcconfig.InstanceOpt{
		mcconfig.WithVersion(req.Version),
		mcconfig.ToDestinationFolder(dest),
		mcconfig.WithMemoryLimit(req.MemoryLimit),
		mcconfig.WithHeadlessConfig(true),
		mcconfig.WithServerFlavour(string(mcconfig.VanillaServerSoftware)),
	}
	if req.Motd != "" {
		instanceOpts = append(instanceOpts, mcconfig.WithServerPropsMotd(req.Motd))
	}
	if req.LevelName != "" {
		instanceOpts = append(instanceOpts, mcconfig.WithServerPropsLevelName(req.LevelName))
	}
	if req.Seed != "" {
		instanceOpts = append(instanceOpts, mcconfig.WithServerPropsSeed(req.Seed))
	}
	if req.RconEnabled {
		instanceOpts = append(instanceOpts, mcconfig.WithServerPropsRconEnabled(rconDefaultPort, req.RconPassword))
	}

	if err := svc.Install(r.Context(), instanceOpts...); err != nil {
		slog.With("error", err, "instance", req.Name).ErrorContext(r.Context(), "install failed")
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("installing server: %v", err))
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"name": req.Name,
		"path": h.sb.Rel(dest),
	})
}

// start answers POST /instances/{id}/start running the instance's
// generated start.sh.
func (h *instanceHandler) start(w http.ResponseWriter, r *http.Request) {
	h.runScript(w, r, provisionerStartScript)
}

// stop answers POST /instances/{id}/stop running the instance's
// generated stop.sh.
func (h *instanceHandler) stop(w http.ResponseWriter, r *http.Request) {
	h.runScript(w, r, provisionerStopScript)
}

const (
	provisionerStartScript = "start.sh"
	provisionerStopScript  = "stop.sh"
)

func (h *instanceHandler) runScript(w http.ResponseWriter, r *http.Request, script string) {
	if !h.requireRepo(w) {
		return
	}
	_, dir, ok := h.findInstance(w, r)
	if !ok {
		return
	}
	scriptPath := filepath.Join(dir, script)
	if _, err := os.Stat(scriptPath); err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("instance has no %s", script))
		return
	}

	// java keeps running after start.sh exits, so its output must go to a
	// file instead of a pipe (a pipe would block cmd.Wait until the
	// server stops).
	logFile, err := os.OpenFile(filepath.Join(dir, "server.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("opening server log: %v", err))
		return
	}
	defer func() { _ = logFile.Close() }()

	ctx, cancel := context.WithTimeout(r.Context(), scriptRunTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", scriptPath)
	cmd.Dir = dir
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Run(); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("running %s: %v (see server.log in the instance folder)", script, err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type backupRequest struct {
	MaxBackups int `json:"max_backups"`
}

// backup answers POST /instances/{id}/backup writing a zip of the
// instance into <home>/backups.
func (h *instanceHandler) backup(w http.ResponseWriter, r *http.Request) {
	if !h.requireRepo(w) {
		return
	}
	var req backupRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
	}
	_, dir, ok := h.findInstance(w, r)
	if !ok {
		return
	}

	svc := minecraft.NewBackupService()
	info, err := svc.Backup(r.Context(), dir, h.sb.BackupsDir())
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("creating backup: %v", err))
		return
	}
	if req.MaxBackups > 0 {
		if err := svc.RolloverBackupFiles(r.Context(), h.sb.BackupsDir(), info.Name, req.MaxBackups); err != nil {
			slog.With("error", err).WarnContext(r.Context(), "backup rollover failed")
		}
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"name":       info.Name,
		"path":       h.sb.Rel(info.Path),
		"created_at": info.Timestamp,
	})
}

// instanceRunning reports whether the instance's server.pid points to a
// live process.
func instanceRunning(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "server.pid"))
	if err != nil {
		return false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}
