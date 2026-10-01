package server

import "embed"

//go:embed assets/dist/*
var assetFS embed.FS

func AssetReady() bool {
	_, err := assetFS.Open("assets/dist/index.html")
	return err == nil
}
