package audioassets

import (
	"fmt"
	"io/fs"
)

type AudioManager struct {
	AudioFS   fs.FS
	mp3Loader *AudioLoader
}

func NewAudioManager(AudioFS fs.FS) (*AudioManager, error) {
	mp3, err := NewAudioLoader()
	if err != nil {
		return nil, fmt.Errorf("NewAudioLoader() でエラーが発生しました。\n%w", err)
	}
	return &AudioManager{
		AudioFS:   AudioFS,
		mp3Loader: mp3,
	}, nil
}
