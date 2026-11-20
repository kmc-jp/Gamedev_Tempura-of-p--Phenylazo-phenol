package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

// HowToPlayCanvasAPI は HowToPlayScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructから注入する。
type HowToPlayCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsHowToPlayCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsHowToPlayCanvas()
}

// HowToPlayScene は遊び方を1枚の画像で見せる画面である。
// タイトル・ポーズのどちらからも重ねて開かれ、閉じると呼び出し元へ戻る。
type HowToPlayScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas HowToPlayCanvasAPI

	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

// NewHowToPlayScene は HowToPlayScene を生成する。
func NewHowToPlayScene(worldDrawer Utils.Factory[WorldDrawer], newCanvas func(*HowToPlayScene) (HowToPlayCanvasAPI, error), Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*HowToPlayScene] {
	return func() (*HowToPlayScene, error) {
		camera, err := worldDrawer()
		if err != nil {
			return &HowToPlayScene{}, fmt.Errorf("worldDrawer() でエラーが発生しました。 \n%w", err)
		}
		if newCanvas == nil {
			return &HowToPlayScene{}, fmt.Errorf("HowToPlaySceneのcanvas生成関数が nil です")
		}

		hs := HowToPlayScene{
			entities:    []Entity{},
			worldDrawer: camera,
			Input:       Input,
			gamedata:    gdata,
			initScenes:  initScenes,
		}
		canvas, err := newCanvas(&hs)
		if err != nil {
			return &HowToPlayScene{}, fmt.Errorf("HowToPlaySceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &HowToPlayScene{}, fmt.Errorf("HowToPlaySceneのcanvas が nil です")
		}
		hs.canvas = canvas
		return &hs, nil
	}
}

func (s HowToPlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *HowToPlayScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
}

// Update はSceneの処理の後でCanvasを更新する。
func (s *HowToPlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	// Esc・決定のどれでも閉じる
	if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) ||
		s.Input.ActionIsJustPressed(inputconfig.ActionSpace) ||
		s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
		return nil, transition.Pop(), nil
	}

	if err := s.canvas.Update(active); err != nil {
		return nil, transition.None(), fmt.Errorf("Canvasの更新に失敗しました。\n%w", err)
	}
	return
}

func (s *HowToPlayScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s HowToPlayScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
