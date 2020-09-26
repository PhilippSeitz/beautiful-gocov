package html

import (
	"github.com/PhilippSeitz/beautiful-gocov/tree"
	"html/template"
	"os"
	"path"
	"strings"
)

type reporter struct {
	listTemplate   *template.Template
	detailTemplate *template.Template
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

func HTML(f *tree.Folder) error {
	listTemplate, err := template.ParseFiles("templates/index.html", "templates/list.html")
	if err != nil {
		return err
	}
	detailTemplate, err := template.ParseFiles("templates/index.html", "templates/detail.html")
	if err != nil {
		return err
	}

	r := reporter{
		listTemplate,
		detailTemplate,
	}

	f.Traverse(r.renderFolder)

	return nil
}

func (r *reporter) renderFolder(f *tree.Folder, parents []string) {
	folders := make([]coverage, 0)
	for key, val := range f.Folders {
		folders = append(folders, coverage{
			Path:     key,
			Coverage: val.Coverage() * 100,
		})
	}

	files := make([]coverage, 0)
	for key, val := range f.Files {
		files = append(files, coverage{
			Path:     key,
			Coverage: val.Coverage() * 100,
		})
	}

	d := data{
		Title:    f.Name,
		Path:     path.Join(base, strings.Join(parents, "/"), f.Name),
		Coverage: f.Coverage() * 100,
		Folders:  folders,
		Files:    files,
	}

	os.MkdirAll(d.Path, os.ModePerm)
	file, _ := os.Create(path.Join(d.Path, "index.html"))
	err := r.listTemplate.Execute(file, d)
	if err != nil {
		panic(err)
	}

	for key, file := range f.Files {
		p := path.Join(d.Path, key+".html")
		detailFile, _ := os.Create(p)
		r.detailTemplate.Execute(detailFile, fileData{Path: key, Title: key, Coverage: file.Coverage() * 100})
	}
	if err != nil {
		panic(err)
	}
}
