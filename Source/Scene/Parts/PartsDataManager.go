package parts

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"math/rand"
)

var Registry []func() *game.Part

const (
	randomRankAProbability = 5
	randomRankBProbability = 20
	randomRankCProbability = 60
	randomRankDProbability = 15
)

func register(builder func() *game.Part) func() *game.Part {
	Registry = append(Registry, builder)
	return builder
}

func LoadRegistryImages(imageManager *imageassets.ImageManager) error {
	for i, builder := range Registry {
		part := builder()
		if part == nil || part.ImgName == "" {
			continue
		}

		imageData, err := imageManager.NewImageData(
			part.ImgName,
			part.ImgName+".png",
			part.ImgName+".json",
		)
		if err != nil {
			return fmt.Errorf("パーツ画像データの作成に失敗しました。パーツ名: %s\n%w", part.Name, err)
		}

		_, image, err := imageManager.Load(*imageData)
		if err != nil {
			return fmt.Errorf("パーツ画像の読み込みに失敗しました。パーツ名: %s\n%w", part.Name, err)
		}

		originalBuilder := builder
		Registry[i] = func() *game.Part {
			loadedPart := originalBuilder()
			loadedPart.Img = image
			return loadedPart
		}
	}

	return nil
}

func RandomPartsFromAll(WithProbability bool) *game.Part {
	if len(Registry) == 0 {
		return nil
	}

	//重みづけ抽選をしない場合はここで抽選
	if WithProbability && !Utils.WithProbability(0.6) {
		return nil
	}

	p := Registry[rand.Intn(len(Registry))]()
	if p.Kind == game.Trash {
		return nil
	}
	if err := p.SetRandomRank(
		randomRankAProbability,
		randomRankBProbability,
		randomRankCProbability,
		randomRankDProbability,
	); err != nil {
		return nil
	}
	return p
}
