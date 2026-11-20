package audioassets

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// ループ再生するためのstreamのラッパーを保存
type BGMLoader struct {
	fs       fs.FS
	audiocxt *audio.Context
	cache    map[string]*audio.InfiniteLoop // キーはパス
}

func NewBGMLoader(fs fs.FS, audiocxt *audio.Context) *BGMLoader {
	return &BGMLoader{
		fs:       fs,
		audiocxt: audiocxt,
		cache:    map[string]*audio.InfiniteLoop{},
	}
}

func (al *BGMLoader) Load(data AudioData) (*audio.Player, error) {
	_, ok := al.cache[data.Path]
	if !ok {
		ad, err := al.fetch(data)
		if err != nil {
			return nil, fmt.Errorf("al.fetch(data.Path) でエラーが発生しました。\n%w", err)
		}
		al.cache[data.Path] = ad
	}

	ad, _ := al.cache[data.Path]
	pl, err := al.mkPlayer(ad)
	if err != nil {
		return nil, fmt.Errorf("al.mkPlayer(ad) でエラーが発生しました。\n%w", err)
	}
	return pl, nil
}

// デコードして []byte にする
func (al BGMLoader) fetch(data AudioData) (*audio.InfiniteLoop, error) {
	rawfile, err := al.fs.Open(data.Path)
	if err != nil {
		return nil, err
	}
	defer rawfile.Close()

	var target *audio.InfiniteLoop

	switch data.ext {
	case ".wav":
		stream, err := wav.DecodeWithSampleRate(al.audiocxt.SampleRate(), rawfile)
		if err != nil {
			return nil, fmt.Errorf("wavファイルの解析に失敗しました。\nPath : %s\n%w", data.Path, err)
		}
		target = audio.NewInfiniteLoop(stream, stream.Length())
	case ".mp3":
		stream, err := mp3.DecodeWithSampleRate(al.audiocxt.SampleRate(), rawfile)
		if err != nil {
			return nil, fmt.Errorf("mp3ファイルの解析に失敗しました。\nPath : %s\n%w", data.Path, err)
		}
		target = audio.NewInfiniteLoop(stream, stream.Length())
	default:
		return nil, fmt.Errorf("対応していないファイル形式です\nPath : %s", data.Path)
	}

	return target, nil
}

// デコードしたデータを *audio.Player にする
func (al BGMLoader) mkPlayer(AudioData *audio.InfiniteLoop) (*audio.Player, error) {
	pl, err := al.audiocxt.NewPlayer(AudioData)
	if err != nil {
		return nil, fmt.Errorf("al.audiocxt.NewPlayer(AudioData) でエラーが発生しました。\n%w", err)
	}
	return pl, nil
}
