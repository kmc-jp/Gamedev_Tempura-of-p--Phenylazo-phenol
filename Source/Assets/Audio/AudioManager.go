package audioassets

import (
	"fmt"
	"io/fs"
)

type AudioManager struct {
	AudioFS   fs.FS
	mp3Loader *MP3Loader
	wavLoader *WAVLoader
}

func NewAudioManager(AudioFS fs.FS) (*AudioManager, error) {
	mp3, err := NewMP3Loader()
	if err != nil {
		return nil, fmt.Errorf("NewMP3Loader() でエラーが発生しました。\n%w", err)
	}
	wav, err := NewWAVLoader()
	if err != nil {
		return nil, fmt.Errorf("NewWAVLoader() でエラーが発生しました。\n%w", err)
	}
	return &AudioManager{
		AudioFS:   AudioFS,
		mp3Loader: mp3,
		wavLoader: wav,
	}, nil
}
