package scene

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

func PartExample() *game.Part {
	return &game.Part{
		Name:        "example",
		Kind:        game.Others,
		Weight:      500,
		Val:         1000,
		Rank:        game.B,
		Performance: 1000,
	}
}

func CPUExample() *game.Part {
	return &game.Part{
		Name:        "CPU_Example",
		Kind:        game.CPU,
		Weight:      100,
		Val:         10000,
		Rank:        game.A,
		Performance: 1000,
	}
}

func GPUExample() *game.Part {
	return &game.Part{
		Name:        "GPU_Example",
		Kind:        game.GPU,
		Weight:      1000,
		Val:         20000,
		Rank:        game.C,
		Performance: 1000,
	}
}

func RAMExample() *game.Part {
	return &game.Part{
		Name:        "RAM_Example",
		Kind:        game.RAM,
		Weight:      100,
		Val:         5000,
		Rank:        game.B,
		Performance: 1000,
	}
}
