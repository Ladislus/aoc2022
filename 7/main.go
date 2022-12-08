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

func newFile(name string, size uint) *File {
	file := new(File)
	file.Name = name
	file.Size = size

	return file
}

type Folder struct {
	Name    string
	Folders []*Folder
	Files   []*File
	Parent  *Folder
}

func newFolder(name string) *Folder {
	folder := new(Folder)
	folder.Name = name
	folder.Folders = make([]*Folder, 0)
	folder.Files = make([]*File, 0)
	folder.Parent = nil

	return folder
}

func (folder *Folder) addFolder(subfolder *Folder) {
	folder.Folders = append(folder.Folders, subfolder)
	subfolder.Parent = folder
}

func (folder *Folder) tryAddFolder(subfolderName string) *Folder {

	for _, currentSubfolder := range folder.Folders {
		if currentSubfolder.Name == subfolderName {
			return currentSubfolder
		}
	}

	newSubfolder := newFolder(subfolderName)
	folder.addFolder(newSubfolder)

	return newSubfolder
}

func (folder *Folder) tryAddFile(name string, size uint) {
	for _, currentFile := range folder.Files {
		if currentFile.Name == name {
			if currentFile.Size != size {
				panic("Found same file with mismatching size")
			} else {
				return
			}
		}
	}

	folder.Files = append(folder.Files, newFile(name, size))
}

func (folder *Folder) size() uint {
	var size uint = 0

	for _, currentSubfolder := range folder.Folders {
		size += currentSubfolder.size()
	}
	for _, currentFile := range folder.Files {
		size += currentFile.Size
	}

	return size
}

func (folder *Folder) compute(maxSize uint) uint {
	var sum uint = 0

	currentFolderSize := folder.size()
	if currentFolderSize < maxSize {
		sum += currentFolderSize
	}

	for _, currentSubfolder := range folder.Folders {
		sum += currentSubfolder.compute(maxSize)
	}

	return sum
}

// TODO: Remove
func (folder *Folder) debug() {

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

func help(programName string) {
	fmt.Fprintf(os.Stderr, "Usage: %s input_file_path\n", programName)
	os.Exit(1)
}

func input_v1(filepath string, maxSize uint) {

	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	commandRegex := regexp.MustCompile("^\\$ (.*)$")
	cdRegex := regexp.MustCompile("^cd ([\\w../]+)$")
	lsRegex := regexp.MustCompile("^ls$")
	fileRegex := regexp.MustCompile("^(\\d+) ([\\w.]+)$")
	folderRegex := regexp.MustCompile("^dir (\\w+)$")

	root := newFolder("/")
	currentFolder := root
	isCurrentlyListing := false

	var lines []string = strings.Split(string(data), "\n")
	for index, line := range lines {

		if match := commandRegex.FindStringSubmatch(line); match != nil {
			if len(match) != 2 {
				panic("Didn't capture encapsulated command")
			}

			currentFolder.debug()
			fmt.Printf("Line %d: \"%s\" is a command\n", index, line)

			line = match[1]
			fmt.Printf("\t Extracted command: \"%s\"\n", line)

			if lsRegex.MatchString(line) {
				isCurrentlyListing = true
			} else if match := cdRegex.FindStringSubmatch(line); match != nil {
				isCurrentlyListing = false

				if len(match) != 2 {
					panic("Didn't properly capture folder name")
				}

				if match[1] == "/" {
					currentFolder = root
					fmt.Printf("\t\tNew current folder: \"%s\"\n", currentFolder.Name)
				} else if match[1] == ".." {
					if currentFolder.Parent == nil {
						panic(fmt.Sprintf("Current folder \"%s\" without parent", currentFolder.Name))
					} else {
						currentFolder = currentFolder.Parent
						fmt.Printf("\t\tNew current folder: \"%s\"\n", currentFolder.Name)
					}
				} else {
					currentFolder = currentFolder.tryAddFolder(match[1])

					fmt.Printf("\t\tNew current folder: \"%s\"\n", currentFolder.Name)
				}

			} else {
				panic(fmt.Sprintf("Unknown command \"%s\"\n", line))
			}

		} else if match := fileRegex.FindStringSubmatch(line); match != nil {

			if !isCurrentlyListing {
				panic("Found file without ls")
			}

			if len(match) != 3 {
				panic("Didn't properly capture file size and name")
			}

			fileSize, fileSizeErr := strconv.ParseUint(match[1], 10, 64)
			if fileSizeErr != nil {
				panic(fmt.Sprintf("Unable to parse file size \"%s\"", match[1]))
			}

			fmt.Printf("\t\tFILE(\"%s\", %d)\n", match[2], fileSize)

			currentFolder.tryAddFile(match[2], uint(fileSize))

		} else if match := folderRegex.FindStringSubmatch(line); match != nil {

			if !isCurrentlyListing {
				panic("Found folder without ls")
			}

			if len(match) != 2 {
				panic("Didn't properly capture directory name")
			}

			dir_name := match[1]

			fmt.Printf("\t\tFOLDER(\"%s\")\n", dir_name)

			currentFolder.tryAddFolder(dir_name)
		} else {
			panic(fmt.Sprintf("Unknown line type \"%s\"", line))
		}
	}

	fmt.Printf("Input v1: %d\n", root.compute(maxSize))
}

func main() {
	if len(os.Args) != 2 {
		help(os.Args[0])
	}

	input_v1(os.Args[1], 100_000)
}
