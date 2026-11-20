package serializetarget

import (
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

// ステージを歩き回るおじさんの設定
// 見た目・歩行に関する設定は PlayerConfig と共通 (Position は Layer のみ使う)
type OjisanConfig struct {
	PlayerConfig

	// 出現人数 (Spawns が空のとき、床からランダムに選ぶ)
	Count int
	// 歩き出すまでの待ちフレーム数の範囲
	MoveIntervalMin int
	MoveIntervalMax int
	// 出現位置。指定した場合は Count を無視してこの座標に出す
	Spawns []Utils.Vec2Int
}

type OjisanConstructData struct {
	PlayerConstructData
	Count           int
	MoveIntervalMin int
	MoveIntervalMax int
	Spawns          []Utils.Vec2Int
}

func (oc OjisanConfig) Construct() (any, error) {
	if oc.MoveIntervalMin < 1 || oc.MoveIntervalMax < oc.MoveIntervalMin {
		return nil, fmt.Errorf("MoveIntervalMin / MoveIntervalMax が不正です。 1 <= Min <= Max を満たしてください。\nMin = %d, Max = %d", oc.MoveIntervalMin, oc.MoveIntervalMax)
	}
	if len(oc.Spawns) == 0 && oc.Count < 0 {
		return nil, fmt.Errorf("Count が負です。 Count = %d", oc.Count)
	}

	pcAny, err := oc.PlayerConfig.Construct()
	if err != nil {
		return nil, fmt.Errorf("PlayerConfig.Construct() でエラーが発生しました。\n%w", err)
	}

	return OjisanConstructData{
		PlayerConstructData: pcAny.(PlayerConstructData),
		Count:               oc.Count,
		MoveIntervalMin:     oc.MoveIntervalMin,
		MoveIntervalMax:     oc.MoveIntervalMax,
		Spawns:              oc.Spawns,
	}, nil
}

// test
var _ serialize.Serializable = OjisanConfig{}
