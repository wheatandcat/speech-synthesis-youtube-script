package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

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

		var sis []sound.Item
		json.Unmarshal([]byte(i), &sis)

		s := sound.New()
		if err := s.Makes(sis); err != nil {
			return err
		}

		m := movie.New(s.ID)
		if err := m.Make(); err != nil {
			return err
		}

		src := "tmp/" + s.ID + "/complex.mp4"
		dist := "output.mp4"

		cmd := "cp " + src + " " + dist
		log.Println(cmd)

		if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
