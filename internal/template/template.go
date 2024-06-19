package template

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/indimeco/jyggalag/internal/timestr"
)

//go:embed templates
var templateFiles embed.FS

var templateValues = map[string]func() string{
	"CURRENT_DATE":  timestr.CurrentDate,
	"CURRENT_YEAR":  timestr.CurrentYear,
	"CURRENT_MONTH": timestr.CurrentMonth,
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

func getLastIdMatching(set []string, r *regexp.Regexp) (int, error) {
	var index = 0
	for _, v := range set {
		matches := r.FindStringSubmatch(v)
		if len(matches) < 2 {
			// no submatch
			continue
		}
		i, _ := strconv.Atoi(matches[1])
		if i > index {
			index = i
		}
	}
	return index, nil
}

func GetLastIdInDir(dir string, r *regexp.Regexp) (int, error) {
	t, err := os.ReadDir(dir)
	if err != nil {
		return -1, fmt.Errorf("Failed to read directory for finding last id %w", err)
	}
	var names []string
	for _, v := range t {
		names = append(names, v.Name())
	}

	return getLastIdMatching(names, r)
}

func CopyTemplate(templatePath string, destination string) error {
	t, err := templateFiles.ReadFile(templatePath)
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
			return fmt.Errorf("Failed to write new template to %v: %w", destination, err)
		}
	}
	// File already exists, so we just continue
	return nil
}
