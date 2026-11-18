package parts

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
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
	A_mb_Blue  //まだassetない
	A_mb_G     //まだassetない

	A_cooler_B //みんげー_CPUクーラー_リザルト用_黒
	A_cooler_W //みんげー_CPUクーラー_リザルト用_白

	A_gpu_B //みんげー_グラボ_リザルト用_黒
	A_gpu_W //みんげー_グラボ_リザルト用_白

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
) *ebiten.Image {
	pc := ebiten.NewImage(72, 64)
	pc.DrawImage(AssembleImages[caseBack], nil)
	pc.DrawImage(AssembleImages[mb], nil)
	pc.DrawImage(AssembleImages[cooler], nil)
	pc.DrawImage(AssembleImages[gpu], nil)
	pc.DrawImage(AssembleImages[caseLid], nil)
	if gaming == A_noAdditionalParts {
		return pc
	} else {
		pc.DrawImage(AssembleImages[gaming], nil)
		return pc
	}
}
