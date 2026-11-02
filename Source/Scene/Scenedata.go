package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type InitializedScenes struct {
	InitTitleScene  Utils.Factory[*TitleScene]
	InitPlayScene   map[Stage]Utils.Factory[*PlayScene]
	InitPauseScene  Utils.Factory[*PauseScene]
	InitTrainScene  Utils.Factory[*TrainScene]
	InitResultScene Utils.Factory[*ResultScene]
}
