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

	l := template.Must(template.ParseFiles("templates/index.html", "templates/list.html"))
	d := template.Must(template.ParseFiles("templates/index.html", "templates/detail.html"))
	if err != nil {
		panic(err)
	}

	root.print(0)
	root.html("", l, d)
}
