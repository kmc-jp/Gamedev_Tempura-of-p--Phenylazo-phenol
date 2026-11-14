# UIアーキテクチャ

## 基本方針

- UIはゲームのModelから分離する。
- 入力は各Scene・Componentの`Update()`で`input.Handler`から取得する。
- UIは入力を取得せず、必要なModelの状態を`Draw()`で参照して表示する。
- `Update()`は時間経過による状態変化やアニメーションを扱い、`Draw()`は現在状態の描画を行う。
- UIは`Update()`や`Draw()`から`*game.Gamedata`を直接参照してよい。ただし、UIからModelを書き換えない。
- UIは必要最小限の構成とし、設計パターンのための不要な層を追加しない。

```text
InputConfig
    ↓ 準備
input.Handler
    ↓ 入力取得
Scene / Component.Update()
    ├── Model更新
    ├── UI構成・状態の管理
    └── シーン遷移
          ↓
        UI.Draw()
```

## アーキテクチャ

MVUを厳密には採用しない。MVUの`Update`に相当する処理は、既存の各Scene・Componentの`Update()`が担当する。

```text
Model
  ↑                 ↓
Scene / Component.Update()
  ↓
受動的なView（UI）
```

MVC・MVVM・MVUの名前を厳密に当てはめるよりも、既存のEbitengineのゲームループに合わせることを優先する。

### Scene・Component.Update()の責務

- `input.Handler`からキーボード等の入力を取得する。
- Modelを更新する。
- 画面遷移やゲーム状態の変更を行う。
- UIの表示状態を更新する。
- 必要に応じてComponentの`Update()`を呼び出す。

### UIの責務

- テキスト、画像、パネルなどを描画する。
- `*game.Gamedata`など、描画に必要なModelの状態を参照する。
- アニメーションなど、時間経過による見た目の状態を更新する。
- レイアウトと描画を担当する。

UIは次の処理を担当しない。

- `input.Handler`へのアクセス
- キー、マウス、ゲームパッド等の入力判定
- Modelの更新
- シーン遷移の判断
- ゲームルールの実装

## 入力方針

入力の入口は、既存の`inputconfig.InputConfig`と`input.Handler`に統一する。
将来的に`input.Handler`を`game.Gamedata`へ移動する場合も、UIへ入力処理を追加しない。

```go
if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
    s.someUI.SetSelected(true)
}
```

UIの`Update(active bool)`は、入力処理のためではなく、アニメーションや表示更新のために使用する。
UIのUpdateが不要な要素は、Scene側から呼び出さず描画だけを行ってもよい。

`active == false`のときのUIの挙動は、UIごとに次のいずれも選択できるものとする。

- アニメーションや内部状態の更新を停止する。
- 非アクティブでも表示更新を継続する。

ただし、入力処理は`active`の値にかかわらずUIでは行わない。

## Ebiten UIの利用方針

