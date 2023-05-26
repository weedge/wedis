package openkv

import "bytes"

const (
	IteratorForward  uint8 = 0
	IteratorBackward uint8 = 1
)

const (
	RangeClose uint8 = 0x00
	RangeLOpen uint8 = 0x01
	RangeROpen uint8 = 0x10
	RangeOpen  uint8 = 0x11
)

// Range min must less or equal than max
//
// range type:
//
//	close: [min, max]
//	open: (min, max)
//	lopen: (min, max]
//	ropen: [min, max)
type Range struct {
	Min []byte
	Max []byte

	Type uint8
}

type Limit struct {
	Offset int
	Count  int
}

type RangeLimitIterator struct {
	Iterator

	r *Range
	l *Limit

	step int

	//0 for IteratorForward, 1 for IteratorBackward
	direction uint8
}

func (it *RangeLimitIterator) Valid() bool {
	if it.l.Offset < 0 {
		return false
	} else if !it.Valid() {
		return false
	} else if it.l.Count >= 0 && it.step >= it.l.Count {
		return false
	}

	if it.direction == IteratorForward {
		if it.r.Max != nil {
			r := bytes.Compare(it.RawKey(), it.r.Max)
			if it.r.Type&RangeROpen > 0 {
				return !(r >= 0)
			} else {
				return !(r > 0)
			}
		}
	} else {
		if it.r.Min != nil {
			r := bytes.Compare(it.RawKey(), it.r.Min)
			if it.r.Type&RangeLOpen > 0 {
				return !(r <= 0)
			} else {
				return !(r < 0)
			}
		}
	}

	return true
}

func (it *RangeLimitIterator) Next() {
	it.step++

	if it.direction == IteratorForward {
		it.Next()
	} else {
		it.Prev()
	}
}

func NewRangeLimitIterator(i *Iterator, r *Range, l *Limit) *RangeLimitIterator {
	return rangeLimitIterator(i, r, l, IteratorForward)
}

func NewRevRangeLimitIterator(i *Iterator, r *Range, l *Limit) *RangeLimitIterator {
	return rangeLimitIterator(i, r, l, IteratorBackward)
}

func NewRangeIterator(i *Iterator, r *Range) *RangeLimitIterator {
	return rangeLimitIterator(i, r, &Limit{0, -1}, IteratorForward)
}

func NewRevRangeIterator(i *Iterator, r *Range) *RangeLimitIterator {
	return rangeLimitIterator(i, r, &Limit{0, -1}, IteratorBackward)
}

func rangeLimitIterator(i *Iterator, r *Range, l *Limit, direction uint8) *RangeLimitIterator {
	it := new(RangeLimitIterator)

	it.r = r
	it.l = l
	it.direction = direction

	it.step = 0

	if l.Offset < 0 {
		return it
	}

	if direction == IteratorForward {
		if r.Min == nil {
			it.SeekToFirst()
		} else {
			it.Seek(r.Min)

			if r.Type&RangeLOpen > 0 {
				if it.Valid() && bytes.Equal(it.RawKey(), r.Min) {
					it.Next()
				}
			}
		}
	} else {
		if r.Max == nil {
			it.SeekToLast()
		} else {
			it.Seek(r.Max)

			if !it.Valid() {
				it.SeekToLast()
			} else {
				if !bytes.Equal(it.RawKey(), r.Max) {
					it.Prev()
				}
			}

			if r.Type&RangeROpen > 0 {
				if it.Valid() && bytes.Equal(it.RawKey(), r.Max) {
					it.Prev()
				}
			}
		}
	}

	for i := 0; i < l.Offset; i++ {
		if it.Valid() {
			if it.direction == IteratorForward {
				it.Next()
			} else {
				it.Prev()
			}
		}
	}

	return it
}
