package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

// PartsListTab は PlayScene から見たパーツ一覧のタブである。
// UIのタブ型とは独立に定義し、Canvasが変換して返す。
type PartsListTab uint8

const (
	// PartsListTabNone は一覧の非表示中を表す。
	PartsListTabNone PartsListTab = iota
	// PartsListTabFound は「発見したパーツ」タブを表す。
	PartsListTabFound
	// PartsListTabOwned は「所持パーツ」タブを表す。
	PartsListTabOwned
	// PartsListTabExit は「店を出る」タブを表す。
	PartsListTabExit
)

// PlayCanvasAPI は PlayScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructからコンストラクタ経由で注入する。
// TimeUp後（Play_Select）のパーツ一覧の表示は、Canvasが PlayScene.SelectListRevision を観測して行う。
type PlayCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsPlayCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsPlayCanvas()

	// ---- 入力用（一覧が非表示のときは何もしない） ----

	// TogglePartsList は、閉じていれば発見・所持パーツで一覧を構築して開き、開いていれば閉じる。
	TogglePartsList() error
	// ShowFoundPartLog は発見したパーツ名をプレイヤー頭上に表示する。
	ShowFoundPartLog(text string)
	// MovePartSelectionDown は購入対象の選択を1つ下へ移す。
	MovePartSelectionDown()
	// MovePartSelectionUp は購入対象の選択を1つ上へ移す。
	MovePartSelectionUp()
	// ShowNextPartsTab は発見／所持タブを次へ切り替える。
	ShowNextPartsTab() error
	// ShowPreviousPartsTab は発見／所持タブを前へ切り替える。
	ShowPreviousPartsTab() error
	// ShowPurchaseConfirmation は選択中の発見パーツの購入確認を表示する。
	ShowPurchaseConfirmation() bool
	MovePurchaseConfirmationLeft()
	MovePurchaseConfirmationRight()
	ConfirmPurchase()
	IsPurchaseConfirmationOpen() bool
	PurchaseConfirmationResult() (bool, bool)
	HidePurchaseConfirmation()

	// ShowExitStoreConfirmation は退店確認を表示する。
	ShowExitStoreConfirmation() bool
	MoveExitStoreConfirmationLeft()
	MoveExitStoreConfirmationRight()
	ConfirmExitStore()
	IsExitStoreConfirmationOpen() bool
	ExitStoreConfirmationResult() (bool, bool)
	HideExitStoreConfirmation()

	SetMessageWindowText(string)
	SetActionLogText(string)

	// ---- 問い合わせ ----

	// SelectedFoundPart は選択中の発見済みパーツを返す。
	// 一覧の非表示中、発見タブ以外の表示中、未選択のときは nil を返す。
	// 購入はSceneが Gamedata.BuyFoundPart で行い、購入後の一覧の再構築はCanvasが観測して行う。
	SelectedFoundPart() *game.Part

	// IsPartsListOpen はパーツ一覧を表示中かを返す。
	IsPartsListOpen() bool

	// CurrentPartsTab は表示中のタブを返す。一覧の非表示中は PartsListTabNone を返す。
	CurrentPartsTab() PartsListTab
}

// メインの遊べるシーン
type PlayScene struct {
	// Stage は、このシーンが遊ばせるステージのある駅。
	Stage       stageinfo.Location
	entities    []Entity
	worldDrawer WorldDrawer

	canvas PlayCanvasAPI
	// selectListRevision は Play_Select に入るたびに増える。Canvasが観測してパーツ一覧を開く。
	selectListRevision int

	timer                     *tick.Timer
	MessageWindowDisplayTimer *tick.Timer

	// text       *ui.Text
	Input      *input.Handler
	Gamedata   *game.Gamedata
	initScenes *InitializedScenes
	// テスト用
	cameraCenter             Utils.Position
	purchaseConfirmationPart *game.Part
	exitStoreFromExplore     bool
}

