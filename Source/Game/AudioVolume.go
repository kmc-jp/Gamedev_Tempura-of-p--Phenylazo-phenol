package game

import (
	audioassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Audio"
	"fmt"
	"weak"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// 音量は 0〜MaxVolume の整数段階で持つ。
const (
	MaxVolume     = 10
	DefaultVolume = 10
)

// audioVolume はシーン間で共有する音量設定である。
type audioVolume struct {
	bgm int
	se  int
	// 再生中のBGMプレイヤー。音量変更を即時反映するために保持する。
	// シーン破棄後に残らないよう弱参照にする。
	bgmPlayers []weak.Pointer[audio.Player]
}

func clampVolume(v int) int {
	return min(max(v, 0), MaxVolume)
}

func volumeRatio(v int) float64 {
	return float64(v) / MaxVolume
}

// BGMVolume はBGM音量(0〜MaxVolume)を返す。
func (g *Gamedata) BGMVolume() int { return g.volume.bgm }

// SEVolume はSE音量(0〜MaxVolume)を返す。
func (g *Gamedata) SEVolume() int { return g.volume.se }

// SetBGMVolume はBGM音量を設定し、読み込み済みのBGMへ即時反映する。
func (g *Gamedata) SetBGMVolume(v int) {
	g.volume.bgm = clampVolume(v)
	alive := g.volume.bgmPlayers[:0]
	for _, wp := range g.volume.bgmPlayers {
		p := wp.Value()
		if p == nil {
			continue
		}
		p.SetVolume(volumeRatio(g.volume.bgm))
		alive = append(alive, wp)
	}
	g.volume.bgmPlayers = alive
}

// StopAllBGM は読み込み済みの全BGMを停止する。
func (g *Gamedata) StopAllBGM() {
	alive := g.volume.bgmPlayers[:0]
	for _, wp := range g.volume.bgmPlayers {
		p := wp.Value()
		if p == nil {
			continue
		}
		p.Pause()
		alive = append(alive, wp)
	}
	g.volume.bgmPlayers = alive
}

// SetSEVolume はSE音量を設定する。次に読み込むSEから反映される。
func (g *Gamedata) SetSEVolume(v int) {
	g.volume.se = clampVolume(v)
}

// LoadBGM はBGMを読み込み、現在の音量を設定したプレイヤーを返す。
func (g *Gamedata) LoadBGM(path string) (*audio.Player, error) {
	pl, err := g.loadAudio(path, audioassets.BGM)
	if err != nil {
		return nil, err
	}
	pl.SetVolume(volumeRatio(g.volume.bgm))
	g.volume.bgmPlayers = append(g.volume.bgmPlayers, weak.Make(pl))
	return pl, nil
}

// LoadSE はSEを読み込み、現在の音量を設定したプレイヤーを返す。
func (g *Gamedata) LoadSE(path string) (*audio.Player, error) {
	pl, err := g.loadAudio(path, audioassets.SE)
	if err != nil {
		return nil, err
	}
	pl.SetVolume(volumeRatio(g.volume.se))
	return pl, nil
}

func (g *Gamedata) loadAudio(path string, t audioassets.AudioType) (*audio.Player, error) {
	ad, err := g.AssetsManager.Audio.NewAudioData(path, t)
	if err != nil {
		return nil, fmt.Errorf("AssetsManager.Audio.NewAudioData() でエラーが発生しました。\npath : %s\n%w", path, err)
	}
	pl, err := g.AssetsManager.Audio.Load(ad)
	if err != nil {
		return nil, fmt.Errorf("AssetsManager.Audio.Load() でエラーが発生しました。\npath : %s\n%w", path, err)
	}
	return pl, nil
}
