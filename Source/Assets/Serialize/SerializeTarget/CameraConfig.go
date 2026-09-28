package serializetarget

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
)

type CameraConfig struct {
	Layers map[string]int
}

// Utils.Factory[*scene.Camera] を生成して返す
func (tmd CameraConfig) Construct() (any, error) {
	return scene.NewCamera(tmd.Layers), nil
}
