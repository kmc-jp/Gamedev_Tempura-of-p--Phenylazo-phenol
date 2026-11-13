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
添字参照は`scene` package内の名前付きフィールドへ置き換え済み（Canvas化は8で実施）。

対象:

- [`InstalledParts.go`](../../Source/UI/ResultUI/InstalledParts.go)
- [`ResultScene.go`](../../Source/Scene/ResultScene.go)

理由:

- UIが`Gamedata.NextState()`を呼び出している
- UIがゲーム状態の遷移を実質的に決めている
- ResultSceneが`elements[1]`、`elements[2]`の添字に依存している

最初に`IsFinished()`だけを導入し、状態遷移をSceneへ戻す。

### 2. PartsListUIの意味ベースAPI追加（UI側のみ）

対象:

- [`PartsListUI.go`](../../Source/UI/PartsListUI.go)
- [`PartsManager.go`](../../Source/Game/PartsManager.go)
- `construct`の`PartsListUI`生成箇所（引数追加のみ）

理由:

- Sceneが選択、タブ、表示状態、インデックスを直接操作している
- 同じUIの操作方法が複数Sceneへ広がっている
- Canvas APIにはPartsListの操作も含まれるため、Canvas API定義より先に意味ベースAPIを確定させる必要がある

`ShowFoundParts`、`ShowParts`、`Show`、`Hide`、`SelectNext`、`SelectPrevious`、`NextTab`、`PreviousTab`、`IsShowingFoundParts`、`BuySelectedPart`を`PartsListUI`へ追加する。
詳細は[`PartsListUI_refactoring.md`](./PartsListUI_refactoring.md)を参照する。

- 旧API（`MoveSelection`、`MoveTab`、`SelectedIndex`、`SetLists`、`SetVisible`）は互換APIとして残す。旧APIを残すことで、`PlayScene`、`TrainScene`、`ResultScene`は変更せずにビルドできる。
- `BuySelectedPart`のため、`game`の購入結果型（現在は非公開の`purchase`）を公開し、`PartsListUI`へ`*game.Gamedata`を渡す。
- Sceneの呼び出し置換と旧APIの削除は、8で行う。

### 3. canvas基盤の導入とUIの`scene`依存除去

対象:

- `canvas` package（新規）
- `Source/UI`配下の各UI（[`ScrollUI.go`](../../Source/UI/ScrollUI.go)など）

理由:

- `Source/UI`の各UIは、`scene.Element`へのアサーションのためだけに`scene` packageをimportしている
- `ScrollUI`の子要素が`scene.Element`型であり、UIがScene packageの型に依存している
- Canvas基盤はSceneを知らない汎用コンテナであり、Scene側API interfaceとは独立に作れる

`canvas.Canvas`と`canvas.CanvasItem`を作成し、UI側のアサーションと`ScrollUI`の子要素の型を`canvas.CanvasItem`へ置き換える。
UIは構造的に`scene.Element`も満たし続けるため、既存Sceneと`construct`の`AddUIElement`は変更しない。

### 4. ButtonUIの入力責務整理

対象:

- [`ButtonUI.go`](../../Source/UI/ButtonUI.go)

理由:

- `ButtonUI`がマウス入力を直接処理している
- `input.Handler`へ入力を集約する方針と一致していない
- 現在`ButtonUI`の呼び出し元がないため、変更がどこにも波及しない

入力をScene・Component側へ集約し、ButtonUIは選択状態などの表示だけを担当する（`SetSelected`、`Selected`）。

### 5. ScrollUIの汎用責務整理（再評価）

対象:

- [`ScrollUI.go`](../../Source/UI/ScrollUI.go)

理由:

- スクロール入力の解釈や子要素の座標変換が、将来複数UIへ漏れる可能性がある
- 現在は主にPartsListUI内部で利用されている

2・3の完了後に、汎用化が必要かを再評価する。

### 6. TitleCanvasの骨組み（`scenecanvas`）

対象:

- `scenecanvas` package（新規）
- [`TitleScene.go`](../../Source/Scene/TitleScene.go)（参照のみ。変更しない）

理由:

- `TitleScene`のUIは静的な`ui.Text`1個だけで、Sceneの`Update()`からUIを操作していない
- Canvas APIは`Update`と`Draw`だけになることが自明であり、API確定前に作っても手戻りが少ない
- Canvas基盤の形を、最も単純なSceneで先に確認できる

`scenecanvas.TitleCanvas`を作成し、`canvas.Canvas`へタイトル表示の`ui.Text`を追加して`Update`、`Draw`を委譲する。
`scene.TitleCanvasAPI`への適合宣言と`construct`からの注入は、TODOとして8へ回す。
`PauseScene`、`SelectStageScene`も同じ構成だが、まずは`TitleScene`だけで形を確認する。

## フェーズB: Scene側Canvas API interfaceに依存する作業

### 7. Scene側Canvas API interfaceの定義（別担当者）

対象:

- `scene` package

`scene.TitleCanvasAPI`、`scene.ResultCanvasAPI`などのCanvas API interfaceを定義する。
フェーズAで確定した`InstalledParts`の公開API（`Start`、`IsFinished`）と、PartsListの意味ベースAPI（`SelectNext`、`NextTab`など）を取り込む。

### 8. Canvas APIの反映とSceneの`elements[index]`依存の全廃

対象:

- `scene` package、`scenecanvas` package、`construct` package
- [`ResultScene.go`](../../Source/Scene/ResultScene.go)
- [`PlayScene.go`](../../Source/Scene/PlayScene.go)
- [`TrainScene.go`](../../Source/Scene/TrainScene.go)

理由:

- UIの追加・削除・順番変更で誤ったUIを更新する
- 実行時まで検出しにくい

現状: `ResultScene`の`elements[index]`による更新は、名前付きフィールドへ暫定移行済み。
具体的なCanvasは`construct`から注入し、SceneはCanvas APIだけを利用する。

7の完了後、次の順で進める。

1. `TitleCanvas`へ`scene.TitleCanvasAPI`の適合宣言を追加し、`construct`から注入する（最初の接続対象）。
2. `ResultCanvas`を実装し、`ResultScene`をCanvas化する。UIの数が少なく、1・2でAPIが固まっているため、本格移行の最初の対象とする。
3. `PlayScene`、`TrainScene`などの他Sceneを、Sceneごとに別の変更として移行する。
4. Scene内のPartsList呼び出しを意味ベースAPIへ置き換え、`scene.partsListUI`を`PartsList`へ変更し、旧APIを削除する（2の残り）。
5. `scene.Element`、Sceneの`elements []Element`、`UIDrawer`、`UICamera`を撤去する。

8の完了までは、名前付きフィールドの暫定構成を維持する。

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

## 現時点の保留事項

- UIのライフサイクルと再生成・リセット方法（UIごとの実装事項として扱い、基本方針では規定しない）

## 実施上の注意

- 一度に全UIを置き換えない。
- まずModel更新・Scene遷移の責務をUIから取り除く。
- 次にSceneからUI内部状態への依存を取り除く。
- Canvas移行では、Scene側のCanvas API interfaceと具体Canvas実装を分離する。
- フェーズAでは`scene` packageを変更せず、別担当者の実装と衝突させない。
- UIの見た目変更と責務変更を同じ変更に詰め込みすぎない。
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
