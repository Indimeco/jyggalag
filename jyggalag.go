package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/indimeco/jyggalag/internal/config"
	"github.com/indimeco/jyggalag/internal/template"
	"github.com/indimeco/jyggalag/internal/timestr"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
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

					journalName := timestr.GetCanonicalDateString()
					journalPath := filepath.Join(c.NotesDir, "journal", journalName+".md")

					err = template.CopyTemplate("./templates/journal.md", journalPath)
					if err != nil {
						return fmt.Errorf("Could not copy template to %v: %w", journalPath, err)
					}

					err = template.OpenEditor(c.Editor, journalPath)
					return nil
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

					journalName := timestr.GetCanonicalDateString()
					journalPath := filepath.Join(c.NotesDir, "journal", journalName+".md")

					err = template.CopyTemplate("./templates/daybook.md", journalPath)
					if err != nil {
						return fmt.Errorf("Could not copy template to %v: %w", journalPath, err)
					}

					err = template.OpenEditor(c.Editor, journalPath)
					return nil
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
					zettelDir := filepath.Join(c.NotesDir, "zettelkasten")
					zettelIdRegex := regexp.MustCompile(`^\[(\d+)\]`)
					zettelId, err := template.GetLastIdInDir(zettelDir, zettelIdRegex)
					zettelId = zettelId + 1
					if err != nil {
						return fmt.Errorf("Could not get new zettel id: %w", err)
					}

					zettelPath := filepath.Join(zettelDir, fmt.Sprintf("[%d] %v.md", zettelId, zettelName))

					err = template.CopyTemplate("./templates/zettelkasten.md", zettelPath)
					if err != nil {
						return err
					}

					err = template.OpenEditor(c.Editor, zettelPath)
					return nil
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

					baseName := timestr.GetCanonicalDateString()
					compositionDir := filepath.Join(c.NotesDir, "composition", timestr.GetCurrentYear())

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
					err = template.CopyTemplate("./templates/composition.md", compositionPath)
					if err != nil {
						return fmt.Errorf("Could not copy template to %v: %w", compositionPath, err)
					}

					err = template.OpenEditor(c.Editor, compositionPath)
					return nil
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
					zettelIdRegex := regexp.MustCompile(`^\[(\d+)\]`)
					zettelId, err := template.GetLastIdInDir(zettelDir, zettelIdRegex)
					zettelId = zettelId + 1
					if err != nil {
						return fmt.Errorf("Could not get new zettel id: %w", err)
					}

					reflectionName := fmt.Sprintf("[%v] Reflection-%v", zettelId, timestr.GetCanonicalDateString())
					zettelPath := filepath.Join(zettelDir, reflectionName)

					err = template.CopyTemplate("./templates/reflection.md", zettelPath)
					if err != nil {
						return err
					}

					err = template.OpenEditor(c.Editor, zettelPath)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
