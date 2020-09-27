package html

import (
	"bytes"
	"github.com/PhilippSeitz/beautiful-gocov/tree"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type reporter struct {
	listTemplate   *template.Template
	detailTemplate *template.Template
}

type coverageChild struct {
	Path     string
	Coverage float64
}

type baseData struct {
	Path     string
	Coverage float64
	Title    string
}

type data struct {
	baseData
	Folders []coverageChild
	Files   []coverageChild
}

type fileData struct {
	baseData
	Code template.HTML
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
	folders := make([]coverageChild, 0)
	for key, val := range f.Folders {
		folders = append(folders, coverageChild{
			Path:     key,
			Coverage: val.Coverage() * 100,
		})
	}

	files := make([]coverageChild, 0)
	for key, val := range f.Files {
		files = append(files, coverageChild{
			Path:     key,
			Coverage: val.Coverage() * 100,
		})
	}

	d := data{
		baseData: baseData{
			Path:     path.Join(base, strings.Join(parents, "/"), f.Name),
			Coverage: f.Coverage() * 100,
			Title:    f.Name,
		},
		Folders: folders,
		Files:   files,
	}

	os.MkdirAll(d.Path, os.ModePerm)
	file, _ := os.Create(path.Join(d.Path, "index.html"))
	err := r.listTemplate.Execute(file, d)
	if err != nil {
		panic(err)
	}

	for key, file := range f.Files {
		r.renderFile(file, key, append(parents, f.Name))
	}
}

func (r *reporter) renderFile(f tree.File, name string, parents []string) {
	p := path.Join(base, strings.Join(parents, "/"), name+".html")
	detailFile, _ := os.Create(p)
	//codeFile, err := os.Open(f.Path)
	//if err != nil {
	//	panic(err)
	//}
	//defer codeFile.Close()
	//
	//scanner := bufio.NewScanner(codeFile)
	//lines := make([]string, 0)
	//for scanner.Scan() {
	//	lines = append(lines, scanner.Text())
	//}
	//if err := scanner.Err(); err != nil {
	//	panic(err)
	//}
	code, err := ioutil.ReadFile(f.Path)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	// Highlight(buf, string(code), "go", "html", "monokai")
	formatter := html.New(
		html.WithLineNumbers(true),
		html.LineNumbersInTable(true),
		html.LinkableLineNumbers(true, "test"),
	)
	it, err := lexers.Get("go").Tokenise(nil, string(code))
	if err != nil {
		panic(err)
	}
	err = formatter.Format(buf, styles.Emacs, it)
	if err != nil {
		panic(err)
	}

	err = r.detailTemplate.Execute(detailFile, fileData{
		baseData: baseData{
			Path:     name,
			Title:    name,
			Coverage: f.Coverage() * 100},
		Code: template.HTML(buf.String())})
	if err != nil {
		panic(err)
	}
}
