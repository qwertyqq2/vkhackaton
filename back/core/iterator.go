package core

import (
	"fmt"

	"github.com/qwertyqq2/filebc/core/types"
)

type iterator struct {
	blocks types.Blocks

	index int
}

func newIterator(blocks types.Blocks) *iterator {
	return &iterator{
		blocks: blocks,
		index:  -1,
	}
}

func (iter *iterator) next() (*types.Block, error) {
	if iter.index+1 >= len(iter.blocks) {
		iter.index = len(iter.blocks)
		return nil, fmt.Errorf("end iterator")
	}
	iter.index++
	return iter.blocks[iter.index], nil
}

func (iter *iterator) prev() *types.Block {
	if iter.index < 1 {
		return nil
	}
	return iter.blocks[iter.index-1]
}

func (iter *iterator) current() *types.Block {
	if iter.index == -1 || iter.index >= len(iter.blocks) {
		return nil
	}
	return iter.blocks[iter.index]
}

func (iter *iterator) first() *types.Block {
	return iter.blocks[0]
}

func (iter *iterator) remaining() int {
	return len(iter.blocks) - iter.index
}

func (iter *iterator) back() {
	iter.index = -1
}

func (iter *iterator) processed() int {
	return iter.index + 1
}
