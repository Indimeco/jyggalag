package state

import (
	"bufio"
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

func ReadRecent() ([]string, error) {
	stateDir, err := createState()
	if err != nil {
		log.Printf("Warning: Failed to create recents file %v", err)
	}

	var lines []string
	file, err := os.Open(path.Join(stateDir, "recents"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if len(lines) >= 5 {
			break
		}
		if !contains(lines, t) {
			lines = append(lines, t)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines, nil
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

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
