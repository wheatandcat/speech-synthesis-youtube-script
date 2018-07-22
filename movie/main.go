package movie

import "os/exec"

type Movie struct {
	ID string
}

func New(id string) *Movie {
	s := &Movie{
		ID: id,
	}

	return s
}

func (m *Movie) Make() error {
	if err := m.toMp4(); err != nil {
		return err
	}

	if err := m.toFilter(); err != nil {
		return err
	}

	return nil
}

func (m *Movie) toMp4() error {
	src := " static/movie/input.mp4"
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
