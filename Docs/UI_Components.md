# UI コンポーネント仕様

## 共通仕様

`Source/UI` の UI 要素は、`canvas` package（`Source/Canvas`）の `CanvasItem` インターフェースを実装する。
位置を持つ UI 要素（`Text`、`AnimatableUI`）は `PositionItem` も実装する。

```go
type CanvasItem interface {
    Update(active bool) error
    Draw(screen *ebiten.Image)
}

type PositionItem interface {
    CanvasItem
    Position() utils.Position
    SetPosition(utils.Position)
}
```

`ui` package は `scene` package を import しない。
UI 要素は Scene 固有 Canvas（`scenecanvas`）が保持し、Scene から直接は参照されない。

## Canvas

`canvas.Canvas` は Scene を知らない汎用 UI コンテナであり、追加された `CanvasItem` の更新・描画・表示状態を管理する。
Scene 固有の構成は `scenecanvas` package の各 Canvas が決める。

- `NewCanvas() *Canvas`: 表示状態の空の Canvas を生成する
- `Add(item CanvasItem) error`: Item を末尾に追加する。`nil` はエラー
- `Update(active bool) error`: 追加順に各 Item の `Update` を呼ぶ。非表示の間は何もしない
- `Draw(screen)`: 追加順に各 Item を描画する（後に追加した Item ほど手前）。非表示の間は何も描画しない
- `Show()` / `Hide()` / `Visible() bool`: Canvas 全体の表示状態（現在は未使用だが、後で使う可能性があるため残している）

`Canvas` 自身も `CanvasItem` を満たすため、Canvas を入れ子にできる。

## Scene 固有 Canvas（`scenecanvas`）

`Source/SceneCanvas` の各 Canvas は、対応する Scene のファイルに定義された Canvas API（`scene.TitleCanvasAPI` など）を実装し、`construct` から Scene へ注入される。

- `TitleCanvas` / `PauseCanvas`: 案内用の `Text` 1個だけを持つ。`NewXxxCanvas(font)` で生成し、`NewXxxScene` の引数で注入する
- `SelectStageCanvas`: 路線図（`RouteMapUI`）と案内 `Text`。`NewSelectStageCanvas(scene, gdata)`。`SelectStageScene.SelectionRevision`（カーソル移動のたびに増える）と `CurrentLocation`（現在地）を毎フレーム観測し、変化したときだけ `SelectedLocation` と合わせて反映する。`StageMap` から駅・路線を組み立てる。行き先選択のルール（隣接駅のみ等）は Scene の `SelectStageRule` を `WithRule` で生成時に渡して決める。ルールを入れないときは制限なし
- `PlayCanvas`: 残り時間（`GameDataText`）、空き重量・所持金（`GameDataUI`）、`PartsListUI`。`NewPlayCanvas(scene, gdata)`。`PlayScene.SelectListRevision` の変化で一覧を開く。購入はSceneが`CurrentPartsTab`で発見タブの表示中かを確かめたうえで、`SelectedFoundPart`と`Gamedata.BuyFoundPart`で行い、購入後の一覧は`Gamedata.PartsRevision`の変化を観測して再構築する
- `TrainCanvas`: 案内 `Text`、`PartsListUI`。`NewTrainCanvas(scene, gdata)`。`TrainScene.PartsListRevision`（または`Gamedata.PartsRevision`）の変化で所持パーツ一覧を構築・表示する
- `ResultCanvas`: 案内 `Text`、ベンチマーク案内（`AnimatableUI`）、`InstalledParts`、`PartsListUI`。`NewResultCanvas(scene, gdata)`。全UIを毎フレーム`Canvas.Update`で更新する。組み立て済みパーツの表示は`Result_Assembled`、ベンチマーク案内のアニメーションは`Result_Benchmark`に変わったときに開始する

Play / Train / Result / SelectStage の Canvas は Scene のポインタや Gamedata を観測して表示を切り替える。そのため、Canvas 本体ではなく Canvas の生成関数（`func(*XxxScene) (XxxCanvasAPI, error)`）を `NewXxxScene` に渡し、Scene のコンストラクタが自身の生成直後に呼んで注入する。生成関数がエラーや `nil` を返したときは Scene の生成エラーになる。

状態変化の観測にはリビジョン番号（変化するたびに増える整数）を使う。書き換える側（Scene・Gamedata）が変更時に番号を増やし、Canvas は反映済みの番号と比べて違っていれば作り直す。

