package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Disk struct {
	Filesystem     []string
	FilePositions  map[string][]int
	FileIds        []int
	SpacePositions []int
}

func (disk *Disk) FilesCount() int {
	return len(disk.FilePositions)
}

func main() {
	disk, err := parseDiskMap("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing disk map:", err)
		os.Exit(1)
	}
	part1(*disk)
	part2(*disk)
}

func part1(disk Disk) {
	freeSpaces := make([]int, len(disk.SpacePositions))
	copy(freeSpaces, disk.SpacePositions)
	fileSystemCompacted := make([]string, len(disk.Filesystem))
	copy(fileSystemCompacted, disk.Filesystem)

	for fileId := len(disk.FileIds) - 1; fileId >= 0; fileId-- {
		filePositions := disk.FilePositions[strconv.Itoa(disk.FileIds[fileId])]

		// iterate over file positions to compact file system with free spaces
		for _, filePosition := range filePositions {
			nextFreeSpace := freeSpaces[0]
			if nextFreeSpace < filePosition {
				fileSystemCompacted[nextFreeSpace] = strconv.Itoa(disk.FileIds[fileId])
				fileSystemCompacted[filePosition] = "."
				freeSpaces = freeSpaces[1:]
				freeSpaces = append(freeSpaces, filePosition)
				sort.Ints(freeSpaces)
			}
		}
	}

	filesystemChecksum := calcFilesystemChecksum(fileSystemCompacted)
	fmt.Println("File system checksum:", filesystemChecksum)
}

func part2(disk Disk) {}

func calcFilesystemChecksum(filesystem []string) int {
	sum := 0
	for i, fileId := range filesystem {
		fileIdInt, _ := strconv.Atoi(fileId)
		sum += i * fileIdInt
	}
	return sum
}

func parseDiskMap(filename string) (*Disk, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	filesystem := []string{}
	fileids := []int{}
	filepositions := map[string][]int{}
	spacepositions := []int{}
	nextFileId := 0
	nextFileSystemIndex := 0
	if scanner.Scan() {
		diskMap := scanner.Text()
		for iChar, char := range diskMap {
			noOfFreeSpaceOrLengthOfFile, err := parseNumber(string(char))
			if err != nil {
				return nil, err
			}

			for nextFileSystemWrite := 0; nextFileSystemWrite < noOfFreeSpaceOrLengthOfFile; nextFileSystemWrite++ {
				if iChar%2 != 0 {
					// free space
					filesystem = append(filesystem, ".")
					spacepositions = append(spacepositions, nextFileSystemIndex)

				} else {
					// file
					filesystem = append(filesystem, strconv.Itoa(nextFileId))
					filepositions[strconv.Itoa(nextFileId)] = append(filepositions[strconv.Itoa(nextFileId)], nextFileSystemIndex)

					if nextFileSystemWrite == noOfFreeSpaceOrLengthOfFile-1 {
						fileids = append(fileids, nextFileId)
						nextFileId++
					}
				}
				nextFileSystemIndex += 1
			}
		}
	}
	return &Disk{Filesystem: filesystem, FilePositions: filepositions, SpacePositions: spacepositions, FileIds: fileids}, nil
}

func parseNumber(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("could not parse %q as number: %w", s, err)
	}
	return i, nil
}
