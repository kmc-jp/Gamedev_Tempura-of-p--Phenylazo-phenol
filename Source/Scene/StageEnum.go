package scene

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
