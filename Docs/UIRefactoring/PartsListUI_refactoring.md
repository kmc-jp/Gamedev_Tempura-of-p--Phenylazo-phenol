# PartsListUI リファクタリング案

## 目的

`PartsListUI`は、一覧表示・タブ・選択・スクロール・パーツカード生成をまとめて担当する複雑なUIである。
Scene側がこのUIの内部操作を直接扱うと、SceneがUIの実装詳細に依存し、状態不整合や変更時の影響範囲が大きくなる。

このリファクタリングでは、`PartsListUI`の実装を`Source/UI`に閉じ込め、Scene側には画面上の意味に対応した最小限のAPIだけを公開する。

## 現状の問題

移行前のScene側の`partsListUI`インターフェースには、UIの内部操作が漏れていた（移行済み）。

```go
type partsListUI interface {
    Element
    SetParts(game.Parts) error
    SetLists(game.Parts, game.Parts) error
    SetVisible(bool)
    Visible() bool
    MoveSelection(int)
    MoveTab(int) error
    IsFoundTab() bool
    SelectedPart() *game.Part
    SelectedIndex() int
}
```

Sceneが以下を直接知っている点が問題である。

- 選択位置を整数で管理している
- タブを数値差分で移動している
- UI内部の`cards`に対応する選択インデックスを参照している
- スクロールや表示状態をUI内部の都合に合わせて操作している

Sceneは「パーツを次に選択する」「発見したパーツを表示する」のような意味を知るべきであり、カードやインデックスの管理方法を知るべきではない。

## 目標構成

```text
Scene
  ├── input.Handlerから入力を取得
  ├── Model / Gamedataを更新
  └── PartsListUIへ意味のある命令を送る
            ↓
      PartsListUI
        ├── 表示状態を保持
        ├── タブを管理
        ├── 選択を管理
        ├── スクロール位置を管理
        ├── PartCardを生成・更新
        └── Draw()
```

SceneはUIの配置と利用目的を知るが、UIの描画方法と内部状態は知らない。

## 表示状態

`PartsListUI`専用の表示状態をUI package内に定義する。
これはModelの代替ではなく、一覧を描画するためのUI内部状態である。

```go
type PartsListViewState struct {
    FoundParts  game.Parts
    OwnedParts  game.Parts
    ActiveTab   PartsListTab
    Selected    int
    Visible     bool
}
```

実装上は、既存の`PartsListUI`がこの状態を直接保持してもよい。
別structに分けるかどうかは、状態更新のテストが必要になった時点で決める。

## Sceneに公開するAPI

Scene側には、UIの内部構造ではなく意味ベースのAPIを公開する。

```go
type PartsList interface {
    ShowFoundParts(found, owned game.Parts) error
    ShowParts(parts game.Parts) error
    Hide()
    Show()
    Visible() bool

    SelectNext()
    SelectPrevious()
    NextTab() error
    PreviousTab() error

    SelectedPart() *game.Part
    IsShowingFoundParts() bool
    RefreshFoundParts(found, owned game.Parts) error
}
```

### APIの責務

- `ShowFoundParts`: 発見済み・所持パーツのタブ一覧を表示する
- `ShowParts`: 単一のパーツ一覧を表示する
- `Show` / `Hide`: 表示状態を変更する
- `SelectNext` / `SelectPrevious`: 選択対象を移動する
- `NextTab` / `PreviousTab`: タブを切り替える
- `SelectedPart`: 現在選択されているModelのパーツを取得する
- `IsShowingFoundParts`: 現在のタブの意味を確認する
- `RefreshFoundParts`: Found / Parts の変化後に、現在のタブと選択位置を維持して一覧を再構築する

`ShowFoundParts` / `ShowParts`は一覧を構築するだけで、表示状態は変更しない。表示するには続けて`Show`を呼ぶ。

`IsShowingFoundParts`は、`ShowFoundParts`によるタブ付き一覧で「発見したパーツ」タブを表示しているときだけ`true`を返す。
`ShowParts`による単一一覧の表示中は`false`を返す。

`PartsListUI`は購入処理（Modelの更新）を行わない。`BuySelectedPart`はコメントアウトした。
購入は、Sceneが`SelectedPart`（Canvas経由では`SelectedFoundPart`）で受け取ったパーツを`Gamedata.BuyFoundPart`へ渡して行う。
`BuyFoundPart`はポインタで`Gamedata.Found`内の位置を解決するため、UIのカード順とModelの順の一致に依存しない。
購入後の一覧は、Canvasが`Gamedata.Found` / `Parts`の変化を観測して`RefreshFoundParts`で再構築する。
`NewPartsListUI`は`*game.Gamedata`を引数に取らない。

`SelectedIndex`、`MoveSelection(delta)`、`MoveTab(delta)`は削除済みである。
購入処理でUIの選択インデックスをSceneへ漏らさないため、Sceneはパーツそのもの（`SelectedPart`）を受け取る。

## 入力処理

入力は引き続きSceneが`input.Handler`から取得する。
`PartsListUI`はキーボード、マウス、ゲームパッドを直接参照しない。

```go
if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
    s.partsList.SelectPrevious()
}
if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
    s.partsList.SelectNext()
}
if s.Input.ActionIsJustPressed(inputconfig.ActionMoveLeft) {
    if err := s.partsList.PreviousTab(); err != nil {
        return err
    }
}
if s.Input.ActionIsJustPressed(inputconfig.ActionMoveRight) {
    if err := s.partsList.NextTab(); err != nil {
        return err
    }
}
```

