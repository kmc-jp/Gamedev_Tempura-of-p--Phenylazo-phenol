package serializetarget

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"fmt"
)

type StagesConfig struct {
	StageMapURL map[string]string `toml:"Stages"`
}

// scene.Stage2URL = map[scene.Stage]string を返す
func (sc StagesConfig) Construct() (any, error) {
	if sc.StageMapURL == nil {
		return nil, fmt.Errorf("シリアライズに失敗したようです。")
	}
	Map := map[scene.Stage]string{}
	for k, v := range sc.StageMapURL {
		Map[scene.String2Stage[k]] = v
	}
	return Map, nil
}
