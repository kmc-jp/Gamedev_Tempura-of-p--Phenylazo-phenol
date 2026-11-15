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

### AssetsManager
- 利用側末端
- キーとパスの接続をやる

### serialize.SerializeManager
- シリアライズ系の末端

### serialize.TomlLoader
- .toml を読み込んで対応する型に変換する

### serializetarget package
- シリアライズ対象

#### 主なシリアライズ対象
- `InputConfig.toml`: アクションとキーの対応（キーマップ）
- `CameraConfig.toml`: 描画レイヤー名とレイヤー番号の対応
- `Player.toml`: プレイヤーの名前・初期位置・歩行時間・アニメーションタグ・タイル定義
- `TileMap.toml`: タイルマップ定義の書式見本（現在はコードから読み込まれていない）
- `StageConfig.toml`: ステージ名とマップ定義ファイルの対応
- `Stages/*.toml`: ステージごとのタイルマップ定義（タイルサイズ、床・壁、タイル、オブジェクト配置。`StageConfig.toml` から参照される）
- `StageMap.toml`: 行き先選択画面の路線図（駅の座標・ラベルと、駅をつなぐ路線）。座標は `Width` x `Height` の基準サイズ上の値で、絵が届いたら絵に合わせて調整する

### imageassets.ImageManager
- Image 系の末端

### imageassets.AsepriteLoader
- .aseprite を読み込んで ebiten.image にする
- キャッシュとかその辺はこいつが担当
