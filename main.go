package main

import (
	"fmt"
	"os"
)

func main() {
	var file = os.Getenv("FILE")
	var _ = os.Getenv("DIR")
	tree, err := buildTree(file)
	if err != nil {
		panic(err)
	}

	root := tree.subFolders["stackit.de"].subFolders["permission"]
	fmt.Println(root)
	root.print(0)


}
