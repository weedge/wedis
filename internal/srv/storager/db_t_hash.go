package storager

import (
	"sync"

	"github.com/weedge/wedis/internal/srv/openkv"
)

type DBHash struct {
	*DB
	batch *Batch
}

func NewDBHash(db *DB) *DBHash {
	batch := NewBatch(db.store, db.IKVStoreDB.NewWriteBatch(),
		&dbBatchLocker{
			l:      &sync.Mutex{},
			wrLock: &db.store.wLock,
		})
	return &DBHash{DB: db, batch: batch}
}

func (db *DBHash) delete(t *Batch, key []byte) (num int64, err error) {
	sk := db.hEncodeSizeKey(key)
	start := db.hEncodeStartKey(key)
	stop := db.hEncodeStopKey(key)

	it := db.IKVStoreDB.RangeLimitIterator(start, stop, openkv.RangeROpen, 0, -1)
	for ; it.Valid(); it.Next() {
		t.Delete(it.Key())
		num++
	}
	it.Close()
	t.Delete(sk)

	return num, nil
}
