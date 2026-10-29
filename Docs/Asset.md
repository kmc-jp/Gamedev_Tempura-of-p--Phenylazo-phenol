# アセット読み込み周り

## 規約
- アセットは /Assets 配下に入れる
  - 絵は /Assets/Images
  - 音は /Assets/Sounds
  - シリアライズは /Assets/Serialize
- アセットは不変であることを前提とする
  - 要するに実行中に増えたり減ったりしない
  - 実行中の増減は反映されない

## 実装

### AssetManager
- 利用側末端
- キーとパスの接続をやる

### serialize.SerializeManager
- シリアライズ系の末端

### serialize.TomlLoader
- .toml を読み込んで対応する型に変換する

### imageassets.ImageManager
- Image 系の末端

### imageassets.AsepriteLoader
- .aseprite を読み込んで ebiten.image にする
- キャッシュとかその辺はこいつが担当
