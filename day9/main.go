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
	fmt.Println("Part 1:File system checksum:", filesystemChecksum)
}

func part2(disk Disk) {
	fileSystemCompacted := make([]string, len(disk.Filesystem))
	copy(fileSystemCompacted, disk.Filesystem)
	// remaining free spaces
	freeSpaces := make([]int, len(disk.SpacePositions))
	copy(freeSpaces, disk.SpacePositions)

	for fileId := len(disk.FileIds) - 1; fileId >= 0; fileId-- {
		filePositions := disk.FilePositions[strconv.Itoa(disk.FileIds[fileId])]
		nextFreeSpaces := foundFreeSpaces(freeSpaces, filePositions)
		if nextFreeSpaces == nil {
			continue
		}

		// iterate over file positions to compact file system with free spaces
		if nextFreeSpaces[0] < filePositions[0] && len(filePositions) <= len(nextFreeSpaces) {
			for filePosI, filePosition := range filePositions {
				fileSystemCompacted[nextFreeSpaces[filePosI]] = strconv.Itoa(disk.FileIds[fileId])
				fileSystemCompacted[filePosition] = "."
				freeSpaces = removeIntFromList(freeSpaces, nextFreeSpaces[filePosI])
				freeSpaces = append(freeSpaces, filePosition)
				sort.Ints(freeSpaces)
			}
		}
	}

	filesystemChecksum := calcFilesystemChecksum(fileSystemCompacted)
	fmt.Println("Part 2: File system checksum:", filesystemChecksum)
}

func calcFilesystemChecksum(filesystem []string) int {
	sum := 0
	for i, fileId := range filesystem {
		fileIdInt, _ := strconv.Atoi(fileId)
		sum += i * fileIdInt
	}
	return sum
}

// iterate over consecutive free spaces to found free spaces for file
func foundFreeSpaces(freeSpaces []int, filePositions []int) []int {
	nextFreeSpaces := []int{}
	tempFreeSpaces := make([]int, len(freeSpaces))
	copy(tempFreeSpaces, freeSpaces)

	// iterate over consecutive free spaces to found free spaces for file
	for len(nextFreeSpaces) < len(filePositions) {
		nextFreeSpaces = []int{}
		nextFreeSpaceIndex := 0

		for len(nextFreeSpaces) == 0 || nextFreeSpaces[nextFreeSpaceIndex-1] == tempFreeSpaces[nextFreeSpaceIndex]-1 {
			nextFreeSpaces = append(nextFreeSpaces, tempFreeSpaces[nextFreeSpaceIndex])
			nextFreeSpaceIndex++
			if nextFreeSpaceIndex >= len(tempFreeSpaces) {
				break
			}
		}
		tempFreeSpaces = tempFreeSpaces[nextFreeSpaceIndex:]
		nextFreeSpaceIndex = 0
		if len(tempFreeSpaces) == 0 {
			break
		}
	}
	if len(nextFreeSpaces) < len(filePositions) {
		return nil
	}
	return nextFreeSpaces
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

func removeIntFromList(list []int, value int) []int {
	for i, v := range list {
		if v == value {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}
