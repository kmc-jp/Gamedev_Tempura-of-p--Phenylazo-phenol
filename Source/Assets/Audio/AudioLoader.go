package audioassets

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type SELoader struct {
	fs       fs.FS
	audiocxt *audio.Context
	cache    map[string][]byte // キーはパス
}

func NewSELoader(fs fs.FS, audiocxt *audio.Context) (*SELoader, error) {
	return &SELoader{
		fs:       fs,
		audiocxt: audiocxt,
		cache:    map[string][]byte{},
	}, nil
}

func (al SELoader) Load(data AudioData) (*audio.Player, error) {
	_, ok := al.cache[data.Path]
	if !ok {
		pl, err := al.fetch(data.Path)
		if err != nil {
			return nil, fmt.Errorf("al.fetch(data.Path) でエラーが発生しました。\n%w", err)
		}
		al.cache[data.Path] = pl
	}

	ad, _ := al.cache[data.Path]
	pl := al.mkPlayer(ad)
	return pl, nil
}

// デコードして []byte にする
func (al SELoader) fetch(path string) ([]byte, error) {
	return nil, fmt.Errorf("Not Implimented!")
}

// デコードしたデータを *audio.Player にする
// 不正なデータでも問答無用で変換してくれる
func (al SELoader) mkPlayer(AudioData []byte) *audio.Player {
	return al.audiocxt.NewPlayerFromBytes(AudioData)
}
