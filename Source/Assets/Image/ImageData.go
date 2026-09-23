package imageassets

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type ImageData struct {
	Name string
	// /Assets/Image からの相対パス
	PngPath  string
	JsonPath string
}

func NewImageData(name string, pngpath string, jsonpath string, ImageDirectory fs.FS) (*ImageData, error) {
	if len(strings.TrimSpace(pngpath)) == 0 {
		return nil, fmt.Errorf("pngpath が空です。")
	}
	if len(strings.TrimSpace(jsonpath)) == 0 {
		return nil, fmt.Errorf("jsonpath が空です。")
	}

	pngcleaned := filepath.Clean(pngpath)
	if pngcleaned == "." {
		return nil, fmt.Errorf("pngpath が不正です。")
	}
	jsoncleaned := filepath.Clean(jsonpath)
	if jsoncleaned == "." {
		return nil, fmt.Errorf("jsonpath が不正です。")
	}
	pngcleaned = strings.TrimPrefix(pngcleaned, "/")
	jsoncleaned = strings.TrimPrefix(jsoncleaned, "/")

	pnginfo, pngerr := fs.Stat(ImageDirectory, pngcleaned)
	if pngerr != nil {
		return nil, fmt.Errorf("pngpath が使用できません。\n%w", pngerr)
	}
	if pnginfo.IsDir() {
		return nil, fmt.Errorf("pngpath はディレクトリです。")
	}
	jsoninfo, jsonerr := fs.Stat(ImageDirectory, jsoncleaned)
	if jsonerr != nil {
		return nil, fmt.Errorf("jsonpath が使用できません。\n%w", jsonerr)
	}
	if jsoninfo.IsDir() {
		return nil, fmt.Errorf("jsonpath はディレクトリです。")
	}

	id := ImageData{
		Name:     name,
		PngPath:  pngcleaned,
		JsonPath: jsoncleaned,
	}
	return &id, nil
}
