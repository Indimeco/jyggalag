package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/indimeco/jyggalag/internal/config"
	"github.com/indimeco/jyggalag/internal/recent"
	"github.com/indimeco/jyggalag/internal/state"
	"github.com/indimeco/jyggalag/internal/stats"
	"github.com/indimeco/jyggalag/internal/template"
	"github.com/indimeco/jyggalag/internal/timestr"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "jyggalag",
		Usage: "manage notes somewhere in the plane of madness",
		Commands: []*cli.Command{
			{
				Name:    "config_note_dir",
				Aliases: []string{},
				Usage:   "set the notes dir",
				Action: func(cCtx *cli.Context) error {
					err := config.SetNotesDir(cCtx.Args().First())
					if err != nil {
						return err
					}

					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					fmt.Printf("Set notes dir to %v", c.NotesDir)
					return nil
				},
			},
			{
				Name:    "config_editor",
				Aliases: []string{},
				Usage:   "set the editor",
				Action: func(cCtx *cli.Context) error {
					err := config.SetEditor(cCtx.Args().First())
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "new_journal",
				Aliases: []string{"nj"},
				Usage:   "create a new journal entry",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					journalName := timestr.CanonicalDateString()
					journalPath := filepath.Join(c.NotesDir, "journal", journalName+".md")

					return createAndOpen(journalPath, "templates/journal.md")
				},
			},
			{
				Name:    "new_daybook",
				Aliases: []string{"nd"},
				Usage:   "create a new daybook entry",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					daybookName := timestr.CanonicalDateString()
					daybookPath := filepath.Join(c.NotesDir, "journal", daybookName+".md")

					return createAndOpen(daybookPath, "templates/daybook.md")
				},
			},
			{
				Name:    "new_zettelkasten",
				Aliases: []string{"nz"},
				Usage:   "create a new zettelkasten",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					zettelName := cCtx.Args().First()
					if zettelName == "" {
						return errors.New("expected a name as first arg")
					}
					zettelDir := filepath.Join(c.NotesDir, "zettelkasten")
					zettelId, err := getNextZettelID(zettelDir)
					if err != nil {
						return err
					}

					zettelPath := filepath.Join(zettelDir, fmt.Sprintf("[%d] %v.md", zettelId, zettelName))
					return createAndOpen(zettelPath, "templates/zettelkasten.md")
				},
			},
			{
				Name:    "new_composition",
				Aliases: []string{"np"},
				Usage:   "create a new composition",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					baseName := timestr.CanonicalDateString()
					compositionDir := filepath.Join(c.NotesDir, "composition", timestr.CurrentYear())

					err = os.MkdirAll(compositionDir, 0777)
					if err != nil {
						return fmt.Errorf("Could not create composition directory %v: %w", compositionDir, err)
					}

					compositionId, err := strconv.Atoi(cCtx.Args().First())
					if err != nil {
						compositionId = 0
					}
					var compositionName string
					if compositionId > 0 {
						compositionName = fmt.Sprintf("%v-%v.md", baseName, compositionId)
					} else {
						compositionName = fmt.Sprintf("%v.md", baseName)
					}

					compositionPath := filepath.Join(compositionDir, compositionName)
					return createAndOpen(compositionPath, "templates/composition.md")
				},
			},
			{
				Name:    "new_reflection",
				Aliases: []string{"nr"},
				Usage:   "create a new reflection",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					zettelDir := filepath.Join(c.NotesDir, "zettelkasten")
					zettelId, err := getNextZettelID(zettelDir)
					if err != nil {
						return err
					}

					reflectionName := fmt.Sprintf("[%v] Reflection-%v", zettelId, timestr.CanonicalDateString())
					zettelPath := filepath.Join(zettelDir, reflectionName)

					return createAndOpen(zettelPath, "templates/reflection.md")
				},
			},
			{
				Name:    "open_recent",
				Aliases: []string{"or"},
				Usage:   "open a recent note",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					recentNote, err := recent.SelectRecent()
					if err != nil {
						return err
					}
					if recentNote == "" {
						return nil
					}

					state.WriteRecent(recentNote)
					return template.OpenEditor(c.Editor, recentNote)
				},
			},
			{
				Name:      "log_view",
				Aliases:   []string{"lv"},
				Args:      true,
				ArgsUsage: "log_view [log_name]",
				Usage:     "open the log",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					logName := cCtx.Args().Get(0)
					if logName == "" {
						return errors.New("Expected a log name")
					}

					logPath := filepath.Join(c.NotesDir, logName+".md")
					return template.OpenEditor(c.Editor, logPath)
				},
			},
			{
				Name:      "log",
				Aliases:   []string{"l"},
				Usage:     "add to the log",
				Args:      true,
				ArgsUsage: "log [log_name] [message]",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					logName := cCtx.Args().Get(0)
					if logName == "" {
						return errors.New("Expected a log name")
					}

					logPath := filepath.Join(c.NotesDir, logName+".md")
					message := strings.Trim(fmt.Sprint(strings.Join(cCtx.Args().Slice()[1:], " ")), " ")
					if message == "" {
						return errors.New("Expected a log message")
					}
					dateStr := timestr.CanonicalDateString()

					f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
					defer f.Close()
					if err != nil {
						return fmt.Errorf("Warning: Failed to open log file %v", err)
					}
					_, err = f.Write([]byte(fmt.Sprintf("%v: %v\n", dateStr, message)))
					if err != nil {
						return fmt.Errorf("Warning: Failed to write log file %v", err)
					}
					return nil
				},
			},
			{
				Name:      "stats",
				Usage:     "show stats for a directory",
				Args:      true,
				ArgsUsage: "stats [dir]",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					dir := filepath.Join(c.NotesDir, cCtx.Args().First())
					result, err := stats.TopWords(dir, 50)
					if err != nil {
						return fmt.Errorf("Failed to get stats: %w", err)
					}

					fmt.Print(result)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func createAndOpen(destination string, templatePath string) error {
	c, err := config.LoadConfig()
	if err != nil {
		return err
	}

	err = template.CopyTemplate(templatePath, destination)
	if err != nil {
		return err
	}

	state.WriteRecent(destination)
	err = template.OpenEditor(c.Editor, destination)
	return err
}

func getNextZettelID(zettelDir string) (int, error) {
	zettelIdRegex := regexp.MustCompile(`^\[(\d+)\]`)
	zettelId, err := template.GetLastIdInDir(zettelDir, zettelIdRegex)
	zettelId = zettelId + 1
	if err != nil {
		return 0, fmt.Errorf("Could not get new zettel id: %w", err)
	}
	return zettelId + 1, nil
}
