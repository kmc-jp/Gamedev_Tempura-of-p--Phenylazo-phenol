package parts

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"strconv"
)

type gpuModelSpec struct {
	name        string
	isWhite     bool
	rarityLevel int
	value       int
	performance int
	detail      string
}

func registerGPU(spec gpuModelSpec) {
	imgName := "Parts/みんげー_グラボ"
	imgScale := GPUBlackImageScale
	imgOffset := &NoOffset
	subcategory := game.GPU_Black
	if spec.isWhite {
		subcategory = game.GPU_White
		imgName = "Parts/みんげー_グラボ_白"
		imgScale = GPUWhiteImageScale
		imgOffset = &GPUWhiteOffset
	}
	rarity := game.S1
	if spec.rarityLevel >= 2 {
		rarity = game.S2
	}
	if spec.rarityLevel >= 3 {
		rarity = game.S3
	}

	register(func() *game.Part {
		return &game.Part{
			Name:        spec.name,
			Kind:        game.GPU,
			Subcategory: subcategory,
			Rarity:      rarity,
			Weight:      1000,
			Val:         spec.value,
			Performance: spec.performance,
			Detail:      spec.detail,
			ImgName:     imgName,
			ImgScale:    imgScale,
			ImgOffset:   imgOffset,
		}
	})
}

func init() {
	manufacturers := []string{"ASAS", "astone", "MSL", "MegaByte", "ruby"}

	rpxGenerations := []struct {
		generation int
		seriesBase int
	}{
		{generation: 20, seriesBase: 2000},
		{generation: 30, seriesBase: 3000},
		{generation: 40, seriesBase: 4000},
		{generation: 50, seriesBase: 5000},
		{generation: 60, seriesBase: 6000},
	}
	rpxGrades := []int{50, 60, 70, 80, 90}

	for generationIndex, generation := range rpxGenerations {
		for gradeIndex, grade := range rpxGrades {
			manufacturer := manufacturers[(generationIndex*len(rpxGrades)+gradeIndex)%len(manufacturers)]
			isWhite := false
			if (generationIndex+gradeIndex)%2 == 1 {
				isWhite = true
			}
			rarityLevel := 1
			if grade >= 70 {
				rarityLevel = 2
			}
			if grade >= 90 {
				rarityLevel = 3
			}
			performance := 2500 + generationIndex*1800 + gradeIndex*900
			value := 12000 + generationIndex*12000 + gradeIndex*7000

			registerGPU(gpuModelSpec{
				name:        manufacturer + " RPX " + strconv.Itoa(generation.seriesBase+grade),
				isWhite:     isWhite,
				rarityLevel: rarityLevel,
				value:       value,
				performance: performance,
				detail:      manufacturer + "製のグラフィックスボード。",
			})
		}
	}

	rideonGenerations := []struct {
		generation int
		seriesBase int
	}{
		{generation: 66, seriesBase: 6600},
		{generation: 67, seriesBase: 6700},
		{generation: 76, seriesBase: 7600},
		{generation: 77, seriesBase: 7700},
		{generation: 78, seriesBase: 7800},
	}
	rideonGrades := []int{00, 00, 50, 70, 90}

	for generationIndex, generation := range rideonGenerations {
		for gradeIndex, gradeOffset := range rideonGrades {
			grade := gradeOffset
			modelNumber := generation.seriesBase
			if grade > 0 {
				modelNumber += grade
			} else if gradeIndex == 1 {
				modelNumber++
			}

			manufacturer := manufacturers[(generationIndex*len(rideonGrades)+gradeIndex+2)%len(manufacturers)]
			isWhite := false
			if (generationIndex+gradeIndex)%2 == 0 {
				isWhite = true
			}
			rarityLevel := 1
			if gradeIndex >= 2 {
				rarityLevel = 2
			}
			if gradeIndex >= 4 {
				rarityLevel = 3
			}
			performance := 2300 + generationIndex*1700 + gradeIndex*850
			value := 10000 + generationIndex*11500 + gradeIndex*6500

			registerGPU(gpuModelSpec{
				name:        manufacturer + " rideon RX " + strconv.Itoa(modelNumber),
				isWhite:     isWhite,
				rarityLevel: rarityLevel,
				value:       value,
				performance: performance,
				detail:      manufacturer + "製のグラフィックスボード。",
			})
		}
	}
}