- シーン固有の変化: Scene が持つ（`PlayScene.SelectListRevision`、`TrainScene.PartsListRevision`、`ResultScene.CandidatesRevision`）
- `Gamedata.Found`・`Parts` の変化: `Gamedata.PartsRevision` が持つ。シーンをまたぐ Model のデータで、Scene を通らない書き換え（Player の拾得）もあるため、書き換える Gamedata のメソッドの中で増やす

既存の `Text` と `DrawText` は維持し、新しい UI はこれらを利用する。

## RichUI

`RichUI` は任意の `*ebiten.Image` と任意の `*Text` を組み合わせて表示する。

- 画像のみ、テキストのみ、画像とテキストの両方に対応する。
- 画像を指定しサイズが `image.Point{}` の場合、画像サイズを境界サイズに使用する。
- 明示したサイズはレイアウトと当たり判定の境界に使用する。
- 画像の読み込みは `AssetsManager` または呼び出し側の責務であり、RichUI は読み込み済み画像を受け取る。

主な API:

- `NewRichUI`
- `Position`
- `Size`
- `Bounds`
- `SetPosition`

## ButtonUI

`ButtonUI` は `RichUI` をラップし、選択状態の表示だけを担当する。入力は扱わない。

- `SetSelected(selected bool)`: 選択状態を設定する。選択中は境界の内側に `uiutils.DrawBorder` でオレンジ枠（`uiutils.SelectionBorderColor`）を描画する
- `Selected() bool`: 選択状態を取得
- `RichUI() *RichUI`: 内包要素を取得

押下・選択の判定は Scene・Component が `input.Handler` で行い、結果を `SetSelected` で渡す。マウスの当たり判定が必要な場合は `RichUI().Bounds()` を使う。
旧マウス判定（`IsMouseOver`、`IsButtonPressed`）は、Scene側の入力移行は完了したが、再利用の可能性があるため理由付きのコメントアウトで残している。

## Confirm

`commonui.Confirm` は、任意メッセージと「はい」「いいえ」の2択を表示する受動的な確認UIである。背景はメッセージとボタン全体を覆い、ボタンの選択状態は既存のオレンジ枠で示す。

- `NewConfirm(ConfirmConfig)`: 位置、サイズ、余白、フォント、色、ボタンサイズ、表示文言を指定して生成する
- `Show()` / `Hide()` / `Visible() bool`: 表示状態を管理する。`Show` は選択を「はい」に戻し、未確定にする
- `SetMessage(string)`: 上部のメッセージを差し替える
- `MoveLeft()` / `MoveRight()`: Scene・Componentから選択を変更する（左が「はい」、右が「いいえ」）。入力は扱わない
- `Confirm()`: 現在の選択を確定する
- `Result() (bool, bool)`: `(結果, 確定済み)`を返す

キー入力はScene・Componentの`input.Handler`で解釈し、左右キーを`MoveLeft` / `MoveRight`、Enterを`Confirm`へ接続する。

## Clock

`commonui.Clock` は`Gamedata.WorldTimer`から取得した世界時刻を`hh:mm`形式で表示する専用UIである。通常の`GameDataText`ではなく、時計表示のデータ取得と描画を専用部品にまとめる。

- `NewClock(position, color, gdata, font)`: 時計UIを生成する
- `Position()` / `SetPosition()`: 表示位置を取得・設定する
- 時刻の計算と文字列化は`tick.WorldTimer`が担当し、Clockは描画時に現在値を参照する

## ScrollUI

`ScrollUI` は複数の `CanvasItem` をコンテンツ用オフスクリーン画像へ描画し、ビューポート部分だけを表示する。

- `ScrollVertical`: 上下方向
- `ScrollHorizontal`: 左右方向
- `AddChild`: 子要素を追加
- `Scroll`: 相対スクロール
- `SetOffset`: 絶対スクロール位置を設定
- `Offset`: 現在位置を取得
- `Bounds`: ビューポートの画面境界を取得

子要素の座標はコンテンツバッファの座標系で指定する。スクロール量は自動的に 0 から最大値までクランプされる。

ScrollUI は表示と子要素の更新だけを担当し、マウスホイールやキー入力の解釈は行わない。オフセットは ScrollUI を持つUI（`PartsListUI` など）が設定する。ボタン等のUIはマウス判定を行わない方針のため、子要素向けの座標変換は提供しない。

## RouteMapUI

