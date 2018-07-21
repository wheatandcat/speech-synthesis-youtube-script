package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/wheatandcat/speech-synthesis-youtube-script/sound"
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

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
		panic(err)
	}

}
