# UIリファクタリング優先順位

## 判断基準

優先順位は、次の観点で決める。

1. UIからModel更新やScene遷移が行われているか
2. SceneがUIの内部実装に依存しているか
3. UIの順番や添字に依存しているか
4. 複数Sceneから利用され、変更の影響範囲が大きいか
5. 現在の開発を止めずに段階移行できるか

## 優先順位

### 1. InstalledPartsUIの責務分離

対象:

- [`InstalledParts.go`](../../Source/UI/ResultUI/InstalledParts.go)
- [`ResultScene.go`](../../Source/Scene/ResultScene.go)

理由:

- UIが`Gamedata.NextState()`を呼び出している
- UIがゲーム状態の遷移を実質的に決めている
- ResultSceneが`elements[1]`、`elements[2]`の添字に依存している

最初に`IsFinished()`だけを導入し、状態遷移をSceneへ戻す。

### 2. PartsListUIのScene向けAPI整理

対象:

- [`PartsListUI.go`](../../Source/UI/PartsListUI.go)
- [`PlayScene.go`](../../Source/Scene/PlayScene.go)
- [`TrainScene.go`](../../Source/Scene/TrainScene.go)
- [`ResultScene.go`](../../Source/Scene/ResultScene.go)

理由:

- Sceneが選択、タブ、表示状態、インデックスを直接操作している
- 同じUIの操作方法が複数Sceneへ広がっている
- UIの内部変更が複数Sceneへ波及する

`SelectNext`、`SelectPrevious`、`NextTab`などの意味ベースAPIへ移行する。
詳細は[`PartsListUI_refactoring.md`](./PartsListUI_refactoring.md)を参照する。

### 3. Sceneの`elements[index]`依存の全廃

対象:

- [`ResultScene.go`](../../Source/Scene/ResultScene.go)

理由:

- UIの追加・削除・順番変更で誤ったUIを更新する
- 実行時まで検出しにくい

`ResultCanvas`などの名前付きstructへ移行する。

### 4. ButtonUIの入力責務整理

対象:

- [`ButtonUI.go`](../../Source/UI/ButtonUI.go)

理由:

- `ButtonUI`がマウス入力を直接処理している
- `input.Handler`へ入力を集約する方針と一致していない

入力をScene・Component側へ集約し、ButtonUIは選択状態などの表示だけを担当する。

### 5. ScrollUIの汎用責務整理

対象:

- [`ScrollUI.go`](../../Source/UI/ScrollUI.go)

理由:

- スクロール入力の解釈や子要素の座標変換が、将来複数UIへ漏れる可能性がある
- 現在は主にPartsListUI内部で利用されている

PartsListUIのリファクタ後に、汎用化が必要かを再評価する。

### 6. Gamedata参照の縮小

対象:

- [`GameDataUI.go`](../../Source/UI/GameDataUI.go)
- [`GameDataText.go`](../../Source/UI/GameDataText.go)
- [`InstalledParts.go`](../../Source/UI/ResultUI/InstalledParts.go)

理由:

- UIがModel全体へ依存している
- Modelの型変更が多くのUIへ波及する

ただし、現時点では参照のみであり、優先度は低い。ゲーム仕様が固まった後に必要な情報だけのReaderや表示用データへ移行する。

## 実施上の注意

- 一度に全UIを置き換えない。
- まずModel更新・Scene遷移の責務をUIから取り除く。
- 次にSceneからUI内部状態への依存を取り除く。
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
