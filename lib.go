package staticassets

import (
	"embed"
	"io/fs"
)

//go:embed data/*.json
var dataFS embed.FS

//go:embed data/products/*.json
var productFS embed.FS

//go:embed hugo/*.json
var hugoFS embed.FS

func Data() fs.FS {
	fsys, err := fs.Sub(dataFS, "data")
	if err != nil {
		panic(err)
	}
	return fsys
}

func Product() fs.FS {
	fsys, err := fs.Sub(productFS, "data/products")
	if err != nil {
		panic(err)
	}
	return fsys
}

func Hugo() fs.FS {
	fsys, err := fs.Sub(hugoFS, "hugo")
	if err != nil {
		panic(err)
	}
	return fsys
}
