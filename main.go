package main

import (
	"html/template"
	"os"
)

func main() {
	var file = os.Getenv("FILE")
	var _ = os.Getenv("DIR")
	tree, err := buildTree(file)
	if err != nil {
		panic(err)
	}

	root := tree.subFolders["stackit.de"]

	t, err := template.ParseFiles("templates/index.html", "templates/list.html")
	root.print(0)
	root.html("", t)
}
