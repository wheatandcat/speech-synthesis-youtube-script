package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli"

	"github.com/wheatandcat/speech-synthesis-youtube-script/movie"
	"github.com/wheatandcat/speech-synthesis-youtube-script/sound"
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

	app := cli.NewApp()
	app.Name = "speech-synthesis-script"
	app.Usage = "make speech synthesis file"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "output file path",
		},
	}

	app.Action = func(c *cli.Context) error {
		i := `
		[
			{
				"body": "こんにちは"
			},
			{
				"body": "よろしく"
			}
		]
		`
		uuid := uuid.NewV4()

		var sis []sound.Item
		json.Unmarshal([]byte(i), &sis)

		s := sound.New(uuid.String())
		if err := s.Makes(sis); err != nil {
			return err
		}

		m := movie.New(s.ID, "Title")
		if err := m.Make(); err != nil {
			return err
		}

		src := "tmp/" + s.ID + "/complex.mp4"
		dist := "tmp/" + s.ID + ".mp4"

		cmd := "cp " + src + " " + dist

		if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
			return err
		}

		if err := exec.Command("sh", "-c", "rm -rf tmp/"+s.ID).Run(); err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
