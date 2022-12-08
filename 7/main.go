package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
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

func (folder *Folder) collectSizes(sizes *[]uint) *[]uint {
	*sizes = append(*sizes, folder.size())

	for _, currentSubfolder := range folder.Folders {
		currentSubfolder.collectSizes(sizes)
	}

	return sizes
}

func help(programName string) {
	fmt.Fprintf(os.Stderr, "Usage: %s input_file_path\n", programName)
	os.Exit(1)
}

func input(filepath string, maxSize uint) {

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
	for _, line := range lines {

		if match := commandRegex.FindStringSubmatch(line); match != nil {
			if len(match) != 2 {
				panic("Didn't capture encapsulated command")
			}

			line = match[1]

			if lsRegex.MatchString(line) {
				isCurrentlyListing = true
			} else if match := cdRegex.FindStringSubmatch(line); match != nil {
				isCurrentlyListing = false

				if len(match) != 2 {
					panic("Didn't properly capture folder name")
				}

				if match[1] == "/" {
					currentFolder = root
				} else if match[1] == ".." {
					if currentFolder.Parent == nil {
						panic(fmt.Sprintf("Current folder \"%s\" without parent", currentFolder.Name))
					} else {
						currentFolder = currentFolder.Parent
					}
				} else {
					currentFolder = currentFolder.tryAddFolder(match[1])
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

			currentFolder.tryAddFile(match[2], uint(fileSize))

		} else if match := folderRegex.FindStringSubmatch(line); match != nil {

			if !isCurrentlyListing {
				panic("Found folder without ls")
			}

			if len(match) != 2 {
				panic("Didn't properly capture directory name")
			}

			currentFolder.tryAddFolder(match[1])

		} else {
			panic(fmt.Sprintf("Unknown line type \"%s\"", line))
		}
	}

	fmt.Printf("Input v1: %d\n", root.compute(maxSize))

	var totalDiskSpace uint = 70_000_000
	var neededDiskSpace uint = 30_000_000
	var freeDiskSpace uint = totalDiskSpace - root.size()

	var sizes []uint = make([]uint, 0)
	root.collectSizes(&sizes)
	sort.Slice(sizes, func(i, j int) bool { return sizes[i] < sizes[j] })

	for _, currentSize := range sizes {
		if (freeDiskSpace + currentSize) >= neededDiskSpace {
			fmt.Printf("Input v2: %d\n", currentSize)
			return
		}
	}

	panic("Didn't find any suitable directory to delete")
}

func main() {
	if len(os.Args) != 2 {
		help(os.Args[0])
	}

	input(os.Args[1], 100_000)
}