`RouteMapUI[N]`（`Source/UI/RouteMapUI.go`）は、ノード（駅）とそれを結ぶ辺を路線図として描画する。`scene` を import しない汎用部品で、入力は扱わない。ノードの ID 型 `N` は呼び出し側が決める（`SelectStageCanvas` は `stageinfo.Location`）。

- `NewRouteMapUI(refSize, nodes, edges, fontSource)`: 基準サイズ・ノード（`RouteMapNode{ID, Pos, Label, Image}`）・辺（ID の組）から生成する。ID の重複、辺の端が未定義のとき、`refSize` が正でないときはエラー
- `SetSelected(id)`: 選択中のノードを設定する（オレンジの脈動するリング）。未知の ID は何もしない
- `SetMarker(id)`: 現在地のノードを設定する（ノード内側の青い円）。未知の ID は何もしない
- `SetBackground(img)`: 路線図の背景画像を設定する。基準サイズ全体に引き伸ばして描く。`nil` で背景なし
- `Update(active)`: 選択リングの脈動のためにフレームを進めるだけ

座標は `refSize`（基準サイズ）上の値で、`Draw` 時に画面サイズへアスペクト比を保ったまま拡大縮小し、中央に描く。画面サイズが変わっても座標を決め打ちしない。
`RouteMapNode.Image` が `nil` のノードは円で描く。画像を指定すると中心を `Pos` に合わせて描くため、絵が届いたらノードごとに差し替えられる。

## selection.Graph

`selection.Graph[N, T]`（`Source/Utils/Selection`）は、選択可能UIの「いま選んでいるノード」と「どの入力でどこへ移るか」を管理する共通基盤で、オートマトン（状態遷移グラフ）として扱う。描画と入力の取得は行わず、ebiten に依存しない。

- `NewGraph(initial)`: 初期ノードを指定して生成する
- `AddTransition(from, trigger, to, cond)`: 遷移を登録する。`cond` が `nil` なら無条件。同じ `(from, trigger)` に複数登録した場合は登録順に評価し、最初に条件を満たしたものを採用する
- `Fire(trigger) bool`: 現在のノードで発火する。遷移したら `true`。遷移先がない・条件を満たさないときは何もせず `false`
- `Next(trigger) (N, bool)`: `Fire` した場合の遷移先を、状態を変えずに返す
- `Current()` / `SetCurrent(n)`: 現在のノード。`SetCurrent` は値が変わったときだけ `Revision` を増やす
- `Revision() int`: 選択が変化するたびに増える整数。UIが変化を観測するために使う
- `SetGuard(guard)`: 全遷移に共通の条件（`guard(from, trigger, to) bool`）を設定する。`nil` で解除。遷移の登録をやり直さずにルールを入れ替えるための差し替え口
- `SetCondition(from, trigger, to, cond) bool`: 登録済みの遷移1本の条件を差し替える。`nil` で無条件に戻る。該当する遷移がなければ `false`

`Direction`（`Left` / `Right` / `Up` / `Down`）と `Directions()` を上下左右移動用のトリガーとして同梱している。

`SelectStageScene` は `stageinfo.StageMap.NewSelector` で路線図から `Graph` を作る。上下左右の移動先は、路線でつながった駅の位置関係から自動で決まる（方向キーの向きと 60° 以内で最も近い路線）。自動で決まるキーを変えたい路線は、`StageMap.toml` の `[[Lines]]` に `FromKey`（From の駅で To へ進むキー）/ `ToKey`（To の駅で From へ戻るキー）を `"Left"` / `"Right"` / `"Up"` / `"Down"` で書いて明示する。明示した向きは自動導出より優先され（その向きは自動では決めない）、書かない向きは自動のまま。同じ駅の同じキーに別の行き先を明示した場合と、どのキーでも選べない路線があるときは起動時エラーになる。路線図の定義は `Assets/Serialize/StageMap.toml`（`StageMapConfig`）。

行き先選択のルールは、`NewSelectStageScene` の末尾の Functional Options `scene.WithRule(scene.SelectStageRule{...})` で生成時に決める。`CanMove(from, to) bool` はカーソルの移動可否（`Graph.SetGuard` として全遷移に効く）、`CanConfirm(loc) bool` は行き先として確定できるかを決める。各関数が `nil` なら制限なしで、`WithRule` を渡さなければルールなし。渡す場所は `SelectStageSceneConstruct.go` で、途中での差し替えは想定していない（参照する `Gamedata` の状態が変われば結果は動的に変わる）。

## PartsListUI

