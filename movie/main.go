package movie

import (
	"os/exec"
)

type Movie struct {
	ID    string
	Title string
}

func New(id string, t string) *Movie {
	s := &Movie{
		ID:    id,
		Title: t,
	}

	return s
}

func (m *Movie) Make() error {
	if err := m.titleImage(); err != nil {
		return err
	}

	if err := m.imageToMp4(); err != nil {
		return err
	}

	if err := m.toMp4(); err != nil {
		return err
	}

	if err := m.toFilter(); err != nil {
		return err
	}

	return nil
}

func (m *Movie) titleImage() error {
	dist := "tmp/" + m.ID + "/title.jpg"
	cmd := `docker run -v $(pwd):/image acleancoder/imagemagick-full convert -size 512x288 xc:transparent -font /usr/share/fonts/truetype/liberation/LiberationSansNarrow-Bold.ttf -pointsize 72 -fill white -stroke black -draw "text 50,144 '` + m.Title + `'" image/` + dist

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}

func (m *Movie) imageToMp4() error {
	src := " tmp/" + m.ID + "/title.jpg"
	dist := " tmp/" + m.ID + "/title.mp4"
	cmd := "docker run -v $(pwd):/tmp/ffmpeg opencoconut/ffmpeg -r 15 -f image2 -i " + src + " -r 15 -an -vcodec libx264 -pix_fmt yuv420p " + dist

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}

func (m *Movie) toMp4() error {
	src := " tmp/" + m.ID + "/title.mp4"
	sound := " tmp/" + m.ID + "/join.m4a"
	dist := " tmp/" + m.ID + "/join.mp4"

	cmd := "docker run -v $(pwd):/tmp/ffmpeg opencoconut/ffmpeg -y -i " + src + " -i " + sound + " -map 0:0 -map 1:0 -movflags faststart -vcodec libx264 -acodec copy " + dist

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}

func (m *Movie) toFilter() error {
	filter := " static/movie/complex.mp4"
	src := " tmp/" + m.ID + "/join.mp4"
	dist := " tmp/" + m.ID + "/complex.mp4"

	cmd := "docker run -v $(pwd):/tmp/ffmpeg opencoconut/ffmpeg -y -i " + src + " -i " + filter + " -filter_complex 'concat=n=2:v=1:a=1' " + dist

	err := exec.Command("sh", "-c", cmd).Run()

	return err
}