[Ebiten UI](https://github.com/ebitenui/ebitenui)は、必要になった標準的なレイアウト・ウィジェットの実装候補とする。
[公式ドキュメント](https://ebitenui.github.io/)

ただし、Ebiten UIのボタン等は通常、マウス・フォーカス・クリックなどの入力処理を内部で行う。
今回の入力方針では、Ebiten UIをインタラクティブな入力担当としてそのまま利用しない。

利用する場合は、次のいずれかに限定する。

- レイアウトや描画部品として利用する。
- 入力処理を無効化または隠蔽できる部品だけを利用する。
- その画面だけ、Sceneの`input.Handler`からEbiten UIへ状態を渡すアダプターを用意する。

時間を優先する間は、既存の自作UIを継続し、Ebiten UIの導入は標準部品が必要になった時点で判断する。

## SceneとUIの接続

各Sceneは「何を表示するか」を知り、Canvasの公開された契約を利用する。
具体的なUI部品の構成・配置・画像・フォント・描画方法は、Sceneから分離する。

SceneごとにScene固有Canvasを持たせ、そこにそのSceneで使用するUI部品をまとめる。
Scene固有Canvasは表示情報のコピーを保持するViewModelではなく、UIの更新・描画・表示状態を管理する構成層である。

Sceneが利用するCanvasの意味ベースAPIは、`scene` packageにinterfaceとして定義する。
Canvas API interfaceは、Sceneと密接に関連するため、対応するScene構造体と同じファイルに定義する。
具体的なCanvas実装は`scenecanvas` packageに置き、`construct` packageからSceneへ注入する。
これにより、Sceneは具体的なCanvas実装や個別UIに依存せず、Canvasの契約だけに依存する。

```text
scene.PlayScene
    ↓ scene.PlayCanvasAPI
scenecanvas.PlayCanvas
    ↓
canvas.Canvas
    ↓
ui
```

### Canvas APIの構成

Canvas APIは`Update`、`Draw`と識別用メソッド（`IsTitleCanvas()`など）を基本とする。
識別用メソッドは処理を持たず、メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもので、`scenecanvas`から実装できるよう公開メソッドとする。
`Update`・`Draw`だけで完結しないSceneは、次の規則でAPIを追加する。

- 入力起点の操作: Sceneが`input.Handler`で入力を判定し、専用の入力用APIを呼ぶ。名前はUI部品の操作名ではなく、そのSceneでの意図に合わせる（例: Playは`MovePartSelectionDown`、Trainは`ScrollPartsListDown`）。入力用APIは、対象UIが非表示なら何もしない。
- シーン固有の情報: CanvasはSceneのポインタを持ち、Sceneの公開メソッドで参照する（例: `ResultScene.CandidateParts`）。ゲームルール（候補の絞り込みなど）はSceneに残す。
- 入力以外の状態変化: Canvasが`Update`でScene・Gamedataを毎フレーム観測し、前回値との差分で反映する（例: `Play_Select`への変化で一覧を開く）。SceneからCanvasへの通知APIは作らない。
- 問い合わせ: UIの状態を戻り値で返すAPIは、個別に合意したものだけにする（現在は`IsPartsListOpen`、`SelectedPart`、`IsInstalledPartsFinished`）。

```go
// PlayScene.go
type PlayCanvasAPI interface {
    Update(active bool) error
    Draw(screen *ebiten.Image)
    IsPlayCanvas()
    // 入力用
    TogglePartsList() error
    MovePartSelectionDown()
    MovePartSelectionUp()
    ShowNextPartsTab() error
    ShowPreviousPartsTab() error
    BuySelectedPart() (game.PurchaseResult, error) // 暫定: UIが購入処理を呼ぶ
    // 問い合わせ
    IsPartsListOpen() bool
}
```

注意点:

- CanvasがSceneのポインタを持つSceneは、Scene生成後に`SetCanvas`で注入する。
- Sceneの`Update`は本体の処理の後で`canvas.Update`を呼び、同じフレームの状態変化をCanvasが観測できるようにする。
- Canvasは`Gamedata.Timer.Expired()`を呼ばない。通知を消費するため、Sceneが満了を検知できなくなる。タイマー判定はSceneが行い、結果を公開メソッドで渡す。
- 購入処理は暫定でUI（`PartsListUI.BuySelectedPart`）がModelの購入処理を呼んでいる。後でSceneが購入する形へ移す（[`UI_refactoring_priority.md`](./UI_refactoring_priority.md)の保留事項）。

Canvasは個別UIを管理し、固定された順序で子UIの`Update()`と`Draw()`を呼び出す。
描画順はSceneごとのCanvasが決定する。Canvas基盤はSceneを知らない汎用コンテナである。

各UIは、描画に必要な`*game.Gamedata`やSceneの表示状態などを参照してよい。
ただし、UIは入力処理、Model更新、Scene遷移、ビジネスルールを担当しない。

```go
type MoneyUI struct {
    Gamedata *game.Gamedata
}

func (u *MoneyUI) Draw(screen *ebiten.Image) {
    // u.Gamedata.Moneyを読み取り、具体的な見た目で描画する
}
```

処理の流れは以下とする。

```text
Model / Gamedata
      ↑
      │ 必要な情報を参照
各UIのDraw()
      ↑
Scene固有Canvas
      ↑
SceneがCanvas APIを利用
```

表示用の整形処理が複雑になった場合だけ、`HogeSceneViewModel`などを追加する。
最初から表示情報のコピーやViewModelを必須にしない。

表示データの更新方式は、単純なUIと複雑なUIで統一しない。
単純なUIは`*game.Gamedata`などのModelを直接参照してよい。
複雑なUIは、Scene固有Canvasの意味ベースAPI、専用の状態型、またはViewModelを必要に応じて利用する。
方式の選択は、表示変換の複雑さとUIが保持する状態の量に応じて決める。

## 状態の永続化方針

UIの生成・再生成・再利用の方法は、UIごとの実装に委ねる。
基本方針では、UIを毎フレーム再生成するか、Scene生成時に生成するかを規定しない。
ただし、UIを破棄して再生成しても失われてよい状態と、UIの外部で保持すべき状態は区別する。

### UIが保持する状態

UIは、表示を成立させるための一時的な状態を保持してよい。

- アニメーションの進行状況
- スクロール位置
- 補間中の座標や透明度
- 描画キャッシュ
- レイアウト計算結果
- 一時的な点滅やホバー表現

これらはUIを再生成したときに初期化されても、ゲームの進行やユーザーの確定結果を変更しない状態とする。

### SceneまたはViewModelが保持する状態

画面内の操作に意味があるが、ゲーム全体のModelやセーブデータには属さない状態は、Sceneまたは必要に応じてViewModelが保持する。

- 選択中の項目
- 現在のタブ
- フォーカス対象
- メニューの開閉状態
- 一時的なUIの表示状態
- 画面内での入力や操作の進行状況

UIはこれらを入力によって直接変更せず、SceneまたはViewModelから渡された状態を表示する。

### Modelが保持する状態

画面を離れた後もゲーム上の意味を持つ状態は、Modelが保持する。

- 確定したパーツ選択や組み立て結果
- プレイヤーの所持品
- ゲームの進行状態
- セーブ対象の設定

選択中の項目のような一時的な状態と、選択を確定した結果のようなゲーム上の状態を混同しない。

### 判断基準

次の問いを状態の所属先を決める基準とする。

> UIを破棄して再生成したとき、その状態が失われるとゲーム上の意味が変わるか。

- 変わらない表示上の状態はUIが保持する。
- 画面内の操作状態はSceneまたはViewModelが保持する。
- ゲーム進行や確定結果に関わる状態はModelが保持する。

## 既存UIとの移行

既存の`Element`、`ButtonUI`、`ScrollUI`は、Canvas導入に合わせて移行した。

- Scene packageの`Element`、`UIDrawer`、`UICamera`による配列ベースのUI管理はCanvasへ置き換え、撤去した。
- Canvas内部で共通処理が必要な場合だけ、Canvas package内の抽象を使用する。
- `ButtonUI`は入力処理を持たず、選択状態の表示のみ担当する（移行済み）。押下判定が必要な画面は、`input.Handler`による入力で行う。
- `ScrollUI`は入力を扱わない表示用コンテナとし、スクロール量は`ScrollUI`を持つUI（`PartsListUI`など）が選択追従などに応じて設定する。Sceneはスクロール量を直接操作しない。
- 入力処理をSceneとUIで二重に実装しない。
- SceneごとのUI部品をScene固有Canvasへまとめる。
- UIは必要なModel参照を受け取り、`Draw()`で状態を取得する。
- 特殊な描画やゲーム固有のアニメーションは自作UIを利用する。

## 実装計画

1. `input.Handler`を入力の唯一の入口として維持する。
2. 既存Sceneの`Update()`から、UIの表示状態を更新する処理を整理する。
3. UI要素からマウス入力などの直接入力処理を取り除く。
4. `scene` packageにSceneごとのCanvas API interfaceを、各Scene構造体と同じファイルに定義する。
5. `scenecanvas` packageに具体的なCanvasを実装し、`construct`から注入する。
6. `canvas` packageで個別UIの更新・描画・表示状態を管理する。
7. 表示変換が複雑な画面だけViewModelまたはPresenterを追加する。
8. 標準的な部品が不足した場合に限り、Ebiten UIの導入を検討する。

実装上のクラス関係は[`UI_architecture_class_diagram.puml`](./UI_architecture_class_diagram.puml)を参照する。

複雑なパーツ一覧UIの分割方針は[`PartsListUI_refactoring.md`](./PartsListUI_refactoring.md)を参照する。

InstalledPartsUIの責務分離は[`InstalledPartsUI_refactoring.md`](./InstalledPartsUI_refactoring.md)を参照する。

全体の実施順は[`UI_refactoring_priority.md`](./UI_refactoring_priority.md)を参照する。
