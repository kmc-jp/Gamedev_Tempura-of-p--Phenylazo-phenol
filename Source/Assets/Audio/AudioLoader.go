package audioassets

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioLoader struct {
	cache map[string]audio.Player
}

func NewAudioLoader() (*AudioLoader, error) {
	return &AudioLoader{
		cache: map[string]audio.Player{},
	}, nil
}
