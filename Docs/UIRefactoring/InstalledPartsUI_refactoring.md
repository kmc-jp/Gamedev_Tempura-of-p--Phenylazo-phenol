# InstalledPartsUI リファクタリング案

## 目的

`InstalledPartsUI`は、組み立てたパーツの一覧表示とスライドインアニメーションを担当するUIである。
現在はUI内部から`Gamedata.NextState()`を呼び出しており、UIがゲーム状態の遷移まで担当している。

このリファクタリングでは、`InstalledPartsUI`を表示とアニメーションに限定し、ゲーム状態の変更とScene遷移を`ResultScene`へ戻す。

## 現状の問題

対象ファイルは[`InstalledParts.go`](../../Source/UI/ResultUI/InstalledParts.go)である。

### UIがModelを変更している

現在の`Update()`は、全アニメーション完了後に次の状態へ遷移させる。

```go
if i.newest == uint(len(i.active)) && i.allAnimationsFinished() {
    if err := i.gdata.NextState(); err != nil {
        return err
    }
}
```

これはUIの責務を越えている。

- UIはアニメーションの完了を検知する
- Sceneがゲーム状態を進める
- Modelは状態遷移のルールを持つ

という分担に変更する。

### SceneがUI配列の添字に依存している

[`ResultScene`](../../Source/Scene/ResultScene.go)は、次のように`elements`の順番でUIを特定している。

```go
s.elements[2].Update(true)
s.elements[1].Update(true)
```

UIの追加・削除・順番変更で別のUIを更新する危険があるため、名前付きのCanvasへ移行する。

## 目標構成

```text
ResultScene
  ├── input.Handlerから入力を取得
  ├── Gamedataの状態を更新
  ├── ResultCanvasのUIを名前で参照
  └── アニメーション完了後にModelの状態を遷移
            ↓
      InstalledPartsUI
        ├── Gamedata.Assembleを表示用に読む
        ├── 表示用UIを生成する
        ├── スライドインを制御する
        ├── 完了状態を公開する
        └── Draw()
```

`InstalledPartsUI`は、ゲームの次の状態が何であるかを知ってはならない。

## UIの公開API

Sceneには、UIの内部配列や`AnimatableUI`を公開しない。

```go
type InstalledParts interface {
    scene.Element

    Start() error
    IsFinished() bool
}
```

### APIの責務

- `Start`: 現在の`Gamedata.Assemble`をもとに表示とアニメーションを初期化する
- `IsFinished`: 全要素の表示アニメーションが完了したかを返す
- `Update`: 表示アニメーションだけを進める
- `Draw`: 表示する

`NextState()`、Scene遷移、ベンチマーク開始などのゲーム処理はAPIに含めない。

## Model参照の扱い

`InstalledPartsUI`が`Gamedata`を参照する構成は当面維持してよい。
ただし、UIが書き換えてよいModelの範囲はゼロとする。

```go
type InstalledParts struct {
    gdata *game.Gamedata // 読み取り専用として扱う
}
```

より厳密に分離する必要が出た場合は、`Assemble`とフォントなどUIに必要な値だけを渡す専用の表示データへ変更する。
最初からコピー用ViewModelを必須にせず、まずは状態遷移の責務をSceneへ戻す。

## ResultScene側の責務

`ResultScene`は、アニメーション完了をゲーム状態の遷移条件として扱う。

```go
if s.Canvas.InstalledParts.IsFinished() {
    if err := s.Gamedata.NextState(); err != nil {
        return nil, transition.None(), fmt.Errorf(
            "InstalledPartsUI完了後の状態遷移に失敗しました: %w",
            err,
        )
    }
}
```

状態遷移の具体的な条件が複雑になった場合でも、判定はSceneまたはModel側に置く。
UIは「表示が終わった」という事実だけを返す。

## ResultCanvas

`ResultScene`のUIは配列の添字ではなく、名前付きのCanvasで管理する。

```go
type ResultCanvas struct {
    Header          *ui.Text
    BenchmarkPrompt *ui.AnimatableUI
    InstalledParts  *resultui.InstalledParts
    PartsList       *ui.PartsListUI
}
```

Sceneからの参照例:

```go
s.Canvas.InstalledParts.Update(active)
s.Canvas.BenchmarkPrompt.Update(active)
```

CanvasはUI部品の構成を保持するだけで、ゲームルールや状態遷移を持たない。

## 内部の分割

まずは既存の`InstalledParts.go`内で責務を整理する。
サイズが大きくなった場合は、次のように分割する。

```text
Source/UI/ResultUI/
├── InstalledParts.go       公開API、全体の更新と描画
├── InstalledPartRecord.go  1パーツ分の表示とアニメーション
├── InstalledPartsState.go  表示開始・完了状態
└── PartCategory.go         パーツ種別名の変換
```

### 分割の基準

- `InstalledPartRecord`: category/nameのAnimatableUIと、その完了判定
- `InstalledPartsState`: 初期化済み、表示開始、完了状態
- `PartCategory`: `game.Part.Kind`から表示名への変換
- `InstalledParts`: レコード生成、順次開始、全体描画

ゲーム状態の変更とScene遷移は、これらのファイルへ移さない。

## 段階的な移行手順

1. `InstalledPartsUI`に`IsFinished()`を追加する。
2. `InstalledPartsUI.Update()`から`Gamedata.NextState()`を削除する。
3. `ResultScene.Update()`で`IsFinished()`を確認し、`Gamedata.NextState()`を呼ぶ。
4. `ResultScene`の`elements[index]`参照を名前付きCanvasへ置き換える。
5. `InstalledPartsUI`の初期化タイミングを`Start()`へ整理する。
6. 必要に応じて内部状態・レコード・カテゴリ変換を別ファイルへ分割する。

移行中は既存のコンストラクタを維持してもよいが、新規コードではCanvasと意味ベースAPIを使用する。

## 完了条件

- `InstalledPartsUI`が`Gamedata.NextState()`を呼び出さない。
- `InstalledPartsUI`がシーン遷移やゲーム状態を決定しない。
- `ResultScene`が`elements[1]`や`elements[2]`でUIを参照しない。
- `InstalledPartsUI`が`scene.Element`を実装する。
- アニメーション完了の検知とゲーム状態の遷移が分離されている。
- UIの追加・削除・順番変更で、別のUIが誤って更新されない。
