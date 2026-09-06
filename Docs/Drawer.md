# 描画系(案)

## game.Game.Draw(screen *ebiten.Image)

### やること
後方のシーンから順に画面全体を描画する

### Scene.Draw(screen *ebiten.Image) に期待すること
そのシーン全体 (裏に何かあるならそれの処理も含む) をいい感じに描画すること

## scene.PlayScene.Draw(screen *ebiten.Image)

### やること
Render.Draw() を呼ぶ

## scene.MockCamera.Draw(screen *ebiten.Image, objects []Entity, center Position)

### やること
受け取った objects の中身のすべてについて、そのワールド座標をそのまま描画する
カメラ位置の調整周りは後から考える

### Entity.Render.Draw() (screen *ebiten.Image) に期待すること
適切なサイズの画像を描いて返してくること

## GameObject.Draw() (screen *ebiten.Image)
自分の画像を返す