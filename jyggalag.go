package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/indimeco/jyggalag/internal/config"
	"github.com/indimeco/jyggalag/internal/template"
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

					journalName := time.Now().Format("2006-01-02")
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
					zettelId := 99
					zettelPath := filepath.Join(zettelDir, fmt.Sprintf("[%d]%v.md", zettelId, zettelName))

					err = template.CopyTemplate("./templates/default.md", zettelPath)
					if err != nil {
						return fmt.Errorf("Could not copy template to %v: %w", zettelPath, err)
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
