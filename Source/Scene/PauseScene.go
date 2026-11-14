package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	stdcolor "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

// PauseCanvasAPI は PauseScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructから注入する。
type PauseCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsPauseCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsPauseCanvas()
}

type PauseScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas PauseCanvasAPI

	// dimImage は背後の画面を暗くする半透明の画像。毎フレーム作らないようキャッシュし、画面サイズが変わったときだけ作り直す。
	dimImage *ebiten.Image

	// text        *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewPauseScene(worldDrawer Utils.Factory[WorldDrawer], canvas PauseCanvasAPI, Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*PauseScene] {
	return func() (*PauseScene, error) {
		camera, err := worldDrawer()
		if err != nil {
			return &PauseScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		if canvas == nil {
			return &PauseScene{}, fmt.Errorf("PauseSceneのcanvas が nil です")
		}

		ps := PauseScene{
			entities:    []Entity{},
			worldDrawer: camera,
			canvas:      canvas,
			// text:        text,
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
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if s.dimImage == nil || s.dimImage.Bounds().Dx() != w || s.dimImage.Bounds().Dy() != h {
		if s.dimImage != nil {
			s.dimImage.Deallocate()
		}
		s.dimImage = ebiten.NewImage(w, h)
		s.dimImage.Fill(stdcolor.RGBA{100, 100, 100, 200})
	}
	screen.DrawImage(s.dimImage, &ebiten.DrawImageOptions{})

	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
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

	if err := s.canvas.Update(active); err != nil {
		return nil, transition.None(), fmt.Errorf("Canvasの更新に失敗しました。\n%w", err)
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
	if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
		return nil, transition.Pop(), nil
	}

	return
}

func (s *PauseScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s PauseScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