`PartsListUI` は、パーツ一覧の表示・タブ・選択・スクロール・パーツカード（`PartCard`）の生成をまとめて担当する。
利用するのは Scene 固有 Canvas（`PlayCanvas`、`TrainCanvas`、`ResultCanvas`）だけで、Scene は Canvas API の入力用API・問い合わせAPIを経由して操作する。

- `NewPartsListUI(position, font, captionFont)`: 生成する。生成直後は非表示
- `ShowFoundParts(found, owned game.Parts) error`: 「発見したパーツ」「所持パーツ」のタブ付き一覧を構築し、発見タブを選ぶ
- `ShowParts(parts game.Parts) error`: タブなしの単一一覧を構築する
- `Show()` / `Hide()` / `Visible() bool`: 表示状態。`ShowFoundParts` / `ShowParts` は表示状態を変えないため、表示するには続けて `Show` を呼ぶ
- `SelectNext()` / `SelectPrevious()`: 選択を1つ移動する。範囲外へはクランプし、選択中のカードが見えるようスクロールを追従させる。非表示中は何もしない
- `NextTab() error` / `PreviousTab() error`: タブを切り替える。単一一覧と非表示中は何もしない
- `SelectedPart() *game.Part`: 選択中のパーツ。非表示中・未選択なら `nil`
- `IsShowingFoundParts() bool`: タブ付き一覧で発見タブを表示中か。単一一覧の表示中は `false`
- `ActiveTab() PartsListTab`: 現在のタブ。表示状態やタブ表示の有無によらず、最後に選ばれたタブを返す
- `RefreshFoundParts(found, owned game.Parts) error`: タブ付き一覧のデータを差し替え、現在のタブと選択位置を維持して再構築する

非表示中に操作を受け付けない判定は `PartsListUI` が持つため、Canvas 側で `Visible()` を確かめてから呼ぶ必要はない。一覧の構築系（`ShowFoundParts` / `ShowParts` / `RefreshFoundParts`）は、非表示のまま構築してから `Show` できるよう、表示状態によらず動く。

`PartsListUI` は購入などの Model 更新を行わない。購入は Scene が `SelectedPart`（Play では Canvas API の `SelectedFoundPart`）で受け取ったパーツを `Gamedata.BuyFoundPart` へ渡して行い、購入後の一覧は Canvas が `Gamedata.PartsRevision` の変化を観測して `RefreshFoundParts` で再構築する。
UI から購入していた `BuySelectedPart` は、UI から Model を書き換えない方針に反し再利用しないため削除した。Scene にスクロール量を操作させないために無効化した `Scroll` は、ホイール等で再利用の可能性があるため理由付きのコメントアウトで残している。

## uiutils

`uiutils`（`Source/UI/UIUtils`）は、UI の描画で共通に使う図形描画のヘルパーである。1×1 の白画像を1枚だけ作って使い回し、拡大と着色で矩形を描くため、描画のたびに `ebiten.NewImage` を呼ばない。

- `FillRect(dst, rect image.Rectangle, c color.Color)`: 矩形を塗りつぶす
- `DrawBorder(dst, rect image.Rectangle, width int, c color.Color)`: 矩形の内側に太さ `width` の枠を描く。矩形が `width*2` より小さいときは何もしない
- `SelectionBorderColor`: 選択中の UI（ボタン、タブ、パーツカード）の枠の色

`ButtonUI` と `PartsListUI`（背景の減光、パネル、タブ、スクロールバー、パーツカード）が利用する。

## InstalledParts

`InstalledParts`（`Source/UI/ResultUI`）は、`Gamedata.Assemble`のパーツを種別名と名前のレコードとして並べ、上から順にスライドインさせる。

- `Start() error`: 現在の`Gamedata.Assemble`から表示を初期化し、アニメーションを開始する。再呼び出しで再初期化される
- `IsFinished() bool`: 全レコードの開始と表示アニメーションが完了したか。`Start()`前は`false`
- `Update(active bool) error`: スライドインの順次開始とアニメーション更新のみ行う。`Start()`前は何もしない
- `Draw(screen)`: 表示する

`Gamedata`は読み取り専用で参照する。`NextState()`などのゲーム状態の変更やScene遷移は行わない。完了後の遷移は`ResultScene`が`IsFinished()`を確認して行う。

## 変更記録

今回の UI 拡張では既存の `Text`、シーンの描画 API、入力設定を破壊的に変更していない。既存シーンのテキスト表示は `Text` 要素を経由するよう変更済みである。
