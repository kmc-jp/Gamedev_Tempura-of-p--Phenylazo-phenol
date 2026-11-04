package serializetarget

import camera "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/Camera"

type CameraConfig struct {
	Layers map[string]int
}

// Utils.Factory[*scene.Camera] を生成して返す
func (tmd CameraConfig) Construct() (any, error) {
	return camera.NewCamera(tmd.Layers), nil
}
