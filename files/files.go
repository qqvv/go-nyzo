package files

import (
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var nyzoPath = AppDataDir()

var mu sync.RWMutex

func AppDataDir() string {
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	name, nameUpper := "nyzo", "Nyzo"

	var path string
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("LOCALAPPDATA")
		if appData == "" {
			appData = os.Getenv("APPDATA")
		}
		if appData != "" {
			path = filepath.Join(appData, nameUpper)
		}
	case "darwin":
		if homeDir != "" {
			path = filepath.Join(homeDir, "Library",
				"Application Support", nameUpper)
		}
	case "plan9":
		if homeDir != "" {
			path = filepath.Join(homeDir, name)
		}
	default:
		if homeDir != "" {
			path = filepath.Join(homeDir, "."+name)
		}
	}

	if path == "" {
		path = filepath.Join("." + name)
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(path, 0777)
			if err != nil {
				panic(err)
			}
		}
	}

	return path
}

func Delete(filename string) error {
	mu.Lock()
	defer mu.Unlock()

	// check if file exists first
	if _, err := os.Stat(filepath.Join(nyzoPath, filename)); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(nyzoPath, filename)); err != nil {
		return err
	}
	return nil
}

func Write(filename string, data []byte) error {
	mu.Lock()
	defer mu.Unlock()

	// first create a temporary file
	f, err := os.Create(filepath.Join(nyzoPath, filename+".tmp"))
	if err != nil {
		return err
	}

	n, err := f.Write(data)
	if err != nil {
		return err
	}
	if n < len(data) {
		return io.ErrShortWrite
	}
	if err = f.Close(); err != nil {
		return err
	}

	err = os.Rename(
		filepath.Join(nyzoPath, filename+".tmp"),
		filepath.Join(nyzoPath, filename),
	)
	return err
}

func ReadString(filename string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	b, err := ioutil.ReadFile(filepath.Join(nyzoPath, filename))
	if err != nil {
		return "", err
	}
	s := strings.TrimSuffix(string(b), "\n")
	return s, nil
}

func ReadBytes(filename string) ([]byte, error) {
	mu.RLock()
	defer mu.RUnlock()

	b, err := ioutil.ReadFile(filepath.Join(nyzoPath, filename))
	if err != nil {
		return b, err
	}
	return b, nil
}

func Exists(filename string) bool {
	if _, err := os.Stat(filepath.Join(nyzoPath, filename)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
