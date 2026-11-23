package web

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eldius/mineserver-manager/internal/utils"
)

// Sandbox confines all filesystem access of the API to a configured home
// folder. Every path the API reads or writes is resolved through it, so
// installs, backups and temp files can never escape the home folder.
type Sandbox struct {
	home string
}

// NewSandbox creates the home folder (and its working subfolders) and
// returns a Sandbox rooted at it.
func NewSandbox(home string) (*Sandbox, error) {
	abs, err := utils.AbsolutePath(home)
	if err != nil {
		return nil, fmt.Errorf("resolving home folder: %w", err)
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf("creating home folder: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return nil, fmt.Errorf("resolving home folder symlinks: %w", err)
	}
	s := &Sandbox{home: resolved}
	for _, d := range []string{s.InstancesDir(), s.BackupsDir(), s.TempDir()} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return nil, fmt.Errorf("creating %s: %w", d, err)
		}
	}
	return s, nil
}

// Home returns the resolved absolute home folder.
func (s *Sandbox) Home() string { return s.home }

// InstancesDir returns the folder server instances are installed into.
func (s *Sandbox) InstancesDir() string { return filepath.Join(s.home, "instances") }

// BackupsDir returns the folder backup files are written to.
func (s *Sandbox) BackupsDir() string { return filepath.Join(s.home, "backups") }

// TempDir returns the folder to be used for temp files, instead of the
// OS default (which would live outside the home folder).
func (s *Sandbox) TempDir() string { return filepath.Join(s.home, "tmp") }

// Resolve maps rel to an absolute path inside the home folder. rel is
// always interpreted as home-relative (even when it looks absolute or
// contains ".."), and symlinks are resolved, so the result can never
// point outside the home folder. It returns ErrPathOutsideHome when a
// symlink inside home points out of it.
func (s *Sandbox) Resolve(rel string) (string, error) {
	p := filepath.Join(s.home, filepath.Clean(string(filepath.Separator)+rel))
	resolved, err := resolveExisting(p)
	if err != nil {
		return "", fmt.Errorf("resolving path %q: %w", rel, err)
	}
	if !s.Contains(resolved) {
		return "", fmt.Errorf("%w: %s", ErrPathOutsideHome, rel)
	}
	return resolved, nil
}

// Contains reports whether p is the home folder or lives inside it. p is
// expected to be absolute and symlink-resolved (see resolveExisting).
func (s *Sandbox) Contains(p string) bool {
	return p == s.home || strings.HasPrefix(p, s.home+string(filepath.Separator))
}

// Rel renders p relative to the home folder when inside it, otherwise it
// returns p unchanged.
func (s *Sandbox) Rel(p string) string {
	if rel, err := filepath.Rel(s.home, p); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return rel
	}
	return p
}

// resolveExisting resolves symlinks of the deepest existing ancestor of p
// and re-appends the non-existing remainder.
func resolveExisting(p string) (string, error) {
	var missing []string
	cur := p
	for {
		if _, err := os.Lstat(cur); err == nil {
			break
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("checking path %s: %w", cur, err)
		}
		missing = append([]string{filepath.Base(cur)}, missing...)
		parent := filepath.Dir(cur)
		if parent == cur {
			return "", fmt.Errorf("no existing ancestor for %s", p)
		}
		cur = parent
	}
	resolved, err := filepath.EvalSymlinks(cur)
	if err != nil {
		return "", fmt.Errorf("evaluating symlinks: %w", err)
	}
	for _, m := range missing {
		resolved = filepath.Join(resolved, m)
	}
	return resolved, nil
}
