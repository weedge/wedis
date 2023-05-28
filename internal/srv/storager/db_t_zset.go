package storager

import (
	"sync"

	"github.com/weedge/wedis/internal/srv/openkv"
)

type DBZSet struct {
	*DB
	batch *Batch
}

func NewDBZSet(db *DB) *DBZSet {
	batch := NewBatch(db.store, db.IKVStoreDB.NewWriteBatch(),
		&dbBatchLocker{
			l:      &sync.Mutex{},
			wrLock: &db.store.wLock,
		})
	return &DBZSet{DB: db, batch: batch}
}

func (db *DBZSet) delete(t *Batch, key []byte) (num int64, err error) {
	num, err = db.zRemRange(t, key, MinScore, MaxScore, 0, -1)
	return
}

func (db *DB) zRemRange(t *Batch, key []byte, min int64, max int64, offset int, count int) (int64, error) {
	if len(key) > MaxKeySize {
		return 0, ErrKeySize
	}

	it := db.zIterator(key, min, max, offset, count, false)
	var num int64
	for ; it.Valid(); it.Next() {
		sk := it.RawKey()
		_, m, _, err := db.zDecodeScoreKey(sk)
		if err != nil {
			continue
		}

		if n, err := db.zDelItem(t, key, m, true); err != nil {
			return 0, err
		} else if n == 1 {
			num++
		}

		t.Delete(sk)
	}
	it.Close()

	if _, err := db.zIncrSize(t, key, -num); err != nil {
		return 0, err
	}

	return num, nil
}

func (db *DB) zIterator(key []byte, min int64, max int64, offset int, count int, reverse bool) *openkv.RangeLimitIterator {
	minKey := db.zEncodeStartScoreKey(key, min)
	maxKey := db.zEncodeStopScoreKey(key, max)

	if !reverse {
		return db.IKVStoreDB.RangeLimitIterator(minKey, maxKey, openkv.RangeClose, offset, count)
	}
	return db.IKVStoreDB.RevRangeLimitIterator(minKey, maxKey, openkv.RangeClose, offset, count)
}

func (db *DB) zDelItem(t *Batch, key []byte, member []byte, skipDelScore bool) (int64, error) {
	ek := db.zEncodeSetKey(key, member)
	if v, err := db.IKVStoreDB.Get(ek); err != nil {
		return 0, err
	} else if v == nil {
		//not exists
		return 0, nil
	} else {
		//exists
		if !skipDelScore {
			//we must del score
			s, err := Int64(v, err)
			if err != nil {
				return 0, err
			}
			sk := db.zEncodeScoreKey(key, member, s)
			t.Delete(sk)
		}
	}

	t.Delete(ek)

	return 1, nil
}

func (db *DB) zIncrSize(t *Batch, key []byte, delta int64) (int64, error) {
	sk := db.zEncodeSizeKey(key)

	size, err := Int64(db.IKVStoreDB.Get(sk))
	if err != nil {
		return 0, err
	}
	size += delta
	if size <= 0 {
		size = 0
		t.Delete(sk)
		db.rmExpire(t, ZSetType, key)
	} else {
		t.Put(sk, PutInt64(size))
	}

	return size, nil
}
