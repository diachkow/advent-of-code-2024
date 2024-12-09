package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BlockType uint8

const (
	TYPE_EMPTY BlockType = iota
	TYPE_FILE
)
const EMPTY_CHAR = '.'

type Block struct {
	Type BlockType
	ID   int // Only valid if Type is TYPE_FILE
}

type DiskMap []Block

func main() {
	dm := ReadDiskMap()
	dm.Print()

	fragmentized := Fragmentize(dm)
	fragmentized.Print()

	fmt.Printf("Checksum of fragmentized disk: %v\n", CalculateChecksum(fragmentized))
}

func ReadDiskMap() DiskMap {
	content, err := os.ReadFile("./disk.input.txt")
	if err != nil {
		panic(err)
	}
	rawDiskMap := strings.TrimSpace(string(content))

	dm := make(DiskMap, 0)
	blockType := TYPE_FILE
	idCounter := 0

	for _, char := range rawDiskMap {
		count := CharToInt(char)
		for i := 0; i < count; i++ {
			if blockType == TYPE_EMPTY {
				dm = append(dm, Block{Type: blockType})
			} else {
				dm = append(dm, Block{Type: blockType, ID: idCounter})
			}
		}

		if blockType == TYPE_FILE {
			idCounter++
		}

		blockType = (blockType + 1) % 2
	}

	return dm
}

func Fragmentize(dm DiskMap) DiskMap {
	dmCopy := make(DiskMap, len(dm))
	copy(dmCopy, dm)
	firstEmptyIdx := dmCopy.BlockIndex(TYPE_EMPTY)

	for i := len(dmCopy) - 1; i > 0; i-- {
		if dmCopy[i].Type == TYPE_EMPTY {
			continue
		}

		dmCopy[i], dmCopy[firstEmptyIdx] = dmCopy[firstEmptyIdx], dmCopy[i]
		if IsFragmentized(dmCopy) {
			return dmCopy
		}

		firstEmptyIdx = dmCopy.BlockIndex(TYPE_EMPTY)
	}

	return dmCopy
}

func IsFragmentized(dm DiskMap) bool {
	lastNumberIdx := dm.RBlockIndex(TYPE_FILE)
	firstEmptyIdx := dm.BlockIndex(TYPE_EMPTY)
	return firstEmptyIdx > lastNumberIdx
}

func CalculateChecksum(dm DiskMap) int {
	checksum := 0
	for i, block := range dm {
		if block.Type == TYPE_FILE {
			checksum += i * block.ID
		}
	}
	return checksum
}

func (dm DiskMap) Print() {
	for _, block := range dm {
		if block.Type == TYPE_EMPTY {
			fmt.Printf(string(EMPTY_CHAR))
		} else {
			fmt.Printf("%v", block.ID)
		}
	}
	fmt.Printf("\n")
}

func (dm DiskMap) BlockIndex(t BlockType) int {
	for i := 0; i < len(dm); i++ {
		if dm[i].Type == t {
			return i
		}
	}
	return -1
}

func (dm DiskMap) RBlockIndex(t BlockType) int {
	for i := len(dm) - 1; i > 0; i-- {
		if dm[i].Type == t {
			return i
		}
	}
	return -1
}

func CharToInt(c rune) int {
	num, err := strconv.Atoi(string(c))
	if err != nil {
		panic(err)
	}
	return num
}
