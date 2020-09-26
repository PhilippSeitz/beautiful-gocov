package main

import (
	"fmt"
	"golang.org/x/tools/cover"
	"strings"
)

type folder struct {
	name       string
	subFiles   map[string][]cover.ProfileBlock
	subFolders map[string]*folder
}

func buildTree(source string) (*folder, error) {
	root := &folder{
		subFiles:   make(map[string][]cover.ProfileBlock, 0),
		subFolders: make(map[string]*folder, 0),
	}

	profiles, err := cover.ParseProfiles(source)
	if err != nil {
		return nil, err
	}
	fmt.Println(profiles)
	for _, profile := range profiles {
		folders := strings.Split(profile.FileName, "/")
		x := root
		for i, f := range folders {
			if len(folders)-1 == i {
				x.subFiles[f] = profile.Blocks
				break
			}
			if val, ok := x.subFolders[f]; ok {
				x = val
			} else {
				t := &folder{
					name:       f,
					subFiles:   make(map[string][]cover.ProfileBlock, 0),
					subFolders: make(map[string]*folder, 0),
				}
				x.subFolders[f] = t
				x = t
			}
		}
	}

	return root, nil
}

func (f *folder) print(depth int) {
	margin := strings.Repeat("  ", depth)
	fmt.Println(margin + "+--" + f.name)
	for val := range f.subFiles {
		fmt.Println(margin + "    " +val)
	}
	for _, val := range f.subFolders {
		val.print(depth + 1)
	}
}
