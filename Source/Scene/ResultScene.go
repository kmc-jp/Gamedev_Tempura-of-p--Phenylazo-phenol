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

// ResultCanvasAPI は ResultScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructからコンストラクタ経由で注入する。
//
// Canvasは次をScene・Gamedataの観測で行う。
//   - Result_SelectParts 中に CandidatesRevision が変わったら、CandidateParts で一覧を作り直して表示する
//   - Result_Assembled に変わったら一覧を閉じ、組み立て済みパーツの表示を開始する
//   - Result_Benchmark に変わったら、ベンチマーク案内のアニメーションを開始する
type ResultCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsResultCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsResultCanvas()

	// ---- 入力用（一覧が非表示のときは何もしない） ----

	// MovePartSelectionDown は組み立てに使うパーツの選択を1つ下へ移す。
	MovePartSelectionDown()
	// MovePartSelectionUp は組み立てに使うパーツの選択を1つ上へ移す。
	MovePartSelectionUp()

	// ---- 問い合わせ ----

	// SelectedPart は選択中のパーツを返す。未選択なら nil。
	SelectedPart() *game.Part
	// IsInstalledPartsFinished は組み立て済みパーツの表示アニメーションが完了したかを返す。
	IsInstalledPartsFinished() bool
}

type ResultScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas        ResultCanvasAPI
	partKindIndex int
	// candidates は現在の選択候補。refreshCandidates で計算してキャッシュする。
	candidates game.Parts
	// candidatesRevision は選択候補が変わるたびに増える。Canvasが観測して一覧を作り直す。
	candidatesRevision int
	// text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

// NewResultScene は ResultScene を生成する。
// CanvasはSceneを参照するため、生成関数 newCanvas を受け取り、Sceneの生成直後に呼んで注入する。
func NewResultScene(NewWorldDrawer Utils.Factory[WorldDrawer], newCanvas func(*ResultScene) (ResultCanvasAPI, error), Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*ResultScene] {
	return func() (*ResultScene, error) {
		if newCanvas == nil {
			return &ResultScene{}, fmt.Errorf("ResultSceneのcanvas生成関数が nil です")
		}
		camera, err := NewWorldDrawer()
		if err != nil {

			return &ResultScene{}, fmt.Errorf("NewDrawer() でエラーが発生しました。 \n%w", err)
		}
		ts := ResultScene{
			entities:    []Entity{},
			worldDrawer: camera,
			// text:        text,
			Input:      Input,
			gamedata:   gdata,
			initScenes: initScenes,
		}
		canvas, err := newCanvas(&ts)
		if err != nil {
			return &ResultScene{}, fmt.Errorf("ResultSceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &ResultScene{}, fmt.Errorf("ResultSceneのcanvas が nil です")
		}
		ts.canvas = canvas
		return &ts, nil
	}
}

func (s ResultScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *ResultScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
	// s.text.Draw(screen)
}

// Update はSceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *ResultScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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

func (s *ResultScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if s.gamedata.InitialProcess {
		//シーン起動時処理
		s.gamedata.InitialProcess = false
		s.gamedata.Assemble = game.Parts{}
		s.partKindIndex = 0
		if err := s.refreshCandidates(); err != nil {
			return nil, transition.None(), err
		}
		// 候補一覧はこのフレームの最後にCanvasが構築するため、入力は次のフレームから受け付ける
		return
	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	if s.gamedata.State == game.Result_SelectParts {
		//組み立てに使うパーツを選択
		if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
			s.canvas.MovePartSelectionUp()
		}
		if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
			s.canvas.MovePartSelectionDown()
		}
		if s.isSelectingArbitraryParts() && s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
			// 一覧を閉じて組み立て済みパーツの表示を開始するのは、Canvasが状態変化を観測して行う
			s.gamedata.State = game.Result_Assembled
			return
		}
		if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) {
			part := s.canvas.SelectedPart()
			if part == nil {
				if s.isSelectingArbitraryParts() {
					return
				}
				return nil, transition.None(), fmt.Errorf("選択されたパーツがありません")
			}
			s.gamedata.Assemble = append(s.gamedata.Assemble, part)
			if !s.isSelectingArbitraryParts() {
				s.partKindIndex++
			}
			if err := s.refreshCandidates(); err != nil {
				return nil, transition.None(), err
			}
		}
	}

	// 組み立てPC観察フェーズ
	if s.gamedata.State == game.Result_Assembled {
		//所持パーツをもとに利用する絵を決定
		//完成PCを表示

		if s.canvas.IsInstalledPartsFinished() {
			if err := s.gamedata.NextState(); err != nil {
				return nil, transition.None(), fmt.Errorf("InstalledPartsUI完了後の状態遷移に失敗しました: %w", err)
			}
		}
	}

	// ベンチマークフェーズ
	if s.gamedata.State == game.Result_Benchmark {
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

// CandidatesRevision は選択候補の版を返す。候補が変わるたびに増える。
// Canvasはこの値の変化を観測し、変化したときだけ CandidateParts で一覧を作り直す。
func (s *ResultScene) CandidatesRevision() int {
	return s.candidatesRevision
}

// CandidateParts は現在のパーツ種別で組み立てに使える候補を返す。
// 候補は refreshCandidates でキャッシュした値で、呼び出しに副作用はない。
func (s *ResultScene) CandidateParts() game.Parts {
	return s.candidates
}

// refreshCandidates は現在のパーツ種別の候補を計算してキャッシュし、CandidatesRevision を進める。
// 候補が変わる時点（シーン開始時、パーツを選んだ後）に1回だけ呼ぶ。
func (s *ResultScene) refreshCandidates() error {
	candidates, err := s.computeCandidates()
	if err != nil {
		return fmt.Errorf("結果画面のパーツ候補の計算に失敗しました。\n%w", err)
	}
	s.candidates = candidates
	s.candidatesRevision++
	return nil
}

// computeCandidates は現在のパーツ種別で組み立てに使える候補を計算する。
// 該当する所持パーツがない種別では、その種別のTrashパーツを1つ生成して返す。
func (s *ResultScene) computeCandidates() (game.Parts, error) {
	if s.partKindIndex < 0 || s.partKindIndex > len(game.AvailableParts) {
		return nil, fmt.Errorf("パーツ種別の選択位置が不正です: %d", s.partKindIndex)
	}

	available := game.Parts{}
	if s.isSelectingArbitraryParts() {
		selected := s.selectedArbitraryParts()
		for _, part := range s.gamedata.Parts {
			if part != nil && !s.isAvailablePart(part) && selected[part] == 0 {
				available = append(available, part)
			}
		}
		return available, nil
	}

	kind := game.AvailableParts[s.partKindIndex]
	for _, part := range s.gamedata.Parts {
		if part != nil && part.Kind == kind {
			available = append(available, part)
		}
	}
	if len(available) == 0 {
		trash := partsdata.Trash()
		if trash == nil {
			return nil, fmt.Errorf("Trashパーツの生成に失敗しました")
		}
		trash.Kind = kind
		available = append(available, trash)
	}
	return available, nil
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
