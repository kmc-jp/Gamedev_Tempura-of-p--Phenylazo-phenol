package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	selection "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Selection"
	"fmt"
	"maps"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

// SelectStageCanvasAPI は SelectStageScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructからコンストラクタ経由で注入する。
// 路線図の表示は、Canvasが SelectStageScene の公開メソッド
// （StageMap / SelectedLocation / CurrentLocation）を観測して行う。
type SelectStageCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsSelectStageCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsSelectStageCanvas()

	SetMessageWindowText(string)
	SetActionLogText(string)
}

// SelectStageRule は行き先選択のルールである。WithRule で NewSelectStageScene の生成時に渡して決める。
// 各関数が nil のとき（ゼロ値を含む）は制限なし。
// 例えば「現在地は選べない」なら次のようにする。
//
//	scene.WithRule(scene.SelectStageRule{
//		CanConfirm: func(loc stageinfo.Location) bool { return loc != gdata.CurrentLocation() },
//	})
type SelectStageRule struct {
	// CanMove は、カーソルを from から to へ動かせるかを返す。
	// 路線でつながっていても、false を返した移動はできなくなる。
	CanMove func(from, to stageinfo.Location) bool
	// CanConfirm は、loc を行き先として確定できるかを返す。
	CanConfirm func(loc stageinfo.Location) bool
}

// SelectStageOption は NewSelectStageScene の生成時に渡すオプションである。
// 何も渡さなければ制限なし（全ての駅へ動けて、全ての駅を確定できる）。
type SelectStageOption func(*SelectStageScene)

// WithRule は行き先選択のルールを設定する。
func WithRule(rule SelectStageRule) SelectStageOption {
	return func(s *SelectStageScene) { s.rule = rule }
}

type SelectStageScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas SelectStageCanvasAPI

	// stageMap は路線図の定義。Canvasが描画に使う。
	stageMap stageinfo.StageMap
	// selector は路線図上のカーソル（選択中の駅）と、方向キーでの移動先を管理する。
	selector *selection.Graph[stageinfo.Location, selection.Direction]
	// rule は行き先選択のルール。ゼロ値は制限なし。
	rule SelectStageRule

	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

