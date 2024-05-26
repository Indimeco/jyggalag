package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
				Name:    "new_journal",
				Aliases: []string{"nj"},
				Usage:   "create a new day book entry",
				Action: func(cCtx *cli.Context) error {
					c, err := config.LoadConfig()
					if err != nil {
						return err
					}

					journalName := time.Now().Format("2006-01-02")
					journalPath := filepath.Join(c.NotesDir, "journal", journalName+".md")

					err = template.CopyTemplate("./templates/journal.md", journalPath)
					if err != nil {
						return fmt.Errorf("Could not copy template to %v: %e", journalPath, err)
					}

					// TODO make this come from config
					cmd := exec.Command("vim", journalPath)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					if err := cmd.Run(); err != nil {
						fmt.Println("Error: ", err)
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}