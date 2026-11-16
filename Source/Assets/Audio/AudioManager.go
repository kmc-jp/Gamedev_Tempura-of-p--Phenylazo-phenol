package audioassets

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioManager struct {
	AudioFS     fs.FS
	audioLoader *SELoader
}

func NewAudioManager(AudioFS fs.FS) (*AudioManager, error) {
	al, err := NewSELoader(AudioFS, nil)
	if err != nil {
		return nil, fmt.Errorf("NewAudioLoader() でエラーが発生しました。\n%w", err)
	}
	return &AudioManager{
		AudioFS:     AudioFS,
		audioLoader: al,
	}, nil
}

func (am AudioManager) Load(data AudioData) (*audio.Player, error) {
	pl, err := am.audioLoader.Load(data)
	if err != nil {
		return nil, fmt.Errorf("am.audioLoader.Load(data) でエラーが発生しました。\n%w", err)
	}
	return pl, nil
}
