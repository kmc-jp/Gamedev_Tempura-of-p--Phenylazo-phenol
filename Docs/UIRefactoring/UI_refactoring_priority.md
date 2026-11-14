# UIリファクタリング優先順位

## 判断基準

優先順位は、次の観点で決める。

1. UIからModel更新やScene遷移が行われているか
2. SceneがUIの内部実装に依存しているか
3. UIの順番や添字に依存しているか
4. 複数Sceneから利用され、変更の影響範囲が大きいか
5. 現在の開発を止めずに段階移行できるか
6. 別担当者が実装するScene側Canvas API interfaceに依存しないか

## 担当の分離

Scene側のCanvas API interface（`scene.XxxCanvasAPI`）の定義と、Sceneへの組み込みは別担当者が実装する。
そのため、実装項目を次の2フェーズに分ける。

- フェーズA: Scene側Canvas API interfaceに依存しない作業。API確定を待たずに先行して実施する。
- フェーズB: Scene側Canvas API interfaceに依存する作業。API確定後に実施する。

フェーズAでは`scene` package（Scene本体と`scene`内のUI用interface）を変更しない。
フェーズAの完了時点で、フェーズBでは「interface適合宣言・`construct`からの注入・Sceneの呼び出し置換」だけを反映すればよい状態にする。

## フェーズA: Scene側Canvas API interfaceに依存しない作業

### 1. InstalledPartsUIの責務分離（完了）

`IsFinished()`と`Start()`を導入し、状態遷移を`ResultScene`へ移した。
添字参照は`scene` package内の名前付きフィールドへ置き換え済み。8で`ResultCanvas`へ移行し、名前付きフィールドも撤去した。

対象:

- [`InstalledParts.go`](../../Source/UI/ResultUI/InstalledParts.go)
- [`ResultScene.go`](../../Source/Scene/ResultScene.go)

理由:

- UIが`Gamedata.NextState()`を呼び出している
- UIがゲーム状態の遷移を実質的に決めている
- ResultSceneが`elements[1]`、`elements[2]`の添字に依存している

最初に`IsFinished()`だけを導入し、状態遷移をSceneへ戻す。

### 2. PartsListUIの意味ベースAPI追加（UI側のみ）（完了）

意味ベースAPIを`PartsListUI`へ追加し、購入結果型を`game.PurchaseResult`として公開した。
旧APIは一時的に`Deprecated`の互換APIとして残し、8でSceneの呼び出しを置換して削除した。

対象:

- [`PartsListUI.go`](../../Source/UI/PartsListUI.go)
- [`PartsManager.go`](../../Source/Game/PartsManager.go)
- `construct`の`PartsListUI`生成箇所（引数追加のみ）

理由:

- Sceneが選択、タブ、表示状態、インデックスを直接操作している
- 同じUIの操作方法が複数Sceneへ広がっている
- Canvas APIにはPartsListの操作も含まれるため、Canvas API定義より先に意味ベースAPIを確定させる必要がある

`ShowFoundParts`、`ShowParts`、`Show`、`Hide`、`SelectNext`、`SelectPrevious`、`NextTab`、`PreviousTab`、`IsShowingFoundParts`、`BuySelectedPart`を`PartsListUI`へ追加した（`BuySelectedPart`は10でコメントアウトし、代わりに`RefreshFoundParts`を追加した。11で削除した）。
現在の公開APIは[`UI_Components.md`](../UI_Components.md)の「PartsListUI」を参照する。

- 旧API（`MoveSelection`、`MoveTab`、`SelectedIndex`、`SetLists`、`SetParts`、`SetVisible`、`IsFoundTab`、`CurrentTab`）は互換APIとして残す。旧APIを残すことで、`PlayScene`、`TrainScene`、`ResultScene`は変更せずにビルドできる。
- `BuySelectedPart`のため、`game`の購入結果型を`PurchaseResult`として公開し、`NewPartsListUI`の引数で`*game.Gamedata`を渡した。（10で`BuySelectedPart`と`Gamedata`引数を撤去し、購入をSceneへ移した。）
- Sceneの呼び出し置換と旧APIの削除は、8で実施済み。

