# PartsListUI リファクタリング案

## 目的

`PartsListUI`は、一覧表示・タブ・選択・スクロール・パーツカード生成をまとめて担当する複雑なUIである。
Scene側がこのUIの内部操作を直接扱うと、SceneがUIの実装詳細に依存し、状態不整合や変更時の影響範囲が大きくなる。

このリファクタリングでは、`PartsListUI`の実装を`Source/UI`に閉じ込め、Scene側には画面上の意味に対応した最小限のAPIだけを公開する。

## 現状の問題

現在のScene側の`partsListUI`インターフェースには、UIの内部操作が漏れている。

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
    scene.Element

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

`SelectedIndex`、`MoveSelection(delta)`、`MoveTab(delta)`は、移行完了後にSceneから削除する。

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
選択されたパーツをModelへ追加する処理は、SceneまたはModel側が担当する。

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
- `PartsListUI`: 上記を組み合わせ、公開APIと`Element`を実装する

ゲームルールやシーン遷移は、これらのファイルへ移さない。

## Scene側の変更

Scene package内の`partsListUI`インターフェースは、意味ベースの`PartsList`インターフェースへ置き換える。
同じインターフェースを`PlayScene`、`TrainScene`、`ResultScene`で共有する。

変更対象の例:

- [`PlayScene`](../Source/Scene/PlayScene.go)
- [`TrainScene`](../Source/Scene/TrainScene.go)
- [`ResultScene`](../Source/Scene/ResultScene.go)
- [`PartsListUI`](../Source/UI/PartsListUI.go)

Sceneの`elements []Element`自体は当面維持してよい。
`HogeSceneCanvas`を導入する場合は、Canvasが`PartsListUI`を保持し、SceneはCanvas経由で参照する。

```go
type PlaySceneCanvas struct {
    PartsList *ui.PartsListUI
}
```

## 段階的な移行手順

1. `PartsListUI`に意味ベースのメソッドを追加する。
2. `PlayScene`、`TrainScene`、`ResultScene`の呼び出しを新APIへ置き換える。
3. Scene側の`SelectedIndex`、`MoveSelection`、`MoveTab`への依存を削除する。
4. `partsListUI`インターフェースを`PartsList`へ変更する。
5. 既存の低レベルAPIを削除する。
6. 必要に応じて`PartCard`、状態、レイアウト計算を別ファイルへ分割する。
7. `HogeSceneCanvas`へUI部品を移す。

移行中は互換メソッドを残してもよいが、新規コードでは意味ベースのAPIだけを使用する。

## 完了条件

- Scene packageが`PartCard`、`cards`、選択インデックス、スクロール量を参照しない。
- Scene packageが`MoveSelection(delta)`や`MoveTab(delta)`を呼び出さない。
- `PartsListUI`が`scene.Element`を実装する。
- `PartsListUI`が入力パッケージや`ebiten/inpututil`を直接参照しない。
- パーツ選択・タブ切り替え・スクロール追従がUI内部で完結する。
- `go test ./...`またはプロジェクトで利用している同等の検証を通過する。
