package audioassets

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type SELoader struct {
	samplingRate int
	fs           fs.FS
	audiocxt     *audio.Context
	cache        map[string][]byte // キーはパス
}

func NewSELoader(fs fs.FS, audiocxt *audio.Context, samplingRate int) (*SELoader, error) {
	return &SELoader{
		samplingRate: samplingRate,
		fs:           fs,
		audiocxt:     audiocxt,
		cache:        map[string][]byte{},
	}, nil
}

func (al SELoader) Load(data AudioData) (*audio.Player, error) {
	_, ok := al.cache[data.Path]
	if !ok {
		pl, err := al.fetch(data)
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
func (al SELoader) fetch(data AudioData) ([]byte, error) {
	rawfile, err := al.fs.Open(data.Path)
	if err != nil {
		return nil, err
	}
	defer rawfile.Close()

	var target []byte

	switch data.ext {
	case ".wav":
		stream, err := wav.DecodeWithSampleRate(al.samplingRate, rawfile)
		if err != nil {
			return nil, fmt.Errorf("wavファイルの解析に失敗しました。\nPath : %s\n%w", data.Path, err)
		}
		target, err = io.ReadAll(stream)
		if err != nil {
			return nil, fmt.Errorf("wavファイルの解析結果が読み込めません。\nPath : %s\n%w", data.Path, err)
		}
	case ".mp3":
		stream, err := mp3.DecodeWithSampleRate(al.samplingRate, rawfile)
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
func (al SELoader) mkPlayer(AudioData []byte) *audio.Player {
	return al.audiocxt.NewPlayerFromBytes(AudioData)
}
