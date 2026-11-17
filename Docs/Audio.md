# 曲を流す

## 全体

- シリアライズ系統 (audio.Context を持たせておく) で直接 audio.Player を生成し、シーン等で流す
  - BGMは audio.Player 実体をキャッシュ (未実装)
    - 一旦バイト列を直接キャッシュする。 重くなってきたら実装変更。
  - 効果音は デコード後の生データを []byte で持っておいて NewPlayerFromBytes() で生成