### 3. canvas基盤の導入とUIの`scene`依存除去（完了）

`canvas.Canvas`、`canvas.CanvasItem`、`canvas.PositionItem`を[`Canvas.go`](../../Source/Canvas/Canvas.go)へ追加した。
`Source/UI`配下のアサーションと`ScrollUI`の子要素の型を`canvas`の型へ置き換え、`ui` packageから`scene`へのimportを除去した。
UIは構造的に`scene.Element`も満たし続けるため、既存Sceneと`construct`は変更していない。

対象:

- `canvas` package（新規）
- `Source/UI`配下の各UI（[`ScrollUI.go`](../../Source/UI/ScrollUI.go)など）

理由:

- `Source/UI`の各UIは、`scene.Element`へのアサーションのためだけに`scene` packageをimportしている
- `ScrollUI`の子要素が`scene.Element`型であり、UIがScene packageの型に依存している
- Canvas基盤はSceneを知らない汎用コンテナであり、Scene側API interfaceとは独立に作れる

`canvas.Canvas`と`canvas.CanvasItem`を作成し、UI側のアサーションと`ScrollUI`の子要素の型を`canvas.CanvasItem`へ置き換える。
UIは構造的に`scene.Element`も満たし続けるため、既存Sceneと`construct`の`AddUIElement`は変更しない。

### 4. ButtonUIの入力責務整理（完了）

マウス判定（`IsMouseOver`、`IsButtonPressed`）をコメントアウトし、`SetSelected`、`Selected`と選択枠の描画を追加した。
Scene側の入力移行は8で完了したが、旧マウス判定は再利用の可能性があるため理由付きのコメントアウトで残す。

対象:

- [`ButtonUI.go`](../../Source/UI/ButtonUI.go)

理由:

- `ButtonUI`がマウス入力を直接処理している
- `input.Handler`へ入力を集約する方針と一致していない
- 現在`ButtonUI`の呼び出し元がないため、変更がどこにも波及しない

入力をScene・Component側へ集約し、ButtonUIは選択状態などの表示だけを担当する（`SetSelected`、`Selected`）。

### 5. ScrollUIの汎用責務整理（再評価）（完了）

再評価の結果、`ScrollUI`の汎用化（入力解釈・子要素の座標変換・追加API）は行わないことにした。
呼び出し元のない`PartsListUI.Scroll`はコメントアウトした。

対象:

- [`ScrollUI.go`](../../Source/UI/ScrollUI.go)
- [`PartsListUI.go`](../../Source/UI/PartsListUI.go)

懸念していた点:

- スクロール入力の解釈や子要素の座標変換が、将来複数UIへ漏れる可能性がある

再評価の結果:

- `ScrollUI`は入力を扱っていない。クランプ付きのオフセットを受け取り、描画と子要素の`Update`だけを行うため、UIの入力方針に適合している。
- 子要素は3で`canvas.CanvasItem`化済みで、`scene`への依存もない。
- 利用箇所は`PartsListUI`の内部だけで、スクロール位置は`ensureSelectedVisible`による選択追従でUI内部に閉じている。
- 子要素のマウス座標変換が必要になる要因だった`ButtonUI`のマウス判定は4で無効化済みのため、座標変換は不要になった。
- `PartsListUI.Scroll`は呼び出し元がなく、Sceneにスクロール量を操作させない方針と矛盾するためコメントアウトした。ホイール等の相対スクロールやヘルパー実装で再利用する可能性があるため、削除はしていない。

`ScrollUI`を利用するUIが2つ目に現れた時点で、選択追従（`ScrollIntoView`相当）や最大オフセットの取得などを`ScrollUI`へ移すか検討する。

### 6. TitleCanvasの骨組み（`scenecanvas`）（完了。7・8-1と同時に実施）

