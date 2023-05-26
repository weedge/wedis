package driver

import (
	"github.com/weedge/wedis/internal/srv/openkv"
)

type IKVStoreDB interface {
	Get(key []byte) ([]byte, error)
	GetSlice(key []byte) (openkv.Slice, error)

	Put(key []byte, value []byte) error
	Delete(key []byte) error

	NewIterator() *openkv.Iterator

	NewWriteBatch() *openkv.WriteBatch

	RangeIterator(min []byte, max []byte, rangeType uint8) *openkv.RangeLimitIterator
	RevRangeIterator(min []byte, max []byte, rangeType uint8) *openkv.RangeLimitIterator
	RangeLimitIterator(min []byte, max []byte, rangeType uint8, offset int, count int) *openkv.RangeLimitIterator
	RevRangeLimitIterator(min []byte, max []byte, rangeType uint8, offset int, count int) *openkv.RangeLimitIterator

	Compact() error
	Close() error
}
