package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Parse parses environment variables into a struct.
func Parse(cfg interface{}) error {
	return env.Parse(cfg)
}

// Load loads environment variables from a file. If no filename is given, it will try to load from
// ".env". If it failed to find the file, it will try to find the file in the parent directory
// recursively.
func Load(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	for _, filename := range filenames {
		if path, ok := find(filename); ok {
			if err := godotenv.Load(path); err != nil {
				return fmt.Errorf("failed to load env file %s: %w", path, err)
			}
		} else {
			return fmt.Errorf("failed to find env file %s", filename)
		}
	}

	return nil
}

func find(filename string) (string, bool) {
	workDir, _ := os.Getwd()
	for {
		if _, err := os.Stat(workDir + "/" + filename); err == nil {
			return workDir + "/" + filename, true
		}
		if workDir == "/" {
			break
		}
		workDir = filepath.Dir(workDir)
	}
	return "", false
}
