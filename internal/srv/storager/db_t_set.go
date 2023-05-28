package storager

import (
	"sync"

	"github.com/weedge/wedis/internal/srv/openkv"
)

type DBSet struct {
	*DB
	batch *Batch
}

func NewDBSet(db *DB) *DBSet {
	batch := NewBatch(db.store, db.IKVStoreDB.NewWriteBatch(),
		&dbBatchLocker{
			l:      &sync.Mutex{},
			wrLock: &db.store.wLock,
		})
	return &DBSet{DB: db, batch: batch}
}

func (db *DBSet) delete(t *Batch, key []byte) (num int64, err error) {
	sk := db.sEncodeSizeKey(key)
	start := db.sEncodeStartKey(key)
	stop := db.sEncodeStopKey(key)

	it := db.IKVStoreDB.RangeLimitIterator(start, stop, openkv.RangeROpen, 0, -1)
	for ; it.Valid(); it.Next() {
		t.Delete(it.RawKey())
		num++
	}

	it.Close()
	t.Delete(sk)

	return num, nil
}
