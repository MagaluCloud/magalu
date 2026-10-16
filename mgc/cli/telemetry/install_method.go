package telemetry

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TODO: habilitar InstallMethodScript quando o script de instalação (MagaluCloud/mgccli#77,
// curl | sh) for publicado. Ele instala em ~/.local/bin ou /usr/local/bin, e só pelo diretório
// um binário copiado à mão para lá também viraria script. Fazer o script deixar uma marca
// (ex.: arquivo ao lado do binário) para identificar esse método sem depender do caminho.
const (
	// InstallMethodScript   = "script"
	InstallMethodHomebrew = "homebrew"
	InstallMethodSystem   = "system"
	InstallMethodSnap     = "snap"
	InstallMethodManual   = "manual"
)

func ResolveExecutablePath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		return resolved
	}
	return exe
}

func DetectInstallMethod(executablePath string) string {
	if executablePath == "" {
		return InstallMethodManual
	}

	p := strings.ToLower(filepath.ToSlash(executablePath))
	p = strings.ReplaceAll(p, `\`, "/")
	dir := path.Dir(p)

	isHomebrew := strings.Contains(p, "/cellar/") || strings.Contains(p, "/homebrew/")
	isHomebrew = isHomebrew || strings.Contains(p, "/linuxbrew/")

	switch {
	case isHomebrew:
		return InstallMethodHomebrew
	case strings.HasPrefix(p, "/snap/"):
		return InstallMethodSnap
	// /usr/bin recebe os pacotes .deb/.rpm, mas também outros gerenciadores e cópias manuais.
	case dir == "/usr/bin":
		return InstallMethodSystem
	}

	return InstallMethodManual
}
