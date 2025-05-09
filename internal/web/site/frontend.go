package site

import (
	"embed"
	"io/fs"
	"log"
	"os"
)

// FrontendFS contains the embedded frontend build
//
//go:embed frontend/**/*
//go:embed frontend/*
var FrontendFS embed.FS

func init() {
	// Print all embedded files but only in debug mode
	if os.Getenv("DEBUG") == "true" {
		fs.WalkDir(FrontendFS, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				log.Printf("Embedded file: %s", path)
			}
			return nil
		})
	}
}

func GetEmbedFrontend() fs.FS {
	frontendFS, err := fs.Sub(FrontendFS, "frontend")
	if err != nil {
		panic(err)
	}
	return frontendFS
}
