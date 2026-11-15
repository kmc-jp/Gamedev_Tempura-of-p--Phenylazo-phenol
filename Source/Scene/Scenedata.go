package scene

import (
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type InitializedScenes struct {
	InitTitleScene       Utils.Factory[*TitleScene]
	InitSelectStageScene Utils.Factory[*SelectStageScene]
	InitPlayScene        map[stageinfo.Location]Utils.Factory[*PlayScene]
	InitPauseScene       Utils.Factory[*PauseScene]
	InitTrainScene       Utils.Factory[*TrainScene]
	InitResultScene      Utils.Factory[*ResultScene]
}
