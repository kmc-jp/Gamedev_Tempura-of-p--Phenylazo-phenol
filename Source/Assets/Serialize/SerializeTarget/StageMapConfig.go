package serializetarget

import (
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Selection"
	"fmt"
)

// StationConfig は StageMap.toml の駅1つ分の定義である。
type StationConfig struct {
	Label string
	X     float64
	Y     float64
}

// LineConfig は StageMap.toml の路線1本分の定義である。
//
// FromKey / ToKey は省略可。"Left" / "Right" / "Up" / "Down" で、その向きに進むキーを明示する。
// 省略すると駅の位置関係から自動で決まる。
//   - FromKey: From の駅で To へ進むキー
//   - ToKey:   To の駅で From へ戻るキー
type LineConfig struct {
	From    string
	To      string
	FromKey string
	ToKey   string
}

var string2Direction = map[string]selection.Direction{
	"Left":  selection.Left,
	"Right": selection.Right,
	"Up":    selection.Up,
	"Down":  selection.Down,
}

// parseKey は空文字なら nil（自動）を返す。
func parseKey(s string) (*selection.Direction, error) {
	if s == "" {
		return nil, nil
	}
	d, ok := string2Direction[s]
	if !ok {
		return nil, fmt.Errorf("未知のキー名です: %s (Left / Right / Up / Down のいずれか)", s)
	}
	return &d, nil
}

// StageMapConfig は StageMap.toml の内容である。
type StageMapConfig struct {
	Width    float64
	Height   float64
	Stations map[string]StationConfig
	Lines    []LineConfig
}

// stageinfo.StageMap を返す
func (c StageMapConfig) Construct() (any, error) {
	stageMap := stageinfo.StageMap{
		Width:    c.Width,
		Height:   c.Height,
		Stations: map[stageinfo.Location]stageinfo.Station{},
	}

	for name, station := range c.Stations {
		loc, ok := stageinfo.String2Location[name]
		if !ok {
			return nil, fmt.Errorf("未知の駅名です: %s", name)
		}
		stageMap.Stations[loc] = stageinfo.Station{
			Label: station.Label,
			Pos:   Utils.Position{X: station.X, Y: station.Y},
		}
	}

	for _, line := range c.Lines {
		from, ok := stageinfo.String2Location[line.From]
		if !ok {
			return nil, fmt.Errorf("路線の駅名が未知です: %s", line.From)
		}
		to, ok := stageinfo.String2Location[line.To]
		if !ok {
			return nil, fmt.Errorf("路線の駅名が未知です: %s", line.To)
		}
		fromKey, err := parseKey(line.FromKey)
		if err != nil {
			return nil, fmt.Errorf("路線 %s - %s の FromKey が不正です。\n%w", line.From, line.To, err)
		}
		toKey, err := parseKey(line.ToKey)
		if err != nil {
			return nil, fmt.Errorf("路線 %s - %s の ToKey が不正です。\n%w", line.From, line.To, err)
		}
		stageMap.Lines = append(stageMap.Lines, stageinfo.Line{From: from, To: to, FromKey: fromKey, ToKey: toKey})
	}

	if err := stageMap.Validate(); err != nil {
		return nil, fmt.Errorf("路線図の検証に失敗しました。\n%w", err)
	}
	return stageMap, nil
}
