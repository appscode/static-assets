package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/appscode/static-assets/api"
)

func main() {
	dirname := "/home/tamal/go/src/github.com/appscode/static-assets/data/products"
	entries, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	for _, fi := range entries {
		if fi.IsDir() || filepath.Ext(fi.Name()) != ".json" {
			continue
		}

		filename := filepath.Join(dirname, fi.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		var prod api.Product
		err = json.Unmarshal(data, &prod)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(prod)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}
}
