# UIアーキテクチャ

## 基本方針

- UIはゲームのModelから分離する。
- 入力は各Scene・Componentの`Update()`で`input.Handler`から取得する。
- UIは入力を取得せず、必要なModelの状態を`Draw()`で参照して表示する。
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
- Scene・Componentから設定された表示状態を保持する。
- アニメーションなど、見た目に関する更新を行う。
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

各Sceneは「何を表示するか」とUIの構成を知る。ただし、具体的な配置・画像・フォント・描画方法は知らない。

Sceneごとに`HogeSceneCanvas`を持たせ、そこにそのSceneで使用するUI部品をまとめる。
`HogeSceneCanvas`は表示情報のコピーを保持するViewModelではなく、UI部品の構成を保持する入れ物である。

```go
type PlayScene struct {
    Gamedata *game.Gamedata
    Canvas   *PlaySceneCanvas
}

type PlaySceneCanvas struct {
    Timer     *ui.TimerUI
    Money     *ui.MoneyUI
    PartsList *ui.PartsListUI
}
```

各UIは、描画に必要な`Gamedata`などの参照を保持し、`Draw()`内で現在の状態を取得して描画する。
UIはModelを変更せず、入力処理やゲームルールも担当しない。

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
HogeSceneCanvas
      ↑
HogeSceneがUI構成を保持
```

表示用の整形処理が複雑になった場合だけ、`HogeSceneViewModel`などを追加する。
最初から表示情報のコピーやViewModelを必須にしない。

## 既存UIとの移行

既存の`Element`、`ButtonUI`、`ScrollUI`は、すぐに全て置き換えず段階的に移行する。

- `Element.Update()`からマウス入力などの処理を削除し、表示更新用にする。
- `ButtonUI`の押下判定を使う画面は、`input.Handler`による入力へ移行する。
- `ScrollUI`のスクロール量はScene・Componentから設定する。
- 入力処理をSceneとUIで二重に実装しない。
- SceneごとのUI部品を`HogeSceneCanvas`へまとめる。
- UIは必要なModel参照を受け取り、`Draw()`で状態を取得する。
- 特殊な描画やゲーム固有のアニメーションは自作UIを利用する。

## 実装計画

1. `input.Handler`を入力の唯一の入口として維持する。
2. 既存Sceneの`Update()`から、UIの表示状態を更新する処理を整理する。
3. UI要素からマウス入力などの直接入力処理を取り除く。
4. 画面ごとの`HogeSceneCanvas`を作り、UI部品の生成と配置を整理する。
5. 各UIへ必要なModel参照を渡し、`Draw()`で状態を取得する。
6. 表示変換が複雑な画面だけViewModelまたはPresenterを追加する。
7. 標準的な部品が不足した場合に限り、Ebiten UIの導入を検討する。

実装上のクラス関係は[`UI_architecture_class_diagram.puml`](./UI_architecture_class_diagram.puml)を参照する。

複雑なパーツ一覧UIの分割方針は[`PartsListUI_refactoring.md`](./PartsListUI_refactoring.md)を参照する。
