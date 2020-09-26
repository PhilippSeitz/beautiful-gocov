package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type field struct {
	numberOfStatements int
	count              int
	position           string
}

type folder struct {
	name       string
	subFiles   map[string][]*field
	subFolders map[string]*folder
}

func buildTree(source string) (*folder, error) {
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	root := &folder{
		subFiles:   make(map[string][]*field, 0),
		subFolders: make(map[string]*folder, 0),
	}
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		position := strings.Split(line[0], ":")
		file := position[0]
		numberOfStatements, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, err
		}
		count, _ := strconv.Atoi(line[2])

		folders := strings.Split(file, "/")
		x := root
		for i, f := range folders {
			if len(folders)-1 == i {
				x.subFiles[f] = append(x.subFiles[f], &field{
					numberOfStatements: numberOfStatements,
					count:              count,
					position:           position[1],
				})
				break
			}
			if val, ok := x.subFolders[f]; ok {
				x = val
			} else {
				t := &folder{
					name:       f,
					subFiles:   make(map[string][]*field, 0),
					subFolders: make(map[string]*folder, 0),
				}
				x.subFolders[f] = t
				x = t
			}
		}

	}

	return root, scanner.Err()
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
