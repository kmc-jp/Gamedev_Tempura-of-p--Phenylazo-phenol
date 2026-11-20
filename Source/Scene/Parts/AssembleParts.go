package parts

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"fmt"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

var AssembleImages [numberOfAssemblePartsAssets]*ebiten.Image

// 組み立て用画像一覧 (意図的に_を使用)
// A: Assemble であることを識別するための識別子 ,B: Black, W: White, S: Side glass, SF: Side and Front glass,
const (
	A_case_B_SF      = iota //みんげー_PCケース_黒_蓋開いてる
	A_case_B_S              //みんげー_PCケース_黒_蓋閉じてる
	A_case_W_SF             //みんげー_PCケース_白_蓋開いてる
	A_case_W_S              //みんげー_PCケース_白_蓋閉じてる
	A_case_cardboard        //みんげー_段ボールPCケース

	A_case_SF_lid        //みんげー_PCケース_蓋_開いてるとき_通常
	A_case_S_lid         //みんげー_PCケース_蓋_閉じてるとき_通常
	A_case_cardboard_lid //みんげー_段ボールPC_蓋

	A_caseSF_lid_gaming //みんげー_PCケース_蓋_開いてるとき_ゲーミング
	A_case_S_lid_gaming //みんげー_PCケース_蓋_閉じてるとき_ゲーミング

	A_gaming_B_SF //みんげー_ゲーミングカラー_追加パーツ_黒_蓋開いてるとき
	A_gaming_B_S  //みんげー_ゲーミングカラー_追加パーツ_黒_蓋閉じてるとき
	A_gaming_W_SF //みんげー_ゲーミングカラー_追加パーツ_白_開いてるとき
	A_gaming_W_S  //みんげー_ゲーミングカラー_追加パーツ_白_閉じてるとき

	A_mb_Black //みんげー_マザボ_リザルト用_黒
	A_mb_W     //みんげー_マザボ_リザルト用_白
	A_mb_Blue  //みんげー_マザボ_リザルト用_青
	A_mb_G     //みんげー_マザボ_リザルト用_緑

	A_cooler_B //みんげー_CPUクーラー_リザルト用_黒
	A_cooler_W //みんげー_CPUクーラー_リザルト用_白

	A_gpu_B //みんげー_グラボ_リザルト用_黒
	A_gpu_W //みんげー_グラボ_リザルト用_白

	A_sticker //みんげー_キャラステッカー_リザルト用

	A_noAdditionalParts //任意パーツの有無を選択するためのもの
)
const numberOfAssemblePartsAssets = A_noAdditionalParts

type AssemblePartsImage uint

func LoadAssemblePartsImages(imageManager *imageassets.ImageManager) error {
	AssembleImages = [numberOfAssemblePartsAssets]*ebiten.Image{}

	loads := []struct {
		part     AssemblePartsImage
		filename string
	}{
		{A_case_B_SF, "みんげー_PCケース_黒_蓋開いてる"},
		{A_case_B_S, "みんげー_PCケース_黒_蓋閉じてる"},
		{A_case_W_SF, "みんげー_PCケース_白_蓋開いてる"},
		{A_case_W_S, "みんげー_PCケース_白_蓋閉じてる"},
		{A_case_cardboard, "みんげー_段ボールPCケース"},

		{A_case_SF_lid, "みんげー_PCケース_蓋_開いてるとき_通常"},
		{A_case_S_lid, "みんげー_PCケース_蓋_閉じてるとき_通常"},
		{A_case_cardboard_lid, "みんげー_段ボールPC_蓋"},

		{A_caseSF_lid_gaming, "みんげー_PCケース_蓋_開いてるとき_ゲーミング"},
		{A_case_S_lid_gaming, "みんげー_PCケース_蓋_閉じてるとき_ゲーミング"},

		{A_gaming_B_SF, "みんげー_ゲーミングカラー_追加パーツ_黒_蓋開いてるとき"},
		{A_gaming_B_S, "みんげー_ゲーミングカラー_追加パーツ_黒_蓋閉じてるとき"},
		{A_gaming_W_SF, "みんげー_ゲーミングカラー_追加パーツ_白_開いてるとき"},
		{A_gaming_W_S, "みんげー_ゲーミングカラー_追加パーツ_白_蓋閉じてるとき"},

		{A_mb_Black, "みんげー_マザボ_リザルト用_黒"},
		{A_mb_W, "みんげー_マザボ_リザルト用_白"},
		{A_mb_Blue, "みんげー_マザボ_リザルト用_青"},
		{A_mb_G, "みんげー_マザボ_リザルト用_緑"},

		{A_cooler_B, "みんげー_CPUクーラー_リザルト用_黒"},
		{A_cooler_W, "みんげー_CPUクーラー_リザルト用_白"},

		{A_gpu_B, "みんげー_グラボ_リザルト用_黒"},
		{A_gpu_W, "みんげー_グラボ_リザルト用_白"},
	}

	for _, load := range loads {
		if err := LoadAssemblePartImage(imageManager, load.part, load.filename); err != nil {
			return err
		}
	}

	return nil
}

