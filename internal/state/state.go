package state

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/user"
	"path"
	"slices"
	"strings"
)

func WriteRecent(recent string) {
	stateDir, err := createState()
	if err != nil {
		log.Printf("Warning: Failed to create recents file %v", err)
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(path.Join(stateDir, "recents"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Printf("Warning: Failed to open recents file %v", err)
	}
	_, err = f.Write([]byte(fmt.Sprintf("%v\n", recent)))
	if err != nil {
		log.Printf("Warning: Failed to write recents file %v", err)
	}
}

func ReadRecent() ([]string, error) {
	stateDir, err := createState()
	if err != nil {
		log.Printf("Warning: Failed to create recents file %v", err)
	}

	file, err := os.Open(path.Join(stateDir, "recents"))
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Failed to stat file: %w", err)
	}
	fileSize := stat.Size()
	bufSize := int(math.Min(float64(fileSize), float64(500)))
	start := stat.Size() - int64(bufSize)
	buf := make([]byte, bufSize)
	_, err = file.ReadAt(buf, start)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}
	lines := strings.Split(string(buf), "\n")
	slices.Reverse(lines)
	result := make([]string, 0, 5)
	for _, line := range lines {
		if len(result) >= 5 {
			break
		}
		if !slices.Contains(result, line) {
			result = append(result, line)
		}
	}

	return result, nil
}

func createState() (string, error) {
	userDir, err := user.Current()
	if err != nil {
		return "", err
	}

	stateDir := path.Join(userDir.HomeDir, ".local/jyggalag")

	err = os.MkdirAll(stateDir, 744)
	if err != nil {
		return "", err
	}

	return stateDir, nil
}
