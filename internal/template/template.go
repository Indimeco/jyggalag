package template

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var templateValues = map[string]func() string{
	"CURRENT_DATE":  getCurrentDate,
	"CURRENT_YEAR":  getCurrentYear,
	"CURRENT_MONTH": getCurrentMonth,
}

const cursorTemplate = "$0"

func OpenEditor(editor string, path string) error {
	const vimToCursor = "-c %s/" + cursorTemplate + "//ge"
	cmd := exec.Command(editor, path, vimToCursor)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to open editor %w", err)
	}
	return nil
}

func getCurrentDate() string {
	time := time.Now()
	return time.Format("02")
}

func getCurrentYear() string {
	time := time.Now()
	return time.Format("2006")
}

func getCurrentMonth() string {
	time := time.Now()
	return time.Format("01")
}

func CopyTemplate(templatePath string, destination string) error {
	t, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("Could not read template %w", err)
	}

	// File doesn't exist, so we create it
	_, err = os.Stat(destination)
	if err != nil {
		contents := string(t[:])
		for key, value := range templateValues {
			contents = strings.ReplaceAll(contents, "${"+key+"}", value())
		}
		err = os.WriteFile(destination, []byte(contents), 0644)
		if err != nil {
			return fmt.Errorf("Failed to write new template: %w", err)
		}
	}
	// File already exists, so we just continue
	return nil
}