func LoadAssemblePartImage(
	imageManager *imageassets.ImageManager,
	part AssemblePartsImage,
	filename string,
) error {
	if part >= AssemblePartsImage(numberOfAssemblePartsAssets) {
		return fmt.Errorf("Assemble Parts定数が不正です: %d", part)
	}

	baseName := strings.TrimSuffix(filename, path.Ext(filename))
	baseName = strings.TrimPrefix(path.Clean(baseName), "Assemble/")
	pngPath := path.Join("Assemble", baseName+".png")
	jsonPath := path.Join("Assemble", baseName+".json")

	imageData, err := imageManager.NewImageData(baseName, pngPath, jsonPath)
	if err != nil {
		return fmt.Errorf("Assemble Parts画像データの作成に失敗しました。定数: %d, ファイル名: %s\n%w", part, filename, err)
	}

	_, image, err := imageManager.Load(*imageData)
	if err != nil {
		return fmt.Errorf("Assemble Parts画像の読み込みに失敗しました。定数: %d, ファイル名: %s\n%w", part, filename, err)
	}

	AssembleImages[part] = image
	return nil
}

func AssembledImage(
	caseBack AssemblePartsImage,
	mb AssemblePartsImage,
	cooler AssemblePartsImage,
	gpu AssemblePartsImage,
	caseLid AssemblePartsImage,
	gaming AssemblePartsImage,
) (*ebiten.Image, error) {
	pc := ebiten.NewImage(72, 64)

	required := []struct {
		name string
		part AssemblePartsImage
		draw bool
	}{
		{"ケース本体", caseBack, true},
		{"マザボ", mb, true},
		{"CPUクーラー", cooler, true},
		{"グラフィックボード", gpu, gpu != AssemblePartsImage(A_noAdditionalParts)},
		{"ケースの蓋", caseLid, caseLid != A_noAdditionalParts},
		{"ゲーミング追加パーツ", gaming, gaming != A_noAdditionalParts},
	}
	for _, image := range required {
		if !image.draw {
			continue
		}
		if image.part >= AssemblePartsImage(numberOfAssemblePartsAssets) {
			return nil, fmt.Errorf("%sに無効な画像定数が指定されています: %d", image.name, image.part)
		}
		if AssembleImages[image.part] == nil {
			return nil, fmt.Errorf("%sの画像が読み込まれていません", image.name)
		}
		if image.part == A_case_cardboard_lid {
			continue
		}
		pc.DrawImage(AssembleImages[image.part], nil)
	}

	if gaming == A_noAdditionalParts {
		return pc, nil
	}
	return pc, nil
}

func findAssemblePart(pts game.Parts, kind uint) *game.Part {
	for _, part := range pts {
		if part != nil && uint(part.Kind) == kind {
			return part
		}
	}
	return nil
}

