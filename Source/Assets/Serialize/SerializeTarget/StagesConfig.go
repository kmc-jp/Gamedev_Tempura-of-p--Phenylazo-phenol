package serializetarget

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"fmt"
)

type StagesConfig struct {
	StageMapURL map[scene.Stage]string `toml:"Stages"`
}

// scene.Stage2URL = map[scene.Stage]string を返す
func (sc StagesConfig) Construct() (any, error) {
	if sc.StageMapURL == nil {
		return nil, fmt.Errorf("シリアライズに失敗したようです。")
	}
	return sc.StageMapURL, nil
}
