package main

import (
	"fmt"
	"golang.org/x/tools/cover"
	"strings"
)

type folder struct {
	name       string
	subFiles   map[string]file
	subFolders map[string]*folder
}

type file struct {
	covered    int64
	total      int64
	percentage float64
	blocks     []cover.ProfileBlock
}

func buildTree(source string) (*folder, error) {
	root := &folder{
		subFiles:   make(map[string]file, 0),
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
				var total, covered int64
				var percentage float64
				for _, b := range profile.Blocks {
					total += int64(b.NumStmt)
					if b.Count > 0 {
						covered += int64(b.NumStmt)
					}
				}

				if total != 0 {
					percentage = float64(covered) / float64(total)
				}

				x.subFiles[f] = file{
					covered:    covered,
					total:      total,
					percentage: percentage,
					blocks:     profile.Blocks,
				}
				break
			}
			if val, ok := x.subFolders[f]; ok {
				x = val
			} else {
				t := &folder{
					name:       f,
					subFiles:   make(map[string]file, 0),
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
	covered, total := f.stats()
	fp := fmt.Sprintf("%f",  float64(covered) / float64(total) *100)
	fmt.Println(margin + "+--" + f.name + ": " + fp)
	for key, val := range f.subFiles {
		sp := fmt.Sprintf("%f", val.percentage*100)
		fmt.Println(margin + "    " + key + ": " + sp)
	}
	for _, val := range f.subFolders {
		val.print(depth + 1)
	}
}

func (f *folder) stats() (int64, int64) {
	var covered, total int64
	for _, val := range f.subFiles {
		total = total + val.total
		covered = covered + val.covered
	}
	for _, val := range f.subFolders {
		fCovered, fTotal := val.stats()
		total = total + fTotal
		covered = covered + fCovered
	}
	return covered, total
}
