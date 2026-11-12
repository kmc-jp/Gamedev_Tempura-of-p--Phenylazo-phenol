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

// installedPartsUI は組み立て済みパーツの表示UI。
// 表示開始と表示完了の通知だけを公開し、ゲーム状態の遷移はSceneが行う。
type installedPartsUI interface {
	Element
	Start() error
	IsFinished() bool
}

type ResultScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	elements        []Element
	uiDrawer        UIDrawer
	partsList       partsListUI
	installedParts  installedPartsUI
	benchmarkPrompt Element
	partKindIndex   int
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
			if s.isSelectingArbitraryParts() && s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
				if s.installedParts == nil {
					return nil, transition.None(), fmt.Errorf("InstalledPartsUIが設定されていません")
				}
				if err := s.installedParts.Start(); err != nil {
					return nil, transition.None(), fmt.Errorf("InstalledPartsUIの開始に失敗しました。\n%w", err)
				}
				s.partsList.SetVisible(false)
				s.gamedata.State = game.Result_Assembled
				return
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) {
				part := s.partsList.SelectedPart()
				if part == nil {
					if s.isSelectingArbitraryParts() {
						return
					}
					return nil, transition.None(), fmt.Errorf("選択されたパーツがありません")
				}
				s.gamedata.Assemble = append(s.gamedata.Assemble, part)
				if s.isSelectingArbitraryParts() {
					if err := s.setCurrentPartsList(); err != nil {
						return nil, transition.None(), fmt.Errorf("任意パーツ一覧の更新に失敗しました。\n%w", err)
					}
					return
				}
				s.partKindIndex++
				if err := s.setCurrentPartsList(); err != nil {
					return nil, transition.None(), fmt.Errorf("結果画面のパーツ一覧の更新に失敗しました。\n%w", err)
				}
			}
		}
	}

	// 組み立てPC観察フェーズ
	if s.gamedata.State == game.Result_Assembled {
		//所持パーツをもとに利用する絵を決定
		//完成PCを表示

		if err := s.installedParts.Update(true); err != nil {
			return nil, transition.None(), fmt.Errorf("InstalledPartsUIの更新に失敗しました。\n%w", err)
		}
		if s.installedParts.IsFinished() {
			if err := s.gamedata.NextState(); err != nil {
				return nil, transition.None(), fmt.Errorf("InstalledPartsUI完了後の状態遷移に失敗しました: %w", err)
			}
		}
	}

	// ベンチマークフェーズ
	if s.gamedata.State == game.Result_Benchmark {
		if err := s.benchmarkPrompt.Update(true); err != nil {
			return nil, transition.None(), fmt.Errorf("ベンチマーク案内UIの更新に失敗しました。\n%w", err)
		}

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

func (s *ResultScene) SetInstalledParts(ip installedPartsUI) {
	s.installedParts = ip
	s.AddUIElement(ip)
}

func (s *ResultScene) SetBenchmarkPrompt(e Element) {
	s.benchmarkPrompt = e
	s.AddUIElement(e)
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
	if s.partKindIndex < 0 || s.partKindIndex > len(game.AvailableParts) {
		return fmt.Errorf("パーツ種別の選択位置が不正です: %d", s.partKindIndex)
	}

	available := game.Parts{}
	if s.isSelectingArbitraryParts() {
		selected := s.selectedArbitraryParts()
		for _, part := range s.gamedata.Parts {
			if part != nil && !s.isAvailablePart(part) && selected[part] == 0 {
				available = append(available, part)
			}
		}
	} else {
		kind := game.AvailableParts[s.partKindIndex]
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
	}

	if err := s.partsList.SetParts(available); err != nil {
		return err
	}
	s.partsList.SetVisible(true)
	return nil
}

func (s *ResultScene) selectedArbitraryParts() map[*game.Part]int {
	selected := map[*game.Part]int{}
	for _, part := range s.gamedata.Assemble {
		if part != nil && !s.isAvailablePart(part) {
			selected[part]++
		}
	}
	return selected
}

func (s *ResultScene) isSelectingArbitraryParts() bool {
	return s.partKindIndex == len(game.AvailableParts)
}

func (s *ResultScene) isAvailablePart(part *game.Part) bool {
	for _, availableKind := range game.AvailableParts {
		if part.Kind == availableKind {
			return true
		}
	}
	return false
}

func (s ResultScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
