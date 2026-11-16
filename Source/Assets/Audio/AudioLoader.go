package audioassets

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioLoader struct {
	fs    fs.FS
	cache map[string][]byte // キーはパス
}

func NewAudioLoader(fs fs.FS) (*AudioLoader, error) {
	return &AudioLoader{
		fs:    fs,
		cache: map[string][]byte{},
	}, nil
}

func (al AudioLoader) Load(data AudioData) (*audio.Player, error) {
	_, ok := al.cache[data.Path]
	if !ok {
		pl, err := al.fetch(data.Path)
		if err != nil {
			return nil, fmt.Errorf("al.fetch(data.Path) でエラーが発生しました。\n%w", err)
		}
		al.cache[data.Path] = pl
	}

	ad, _ := al.cache[data.Path]
	pl, err := al.mkPlayer(ad)
	if err != nil {
		return nil, fmt.Errorf("al.mkPlayer(ad) でエラーが発生しました。\nPath : %s\n%w", data.Path, err)
	}
	return pl, nil
}

func (al AudioLoader) fetch(path string) ([]byte, error) {
	return nil, fmt.Errorf("Not Implimented!")
}

func (al AudioLoader) mkPlayer(AudioData []byte) (*audio.Player, error) {
	return nil, fmt.Errorf("Not Implimented!")
}
