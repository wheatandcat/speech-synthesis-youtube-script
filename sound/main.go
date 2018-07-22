package sound

import (
	"fmt"
	"os"
	"os/exec"

	uuid "github.com/satori/go.uuid"
)

type Sound struct {
	ID string
}

type Item struct {
	Body string `json:"body"`
}

func New() *Sound {
	uuid := uuid.NewV4()

	s := &Sound{
		ID: uuid.String(),
	}

	return s
}

func (s *Sound) Makes(is []Item) error {
	if err := os.Mkdir("tmp/"+s.ID, 0777); err != nil {
		return err
	}

	if err := os.Mkdir("tmp/"+s.ID+"/words", 0777); err != nil {
		return err
	}

	for i, v := range is {
		if err := s.make(v, fileName(i)); err != nil {
			return err
		}
	}

	if err := s.join(is); err != nil {
		return err
	}

	if err := s.toMp3(); err != nil {
		return err
	}

	if err := s.toMp4a(); err != nil {
		return err
	}

	if err := s.remove(); err != nil {
		return err
	}

	return nil
}

func (s *Sound) make(i Item, f string) error {
	cmd := "echo " + i.Body + " | docker run -i --rm u6kapps/open_jtalk > tmp/" + s.ID + "/words/" + f + ".wav"

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}

func fileName(index int) string {
	return fmt.Sprintf("%04d", index)
}

func (s *Sound) join(is []Item) error {
	cmd := "docker run -v $(pwd):/work --rm bigpapoo/sox mysox"

	for i := range is {
		cmd += " tmp/" + s.ID + "/words/" + fileName(i) + ".wav"
	}

	cmd += " tmp/" + s.ID + "/join.wav"

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}

func (s *Sound) toMp3() error {
	src := " tmp/" + s.ID + "/join.wav"
	dist := " tmp/" + s.ID + "/join.mp3"

	cmd := "docker run -v $(pwd):/mp3 renyufu/lame -V2 " + src + " " + dist

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}

func (s *Sound) toMp4a() error {
	src := " tmp/" + s.ID + "/join.mp3"
	dist := " tmp/" + s.ID + "/join.m4a"

	cmd := "docker run -v $(pwd):/tmp/ffmpeg opencoconut/ffmpeg -y -i " + src + " -vn -ac 2 -vol 256 -ab 112k" + dist

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}

func (s *Sound) remove() error {
	r := "tmp/" + s.ID

	if err := os.Remove(r + "/join.wav"); err != nil {
		return err
	}

	if err := os.Remove(r + "/join.mp3"); err != nil {
		return err
	}

	if err := os.RemoveAll(r + "/words"); err != nil {
		return err
	}

	return nil
}