入力の解釈はScene、選択範囲のクランプやスクロール追従は`PartsListUI`が担当する。

## データの扱い

パーツ一覧の元データは`Gamedata`をModelとして扱う。
一覧を開くタイミングや、どの一覧を表示するかはSceneが決める。

```go
if err := s.partsList.ShowFoundParts(
    s.Gamedata.Found,
    s.Gamedata.Parts,
); err != nil {
    return nil, transition.None(), err
}
s.partsList.Show()
```

`PartsListUI`は受け取ったパーツを表示用に利用するが、パーツのゲームルール上の状態を変更しない。
選択されたパーツの購入処理は、Sceneが行う。
Sceneは選択中のパーツを問い合わせ、Modelの購入処理（`Gamedata.BuyFoundPart`）へ渡す。
購入結果の扱いとエラー表示もSceneが担当する。

```go
part := s.canvas.SelectedFoundPart()
switch s.Gamedata.BuyFoundPart(part) {
case game.Buy:
    // 購入成功時の処理
case game.NoItem, game.Lack_money, game.Lack_weight:
    // 購入できない場合の通知
}
```

## `PartsListUI`内部の分割

最初から多数のファイルへ分割せず、現在の`PartsListUI.go`を保ったまま責務を明確にする。
必要になった場合は、次の順で分割する。

```text
Source/UI/
├── PartsListUI.go       公開API、表示状態、全体の更新と描画
├── PartsListLayout.go   サイズ・位置・タブ・スクロールバーの計算
├── PartCard.go          1パーツ分の表示
└── PartsListState.go    表示状態とタブ状態
```

### 分割の基準

- `PartCard`: パーツ1件の表示と選択表示だけを担当する
- `PartsListState`: 選択中の要素、タブ、表示状態だけを保持する
- `PartsListLayout`: 座標、サイズ、スクロール範囲を計算する
- `PartsListUI`: 上記を組み合わせ、公開APIを実装する

ゲームルールやシーン遷移は、これらのファイルへ移さない。

## Scene側の変更（完了）

Sceneは`PartsListUI`を直接持たない。`PartsListUI`は各Scene固有Canvas（`scenecanvas.PlayCanvas`、`TrainCanvas`、`ResultCanvas`）が保持する。
Sceneは各Canvas APIの、意図に合わせた名前のAPIだけを利用する。

| Scene | Canvas API | 内部で呼ぶ`PartsListUI`のAPI |
|---|---|---|
| Play | `TogglePartsList`、`MovePartSelectionDown/Up`、`ShowNextPartsTab`、`ShowPreviousPartsTab`、`IsPartsListOpen`、`SelectedFoundPart` | `ShowFoundParts`、`RefreshFoundParts`、`Show`、`Hide`、`SelectNext/Previous`、`NextTab/PreviousTab`、`IsShowingFoundParts`、`SelectedPart`、`Visible` |
| Train | `ScrollPartsListDown/Up` | `ShowParts`、`Show`、`SelectNext/Previous` |
| Result | `MovePartSelectionDown/Up`、`SelectedPart` | `ShowParts`、`Show`、`Hide`、`SelectNext/Previous`、`SelectedPart` |

一覧の構築・表示は、CanvasがScene・Gamedataの状態変化を観測して行う（Playの`Play_Select`への変化、Trainの`IsPartsListRevealed`、Resultの`CandidatesRevision`）。

購入はSceneが`Gamedata.BuyFoundPart`で行い、購入後の一覧の再構築はCanvasが`Gamedata.Found`・`Parts`の変化を観測して行う。

## 段階的な移行手順

1. （完了）`PartsListUI`に意味ベースのメソッドを追加する。
2. （完了）`PlayScene`、`TrainScene`、`ResultScene`の呼び出しをCanvas API経由へ置き換える。
3. （完了）Scene側の`SelectedIndex`、`MoveSelection`、`MoveTab`への依存を削除する。
4. （完了）購入処理をSceneへ移す（`SelectedFoundPart`と`Gamedata.BuyFoundPart`）。
5. （完了）Scene側の`partsListUI`インターフェースを削除し、各Canvas APIへ置き換える。
6. （完了）互換API（`SetParts`、`SetLists`、`SetVisible`、`MoveTab`、`MoveSelection`、`SelectedIndex`、`CurrentTab`、`IsFoundTab`）を削除する。
7. （未実施）必要に応じて`PartCard`、状態、レイアウト計算を別ファイルへ分割する。
8. （完了）Scene固有CanvasへUI部品を移す。
9. （完了）`scene.Element`、Sceneの`elements []Element`、`UIDrawer`、`UICamera`などの配列ベースUI管理を削除する。
10. （完了）ソース全体を検索し、廃止対象の型・フィールド・メソッドが残っていないことを確認する。

## 完了条件

- Scene packageが`PartCard`、`cards`、選択インデックス、スクロール量を参照しない。
- Scene packageが`MoveSelection(delta)`や`MoveTab(delta)`を呼び出さない。
- Scene packageが`SelectedIndex()`を参照しない。
- Scene packageが購入処理のために選択インデックスをModelへ渡さない。
- Canvas移行完了後、`scene.Element`、Sceneの`elements []Element`、`UIDrawer`、`UICamera`が
  コメントアウトまたは削除され、ソース全体の検索で利用箇所がない。
- `PartsListUI`が入力パッケージや`ebiten/inpututil`を直接参照しない。
- パーツ選択・タブ切り替え・スクロール追従がUI内部で完結する。
- `go test ./...`またはプロジェクトで利用している同等の検証を通過する。
