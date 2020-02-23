package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sandbox/indexer/document"
	"sandbox/indexer/index"
)

func main() {
	root := flag.String("root", ".", "source path")
	flag.Parse()
	ix := index.New()
	err := filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err.Error())
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		doc := document.New(path, file)
		if err := ix.Process(doc); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(ix)
}
