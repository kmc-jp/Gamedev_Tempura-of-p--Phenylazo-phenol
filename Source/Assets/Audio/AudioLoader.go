package audioassets

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioLoader struct {
	cache map[string]*audio.Player // キーはパス
}

func NewAudioLoader() (*AudioLoader, error) {
	return &AudioLoader{
		cache: map[string]*audio.Player{},
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

	pl, _ := al.cache[data.Path]
	return pl, nil
}

func (al AudioLoader) fetch(path string) (*audio.Player, error) {
	return nil, fmt.Errorf("Not Implimented!")
}
