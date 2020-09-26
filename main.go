package main

import (
	"github.com/PhilippSeitz/beautiful-gocov/report/html"
	"github.com/PhilippSeitz/beautiful-gocov/tree"
	"os"
)

func main() {
	var file = os.Getenv("FILE")
	var dir = os.Getenv("DIR")
	var mod = os.Getenv("MOD")
	t, err := tree.BuildTree(file, dir, mod)
	if err != nil {
		panic(err)
	}

	root := t.Folders["stackit.de"]
	root.Print(0)

	err = html.HTML(root)
	if err != nil {
		panic(err)
	}

}
