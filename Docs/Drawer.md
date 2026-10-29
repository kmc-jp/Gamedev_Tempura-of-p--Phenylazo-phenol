# 描画系

## game.Game.Draw(screen *ebiten.Image)

### やること
後方のシーンから順に画面全体を描画する

### Scene.Draw(screen *ebiten.Image) に期待すること
そのシーン全体 (裏に何かあるならそれの処理も含む) をいい感じに描画すること

## scene.PlayScene.Draw(screen *ebiten.Image) (ほかのシーンも概ね同様 : 不要なら片方だけになる)

### やること
シーンが保持するワールド要素とUI要素を描画する。

- WorldDrawer に Entity の描画を委譲する
- UIDrawer に Element の描画を委譲する

## scene.Camera.Draw(screen *ebiten.Image, objects []Entity) : WorldDrawer

### やること
受け取った objects の中身を描画する
画面中心は center に合わせる

### Entity.Draw() (image *ebiten.Image, Layer string) に期待すること
適切なサイズの画像を描いてレイヤーと一緒に返してくること

## scene.Camera.SetCenter(center)

中心座標の設定

## GameObject.Draw() (image *ebiten.Image)
自分の画像を返す
(実装はRenderが持っている)

## scene.UICamera.Draw(screen *ebiten.Image, elements []Element) : UIDrawer

### やること
受け取った elements の中身の .Draw() を呼ぶ

### Element.Draw(screen *ebiten.Image) に期待すること
いい感じの場所に自分のUIを設置すること