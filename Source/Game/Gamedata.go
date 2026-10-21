package game

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// シーン間共有データ
type Gamedata struct {
	//文字描画用の基本フォント格納
	Font *text.GoTextFaceSource

	AssetsManager *assets.AssetsManager
}

func NewGamedata() *Gamedata {
	return &Gamedata{}
}
