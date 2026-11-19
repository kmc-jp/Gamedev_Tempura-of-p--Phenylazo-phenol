package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var (
	CaseBlackSideGlass = register(func() *game.Part {
		return &game.Part{
			Name: "黒のPCケース", Kind: game.Case, Subcategory: game.Case_Black_SideGlass,
			Weight: 2000, Val: 2000, Rank: game.B, Performance: 1000,
			Detail: "側面がガラスの黒PCケース", ImgName: "Parts/みんげー_PCケース_黒_蓋閉じてる",
			ImgScale: CaseImageScale, ImgOffset: &NoOffset,
		}
	})
	CaseBlackSideFrontGlass = register(func() *game.Part {
		return &game.Part{
			Name: "黒のピラーレスケース", Kind: game.Case, Subcategory: game.Case_Black_SideFrontGlass,
			Weight: 2000, Val: 2000, Rank: game.B, Performance: 1000,
			Detail: "側面と前面がガラスの黒PCケース", ImgName: "Parts/みんげー_PCケース_黒_蓋開いてる",
			ImgScale: CaseImageScale, ImgOffset: &NoOffset,
		}
	})
	CaseWhiteSideGlass = register(func() *game.Part {
		return &game.Part{
			Name: "白のPCケース", Kind: game.Case, Subcategory: game.Case_White_SideGlass,
			Weight: 2000, Val: 2000, Rank: game.B, Performance: 1000,
			Detail: "側面がガラスの白PCケース", ImgName: "Parts/みんげー_PCケース_白_蓋閉じてる",
			ImgScale: CaseImageScale, ImgOffset: &NoOffset,
		}
	})
	CaseWhiteSideFrontGlass = register(func() *game.Part {
		return &game.Part{
			Name: "白のピラーレスケース", Kind: game.Case, Subcategory: game.Case_White_SideFrontGlass,
			Weight: 2000, Val: 2000, Rank: game.B, Performance: 1000,
			Detail: "側面と前面がガラスの白PCケース", ImgName: "Parts/みんげー_PCケース_白_蓋開いてる",
			ImgScale: CaseImageScale, ImgOffset: &NoOffset,
		}
	})
)

func init() {
	builders := []func() *game.Part{
		CaseBlackSideGlass,
		CaseBlackSideFrontGlass,
		CaseWhiteSideGlass,
		CaseWhiteSideFrontGlass,
	}

	for _, builder := range builders {
		for range 11 {
			register(builder)
		}
	}
}
