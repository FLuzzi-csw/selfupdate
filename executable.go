package selfupdate

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fynelabs/selfupdate/internal/osext"
)

var (
	exePath           string
	defaultOldExePath string
	newExePath        string
	exeErr            error
	once              sync.Once
)

func ExecutableRealPath() (string, error) {
	if LoadPath() != nil {
		return "", exeErr
	}
	return exePath, nil
}

func ExecutableDefaultOldPath() (string, error) {
	if LoadPath() != nil {
		return "", exeErr
	}
	return defaultOldExePath, nil
}

func ExecutableNewPath() (string, error) {
	if LoadPath() != nil {
		return "", exeErr
	}
	return newExePath, nil
}

func LoadPath() error {
	once.Do(func() {
		exePath, defaultOldExePath, newExePath, exeErr = loadPath()
	})
	return exeErr
}

func loadPath() (string, string, string, error) {
	exePath, err := getExecutableRealPath()
	if err != nil {
		return "", "", "", err
	}
	// get the directory the executable exists in
	updateDir := filepath.Dir(exePath)
	filename := filepath.Base(exePath)

	// get file paths to new and old executable file paths
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))
	return exePath, oldPath, newPath, nil
}

func lastModifiedExecutable() (time.Time, error) {
	exe, err := ExecutableRealPath()
	if err != nil {
		return time.Time{}, err
	}

	fi, err := os.Stat(exe)
	if err != nil {
		return time.Time{}, err
	}

	return fi.ModTime(), nil
}

func getExecutableRealPath() (string, error) {
	exe, err := osext.Executable()
	if err != nil {
		return "", err
	}

	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}

	return exe, nil
}
