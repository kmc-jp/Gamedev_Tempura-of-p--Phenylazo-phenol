package parts

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"strconv"
)

type cpuModelSpec struct {
	name        string
	grade       string
	performance int
	value       int
	detail      string
}

func registerCPU(spec cpuModelSpec) {
	rarity := game.S1
	if spec.grade == "i5" {
		rarity = game.S2
	}
	if spec.grade == "i7" {
		rarity = game.S2
	}
	if spec.grade == "i9" {
		rarity = game.S3
	}
	register(func() *game.Part {
		return &game.Part{
			Name:        spec.name,
			Kind:        game.CPU,
			Rarity:      rarity,
			Weight:      100,
			Val:         spec.value,
			Performance: spec.performance,
			Detail:      spec.detail,
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})
}

func init() {
	nitelGenerations := []int{2100, 3100, 4100, 5100, 6100}
	nitelBasePerformance := map[string]int{
		"i3": 2200,
		"i5": 4000,
		"i7": 5800,
		"i9": 8000,
	}
	nitelBaseValue := map[string]int{
		"i3": 5000,
		"i5": 9000,
		"i7": 14000,
		"i9": 22000,
	}

	for generationIndex, generation := range nitelGenerations {
		specs := []struct {
			grade       string
			suffix      string
			performance int
			value       int
		}{
			{grade: "i3", suffix: "", performance: 0, value: 0},
			{grade: "i5", suffix: "", performance: 0, value: 0},
			{grade: "i5", suffix: "T", performance: 100, value: 1500},
			{grade: "i7", suffix: "", performance: 0, value: 0},
			{grade: "i7", suffix: "T", performance: 150, value: 2000},
			{grade: "i9", suffix: "S", performance: 350, value: 3500},
		}

		for _, spec := range specs {
			name := "Nitel Pure " + spec.grade + " " + strconv.Itoa(generation) + spec.suffix
			registerCPU(cpuModelSpec{
				name:        name,
				grade:       spec.grade,
				performance: nitelBasePerformance[spec.grade] + generationIndex*700 + spec.performance,
				value:       nitelBaseValue[spec.grade] + generationIndex*2500 + spec.value,
				detail:      "Nitel社の" + spec.grade + "世代CPU。",
			})
		}
	}

	ryzinGenerations := []int{2300, 3300, 4300, 5300, 6300}
	ryzinBasePerformance := map[string]int{
		"i3": 2600,
		"i5": 4400,
		"i7": 6200,
		"i9": 8500,
	}
	ryzinBaseValue := map[string]int{
		"i3": 6000,
		"i5": 10000,
		"i7": 15000,
		"i9": 23000,
	}

	for generationIndex, generation := range ryzinGenerations {
		specs := []struct {
			grade       string
			suffix      string
			performance int
			value       int
		}{
			{grade: "i3", suffix: "", performance: 0, value: 0},
			{grade: "i5", suffix: "U", performance: 500, value: 1800},
			{grade: "i7", suffix: "U", performance: 550, value: 2200},
			{grade: "i9", suffix: "X", performance: 900, value: 4000},
		}

		for _, spec := range specs {
			name := "Ryzin" + spec.grade[1:] + " " + strconv.Itoa(generation) + spec.suffix
			registerCPU(cpuModelSpec{
				name:        name,
				grade:       spec.grade,
				performance: ryzinBasePerformance[spec.grade] + generationIndex*700 + spec.performance,
				value:       ryzinBaseValue[spec.grade] + generationIndex*2800 + spec.value,
				detail:      "Ryzin社の" + spec.grade + "世代CPU。",
			})
		}
	}
}
