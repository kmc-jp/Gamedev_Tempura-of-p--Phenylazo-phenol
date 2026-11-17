package audioassets

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioManager struct {
	AudioFS   fs.FS
	seLoader  *SELoader
	bgmLoader *BGMLoader
}

var audiocxtCreated bool = false

func NewAudioManager(AudioFS fs.FS) (*AudioManager, error) {
	if audiocxtCreated {
		return nil, fmt.Errorf("AudioManager は既に構築されています。重複して生成できません。")
	}
	audiocxtCreated = true
	SampleRate := 44100
	cxt := audio.NewContext(SampleRate)

	sl, err := NewSELoader(AudioFS, cxt)
	if err != nil {
		return nil, fmt.Errorf("NewSELoader() でエラーが発生しました。\n%w", err)
	}
	bl, err := NewBGMLoader(AudioFS, cxt)
	if err != nil {
		return nil, fmt.Errorf("NewBGMLoader() でエラーが発生しました。\n%w", err)
	}
	return &AudioManager{
		AudioFS:   AudioFS,
		seLoader:  sl,
		bgmLoader: bl,
	}, nil
}

func (am AudioManager) Load(data AudioData) (*audio.Player, error) {
	switch data.Type {
	case SE:
		pl, err := am.seLoader.Load(data)
		if err != nil {
			return nil, fmt.Errorf("am.seLoader.Load(data) でエラーが発生しました。\n%w", err)
		}
		return pl, nil
	case BGM:
		pl, err := am.bgmLoader.Load(data)
		if err != nil {
			return nil, fmt.Errorf("am.bgmLoader.Load(data) でエラーが発生しました。\n%w", err)
		}
		return pl, nil
	default:
		return nil, fmt.Errorf("data.Type が不正です。\n%d", data.Type)
	}
}
