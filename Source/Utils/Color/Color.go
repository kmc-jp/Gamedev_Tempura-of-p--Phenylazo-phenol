package color

// RGBA type
type Color struct {
	R float32
	G float32
	B float32
	A float32
}

func Red() *Color {
	return &Color{1, 0, 0, 1}
}

func Yellow() *Color {
	return &Color{1, 1, 0, 1}
}

func Green() *Color {
	return &Color{0, 1, 0, 1}
}

func Cyan() *Color {
	return &Color{0, 1, 1, 1}
}

func Blue() *Color {
	return &Color{0, 0, 1, 1}
}

func Magenta() *Color {
	return &Color{1, 0, 1, 1}
}

func Black() *Color {
	return &Color{0, 0, 0, 1}
}

func White() *Color {
	return &Color{1, 1, 1, 1}
}

func Gray() *Color {
	return &Color{0.5, 0.5, 0.5, 1}
}

func LightGray() *Color {
	return &Color{0.8, 0.8, 0.8, 1}
}

func DarkGray() *Color {
	return &Color{0.2, 0.2, 0.2, 1}
}

func Orange() *Color {
	return &Color{1, 0.5, 0, 1}
}

func Pink() *Color {
	return &Color{1, 0.75, 0.8, 1}
}

func Purple() *Color {
	return &Color{0.5, 0, 0.5, 1}
}

func Brown() *Color {
	return &Color{0.6, 0.3, 0, 1}
}

func Navy() *Color {
	return &Color{0, 0, 0.5, 1}
}

func Teal() *Color {
	return &Color{0, 0.5, 0.5, 1}
}

func Transparent() *Color {
	return &Color{0, 0, 0, 0}
}

func SemiTransparentBlack() *Color {
	return &Color{0, 0, 0, 0.5}
}

func SemiTransparentWhite() *Color {
	return &Color{1, 1, 1, 0.5}
}

// 赤オレンジ / 朱色 (#FF4500)
func RedOrange() *Color {
	return &Color{1, 0.27, 0, 1}
}

// クリムゾン / 深紅 (#DC143C)
func Crimson() *Color {
	return &Color{0.86, 0.08, 0.24, 1}
}

// マルーン / えんじ色 (#800000)
func Maroon() *Color {
	return &Color{0.5, 0, 0, 1}
}

// コーラル / 珊瑚色 (#FF7F50)
func Coral() *Color {
	return &Color{1, 0.5, 0.31, 1}
}

// ゴールド / 金色 (#FFD700)
func Gold() *Color {
	return &Color{1, 0.84, 0, 1}
}

// 深緑 / ダークグリーン (#006400)
func DarkGreen() *Color {
	return &Color{0, 0.39, 0, 1}
}

// フォレストグリーン / 森林の緑 (#228B22)
func ForestGreen() *Color {
	return &Color{0.13, 0.55, 0.13, 1}
}

// ライム / 黄緑 (#32CD32)
func Lime() *Color {
	return &Color{0.2, 0.8, 0.2, 1}
}

// ミントグリーン (#98FF98)
func MintGreen() *Color {
	return &Color{0.6, 1, 0.6, 1}
}

// オリーブ (#808000)
func Olive() *Color {
	return &Color{0.5, 0.5, 0, 1}
}

// 水色 / スカイブルー (#87CEEB)
func SkyBlue() *Color {
	return &Color{0.53, 0.81, 0.98, 1}
}

// ライトブルー / 薄水色 (#ADD8E6)
func LightBlue() *Color {
	return &Color{0.68, 0.85, 0.9, 1}
}

// ロイヤルブルー / 鮮やかな濃青 (#4169E1)
func RoyalBlue() *Color {
	return &Color{0.25, 0.41, 0.88, 1}
}

// インディゴ / 藍色 (#4B0082)
func Indigo() *Color {
	return &Color{0.29, 0, 0.51, 1}
}

// バイオレット / 菫色 (#8A2BE2)
func Violet() *Color {
	return &Color{0.54, 0.17, 0.89, 1}
}

// ラベンダー (#E6E6FA)
func Lavender() *Color {
	return &Color{0.9, 0.9, 0.98, 1}
}

// ホットピンク (#FF69B4)
func HotPink() *Color {
	return &Color{1, 0.41, 0.71, 1}
}

// ローズ / 薔薇色 (#FF007F)
func Rose() *Color {
	return &Color{1, 0, 0.5, 1}
}

// ベージュ (#F5F5DC)
func Beige() *Color {
	return &Color{0.96, 0.96, 0.86, 1}
}

// カーキ (#C3B091)
func Khaki() *Color {
	return &Color{0.76, 0.69, 0.57, 1}
}
