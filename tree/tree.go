package tree

import (
	"fmt"
	"golang.org/x/tools/cover"
	"strings"
)

type Folder struct {
	Name    string
	Files   map[string]File
	Folders map[string]*Folder
}

type File struct {
	Covered int64
	Total   int64
	Blocks  []cover.ProfileBlock
	Path    string
}

func BuildTree(source, path, mod string) (*Folder, error) {
	root := &Folder{
		Files:   make(map[string]File, 0),
		Folders: make(map[string]*Folder, 0),
	}

	profiles, err := cover.ParseProfiles(source)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		fileOrFolder := strings.Split(profile.FileName, "/")
		folder := root
		for i, subFolder := range fileOrFolder {
			if len(fileOrFolder)-1 == i {
				var total, covered int64
				for _, b := range profile.Blocks {
					total += int64(b.NumStmt)
					if b.Count > 0 {
						covered += int64(b.NumStmt)
					}
				}

				folder.Files[subFolder] = File{
					Covered: covered,
					Total:   total,
					Blocks:  profile.Blocks,
					Path:    origin(profile.FileName, path, mod),
				}
				break
			}
			if val, ok := folder.Folders[subFolder]; ok {
				folder = val
			} else {
				t := &Folder{
					Name:    subFolder,
					Files:   make(map[string]File, 0),
					Folders: make(map[string]*Folder, 0),
				}
				folder.Folders[subFolder] = t
				folder = t
			}
		}
	}

	return root, nil
}

func origin(file, origin, mod string) string {
	return strings.Replace(file, mod, origin, -1)
}

func (f *Folder) Print(depth int) {
	margin := strings.Repeat("  ", depth)
	fp := fmt.Sprintf("%f", f.Coverage())
	fmt.Println(margin + "+--" + f.Name + ": " + fp)
	for key, val := range f.Files {
		sp := fmt.Sprintf("%f", val.Coverage())
		fmt.Println(margin + "    " + key + ": " + sp + " --> " + val.Path)
	}
	for _, val := range f.Folders {
		val.Print(depth + 1)
	}
}

func (f *Folder) Stats() (covered, total int64) {
	for _, val := range f.Files {
		total = total + val.Total
		covered = covered + val.Covered
	}
	for _, val := range f.Folders {
		fCovered, fTotal := val.Stats()
		total = total + fTotal
		covered = covered + fCovered
	}
	return covered, total
}

func (f *File) Coverage() (coverage float64) {
	if f.Total != 0 {
		coverage = float64(f.Covered) / float64(f.Total)
	}
	return coverage
}

func (f *Folder) Coverage() (coverage float64) {
	covered, total := f.Stats()
	if total != 0 {
		coverage = float64(covered) / float64(total)
	}
	return coverage
}

type Handler func(f *Folder, parents []string)

func (f *Folder) Traverse(handler Handler) {
	f.traverse(handler, make([]string, 0))
}

func (f *Folder) traverse(handler Handler, parents []string) {
	handler(f, parents)
	for _, s := range f.Folders {
		s.traverse(handler, append(parents, f.Name))
	}
}
