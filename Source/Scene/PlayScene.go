package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type partsListUI interface {
	Element
	SetParts(game.Parts) error
	SetVisible(bool)
	Visible() bool
	MoveSelection(int)
	SelectedPart() *game.Part
}

// メインの遊べるシーン
type PlayScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	elements  []Element
	uiDrawer  UIDrawer
	partsList partsListUI

	// text       *ui.Text
	Input      *input.Handler
	Gamedata   *game.Gamedata
	initScenes *InitializedScenes
	// テスト用
	cameraCenter Utils.Position
}

func NewPlayScene(worldDrawer Utils.Factory[WorldDrawer], uiDrawer Utils.Factory[UIDrawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*PlayScene] {
	return func() (*PlayScene, error) {
		camera, err := worldDrawer()
		if err != nil {
			return &PlayScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		uiCamera, err := uiDrawer()
		if err != nil {
			return &PlayScene{}, fmt.Errorf("uiDrawer() でエラーが発生しました。 \n%w", err)
		}
		ps := PlayScene{
			entities:    []Entity{},
			worldDrawer: camera,
			elements:    []Element{},
			uiDrawer:    uiCamera,
			// text:       text,
			Input:      Input,
			Gamedata:   gdata,
			initScenes: initScenes,
		}
		return &ps, nil
	}
}

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.uiDrawer.Draw(screen, s.elements)
}

// Update implements [Scene].
func (s *PlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	if s.Gamedata.InitialProcess {
		s.Gamedata.InitialProcess = false
		s.Gamedata.Text = "10"
		// s.text.SetText("探索フェーズ 残り時間: " + s.gamedata.Text)

	}

	if s.Gamedata.State == game.Play_Explore {
		if s.partsList != nil && s.Input.ActionIsJustPressed(inputconfig.ActionPartsList) {
			if err := s.partsList.SetParts(s.Gamedata.Found); err != nil {
				return nil, transition.None(), fmt.Errorf("パーツ一覧の構築に失敗しました。\n%w", err)
			}
			s.partsList.SetVisible(true)
		}

		if s.partsList != nil && s.partsList.Visible() {
			if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
				s.partsList.SetVisible(false)
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
				s.partsList.MoveSelection(-1)
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
				s.partsList.MoveSelection(1)
			}
			return
		}

		for _, entity := range s.entities {
			err := entity.Update(active)
			if err != nil {
				return nil, transition.None(), fmt.Errorf("%s.Update(%s) でエラーが発生しました。\n%w", entity.Name(), strconv.FormatBool(active), err)
			}
		}

		s.Gamedata.Timer.Set(10)

		//制限時間内プレイ中の制御
		if s.Gamedata.Timer.Expired() {
			num, _ := strconv.Atoi(s.Gamedata.Text)
			num--
			s.Gamedata.Text = strconv.Itoa(num)

			if num <= 0 {
				s.Gamedata.Timer.Set(0)
				s.Gamedata.State = game.Play_Timeup

				s.Gamedata.Text = "0\nTimeUp 次の店に移動します"
				// s.text.SetText("探索フェーズ 残り時間: " + s.gamedata.Text)
				fmt.Print("みつけたもの\n", s.Gamedata.ExtractFoundPartsName(), "\n")
			}
		}
		// s.text.SetText("探索フェーズ 残り時間: " + s.gamedata.Text)

		// if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
			return func() (game.Scene, error) {
				pauseScene, err := s.initScenes.InitPauseScene()
				if err != nil {
					return nil, err
				}
				pauseScene.gamedata = s.Gamedata
				//pauseScene.gamedata.ForNewScene(0)
				pauseScene.initScenes = s.initScenes
				return pauseScene, nil
			}, transition.Push(), nil
		}
	}

	//TimeUp後の処理
	if s.Gamedata.State == game.Play_Timeup {
		s.Gamedata.Timer.Set(50)
		if s.Gamedata.Timer.Expired() {
			s.Gamedata.Timer.Set(0)
			s.Gamedata.State = game.Play_Select

			if err := s.partsList.SetParts(s.Gamedata.Found); err != nil {
				return nil, transition.None(), fmt.Errorf("パーツ一覧の構築に失敗しました。\n%w", err)
			}
			s.partsList.SetVisible(true)
		}

		return
	}

	if s.Gamedata.State == game.Play_Select {

		if s.partsList != nil && s.partsList.Visible() {
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
				s.partsList.MoveSelection(-1)
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
				s.partsList.MoveSelection(1)
			}
		}

		if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
			return func() (game.Scene, error) {
				trainScene, err := s.initScenes.InitTrainScene()
				if err != nil {
					return nil, err
				}
				trainScene.gamedata = s.Gamedata
				trainScene.gamedata.ForNewScene(game.Train_Default)
				trainScene.initScenes = s.initScenes
				return trainScene, nil
			}, transition.NewScene(), nil
		}
	}

	return
}

func (s *PlayScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s *PlayScene) AddUIElement(e Element) {
	s.elements = append(s.elements, e)
}

func (s *PlayScene) SetPartsList(partsList partsListUI) {
	s.partsList = partsList
	s.AddUIElement(partsList)
}

func (s PlayScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
