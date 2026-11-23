package web

import (
	"fmt"
	"net/http"

	"github.com/eldius/mineserver-manager/internal/buildinfo"
	cfg "github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/installer"
	"github.com/eldius/mineserver-manager/internal/model"
	"github.com/eldius/mineserver-manager/internal/mojang"
)

type versionsHandler struct{}

// version answers GET /version with the API binary build info.
func (versionsHandler) version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, model.CliVersion{
		Version:   buildinfo.Version,
		Commit:    buildinfo.Commit,
		BuildDate: buildinfo.BuildDate,
	})
}

// list answers GET /versions with the Minecraft versions available for
// installation (fetched live from Mojang).
func (versionsHandler) list(w http.ResponseWriter, r *http.Request) {
	flavor := installer.NewVanillaFlavor(mojang.NewClient(mojang.WithTimeout(cfg.GetMinecraftApiTimeout())))
	versions, err := flavor.ListVersions(r.Context())
	if err != nil {
		writeError(w, http.StatusBadGateway, fmt.Sprintf("listing available versions: %v", err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"versions": versions})
}
