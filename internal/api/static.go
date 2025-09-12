package api

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/*
var webFS embed.FS

func WebHandler() http.Handler {
	// 将嵌入的文件系统的 web 目录做为子目录
	subFS, err := fs.Sub(webFS, "web")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(subFS))
}
