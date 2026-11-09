package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	partsdata "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/Parts"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type ResultScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	elements      []Element
	uiDrawer      UIDrawer
	partsList     partsListUI
	partKindIndex int
	// text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewResultScene(NewWorldDrawer Utils.Factory[WorldDrawer], uiDrawer Utils.Factory[UIDrawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*ResultScene] {
	return func() (*ResultScene, error) {
		camera, err := NewWorldDrawer()
		if err != nil {

			return &ResultScene{}, fmt.Errorf("NewDrawer() でエラーが発生しました。 \n%w", err)
		}
		uiCamera, err := uiDrawer()
		if err != nil {
			return &ResultScene{}, fmt.Errorf("NewDrawer() でエラーが発生しました。\n%w", err)
		}
		ts := ResultScene{
			entities:    []Entity{},
			worldDrawer: camera,
			uiDrawer:    uiCamera,
			// text:        text,
			Input:      Input,
			gamedata:   gdata,
			initScenes: initScenes,
		}
		// なにもしない
		return &ts, nil
	}
}

func (s ResultScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *ResultScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.uiDrawer.Draw(screen, s.elements)
	// s.text.Draw(screen)
}

func (s *ResultScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	if s.gamedata.InitialProcess {
		//シーン起動時処理
		s.gamedata.InitialProcess = false
		s.gamedata.Assemble = game.Parts{}
		s.partKindIndex = 0
		if err := s.setCurrentPartsList(); err != nil {
			return nil, transition.None(), fmt.Errorf("結果画面のパーツ一覧の構築に失敗しました。\n%w", err)
		}
	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	if s.gamedata.State == game.Result_SelectParts {
		//組み立てに使うパーツを選択
		if s.partsList != nil && s.partsList.Visible() {
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
				s.partsList.MoveSelection(-1)
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
				s.partsList.MoveSelection(1)
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) {
				part := s.partsList.SelectedPart()
				if part == nil {
					return nil, transition.None(), fmt.Errorf("選択されたパーツがありません")
				}
				s.gamedata.Assemble = append(s.gamedata.Assemble, part)
				s.partKindIndex++
				if s.partKindIndex >= len(game.AvailableParts) {
					s.partsList.SetVisible(false)
					s.gamedata.State = game.Result_Assembled
				} else if err := s.setCurrentPartsList(); err != nil {
					return nil, transition.None(), fmt.Errorf("結果画面のパーツ一覧の更新に失敗しました。\n%w", err)
				}
			}
		}
	}

	// 組み立てPC観察フェーズ
	if s.gamedata.State == game.Result_Assembled {
		//所持パーツをもとに利用する絵を決定
		//完成PCを表示

		s.elements[2].Update(true)
	}

	// ベンチマークフェーズ
	if s.gamedata.State == game.Result_Benchmark {
		s.elements[1].Update(true)

		if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
			var score float32
			for _, v := range s.gamedata.Assemble {
				score += v.Score()
			}
			fmt.Print("ベンチマーク計測結果", int(score), "\n")
			s.gamedata.State = game.Result_QuitGame
			return
		}
	}

	if s.gamedata.State == game.Result_QuitGame {

		if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
			return func() (game.Scene, error) {
				nextScene, err := s.initScenes.InitTitleScene()
				if err != nil {
					return nil, err
				}
				nextScene.gamedata = s.gamedata
				nextScene.gamedata.ForNewScene(game.Title_Default)
				nextScene.gamedata.ResetGame()
				nextScene.initScenes = s.initScenes
				return nextScene, nil
			}, transition.NewScene(), nil
		}
	}

	return
}

func (s *ResultScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s *ResultScene) AddUIElement(e Element) {
	s.elements = append(s.elements, e)
}

func (s *ResultScene) SetPartsList(partsList partsListUI) {
	s.partsList = partsList
	partsList.SetVisible(false)
	s.AddUIElement(partsList)
}

func (s *ResultScene) setCurrentPartsList() error {
	if s.partsList == nil {
		return fmt.Errorf("パーツ一覧が設定されていません")
	}
	if s.partKindIndex < 0 || s.partKindIndex >= len(game.AvailableParts) {
		return fmt.Errorf("パーツ種別の選択位置が不正です: %d", s.partKindIndex)
	}

	kind := game.AvailableParts[s.partKindIndex]
	available := game.Parts{}
	for _, part := range s.gamedata.Parts {
		if part != nil && part.Kind == kind {
			available = append(available, part)
		}
	}
	if len(available) == 0 {
		trash := partsdata.Trash()
		if trash == nil {
			return fmt.Errorf("Trashパーツの生成に失敗しました")
		}
		trash.Kind = kind
		available = append(available, trash)
	}

	if err := s.partsList.SetParts(available); err != nil {
		return err
	}
	s.partsList.SetVisible(true)
	return nil
}

func (s ResultScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
