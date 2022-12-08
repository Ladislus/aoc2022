package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type File struct {
	Name string
	Size uint
}

func make_file(name string, size uint) File {
	return File{Name: name, Size: size}
}

type Folder struct {
	Name    string
	Folders []Folder
	Files   []File
	Parent  *Folder
}

func make_folder(name string) *Folder {
	return &Folder{Name: name, Folders: make([]Folder, 0), Files: make([]File, 0), Parent: nil}
}

func (folder *Folder) add_folder(subfolder *Folder) {
	folder.Folders = append(folder.Folders, *subfolder)
	(*subfolder).Parent = folder
}

func (folder *Folder) try_add_folder(subfolder_name string) Folder {
	for _, current_subfolder := range folder.Folders {
		if current_subfolder.Name == subfolder_name {
			return current_subfolder
		}
	}

	new_folder := make_folder(subfolder_name)
	folder.add_folder(new_folder)

	return *new_folder
}

func (folder *Folder) try_add_file(name string, size uint) {
	for _, current_file := range folder.Files {
		if current_file.Name == name {
			if current_file.Size != size {
				panic("Found same file with mismatching size")
			} else {
				return
			}
		}
	}

	folder.Files = append(folder.Files, make_file(name, size))
}

func (folder Folder) size() uint {
	var size uint = 0

	for _, current_folder := range folder.Folders {
		size += current_folder.size()
	}
	for _, current_file := range folder.Files {
		size += current_file.Size
	}

	return size
}

// TODO: Remove
func (folder Folder) debug() {

	foln := "[ "
	fols := len(folder.Folders) - 1
	filn := "[ "
	fils := len(folder.Files) - 1

	for index, curr := range folder.Folders {
		foln += curr.Name
		if index < fols {
			foln += ", "
		}
	}
	foln += " ]"

	for index, curr := range folder.Files {
		filn += curr.Name
		if index < fils {
			filn += ", "
		}
	}
	filn += " ]"

	var pn string
	if folder.Parent == nil {
		pn = "nil"
	} else {
		pn = folder.Parent.Name
	}

	fmt.Printf("Folder [ \"%s\", %v, %v, \"%s\" ]\n", folder.Name, foln, filn, pn)
}

func help(program_name string) {
	fmt.Fprintf(os.Stderr, "Usage: %s input_file_path\n", program_name)
	os.Exit(1)
}

func input_v1(filepath string, size_cap uint) {

	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	command_regex := regexp.MustCompile("^\\$ (.*)$")
	cd_regex := regexp.MustCompile("^cd ([\\w../]+)$")
	ls_regex := regexp.MustCompile("^ls$")
	file_regex := regexp.MustCompile("^(\\d+) ([\\w.]+)$")
	dir_regex := regexp.MustCompile("^dir (\\w+)$")

	root := *make_folder("/")
	current_folder := root
	isCurrentlyListing := false

	var lines []string = strings.Split(string(data), "\n")
	for index, line := range lines {

		if match := command_regex.FindStringSubmatch(line); match != nil {
			if len(match) != 2 {
				panic("Didn't capture encapsulated command")
			}

			current_folder.debug()
			fmt.Printf("Line %d: \"%s\" is a command\n", index, line)

			line = match[1]
			fmt.Printf("\t Extracted command: \"%s\"\n", line)

			if ls_regex.MatchString(line) {
				isCurrentlyListing = true
			} else if match := cd_regex.FindStringSubmatch(line); match != nil {
				isCurrentlyListing = false

				if len(match) != 2 {
					panic("Didn't properly capture folder name")
				}

				if match[1] == "/" {
					current_folder = root
					fmt.Printf("\t\tNew current folder: \"%s\"\n", current_folder.Name)
				} else if match[1] == ".." {
					if current_folder.Parent == nil {
						panic(fmt.Sprintf("Current folder \"%s\" without parent", current_folder.Name))
					} else {
						current_folder = *current_folder.Parent
						fmt.Printf("\t\tNew current folder: \"%s\"\n", current_folder.Name)
					}
				} else {
					current_folder = current_folder.try_add_folder(match[1])

					fmt.Printf("\t\tNew current folder: \"%s\"\n", current_folder.Name)
				}

			} else {
				panic(fmt.Sprintf("Unknown command \"%s\"\n", line))
			}

		} else if match := file_regex.FindStringSubmatch(line); match != nil {

			if !isCurrentlyListing {
				panic("Found file without ls")
			}

			if len(match) != 3 {
				panic("Didn't properly capture file size and name")
			}

			file_size, file_size_err := strconv.ParseUint(match[1], 10, 64)
			if file_size_err != nil {
				panic(fmt.Sprintf("Unable to parse file size \"%s\"", match[1]))
			}

			fmt.Printf("\t\tFILE(\"%s\", %d)\n", match[2], file_size)

			current_folder.try_add_file(match[2], uint(file_size))

		} else if match := dir_regex.FindStringSubmatch(line); match != nil {

			if !isCurrentlyListing {
				panic("Found folder without ls")
			}

			if len(match) != 2 {
				panic("Didn't properly capture directory name")
			}

			dir_name := match[1]

			fmt.Printf("\t\tFOLDER(\"%s\")\n", dir_name)

			current_folder.try_add_folder(dir_name)
		} else {
			panic(fmt.Sprintf("Unknown line type \"%s\"", line))
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		help(os.Args[0])
	}

	input_v1(os.Args[1], 100_000)
}
