package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type PauseScene struct {
	entities   []Entity
	drawer     WorldDrawer
	text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewPauseScene(NewDrawer Utils.Factory[WorldDrawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*PauseScene] {
	return func() (*PauseScene, error) {
		camera, err := NewDrawer()
		if err != nil {
			return &PauseScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		text, err := ui.NewText(
			Utils.Position{X: 0, Y: 50},
			color.Orange(),
			"Pause\nPress ESC Key to Return the Game",
			ui.FontFace(gdata.Font, 15),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		ps := PauseScene{
			entities:   []Entity{},
			drawer:     camera,
			text:       text,
			Input:      Input,
			gamedata:   gdata,
			initScenes: initScenes,
		}
		return &ps, nil
	}
}

func (s PauseScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *PauseScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.entities)
	s.text.Draw(screen)
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

	// if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
	if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
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
