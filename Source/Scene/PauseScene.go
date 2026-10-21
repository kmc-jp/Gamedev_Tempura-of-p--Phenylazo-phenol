package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	input "github.com/quasilyte/ebitengine-input"
)

type PauseScene struct {
	entities   []Entity
	drawer     Drawer
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewPauseScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*PauseScene] {
	camera, err := NewDrawer()
	if err != nil {
		return func() (*PauseScene, error) {
			return &PauseScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
	}
	ps := PauseScene{
		entities:   []Entity{},
		drawer:     camera,
		Input:      Input,
		gamedata:   gdata,
		initScenes: initScenes,
	}
	// なにもしない
	return func() (*PauseScene, error) {
		return &ps, nil
	}
}

func (s PauseScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *PauseScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.entities)
	ui.DrawText(Utils.Position{}, color.Orange(),
		"Pause\nPress ESC Key to Return the Game",
		ui.FontFace(s.gamedata.Font, 30), screen)
}

func (s *PauseScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	if s.gamedata.InitialProcess {
		//シーン起動時処理
		s.gamedata.InitialProcess = false

	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return nil, transition.Pop(), nil
	}

	return
}

func (s PauseScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s PauseScene) SetCameraCenter(center Utils.Position) {
	s.drawer.SetCenter(center)
}