// NewSelectStageScene は SelectStageScene を生成する。
// CanvasはSceneを参照するため、生成関数 newCanvas を受け取り、Sceneの生成直後に呼んで注入する。
// 行き先選択のルールは opts で生成時に決める（WithRule。渡さなければ制限なし）。
func NewSelectStageScene(
	worldDrawer Utils.Factory[WorldDrawer],
	stageMap stageinfo.StageMap,
	newCanvas func(*SelectStageScene) (SelectStageCanvasAPI, error),
	Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes,
	opts ...SelectStageOption,
) Utils.Factory[*SelectStageScene] {
	return func() (*SelectStageScene, error) {
		if newCanvas == nil {
			return &SelectStageScene{}, fmt.Errorf("SelectStageSceneのcanvas生成関数が nil です")
		}
		if gdata == nil {
			return &SelectStageScene{}, fmt.Errorf("SelectStageSceneのgdata が nil です")
		}
		camera, err := worldDrawer()
		if err != nil {
			return &SelectStageScene{}, fmt.Errorf("worldDrawer() でエラーが発生しました。 \n%w", err)
		}
		if initScenes == nil {
			return &SelectStageScene{}, fmt.Errorf("SelectStageSceneのinitScenes が nil です")
		}
		// 自宅以外の駅は、行き先として確定するとそのステージへ向かう。ステージが無い駅は入った時点で気付けるようにする。
		for _, loc := range slices.Sorted(maps.Keys(stageMap.Stations)) {
			if loc == stageinfo.Departure {
				continue
			}
			if _, ok := initScenes.InitPlayScene[loc]; !ok {
				return &SelectStageScene{}, fmt.Errorf("駅 %s に対応するステージがありません", stageinfo.Location2String[loc])
			}
		}
		selector, err := stageMap.NewSelector(gdata.CurrentLocation())
		if err != nil {
			return &SelectStageScene{}, fmt.Errorf("stageMap.NewSelector() でエラーが発生しました。\n%w", err)
		}
		ts := SelectStageScene{
			entities:    []Entity{},
			worldDrawer: camera,
			stageMap:    stageMap,
			selector:    selector,
			Input:       Input,
			gamedata:    gdata,
			initScenes:  initScenes,
		}
		for _, opt := range opts {
			if opt != nil {
				opt(&ts)
			}
		}
		if ts.rule.CanMove != nil {
			canMove := ts.rule.CanMove
			selector.SetGuard(func(from stageinfo.Location, _ selection.Direction, to stageinfo.Location) bool {
				return canMove(from, to)
			})
		}
		canvas, err := newCanvas(&ts)
		if err != nil {
			return &SelectStageScene{}, fmt.Errorf("SelectStageSceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &SelectStageScene{}, fmt.Errorf("SelectStageSceneのcanvas が nil です")
		}
		ts.canvas = canvas
		return &ts, nil
	}
}

func (s SelectStageScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// StageMap は路線図の定義を返す。
func (s *SelectStageScene) StageMap() stageinfo.StageMap {
	return s.stageMap
}

// SelectedLocation はカーソルが指している駅を返す。
func (s *SelectStageScene) SelectedLocation() stageinfo.Location {
	return s.selector.Current()
}

// SelectionRevision はカーソルが動くたびに増える整数を返す。
// Canvasはこの値の変化を観測し、変化したときに表示を更新する。
func (s *SelectStageScene) SelectionRevision() int {
	return s.selector.Revision()
}

// CurrentLocation はプレイヤーが今いる駅（直前に行き先として確定した駅）を返す。
func (s *SelectStageScene) CurrentLocation() stageinfo.Location {
	return s.gamedata.CurrentLocation()
}

// CanConfirm は loc を行き先として確定できるかを返す。副作用はない。
func (s *SelectStageScene) CanConfirm(loc stageinfo.Location) bool {
	if s.rule.CanConfirm == nil {
		return true
	}
	return s.rule.CanConfirm(loc)
}

func (s *SelectStageScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
}

// Update はSceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *SelectStageScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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

func (s *SelectStageScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if s.gamedata.InitialProcess {
		//シーン起動時処理
		s.gamedata.InitialProcess = false
		// 遷移元が現在地を更新している場合に備え、カーソルを現在地に合わせる
		s.selector.SetCurrent(s.gamedata.CurrentLocation())

		// 時間切れ、または行ける場所がないときは、行先選択を出さずに帰宅する
		if !s.hasDestination() {
			return s.mkResultScene(), transition.NewScene(), nil
		}
	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	// カーソル移動（路線でつながっていない方向は何も起きない）
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveLeft) {
		s.selector.Fire(selection.Left)
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveRight) {
		s.selector.Fire(selection.Right)
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
		s.selector.Fire(selection.Up)
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
		s.selector.Fire(selection.Down)
	}

	dst, time, _ := s.destination()
	s.canvas.SetMessageWindowText(fmt.Sprintf("%s: 移動時間 %d 分  Enterで出発", stageinfo.Location2String[dst], time))

	// 決定
	if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
		dst, min, err := s.destination()
		if !s.CanConfirm(dst) {
			s.canvas.SetActionLogText("今いる場所には移動できません")
			return nil, transition.None(), nil
		}
		if dst == stageinfo.Departure {
			return s.confirm()
		}
		if err != nil {
			return nil, transition.None(), fmt.Errorf("s.destination() でエラーが発生しました\n%w\n", err)
		}
		if !s.gamedata.WorldTimer.EnoughTimeLeft(min) {
			s.canvas.SetActionLogText("時間が足りません")
			return nil, transition.None(), nil
		}
		s.gamedata.WorldTimer.Skip(uint16(min))
		return s.confirm()
	}

	return
}

// destination はカーソルが指している駅（移動先）と、現在地からそこまでの最短所要時間（分）を返す。
func (s *SelectStageScene) destination() (loc stageinfo.Location, minutes uint64, err error) {
	loc = s.selector.Current()
	minutes, err = s.stageMap.ShortestMinutes(s.gamedata.CurrentLocation(), loc)
	if err != nil {
		return loc, 0, fmt.Errorf("最短所要時間の計算に失敗しました。\n%w", err)
	}
	return loc, minutes, nil
}

// hasDestination は、自宅と現在地を除いて、残り時間で行ける駅があるかを返す。
// 時間切れ（残り0分）のときは常に false になる。
func (s *SelectStageScene) hasDestination() bool {
	if s.gamedata.WorldTimer.TimeLeft() == 0 {
		return false
	}
	current := s.gamedata.CurrentLocation()
	for loc := range s.stageMap.Stations {
		if loc == stageinfo.Departure || loc == current || !s.CanConfirm(loc) {
			continue
		}
		minutes, err := s.stageMap.ShortestMinutes(current, loc)
		if err != nil {
			continue
		}
		if s.gamedata.WorldTimer.EnoughTimeLeft(minutes) {
			return true
		}
	}
	return false
}

// confirm はカーソルが指している駅を行き先として確定する。
// 自宅なら帰宅（Result）へ、それ以外ならその駅のステージへ向かう（Train を経由する）。
func (s *SelectStageScene) confirm() (Utils.Factory[game.Scene], transition.Type, error) {
	loc, minutes, err := s.destination()
	if err != nil {
		return nil, transition.None(), err
	}
	if !s.CanConfirm(loc) {
		return nil, transition.None(), nil
	}
	fmt.Printf("所要時間 : %d\n", minutes)
	// TODO: 移動の所要時間（分）として使う

	if loc == stageinfo.Departure {
		return s.mkResultScene(), transition.NewScene(), nil
	}

	if _, ok := s.initScenes.InitPlayScene[loc]; !ok {
		return nil, transition.None(), fmt.Errorf("駅 %s に対応するステージがありません", stageinfo.Location2String[loc])
	}
	s.gamedata.SetNextStage(s.mkStage(loc))
	s.gamedata.SetCurrentLocation(loc)
	return s.mkTrainScene(), transition.NewScene(), nil
}

func (s *SelectStageScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s SelectStageScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}

func (s SelectStageScene) mkStage(stage stageinfo.Location) Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		nextScene, err := s.initScenes.InitPlayScene[stage]()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitPlayScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.Gamedata = s.gamedata
		nextScene.Gamedata.ForNewScene(game.Play_EnterStore)
		nextScene.initScenes = s.initScenes
		return nextScene, nil
	}
}

func (s SelectStageScene) mkTrainScene() Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		nextScene, err := s.initScenes.InitTrainScene()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitTrainScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.Gamedata = s.gamedata
		nextScene.Gamedata.ForNewScene(game.Train_Default)
		nextScene.initScenes = s.initScenes
		return nextScene, nil
	}
}

func (s SelectStageScene) mkResultScene() Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		nextScene, err := s.initScenes.InitResultScene()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitResultScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.gamedata = s.gamedata
		nextScene.gamedata.ForNewScene(game.Result_SelectParts)
		nextScene.initScenes = s.initScenes
		return nextScene, nil
	}
}
