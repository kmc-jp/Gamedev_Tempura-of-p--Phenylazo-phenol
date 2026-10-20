package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	input "github.com/quasilyte/ebitengine-input"
)

type TitleScene struct {
	entities []Entity
	drawer   Drawer
	Input    *input.Handler
}

func NewTitleScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler) Utils.Factory[*PlayScene] {
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
	}
	// なにもしない
	return func() (*PlayScene, error) {
		return &ps, nil
	}
}
