package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

// メインの遊べるシーン
type PlayScene struct {
	entities   []Entity
	drawer     Drawer
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
	// テスト用
	cameraCenter Utils.Position
}

func NewPlayScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*PlayScene] {
	camera, err := NewDrawer()
	if err != nil {
		return func() (*PlayScene, error) {
			return &PlayScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
	}
	ps := PlayScene{
		entities:   []Entity{},
		drawer:     camera,
		Input:      Input,
		gamedata:   gdata,
		initScenes: initScenes,
	}
	// なにもしない
	return func() (*PlayScene, error) {
		return &ps, nil
	}
}

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.entities)
	ui.DrawText(Utils.Position{X: 0, Y: 0}, color.Red(),
		"探索フェーズ 残り時間: "+s.gamedata.Text, ui.FontFace(s.gamedata.Font, 30), screen)
}

// Update implements [Scene].
func (s *PlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	if s.gamedata.InitialProcess {
		s.gamedata.InitialProcess = false
		s.gamedata.Text = "10"

		s.gamedata.Timer = append(s.gamedata.Timer, tick.NewTimer()) //OneSec
		s.gamedata.Timer = append(s.gamedata.Timer, tick.NewTimer()) //FiveSec

		s.gamedata.Flag = append(s.gamedata.Flag, false)
	}

	//TimeUp後の処理
	if s.gamedata.Flag[game.TimeUp] {
		if s.gamedata.Timer[game.FiveSec].Count() {
			s.gamedata.Timer[game.FiveSec].Set(0)

			trainScene, err := s.initScenes.InitTrainScene()
			if err != nil {
				return nil, transition.None(), err
			}
			trainScene.gamedata = s.gamedata
			trainScene.gamedata.InitialProcess = true
			trainScene.initScenes = s.initScenes
			return func() (game.Scene, error) { return trainScene, nil }, transition.NewScene(), nil
		}
		return
	}

	for _, entity := range s.entities {
		err := entity.Update(active)
		if err != nil {
			return nil, transition.None(), fmt.Errorf("%s.Update(%s) でエラーが発生しました。\n%w", entity.Name(), strconv.FormatBool(active), err)
		}
	}

	// 何か軽いテストはここでやりましょう

	s.gamedata.Timer[game.OneSec].Set(10)
	s.gamedata.Timer[game.FiveSec].Set(50)

	//制限時間内プレイ中の制御
	if cnt := s.gamedata.Timer[game.OneSec].Count(); cnt {
		num, _ := strconv.Atoi(s.gamedata.Text)
		num--
		s.gamedata.Text = strconv.Itoa(num)

		if num <= 0 {
			s.gamedata.Timer[game.OneSec].Set(0)
			s.gamedata.Flag[game.TimeUp] = true

			s.gamedata.Text = "0\nTimeUp 次の店に移動します"
		}
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
	if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
		pauseScene, err := s.initScenes.InitPauseScene()
		if err != nil {
			return nil, transition.None(), err
		}
		pauseScene.gamedata = s.gamedata
		pauseScene.gamedata.InitialProcess = true
		pauseScene.initScenes = s.initScenes
		return func() (game.Scene, error) { return pauseScene, nil }, transition.Push(), nil
	}

	return
}

func (s *PlayScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s PlayScene) SetCameraCenter(center Utils.Position) {
	s.drawer.SetCenter(center)
}
