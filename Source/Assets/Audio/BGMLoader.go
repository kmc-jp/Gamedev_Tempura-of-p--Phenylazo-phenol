package audioassets

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// BGMは長いのでキャッシュが重くなりすぎたら方式を変える
// 現在はバイナリ音源を全キャッシュ (SEと同じ)
type BGMLoader struct {
	fs       fs.FS
	audiocxt *audio.Context
	cache    map[string][]byte // キーはパス
}

func NewBGMLoader(fs fs.FS, audiocxt *audio.Context) *BGMLoader {
	return &BGMLoader{
		fs:       fs,
		audiocxt: audiocxt,
		cache:    map[string][]byte{},
	}
}

func (al BGMLoader) Load(data AudioData) (*audio.Player, error) {
	_, ok := al.cache[data.Path]
	if !ok {
		ad, err := al.fetch(data)
		if err != nil {
			return nil, fmt.Errorf("al.fetch(data.Path) でエラーが発生しました。\n%w", err)
		}
		al.cache[data.Path] = ad
	}

	ad, _ := al.cache[data.Path]
	pl := al.mkPlayer(ad)
	return pl, nil
}

// デコードして []byte にする
func (al BGMLoader) fetch(data AudioData) ([]byte, error) {
	rawfile, err := al.fs.Open(data.Path)
	if err != nil {
		return nil, err
	}
	defer rawfile.Close()

	var target []byte

	switch data.ext {
	case ".wav":
		stream, err := wav.DecodeWithSampleRate(al.audiocxt.SampleRate(), rawfile)
		if err != nil {
			return nil, fmt.Errorf("wavファイルの解析に失敗しました。\nPath : %s\n%w", data.Path, err)
		}
		target, err = io.ReadAll(stream)
		if err != nil {
			return nil, fmt.Errorf("wavファイルの解析結果が読み込めません。\nPath : %s\n%w", data.Path, err)
		}
	case ".mp3":
		stream, err := mp3.DecodeWithSampleRate(al.audiocxt.SampleRate(), rawfile)
		if err != nil {
			return nil, fmt.Errorf("mp3ファイルの解析に失敗しました。\nPath : %s\n%w", data.Path, err)
		}
		target, err = io.ReadAll(stream)
		if err != nil {
			return nil, fmt.Errorf("mp3ファイルの解析結果が読み込めません。\nPath : %s\n%w", data.Path, err)
		}
	default:
		return nil, fmt.Errorf("対応していないファイル形式です\nPath : %s", data.Path)
	}

	return target, nil
}

// デコードしたデータを *audio.Player にする
// 不正なデータでも問答無用で変換してくれる
func (al BGMLoader) mkPlayer(AudioData []byte) *audio.Player {
	return al.audiocxt.NewPlayerFromBytes(AudioData)
}
