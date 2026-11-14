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

// PlayCanvasAPI は PlayScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructから注入する。
// TimeUp後（Play_Select）のパーツ一覧の表示は、CanvasがGamedata.Stateを観測して行う。
type PlayCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsPlayCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsPlayCanvas()

	// ---- 入力用（一覧が非表示のときは何もしない） ----

	// TogglePartsList は、閉じていれば発見・所持パーツで一覧を構築して開き、開いていれば閉じる。
	TogglePartsList() error
	// MovePartSelectionDown は購入対象の選択を1つ下へ移す。
	MovePartSelectionDown()
	// MovePartSelectionUp は購入対象の選択を1つ上へ移す。
	MovePartSelectionUp()
	// ShowNextPartsTab は発見／所持タブを次へ切り替える。
	ShowNextPartsTab() error
	// ShowPreviousPartsTab は発見／所持タブを前へ切り替える。
	ShowPreviousPartsTab() error

	// ---- 問い合わせ ----

	// SelectedFoundPart は選択中の発見済みパーツを返す。
	// 一覧の非表示中、発見タブ以外の表示中、未選択のときは nil を返す。
	// 購入はSceneが Gamedata.BuyFoundPart で行い、購入後の一覧の再構築はCanvasが観測して行う。
	SelectedFoundPart() *game.Part

	// IsPartsListOpen はパーツ一覧を表示中かを返す。
	IsPartsListOpen() bool
}

// メインの遊べるシーン
type PlayScene struct {
	Stage       Stage
	entities    []Entity
	worldDrawer WorldDrawer

	canvas PlayCanvasAPI

	// text       *ui.Text
	Input      *input.Handler
	Gamedata   *game.Gamedata
	initScenes *InitializedScenes
	// テスト用
	cameraCenter Utils.Position
}

func NewPlayScene(worldDrawer Utils.Factory[WorldDrawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes, stage Stage) Utils.Factory[*PlayScene] {
	return func() (*PlayScene, error) {
		camera, err := worldDrawer()
		if err != nil {
			return &PlayScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		ps := PlayScene{
			entities:    []Entity{},
			worldDrawer: camera,
			// text:       text,
			Input:      Input,
			Gamedata:   gdata,
			initScenes: initScenes,
			Stage:      stage,
		}
		return &ps, nil
	}
}

// SetCanvas はCanvasを注入する。CanvasがSceneを参照できるよう、Sceneの生成後に呼ぶ。
func (s *PlayScene) SetCanvas(canvas PlayCanvasAPI) {
	s.canvas = canvas
}

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T(Stage : %s)", s, Stage2String[s.Stage])
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	if s.canvas != nil {
		s.canvas.Draw(screen)
	}
}

// Update implements [Scene].
// Sceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *PlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}
	if s.canvas == nil {
		return nil, transition.None(), fmt.Errorf("PlaySceneのcanvas が設定されていません")
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

func (s *PlayScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if s.Gamedata.InitialProcess {
		s.Gamedata.InitialProcess = false
		s.Gamedata.Text = "10"
		// s.text.SetText("探索フェーズ 残り時間: " + s.gamedata.Text)

	}

	if s.Gamedata.State == game.Play_Explore {
		if s.Input.ActionIsJustPressed(inputconfig.ActionPartsList) {
			if err := s.canvas.TogglePartsList(); err != nil {
				return nil, transition.None(), fmt.Errorf("パーツ一覧の開閉に失敗しました。\n%w", err)
			}
		}

		// 一覧の表示中は一覧の操作だけを受け付け、探索を止める
		if s.canvas.IsPartsListOpen() {
			if err := s.handlePartsListInput(); err != nil {
				return nil, transition.None(), err
			}
			return
		}

		for _, entity := range s.entities {
			err := entity.Update(active)
			if err != nil {
				return nil, transition.None(), fmt.Errorf("%s.Update(%s) でエラーが発生しました。\n%w", entity.Name(), strconv.FormatBool(active), err)
			}
		}
		s.removeExpiredEntities()

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
			// パーツ一覧はCanvasがPlay_Selectへの変化を観測して表示する
			s.Gamedata.State = game.Play_Select
		}

		return
	}

	if s.Gamedata.State == game.Play_Select {
		if err := s.handlePartsListInput(); err != nil {
			return nil, transition.None(), err
		}

		if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) {
			part := s.canvas.SelectedFoundPart()
			switch s.Gamedata.BuyFoundPart(part) {
			case game.NoItem:
				fmt.Print("購入可能な商品がありません")
			case game.Lack_money:
				fmt.Print("所持金が不足しています")
			case game.Lack_weight:
				fmt.Print("重量オーバーです")
			case game.Buy:
				fmt.Printf("%s を購入した", part.Name)
			default:
				fmt.Print("不正な値がBuyFoundPart()から渡されました")
			}
		}

		if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
			return func() (game.Scene, error) {
				nextScene, err := s.initScenes.InitSelectStageScene()
				if err != nil {
					return nil, err
				}
				nextScene.gamedata = s.Gamedata
				nextScene.gamedata.ForNewScene(game.SelectStage_Default)
				nextScene.initScenes = s.initScenes
				return nextScene, nil
			}, transition.NewScene(), nil
		}
	}

	return
}

// handlePartsListInput はパーツ一覧のタブ切り替えと選択移動の入力を処理する。
func (s *PlayScene) handlePartsListInput() error {
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveLeft) {
		if err := s.canvas.ShowPreviousPartsTab(); err != nil {
			return fmt.Errorf("パーツ一覧タブの切り替えに失敗しました。\n%w", err)
		}
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveRight) {
		if err := s.canvas.ShowNextPartsTab(); err != nil {
			return fmt.Errorf("パーツ一覧タブの切り替えに失敗しました。\n%w", err)
		}
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
		s.canvas.MovePartSelectionUp()
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
		s.canvas.MovePartSelectionDown()
	}
	return nil
}

func (s *PlayScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s *PlayScene) ShowPartPopup(part *game.Part, position Utils.Position) {
	if part == nil || part.Img == nil {
		return
	}

	s.AddEntity(NewPartPopup(part.Img, position))
}

func (s *PlayScene) removeExpiredEntities() {
	active := s.entities[:0]
	for _, entity := range s.entities {
		if popup, ok := entity.(interface{ IsExpired() bool }); ok && popup.IsExpired() {
			continue
		}
		active = append(active, entity)
	}
	s.entities = active
}

func (s PlayScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
