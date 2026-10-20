package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	input "github.com/quasilyte/ebitengine-input"
)

type TitleScene struct {
	entities []Entity
	drawer   Drawer
	Input    *input.Handler
	gamedata game.Gamedata
}

func NewTitleScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler, gdata *game.Gamedata) Utils.Factory[*PlayScene] {
	camera, err := NewDrawer()
	if err != nil {
		return func() (*PlayScene, error) {
			return &PlayScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
	}
	ps := PlayScene{
		entities: []Entity{},
		drawer:   camera,
		Input:    Input,
		gamedata: gdata,
	}
	// なにもしない
	return func() (*PlayScene, error) {
		return &ps, nil
	}
}
