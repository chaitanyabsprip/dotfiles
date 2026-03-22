package icon

import (
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrNoChildProcess = exec.ErrNotFound

func ResolveCommand(panePid string) (string, error) {
	pgrepCmd := exec.Command("pgrep", "-P", panePid)
	pgrepOut, err := pgrepCmd.Output()
	if err != nil {
		return "", ErrNoChildProcess
	}

	lines := strings.Split(strings.TrimSpace(string(pgrepOut)), "\n")
	if len(lines) == 0 || lines[len(lines)-1] == "" {
		return "", ErrNoChildProcess
	}

	lastChildPid := strings.TrimSpace(lines[len(lines)-1])
	if lastChildPid == "" {
		return "", ErrNoChildProcess
	}

	psCmd := exec.Command("ps", "-o", "args=", "-p", lastChildPid)
	psOut, err := psCmd.Output()
	if err != nil {
		return "", ErrNoChildProcess
	}

	raw := strings.TrimSpace(string(psOut))
	if raw == "" {
		return "", ErrNoChildProcess
	}

	return normalizeCommand(raw), nil
}

func normalizeCommand(raw string) string {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return ""
	}

	basename := filepath.Base(parts[0])
	if len(parts) == 1 {
		return basename
	}

	return basename + " " + strings.Join(parts[1:], " ")
}