// NewPlayScene は PlayScene を生成する。
// CanvasはSceneを参照するため、生成関数 newCanvas を受け取り、Sceneの生成直後に呼んで注入する。
func NewPlayScene(worldDrawer Utils.Factory[WorldDrawer], newCanvas func(*PlayScene) (PlayCanvasAPI, error), Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes, stage stageinfo.Location) Utils.Factory[*PlayScene] {
	return func() (*PlayScene, error) {
		if newCanvas == nil {
			return &PlayScene{}, fmt.Errorf("PlaySceneのcanvas生成関数が nil です")
		}
		camera, err := worldDrawer()
		if err != nil {
			return &PlayScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		ps := PlayScene{
			entities:    []Entity{},
			worldDrawer: camera,
			// text:       text,
			Input:                     Input,
			timer:                     tick.NewTimer(),
			MessageWindowDisplayTimer: tick.NewTimer(),
			Gamedata:                  gdata,
			initScenes:                initScenes,
			Stage:                     stage,
		}
		canvas, err := newCanvas(&ps)
		if err != nil {
			return &PlayScene{}, fmt.Errorf("PlaySceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &PlayScene{}, fmt.Errorf("PlaySceneのcanvas が nil です")
		}
		ps.canvas = canvas
		return &ps, nil
	}
}

// SelectListRevision は Play_Select に入った回数を返す。
// Canvasはこの値の変化を観測し、変化したときにパーツ一覧を開く。
func (s *PlayScene) SelectListRevision() int {
	return s.selectListRevision
}

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T(Stage : %s)", s, stageinfo.Location2String[s.Stage])
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
}

// Update implements [Scene].
// Sceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *PlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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

func (s *PlayScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if s.Gamedata.State == game.Play_Explore {
		if s.canvas.IsPartsListOpen() {
			s.Gamedata.DisableWorldTimeProgressing()
		} else {
			s.Gamedata.EnableWorldTimeProgressing()
		}
	} else {
		s.Gamedata.DisableWorldTimeProgressing()
	}

	if s.MessageWindowDisplayTimer.Expired() {
		s.canvas.SetMessageWindowText("")
	}

	if s.Gamedata.InitialProcess {
		s.Gamedata.InitialProcess = false

	}

	if s.Gamedata.State == game.Play_EnterStore {
		// 入店時処理

		// 最初の店(または売却予定パーツ無し)なら売却処理をスキップ
		if s.Gamedata.Visited() <= 1 || !s.Gamedata.IsSellPartExists() {
			s.Gamedata.State = game.Play_Explore
		} else {
			s.Gamedata.State = game.Play_Sell
		}
		return nil, transition.None(), nil
	}

	if s.Gamedata.State == game.Play_Sell {
		// TrainSceneで保存したパーツを売却し、探索へ進む。
		s.timer.Set(50)
		if s.Gamedata.IsSellPartExists() {
			v, w, n := s.Gamedata.SellPart()
			s.canvas.SetMessageWindowText(fmt.Sprintf("%s を %d 円で売って %d 持ち運べる量が増えた", n, v, w))
		}
		if s.timer.Expired() {
			s.Gamedata.State = game.Play_Explore
		}

		return nil, transition.None(), nil
	}

	if s.Gamedata.State == game.Play_Explore {
		if s.Input.ActionIsJustPressed(inputconfig.ActionPartsList) {
			if err := s.canvas.TogglePartsList(); err != nil {
				return nil, transition.None(), fmt.Errorf("パーツ一覧の開閉に失敗しました。\n%w", err)
			}
			if s.canvas.IsPartsListOpen() {
				s.Gamedata.DisableWorldTimeProgressing()
			} else {
				s.Gamedata.EnableWorldTimeProgressing()
			}
		}

		// 一覧の表示中は一覧の操作だけを受け付け、探索を止める
		if s.canvas.IsPartsListOpen() {
			if err := s.handlePartsListInput(); err != nil {
				return nil, transition.None(), err
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) &&
				s.canvas.CurrentPartsTab() == PartsListTabExit {
				s.Gamedata.State = game.Play_ExitStore
				s.exitStoreFromExplore = true
				s.canvas.ShowExitStoreConfirmation()
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

		if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
			return func() (game.Scene, error) {
				pauseScene, err := s.initScenes.InitPauseScene()
				if err != nil {
					return nil, err
				}
				pauseScene.gamedata = s.Gamedata
				pauseScene.gamedata.DisableWorldTimeProgressing()
				//pauseScene.gamedata.ForNewScene(0)
				pauseScene.initScenes = s.initScenes
				return pauseScene, nil
			}, transition.Push(), nil
		}
	}

	//TimeUp後の処理
	if s.Gamedata.State == game.Play_Register {
		// パーツ一覧はCanvasが selectListRevision の変化を観測して表示する
		s.Gamedata.State = game.Play_Select
		s.selectListRevision++

		return
	}

	if s.Gamedata.State == game.Play_Select {
		s.canvas.SetMessageWindowText("Space で買うパーツを選択  Enter で決定")
		if s.canvas.IsPurchaseConfirmationOpen() {
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveLeft) {
				s.canvas.MovePurchaseConfirmationLeft()
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveRight) {
				s.canvas.MovePurchaseConfirmationRight()
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
				s.canvas.ConfirmPurchase()
			}
			accepted, confirmed := s.canvas.PurchaseConfirmationResult()
			if confirmed {
				part := s.purchaseConfirmationPart
				s.purchaseConfirmationPart = nil
				s.canvas.HidePurchaseConfirmation()
				if !accepted || part == nil {
					return
				}
				switch s.Gamedata.BuyFoundPart(part) {
				case game.NoItem:
					s.canvas.SetActionLogText("購入可能な商品がありません")
				case game.Lack_money:
					s.canvas.SetActionLogText("所持金が不足しています")
				case game.Lack_weight:
					s.canvas.SetActionLogText("重量オーバーです")
				case game.Buy:
					s.canvas.SetActionLogText(fmt.Sprintf("%s を購入した", part.Name))
				default:
					return nil, transition.None(), fmt.Errorf("不正な値がBuyFoundPart()から渡されました")
				}
			}
			return
		}

		if err := s.handlePartsListInput(); err != nil {
			return nil, transition.None(), err
		}

		// 購入できるのは発見タブの表示中だけ
		if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) && s.canvas.CurrentPartsTab() == PartsListTabFound {
			part := s.canvas.SelectedFoundPart()
			if s.Gamedata.Buyable(part) == game.Buyable && part != nil && s.canvas.ShowPurchaseConfirmation() {
				s.purchaseConfirmationPart = part
			}

		}

		if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) &&
			s.canvas.CurrentPartsTab() == PartsListTabExit {
			s.Gamedata.State = game.Play_ExitStore
			s.exitStoreFromExplore = false
			s.canvas.ShowExitStoreConfirmation()
			return
		}

	}

	if s.Gamedata.State == game.Play_ExitStore {
		if s.canvas.IsExitStoreConfirmationOpen() {
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveLeft) {
				s.canvas.MoveExitStoreConfirmationLeft()
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionMoveRight) {
				s.canvas.MoveExitStoreConfirmationRight()
			}
			if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
				s.canvas.ConfirmExitStore()
			}
			accepted, confirmed := s.canvas.ExitStoreConfirmationResult()
			if confirmed {
				if accepted {
					return func() (game.Scene, error) {
						nextScene, err := s.initScenes.InitSelectStageScene()
						if err != nil {
							return nil, err
						}
						nextScene.gamedata = s.Gamedata
						nextScene.gamedata.ForNewScene(game.SelectStage_Select)
						nextScene.initScenes = s.initScenes
						return nextScene, nil
					}, transition.NewScene(), nil
				} else {
					s.canvas.HideExitStoreConfirmation()
					if s.exitStoreFromExplore {
						s.Gamedata.State = game.Play_Explore
						s.Gamedata.DisableWorldTimeProgressing()
					} else {
						s.Gamedata.State = game.Play_Select
					}
					s.exitStoreFromExplore = false
				}
			}
		}

	}

	if s.Gamedata.State == game.Play_Timeup {
		// worldtime時間切れの処理
		return
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

func (s *PlayScene) ShowFoundPartLog(part *game.Part) {
	if part == nil {
		return
	}
	s.canvas.ShowFoundPartLog(fmt.Sprintf("%s を見つけた", part.Name))
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

func (s *PlayScene) Canvas() PlayCanvasAPI {
	return s.canvas
}
