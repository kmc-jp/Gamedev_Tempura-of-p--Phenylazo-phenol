package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	ranking "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Ranking"
	parts "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/Parts"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	selection "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Selection"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	input "github.com/quasilyte/ebitengine-input"
)

// RankingCanvasAPI は RankingScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructから注入する。
type RankingCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsRankingCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsRankingCanvas()
}

// RankingScene はスコア上位の記録と、その時のPC構成・見た目を見せる画面である。
// リザルトの後に開かれ、閉じるとゲームをリセットしてタイトルへ戻る。
type RankingScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas RankingCanvasAPI
	store  RankingStore

	selector *selection.Graph[int, selection.Direction]
	// ranking は表示中の記録。シーン開始時に読み込む。
	records []ranking.Entry
	// highlightIndex は今回のスコアが入った順位(0始まり)。圏外は -1。
	highlightIndex int
	// selectedImage は選択中の記録のPC画像。nil なら画像なしで表示する。
	selectedImage *ebiten.Image
	imageIndex    int

	Input     *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes

	// resultBGM はリザルトから引き継いだBGM。タイトルへ戻るときに止める。
	resultBGM *audio.Player
}

// NewRankingScene は RankingScene を生成する。
// CanvasはSceneを参照するため、生成関数 newCanvas を受け取り、Sceneの生成直後に呼んで注入する。
func NewRankingScene(worldDrawer Utils.Factory[WorldDrawer], newCanvas func(*RankingScene) (RankingCanvasAPI, error), Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes, store RankingStore) Utils.Factory[*RankingScene] {
	return func() (*RankingScene, error) {
		camera, err := worldDrawer()
		if err != nil {
			return &RankingScene{}, fmt.Errorf("worldDrawer() でエラーが発生しました。 \n%w", err)
		}
		if newCanvas == nil {
			return &RankingScene{}, fmt.Errorf("RankingSceneのcanvas生成関数が nil です")
		}
		if store == nil {
			return &RankingScene{}, fmt.Errorf("RankingSceneのstoreが nil です")
		}

		rs := RankingScene{
			entities:       []Entity{},
			worldDrawer:    camera,
			store:          store,
			selector:       selection.NewGraph[int, selection.Direction](0),
			highlightIndex: -1,
			imageIndex:     -1,
			Input:          Input,
			gamedata:       gdata,
			initScenes:     initScenes,
		}
		canvas, err := newCanvas(&rs)
		if err != nil {
			return &RankingScene{}, fmt.Errorf("RankingSceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &RankingScene{}, fmt.Errorf("RankingSceneのcanvas が nil です")
		}
		rs.canvas = canvas
		return &rs, nil
	}
}

// Records は表示中の記録を順位順に返す。
func (s *RankingScene) Records() []ranking.Entry { return s.records }

// SelectedIndex はカーソルが指している順位(0始まり)を返す。
func (s *RankingScene) SelectedIndex() int { return s.selector.Current() }

// SelectionRevision はカーソルが動くたびに増える整数を返す。Canvasが変化を観測するために使う。
func (s *RankingScene) SelectionRevision() int { return s.selector.Revision() }
// HighlightIndex は今回のスコアが入った順位(0始まり)を返す。圏外は -1。
func (s *RankingScene) HighlightIndex() int { return s.highlightIndex }

// SelectedImage は選択中の記録のPC画像を返す。作れなかった場合は nil。
func (s *RankingScene) SelectedImage() *ebiten.Image { return s.selectedImage }

func (s RankingScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *RankingScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
}

// Update はSceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *RankingScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	nextScene, transitionType, err = s.updateScene(active)
	if err != nil {
		return nil, transition.None(), err
	}
	if err := s.canvas.Update(active); err != nil {
		return nil, transition.None(), fmt.Errorf("Canvasの更新に失敗しました。\n%w", err)
	}
	return
}

func (s *RankingScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if s.gamedata.InitialProcess {
		s.gamedata.InitialProcess = false
		if err := s.loadRecords(); err != nil {
			return nil, transition.None(), err
		}
		// 次のフレームから入力を受け付ける
		return
	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	if s.resultBGM != nil {
		s.resultBGM.Play()
	}

	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
		s.selector.Fire(selection.Up)
		s.refreshImage()
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
		s.selector.Fire(selection.Down)
		s.refreshImage()
	}

	if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) || s.Input.ActionIsJustPressed(inputconfig.ActionSpace) {
		if s.resultBGM != nil {
			s.resultBGM.Pause()
		}
		return backToTitleFactory(s.gamedata, s.initScenes), transition.NewScene(), nil
	}

	return
}

// loadRecords は保存済みの記録を読み込み、カーソルを初期位置へ置く。
// 読み込みに失敗しても空の一覧で続行する。
func (s *RankingScene) loadRecords() error {
	records, err := s.store.Entries()
	if err != nil {
		log.Printf("ランキングの読み込みに失敗しました: %v", err)
	}
	s.records = records
	s.imageIndex = -1

	// 上下に隣り合う順位へ移る。端では止まる（ループしない）。
	count := max(len(records), 1)
	start := 0
	if s.highlightIndex >= 0 && s.highlightIndex < count {
		start = s.highlightIndex
	}
	s.selector = selection.NewGraph[int, selection.Direction](start)
	for i := 0; i+1 < count; i++ {
		s.selector.AddTransition(i, selection.Down, i+1, nil)
		s.selector.AddTransition(i+1, selection.Up, i, nil)
	}
	s.refreshImage()
	return nil
}

// refreshImage は選択中の記録のPC画像を作り直す。選択が変わっていなければ何もしない。
func (s *RankingScene) refreshImage() {
	index := s.selector.Current()
	if index == s.imageIndex {
		return
	}
	s.imageIndex = index
	s.selectedImage = nil
	if index < 0 || index >= len(s.records) {
		return
	}
	var appearance parts.Appearance
	for i, v := range s.records[index].Appearance {
		appearance[i] = parts.AssemblePartsImage(v)
	}
	img, err := parts.AssembledImageFromAppearance(appearance)
	if err != nil {
		log.Printf("ランキングのPC画像の構築に失敗しました: %v", err)
		return
	}
	s.selectedImage = img
}

func (s *RankingScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s RankingScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