func DeterminePartsConfiguration(pts game.Parts) (AssemblePartsImage, AssemblePartsImage, AssemblePartsImage, AssemblePartsImage, AssemblePartsImage, AssemblePartsImage, error) {
	var (
		caseback AssemblePartsImage
		mb       AssemblePartsImage
		cooler   AssemblePartsImage
		gpu      AssemblePartsImage
		caselid  AssemblePartsImage
		gaming   = AssemblePartsImage(A_noAdditionalParts)
	)

	gpuPart := findAssemblePart(pts, uint(game.GPU))
	if gpuPart == nil {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("組み立てパーツにGPUがありません")
	}
	if gpuPart.IsSpare {
		gpu = AssemblePartsImage(A_noAdditionalParts)
	} else if gpuPart.Subcategory == game.GPU_Black {
		gpu = A_gpu_B
	} else if gpuPart.Subcategory == game.GPU_White {
		gpu = A_gpu_W
	} else {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("GPUに対して無効なSubcategoryです")
	}

	mbPart := findAssemblePart(pts, uint(game.MB))
	if mbPart == nil {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("組み立てパーツにマザボがありません")
	}
	if mbPart.IsSpare {
		mb = A_mb_G
		cooler = A_cooler_B
	} else {
		switch mbPart.Subcategory {
		case game.MB_Black:
			mb = A_mb_Black
			cooler = A_cooler_B
		case game.MB_White:
			mb = A_mb_W
			cooler = A_cooler_W
		case game.MB_Blue:
			mb = A_mb_Blue
			cooler = A_cooler_B
		case game.MB_Green:
			mb = A_mb_G
			cooler = A_cooler_B
		default:
			return 0, 0, 0, 0, 0, 0, fmt.Errorf("マザボに対して無効なSubcategoryです")
		}
	}

	casePart := findAssemblePart(pts, uint(game.Case))
	if casePart == nil {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("組み立てパーツにケースがありません")
	}
	if casePart.IsSpare {
		caseback = A_case_cardboard
		caselid = A_case_cardboard_lid
		gaming = AssemblePartsImage(A_noAdditionalParts)
		return caseback, mb, cooler, gpu, caselid, gaming, nil
	}
	_, hasLED := pts.Contain(game.LED)
	switch casePart.Subcategory {
	case game.Case_Black_SideFrontGlass:
		caseback = A_case_B_SF
		caselid = A_case_SF_lid
		if hasLED {
			caselid = A_caseSF_lid_gaming
			gaming = A_gaming_B_SF
		}
	case game.Case_Black_SideGlass:
		caseback = A_case_B_S
		caselid = A_case_S_lid
		if hasLED {
			caselid = A_case_S_lid_gaming
			gaming = A_gaming_B_S
		}
	case game.Case_White_SideFrontGlass:
		caseback = A_case_W_SF
		caselid = A_case_SF_lid
		if hasLED {
			caselid = A_caseSF_lid_gaming
			gaming = A_gaming_W_SF
		}
	case game.Case_White_SideGlass:
		caseback = A_case_W_S
		caselid = A_case_S_lid
		if hasLED {
			caselid = A_case_S_lid_gaming
			gaming = A_gaming_W_S
		}
	default:
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("ケースに対して無効なSubcategoryです")
	}
	return caseback, mb, cooler, gpu, caselid, gaming, nil
}

func CheckAssemblePartsInfo(pts game.Parts) error {
	gpuPart := findAssemblePart(pts, uint(game.GPU))
	if gpuPart == nil {
		return fmt.Errorf("組み立てパーツにGPUがありません")
	}
	if !gpuPart.IsSpare && gpuPart.Subcategory != game.GPU_Black && gpuPart.Subcategory != game.GPU_White {
		return fmt.Errorf("GPUに対して無効なSubcategoryです")
	}

	mbPart := findAssemblePart(pts, uint(game.MB))
	if mbPart == nil {
		return fmt.Errorf("組み立てパーツにマザボがありません")
	}
	if !mbPart.IsSpare {
		switch mbPart.Subcategory {
		case game.MB_Black:
		case game.MB_White:
		case game.MB_Blue:
		case game.MB_Green:
		default:
			return fmt.Errorf("マザボに対して無効なSubcategoryです")
		}
	}

	casePart := findAssemblePart(pts, uint(game.Case))
	if casePart == nil {
		return fmt.Errorf("組み立てパーツにケースがありません")
	}
	if !casePart.IsSpare {
		switch casePart.Subcategory {
		case game.Case_Black_SideFrontGlass:
		case game.Case_Black_SideGlass:
		case game.Case_White_SideFrontGlass:
		case game.Case_White_SideGlass:
		default:
			return fmt.Errorf("ケースに対して無効なSubcategoryです")
		}
	}

	return nil
}
