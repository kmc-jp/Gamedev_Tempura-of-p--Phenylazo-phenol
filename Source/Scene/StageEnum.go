package scene

import stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"

type Stage int
type Stage2URL = map[Stage]string

const (
	Red Stage = iota
	Blue
	Green
	Orange
	Purple
)

const Stages = 5 //Stage数

var Stage2String map[Stage]string = map[Stage]string{
	Red:    "Red",
	Blue:   "Blue",
	Green:  "Green",
	Orange: "Orange",
	Purple: "Purple",
}

var String2Stage map[string]Stage = map[string]Stage{
	"Red":    Red,
	"Blue":   Blue,
	"Green":  Green,
	"Orange": Orange,
	"Purple": Purple,
}

var location2Stage = map[stageinfo.Location]Stage{
	stageinfo.Red:    Red,
	stageinfo.Blue:   Blue,
	stageinfo.Green:  Green,
	stageinfo.Orange: Orange,
	stageinfo.Purple: Purple,
}

// StageOfLocation は路線図上の駅に対応するステージを返す。
// 自宅（stageinfo.Departure）などステージのない駅は false。
func StageOfLocation(loc stageinfo.Location) (Stage, bool) {
	stage, ok := location2Stage[loc]
	return stage, ok
}

// // Stage のデシリアライズに使えるようにする (できなかった)
// func (s *Stage) UnmarshalText(charArray []byte) error {
// 	text := string(charArray)
// 	ok := false
// 	for stage, str := range Stage2String {
// 		if text == str {
// 			*s = stage
// 			ok = true
// 		}
// 	}
// 	if !ok {
// 		return fmt.Errorf("invalid Stage value: %s", text)
// 	}
// 	return nil
// }
