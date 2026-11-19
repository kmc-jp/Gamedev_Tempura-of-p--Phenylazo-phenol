package parts

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"strconv"
)

type mbModelSpec struct {
	name        string
	color       string
	performance int
	value       int
	rarityLevel int
}

func registerMB(spec mbModelSpec) {
	subcategory := game.MB_Black
	imgName := "Parts/みんげー_マザボ"
	switch spec.color {
	case "White":
		subcategory = game.MB_White
		imgName = "Parts/みんげー_マザボ_白"
	case "Blue":
		subcategory = game.MB_Blue
		imgName = "Parts/みんげー_マザボ_青"
	case "Green":
		subcategory = game.MB_Green
		imgName = "Parts/みんげー_マザボ_緑"
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
			Kind:        game.MB,
			Subcategory: subcategory,
			Rarity:      rarity,
			Weight:      500,
			Val:         spec.value,
			Rank:        game.C,
			Performance: spec.performance,
			Detail:      spec.color + "のマザーボード。",
			ImgName:     imgName,
			ImgScale:    MBImageScale,
			ImgOffset:   &NoOffset,
		}
	})
}

func init() {
	numbers := []int{350, 370, 390, 410, 430, 450, 470, 490, 510, 530, 550, 570}
	series := []string{"H", "B"}
	lightManufacturers := []string{"ASAS", "astone", "MSL", "MegaByte"}
	accentManufacturers := []string{"MegaByte", "ASAS"}

	for seriesIndex, seriesName := range series {
		for numberIndex, number := range numbers {
			generationBonus := numberIndex * 350
			seriesBonus := 0
			if seriesName == "H" {
				seriesBonus = 1000
			}

			rarityLevel := 1
			if numberIndex >= 4 {
				rarityLevel = 2
			}
			if numberIndex >= 9 {
				rarityLevel = 3
			}

			performance := 3000 + seriesBonus + generationBonus
			value := 3000 + seriesIndex*1000 + generationBonus*8

			blackManufacturer := lightManufacturers[numberIndex%len(lightManufacturers)]
			registerMB(mbModelSpec{
				name:        blackManufacturer + " " + seriesName + strconv.Itoa(number),
				color:       "Black",
				performance: performance,
				value:       value,
				rarityLevel: rarityLevel,
			})

			whiteManufacturer := lightManufacturers[(numberIndex+2)%len(lightManufacturers)]
			registerMB(mbModelSpec{
				name:        whiteManufacturer + " " + seriesName + strconv.Itoa(number),
				color:       "White",
				performance: performance + 100,
				value:       value + 500,
				rarityLevel: rarityLevel,
			})

			blueManufacturer := accentManufacturers[numberIndex%len(accentManufacturers)]
			registerMB(mbModelSpec{
				name:        blueManufacturer + " " + seriesName + strconv.Itoa(number),
				color:       "Blue",
				performance: performance + 50,
				value:       value + 300,
				rarityLevel: rarityLevel,
			})

			greenManufacturer := accentManufacturers[(numberIndex+1)%len(accentManufacturers)]
			registerMB(mbModelSpec{
				name:        greenManufacturer + " " + seriesName + strconv.Itoa(number),
				color:       "Green",
				performance: performance + 50,
				value:       value + 300,
				rarityLevel: rarityLevel,
			})
		}
	}
}
