# InstalledPartsUI リファクタリング案

## 目的

`InstalledPartsUI`は、組み立てたパーツの一覧表示とスライドインアニメーションを担当するUIである。
現在はUI内部から`Gamedata.NextState()`を呼び出しており、UIがゲーム状態の遷移まで担当している。

このリファクタリングでは、`InstalledPartsUI`を表示とアニメーションに限定し、ゲーム状態の変更とScene遷移を`ResultScene`へ戻す。

## 実施状況

移行手順1〜5とCanvasへの移行は実施済みである。手順6（ファイル分割）は未実施である。

対象ファイルは[`InstalledParts.go`](../../Source/UI/ResultUI/InstalledParts.go)と[`ResultScene.go`](../../Source/Scene/ResultScene.go)である。

### 解消した問題

- UIが`Update()`内で`Gamedata.NextState()`を呼んでいた。現在は`IsFinished()`で完了を返すだけで、`NextState()`は`ResultScene`が呼ぶ。
- UIの初期化が初回`Update()`の遅延初期化だった。現在は`ResultScene`が`Start()`を呼ぶ。
- `ResultScene`が`elements[1]`、`elements[2]`の添字でUIを特定していた。現在は`scene` package内の名前付きフィールド（`installedParts`、`benchmarkPrompt`）で参照する。

### Canvas化

`InstalledParts`は`scenecanvas.ResultCanvas`が保持する。
`ResultScene`は具体型を知らず、`ResultCanvasAPI`の問い合わせAPI`IsInstalledPartsFinished()`で完了を確認する。
`Start()`は、`Gamedata.State`が`Result_Assembled`へ変わったことをCanvasが観測して呼ぶ。
`Update()`は`Result_Assembled`中だけCanvasが呼ぶ。

## 目標構成

```text
ResultScene
  ├── input.Handlerから入力を取得
  ├── Gamedataの状態を更新
  ├── ResultCanvasAPIでUIを利用
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
`InstalledParts`の公開APIは次のとおりで、これを使うのは`ResultCanvas`だけである。
Sceneは`ResultCanvasAPI.IsInstalledPartsFinished()`だけを利用する。

```go
func (ips *InstalledParts) Start() error
func (ips *InstalledParts) IsFinished() bool
func (ips *InstalledParts) Update(active bool) error
func (ips *InstalledParts) Draw(screen *ebiten.Image)
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
if s.canvas.IsInstalledPartsFinished() {
    if err := s.gamedata.NextState(); err != nil {
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

Canvas導入により、名前付きフィールド（`installedParts`、`benchmarkPrompt`、`partsList`）は`ResultScene`から撤去した。
UI部品は`scenecanvas.ResultCanvas`が保持し、描画順もCanvasが決める。
CanvasはUI部品の構成と、状態の観測による表示の切り替えだけを持ち、ゲームルールや状態遷移を持たない。

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

1. （実施済み）`InstalledPartsUI`に`IsFinished()`を追加する。
2. （実施済み）`InstalledPartsUI.Update()`から`Gamedata.NextState()`を削除する。
3. （実施済み）`ResultScene.Update()`で`IsFinished()`を確認し、`Gamedata.NextState()`を呼ぶ。
4. （実施済み）`ResultScene`の`elements[index]`参照を名前付きフィールドへ置き換え、その後`ResultCanvas`へ移行した。
5. （実施済み）`InstalledPartsUI`の初期化タイミングを`Start()`へ整理する。
6. （未実施）必要に応じて内部状態・レコード・カテゴリ変換を別ファイルへ分割する。

移行中は既存のコンストラクタを維持してもよいが、新規コードではCanvasと意味ベースAPIを使用する。

## 完了条件

- 達成: `InstalledPartsUI`が`Gamedata.NextState()`を呼び出さない。
- 達成: `InstalledPartsUI`がシーン遷移やゲーム状態を決定しない。
- 達成: `ResultScene`が`elements[1]`や`elements[2]`でUIを参照しない。
- 達成: `InstalledPartsUI`が`scene`の共通UI要素インターフェースに依存しない（`canvas.CanvasItem`として`ResultCanvas`が管理する）。
- アニメーション完了の検知とゲーム状態の遷移が分離されている。
- UIの追加・削除・順番変更で、別のUIが誤って更新されない。
