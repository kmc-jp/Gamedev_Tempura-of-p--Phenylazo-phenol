# UI コンポーネント仕様

## 共通仕様

`Source/UI` の UI 要素は、次の `Element` インターフェースを実装する。

```go
type Element interface {
    Update(active bool) error
    Draw(screen *ebiten.Image)
}
```

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

`ButtonUI` は `RichUI` をラップし、画面座標上のマウスポインタと左クリックを判定する。

- `IsMouseOver() bool`: 現在ポインタが境界内にあるか
- `IsButtonPressed() bool`: `Update(true)` を呼んだフレームで左クリックが開始されたか
- `RichUI() *RichUI`: 内包要素を取得

押下判定は `inpututil.IsMouseButtonJustPressed` を利用するため、押し続けても押下開始フレームのみ `true` になる。判定は `Update(true)` 後に行う。`Update(false)` では状態を false にリセットし、入力を処理しない。

## ScrollUI

`ScrollUI` は複数の `Element` をコンテンツ用オフスクリーン画像へ描画し、ビューポート部分だけを表示する。

- `ScrollVertical`: 上下方向
- `ScrollHorizontal`: 左右方向
- `AddChild`: 子要素を追加
- `Scroll`: 相対スクロール
- `SetOffset`: 絶対スクロール位置を設定
- `Offset`: 現在位置を取得
- `Bounds`: ビューポートの画面境界を取得

子要素の座標はコンテンツバッファの座標系で指定する。スクロール量は自動的に 0 から最大値までクランプされる。

現時点の ScrollUI は表示と子要素の更新を担当し、マウスホイールやキー入力の解釈は呼び出し側に委譲する。また、スクロール中の子ボタンのマウス座標変換は未実装であり、表示用コンテナとして利用する。

## 変更記録

今回の UI 拡張では既存の `Text`、シーンの描画 API、入力設定を破壊的に変更していない。既存シーンのテキスト表示は `Text` 要素を経由するよう変更済みである。
