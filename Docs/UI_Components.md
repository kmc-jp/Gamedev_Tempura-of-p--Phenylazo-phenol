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

- `TitleCanvas` / `PauseCanvas` / `SelectStageCanvas`: 案内用の `Text` 1個だけを持つ。`NewXxxCanvas(font)` で生成し、`NewXxxScene` の引数で注入する
- `PlayCanvas`: 残り時間（`GameDataText`）、空き重量・所持金（`GameDataUI`）、`PartsListUI`。`NewPlayCanvas(scene, gdata)`。`PlayScene.SelectListRevision` の変化で一覧を開く。購入はSceneが`CurrentPartsTab`で発見タブの表示中かを確かめたうえで、`SelectedFoundPart`と`Gamedata.BuyFoundPart`で行い、購入後の一覧は`Gamedata.PartsRevision`の変化を観測して再構築する
- `TrainCanvas`: 案内 `Text`、`PartsListUI`。`NewTrainCanvas(scene, gdata)`。`TrainScene.PartsListRevision`（または`Gamedata.PartsRevision`）の変化で所持パーツ一覧を構築・表示する
- `ResultCanvas`: 案内 `Text`、ベンチマーク案内（`AnimatableUI`）、`InstalledParts`、`PartsListUI`。`NewResultCanvas(scene, gdata)`。全UIを毎フレーム`Canvas.Update`で更新する。組み立て済みパーツの表示は`Result_Assembled`、ベンチマーク案内のアニメーションは`Result_Benchmark`に変わったときに開始する

Play / Train / Result の Canvas は Scene のポインタや Gamedata を観測して表示を切り替える。そのため、Canvas 本体ではなく Canvas の生成関数（`func(*XxxScene) (XxxCanvasAPI, error)`）を `NewXxxScene` に渡し、Scene のコンストラクタが自身の生成直後に呼んで注入する。生成関数がエラーや `nil` を返したときは Scene の生成エラーになる。

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