Canvas APIを`Update`・`Draw`・識別用メソッドだけで構成する方針が決まったため、7・8-1と同時に実施した。
`TitleScene`と同じく静的な`ui.Text`1個だけの`PauseScene`、`SelectStageScene`も同時に移行した。

対象:

- `scenecanvas` package（新規、[`Source/SceneCanvas`](../../Source/SceneCanvas/)）
- [`TitleScene.go`](../../Source/Scene/TitleScene.go)、[`PauseScene.go`](../../Source/Scene/PauseScene.go)、[`SelectStageScene.go`](../../Source/Scene/SelectStageScene.go)
- `construct`の各Scene生成箇所

`scenecanvas.TitleCanvas`、`PauseCanvas`、`SelectStageCanvas`を作成し、`canvas.Canvas`へ案内用の`ui.Text`を追加して`Update`、`Draw`を委譲する。
3つのCanvasは共通ヘルパー`newTextCanvas`で構築する。

## フェーズB: Scene側Canvas API interfaceに依存する作業

### 7. Scene側Canvas API interfaceの定義（完了）

対象:

- `scene` package

各Canvas API interfaceは、対応するScene構造体と同じファイルに定義する。
基本は`Update`、`Draw`と識別用メソッドだけで構成する。
識別用メソッド（`IsTitleCanvas()`など）は処理を持たず、メソッド集合が同じ他SceneのCanvasを誤って注入しないためのものである。
`scenecanvas` packageから実装できるよう、識別用メソッドは公開メソッドとする（非公開メソッドは別packageから実装できない）。

```go
// TitleScene.go
type TitleCanvasAPI interface {
    Update(active bool) error
    Draw(screen *ebiten.Image)
    IsTitleCanvas()
}
```

`Update`・`Draw`だけで完結しないSceneは、次の規則でAPIを追加する。

- 入力起点の操作は、そのSceneでの意図に合わせた名前の入力用APIにする（例: `MovePartSelectionDown`、`ScrollPartsListDown`）。UI部品の操作名（`SelectNext`など）をそのまま出さない。
- シーン固有の情報が必要な場合、CanvasはSceneのポインタを持ち、Sceneの公開メソッドで参照する。
- 入力以外の状態変化（TimeUp、タイマー満了、パーツ種別の進行など）は、Canvasが`Update`でScene・Gamedataを観測して反映する。SceneからCanvasへの通知APIは作らない。
- 戻り値でUIの状態を返す問い合わせAPIは、個別に合意したものだけにする。

| Canvas API | 入力用 | 問い合わせ | Canvasが観測するもの |
|---|---|---|---|
| `TitleCanvasAPI` / `PauseCanvasAPI` / `SelectStageCanvasAPI` | なし | なし | なし |
| `PlayCanvasAPI` | `TogglePartsList`、`MovePartSelectionDown/Up`、`ShowNextPartsTab`、`ShowPreviousPartsTab` | `IsPartsListOpen`、`SelectedFoundPart` | `PlayScene.SelectListRevision`、`Gamedata.PartsRevision` |
| `TrainCanvasAPI` | `ScrollPartsListDown/Up` | なし | `TrainScene.PartsListRevision`、`Gamedata.PartsRevision` |
| `ResultCanvasAPI` | `MovePartSelectionDown/Up` | `SelectedPart`、`IsInstalledPartsFinished` | `ResultScene.CandidatesRevision`・`CandidateParts`、`Gamedata.State` |

（観測対象は11でリビジョン番号にそろえた。）

### 8. Canvas APIの反映とSceneの`elements[index]`依存の全廃（完了）

対象:

- `scene` package、`scenecanvas` package、`construct` package
- 全Scene

理由:

- UIの追加・削除・順番変更で誤ったUIを更新する
- 実行時まで検出しにくい

実施内容:

