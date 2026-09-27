# アセット読み込み周り

## 規約
- アセットは /Asset 配下に入れる
  - 絵は /Asset/Images
  - 音は /Asset/Sounds
  - シリアライズは /Asset/Serialize
  - 設定ファイルは場所未定
- アセットは不変であることを前提とする
  - 要するに実行中に増えたり減ったりしない
  - 実行中の増減は反映されない

## 実装

### AssetManager
- 利用側末端
- キーとパスの接続をやる

### SerializeManager
- シリアライズ系の末端

### TomlLoader
- .toml を読み込んで対応する型に変換する

### ImageManager
- Image 系の末端

### AsepriteLoader
- .aseprite を読み込んで ebiten.image にする
- キャッシュとかその辺はこいつが担当
