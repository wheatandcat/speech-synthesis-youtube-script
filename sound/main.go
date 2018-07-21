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
	uuid, _ := uuid.NewV4()

	s := &Sound{
		ID: uuid.String(),
	}

	return s
}

func (s *Sound) Makes(is []Item) error {
	if err := os.Mkdir("tmp/"+s.ID, 0777); err != nil {
		return err
	}

	for i, v := range is {
		if err := s.make(v, fileName(i)); err != nil {
			return err
		}
	}

	return s.join(is)
}

func (s *Sound) make(i Item, f string) error {
	cmd := "echo " + i.Body + " | docker run -i --rm u6kapps/open_jtalk > tmp/" + s.ID + "/" + f + ".wav"

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}

func fileName(index int) string {
	return fmt.Sprintf("%04d", index)
}

func (s *Sound) join(is []Item) error {
	cmd := "docker run -v $(pwd):/work --rm bigpapoo/sox mysox"

	for i := range is {
		cmd += " tmp/" + s.ID + "/" + fileName(i) + ".wav"
	}

	cmd += " tmp/" + s.ID + "/join.wav"

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}
