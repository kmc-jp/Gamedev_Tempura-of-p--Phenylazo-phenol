package serializetarget

import (
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	"fmt"
)

type StagesConfig struct {
	StageMapURL map[string]string `toml:"Stages"`
}

// stageinfo.Stage2URL = map[stageinfo.Location]string を返す
func (sc StagesConfig) Construct() (any, error) {
	if sc.StageMapURL == nil {
		return nil, fmt.Errorf("シリアライズに失敗したようです。")
	}
	Map := stageinfo.Stage2URL{}
	for k, v := range sc.StageMapURL {
		loc, ok := stageinfo.String2Location[k]
		if !ok {
			return nil, fmt.Errorf("未知のステージ名です: %s", k)
		}
		if loc == stageinfo.Departure {
			return nil, fmt.Errorf("自宅 (%s) にはステージを設定できません", k)
		}
		Map[loc] = v
	}
	return Map, nil
}
