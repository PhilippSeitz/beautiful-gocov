package main

import (
	"fmt"
	"golang.org/x/tools/cover"
	"html/template"
	"os"
	"path"
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
	fp := fmt.Sprintf("%f", float64(covered)/float64(total)*100)
	fmt.Println(margin + "+--" + f.name + ": " + fp)
	for key, val := range f.subFiles {
		sp := fmt.Sprintf("%f", val.percentage*100)
		fmt.Println(margin + "    " + key + ": " + sp)
	}
	for _, val := range f.subFolders {
		val.print(depth + 1)
	}
}

func (f *folder) stats() (covered, total int64) {
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

type data struct {
	Title    string
	Path     string
	Coverage float64
	Folders  []coverage
	Files    []coverage
}

type coverage struct {
	Path     string
	Coverage float64
}

type fileData struct {
	Path     string
	Title    string
	Coverage float64
}

const base = "out"

func (f *folder) html(p string, list *template.Template, detail *template.Template) {
	covered, total := f.stats()
	folders := make([]coverage, 0)
	for key, val := range f.subFolders {
		covered, total := val.stats()
		folders = append(folders, coverage{
			Path:     key,
			Coverage: float64(covered) / float64(total) * 100,
		})
	}
	files := make([]coverage, 0)
	for key, val := range f.subFiles {
		files = append(files, coverage{
			Path:     key,
			Coverage: val.percentage * 100,
		})
	}
	d := data{
		Title:    f.name,
		Path:     path.Join(base, p),
		Coverage: float64(covered) / float64(total) * 100,
		Folders:  folders,
		Files:    files,
	}
	os.MkdirAll(d.Path, os.ModePerm)
	file, _ := os.Create(path.Join(d.Path, "index.html"))
	err := list.Execute(file, d)

	if err != nil {
		panic(err)
	}

	for key, val := range f.subFolders {
		val.html(path.Join(p, key), list, detail)
	}

	for key, file := range f.subFiles {
		p := path.Join(d.Path, key+".html")
		detailFile, _ := os.Create(p)
		detail.Execute(detailFile, fileData{Path: key, Title: key, Coverage: file.percentage * 100})
	}
}
