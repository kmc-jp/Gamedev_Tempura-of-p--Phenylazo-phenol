package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	serializetarget "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize/SerializeTarget"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"reflect"
)

// カメラ設定
func MakeWorldDrawerFactory(assetsManager *assets.AssetsManager) (Utils.Factory[scene.WorldDrawer], error) {
	cameraConfigData, err := serialize.NewSerializeData("CameraConfig", reflect.TypeFor[serializetarget.CameraConfig](), "CameraConfig.toml", assetsManager.Serialize)
	if err != nil {
		return nil, fmt.Errorf("serialize.NewSerializeData(\"CameraConfig\", reflect.TypeFor[serializetarget.CameraConfig](), \"CameraConfig.toml\", assetsManager.Serialize) でエラーが発生しました。\n%w", err)
	}
	cameraConfigAny, err := assetsManager.Serialize.Load(*cameraConfigData)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.Load(*cameraConfigData) でエラーが発生しました。\n%w", err)
	}

	cameraConfig := cameraConfigAny.(serializetarget.CameraConfig)
	cameraFactory, err := cameraConfig.Construct()
	if err != nil {
		return nil, fmt.Errorf("cameraConfig.Construct() でエラーが発生しました。\n%w", err)
	}

	DrawerFactory := Utils.CastFactory[*scene.Camera, scene.WorldDrawer](cameraFactory.(Utils.Factory[*scene.Camera]))
	return DrawerFactory, nil
}

func MakeUIDrawerFactory() (Utils.Factory[scene.UIDrawer], error) {
	DrawerFactory := Utils.CastFactory[*scene.UICamera, scene.UIDrawer](scene.NewUICamera())
	return DrawerFactory, nil
}
