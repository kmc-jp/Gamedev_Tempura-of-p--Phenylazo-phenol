package parts

import "github.com/hajimehoshi/ebiten/v2"

var AssembleImages [numberOfAssemblePartsAssets]*ebiten.Image

// 組み立て用画像一覧 (意図的に_を使用)
// A: Assemble であることを識別するための識別子 ,B: Black, W: White, S: Side glass, SF: Side and Front glass,
const (
	A_case_B_SF = iota
	A_case_B_S
	A_case_W_SF
	A_case_W_S
	A_case_cardboard

	A_case_B_SF_lid
	A_case_B_S_lid
	A_case_W_SF_lid
	A_case_W_S_lid
	A_case_cardboard_lid

	A_gaming_B_SF
	A_gaming_B_S
	A_gaming_W_SF
	A_gaming_W_S

	A_mb_Black
	A_mb_W
	A_mb_Blue //まだassetない
	A_mb_G    //まだassetない

	A_cooler_B
	A_cooler_W

	A_gpu_B
	A_gpu_W

	A_noAdditionalParts
)
const numberOfAssemblePartsAssets = A_noAdditionalParts

type AssemblePartsImage uint

func LoadAssemblePartsImages() {

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