1. Title・Pause・SelectStageは、`construct`で生成したCanvasを`NewXxxScene`の引数で注入する。
2. Play・Train・Resultは、CanvasがSceneのポインタを持つため、Scene生成後に`SetCanvas`で注入する。（11で、Canvasの生成関数をコンストラクタに渡す方式に変え、`SetCanvas`を削除した。）
3. Play・Train・Resultの`Update`は本体を`updateScene`へ移し、その後で必ず`canvas.Update`を呼ぶ。Sceneの途中returnがあっても、同じフレームの状態変化をCanvasが観測できる。
4. CanvasはGamedata.Timerの`Expired()`を呼ばない（通知を消費するため）。タイマー判定はSceneが行い、結果を公開メソッドで渡す。
5. `PartsListUI`の互換API（`SetParts`、`SetLists`、`SetVisible`、`MoveTab`、`MoveSelection`、`SelectedIndex`、`CurrentTab`、`IsFoundTab`）を削除した。
6. `scene.Element`、`PositionElement`、`UIDrawer`、`camera.UICamera`、`construct.MakeUIDrawerFactory`、Sceneの`elements`・`AddUIElement`・`partsListUI`・`installedPartsUI`を削除した。

### 9. Gamedata参照の縮小（保留）

対象:

- [`GameDataUI.go`](../../Source/UI/GameDataUI.go)
- [`GameDataText.go`](../../Source/UI/GameDataText.go)
- [`InstalledParts.go`](../../Source/UI/ResultUI/InstalledParts.go)

理由:

- UIがModel全体へ依存している
- Modelの型変更が多くのUIへ波及する

UIが`*game.Gamedata`を直接参照することを許容するため、現時点では実施しない。
Model参照をReaderや表示用データへ縮小する必要が生じた場合に、改めて検討する。

### 10. 購入処理のScene移行（完了）

購入処理をUIからSceneへ移した。`ResultScene`の`SelectedPart`と同じ形にそろえた。

- `PlayCanvasAPI.BuySelectedPart`を削除し、問い合わせAPI`SelectedFoundPart() *game.Part`を追加した。一覧の非表示中、発見タブ以外の表示中、未選択のときは`nil`を返す。
- `PlayScene`が`SelectedFoundPart`で受け取ったパーツを`Gamedata.BuyFoundPart`へ渡して購入する。`BuyFoundPart`はポインタで`Found`内の位置を解決し、`BuyPart(index)`へ委譲する。
- `PartsListUI.BuySelectedPart`はコメントアウトし、再構築用の`RefreshFoundParts`を追加した（11で削除した）。`NewPartsListUI`の`Gamedata`引数（2で追加したもの）は不要になったため外した。
- 購入後の一覧の再構築は、`PlayCanvas`が`Gamedata.Found`・`Parts`の変化を`Update`で観測して行う（11で`Gamedata.PartsRevision`の観測に変えた）。

### 11. PR #127 レビュー指摘の反映（完了）

- Canvasの注入: Play・Train・Resultは、Canvasの生成関数（`func(*XxxScene) (XxxCanvasAPI, error)`）を`NewXxxScene`の引数で受け取り、Sceneの生成直後に呼んで注入する。注入忘れや`nil`はScene生成時のエラーになる。`SetCanvas`と、Update・Drawでのcanvasの`nil`判定を削除した。CanvasはSceneの具体型に依存したままにする（別Sceneを誤って受け取らないため）。`PlayCanvas`も`PlayScene`のポインタを持つようにした。
- 候補のキャッシュ: `ResultScene.CandidateParts`は副作用のない取得APIにした。候補の計算（Trashの生成を含む）は、`refreshCandidates`が候補の変わる時点（シーン開始時、パーツ選択後）に1回だけ行う。
- 状態変化の検知をリビジョン番号にそろえた。
  - `Gamedata.PartsRevision`: `Found`・`Parts`を書き換えるGamedataのメソッド（`AddPartToFoundList`、`AddPartToPartsList`、`BuyPart`、`ResetGame`）の中で増やす。Player の拾得のように、Sceneを通らない書き換えでも増やし忘れない。
  - `PlayScene.SelectListRevision`: `Play_Select`に入ったときに増やす。`PlayCanvas`は`slices.Equal`による比較と`slices.Clone`による複製をやめた。
  - `TrainScene.PartsListRevision`: `IsPartsListRevealed`の代わり。到着演出のタイマー満了時に1回だけ増やす。
