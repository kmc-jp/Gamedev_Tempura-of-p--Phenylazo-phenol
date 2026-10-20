package game

import "github.com/hajimehoshi/ebiten/v2/text/v2"

//シーン間共有データ
type Gamedata struct {
	font *text.GoTextFaceSource
}
