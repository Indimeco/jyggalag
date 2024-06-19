package state

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
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
	if _, err := f.Write([]byte(fmt.Sprintf("%v\n", recent))); err != nil {
		log.Printf("Warning: Failed to write recents file %v", err)
	}
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
