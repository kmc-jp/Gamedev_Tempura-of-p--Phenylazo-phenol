package scene

import "fmt"

type Stage int
type Stage2URL = map[Stage]string

const (
	Stage1 Stage = iota
	Stage2
	Stage3
)

var Stage2String map[Stage]string = map[Stage]string{
	Stage1: "Stage1",
	Stage2: "Stage2",
	Stage3: "Stage3",
}

// Stage のデシリアライズに使えるようにする
func (s *Stage) UnmarshalText(charArray []byte) error {
	text := string(charArray)
	ok := false
	for stage, str := range Stage2String {
		if text == str {
			*s = stage
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("invalid Stage value: %s", text)
	}
	return nil
}
