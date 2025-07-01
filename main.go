package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver" // 必须添加这个导入
)

//go:embed frontend/dist/*
var assets embed.FS

func main() {
	launcher := NewLauncher()

	err := wails.Run(&options.App{
		Title:  "雀魂工具箱",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{ // 使用正确的类型
			Assets: assets,
		},
		OnStartup: launcher.OnStartup,
		Bind:      []interface{}{launcher},
	})

	if err != nil {
		log.Fatal(err)
	}
}