- 非表示中の操作の判定は`PartsListUI`の`SelectNext`・`SelectPrevious`・`NextTab`・`PreviousTab`・`SelectedPart`に入れ、Canvas側の重複した判定を外した。Resultでも同じ判定が効く。
- `ResultCanvas`は毎フレーム`Canvas.Update`で全UIを更新する。ベンチマーク案内のアニメーションは、`Result_Benchmark`に変わったときに開始する。`Canvas.Show/Hide/Visible`と`PositionItem`は、後で使う可能性があるため残した。
- `uiutils`（`Source/UI/UIUtils`）に`FillRect`・`DrawBorder`・`SelectionBorderColor`を追加した。`ButtonUI`と`PartsListUI`の描画で、毎フレームの`ebiten.NewImage`をやめた。
- 不要になったコードは、方針（実施上の注意を参照）に従い、`PartsListUI.BuySelectedPart`を削除した。`ButtonUI`のマウス判定と`PartsListUI.Scroll`はコメントアウトで残す。

## 現時点の保留事項

- UIのライフサイクルと再生成・リセット方法（UIごとの実装事項として扱い、基本方針では規定しない）
- `PartsListUI`（`PartCard`、選択・タブの状態、レイアウト計算）と`InstalledParts`（1パーツ分のレコード、種別名の変換）のファイル分割。ファイルが大きくなった時点で検討する。分割後もゲームルールやScene遷移はこれらへ移さない。

## 実施上の注意

- 一度に全UIを置き換えない。
- まずModel更新・Scene遷移の責務をUIから取り除く。
- 次にSceneからUI内部状態への依存を取り除く。
- Canvas移行では、Scene側のCanvas API interfaceと具体Canvas実装を分離する。
- フェーズAでは`scene` packageを変更せず、別担当者の実装と衝突させない。
- UIの見た目変更と責務変更を同じ変更に詰め込みすぎない。
- リファクタで不要になったコードは、再利用の可能性がある（マウス操作の復活など）なら理由付きでコメントアウトし、そうでなければ削除する。
- 各移行の完了時に、関連ドキュメントが実装と設計方針に一致しているか確認する。
- ドキュメントの陳腐化や、実装済みの内容が移行案のまま残っていないか確認する。
- 方針変更・API変更・責務変更があった場合は、`UI.md`、`UI_Components.md`、該当するリファクタリング文書、クラス図を同じ変更で更新する。
- 各移行の検証は、AIによる静的確認と手動による実機確認に分ける。

### AIが実施する検証

- `go vet ./...`を実行する。
- `gofmt`を適用し、フォーマット差分が残っていないことを確認する。
- 変更したAPI、型、責務に関する依存箇所や旧APIの使用箇所を検索する。
- テストが存在する場合は`go test ./...`を実行する。現時点ではテスト未導入のため、必須手順とはしない。

### 手動で実施する実機確認

- ゲームを実機で起動する。
- 対象Sceneへ遷移し、対象UIが表示されることを確認する。
- 入力、選択、タブ切り替え、スクロール、アニメーション、画面遷移を操作して確認する。
- リファクタ前後で、対象機能の操作結果と表示が変わっていないことを確認する。

### 移行ごとのドキュメント確認

各リファクタリングを完了したら、次の順でドキュメントを確認する。

1. 実装とリファクタリング案の差分を確認する。
2. 実装済みの方針を`UI.md`へ反映する。
3. 公開APIや挙動を`UI_Components.md`へ反映する。
4. 完了した移行手順や古い設計案を、該当する`UIRefactoring`文書から更新または削除する。
5. クラス図と実装上の型・依存関係が一致しているか確認する。
