package main

import (
	"fmt"
	"iter"
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
	Len  int
	ID   int // Only valid if Type is TYPE_FILE
}
type BlockItem struct {
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
		dm = append(dm, Block{Type: blockType, Len: count, ID: idCounter})

		if blockType == TYPE_FILE {
			idCounter++
		}

		blockType = (blockType + 1) % 2
	}

	return dm
}

func Fragmentize(dm DiskMap) DiskMap {
	res := make(DiskMap, len(dm))
	copy(res, dm)

	for i := len(res) - 1; i > 0; i-- {
		replaceBlock := res[i]

		if replaceBlock.Type != TYPE_FILE {
			continue
		}

		for j := 0; j < i; j++ {
			if replacedBlock := res[j]; replacedBlock.Type == TYPE_EMPTY && replacedBlock.Len >= replaceBlock.Len {
				leftover := replacedBlock.Len - replaceBlock.Len

				res[j] = replaceBlock
				res[i] = Block{Type: TYPE_EMPTY, Len: replaceBlock.Len}

				if leftover > 0 {
					leftoverIdx := j + 1
					res = append(res[:leftoverIdx+1], res[leftoverIdx:]...)
					res[leftoverIdx] = Block{Type: TYPE_EMPTY, Len: leftover}
				}

				break
			}
		}
	}

	return res
}

func CalculateChecksum(dm DiskMap) int {
	checksum := 0
	for i, item := range dm.Flatten() {
		if item.Type == TYPE_FILE {
			checksum += i * item.ID
		}
	}
	return checksum
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

func (dm DiskMap) Flatten() iter.Seq2[int, BlockItem] {
	return func(yield func(int, BlockItem) bool) {
		cnt := 0
		for _, block := range dm {
			bi := BlockItem{Type: block.Type, ID: block.ID}
			for i := 0; i < block.Len; i++ {
				if !yield(cnt, bi) {
					return
				}
				cnt++
			}
		}
	}
}

func (dm DiskMap) Print() {
	for _, item := range dm.Flatten() {
		if item.Type == TYPE_EMPTY {
			fmt.Printf(string(EMPTY_CHAR))
		} else {
			fmt.Printf("%v", item.ID)
		}
	}
	fmt.Printf("\n")
}

func CharToInt(c rune) int {
	num, err := strconv.Atoi(string(c))
	if err != nil {
		panic(err)
	}
	return num
}
