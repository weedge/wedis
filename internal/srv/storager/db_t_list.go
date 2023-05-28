package storager

import (
	"encoding/binary"
	"sync"

	"github.com/weedge/wedis/internal/srv/openkv"
)

type DBList struct {
	*DB
	batch *Batch
}

func NewDBList(db *DB) *DBList {
	batch := NewBatch(db.store, db.IKVStoreDB.NewWriteBatch(),
		&dbBatchLocker{
			l:      &sync.Mutex{},
			wrLock: &db.store.wLock,
		})
	return &DBList{DB: db, batch: batch}
}

func (db *DBList) delete(t *Batch, key []byte) (num int64, err error) {
	it := db.IKVStoreDB.NewIterator()
	defer it.Close()

	mk := db.lEncodeMetaKey(key)
	headSeq, tailSeq, _, err := db.lGetMeta(it, mk)
	if err != nil {
		return
	}

	startKey := db.lEncodeListKey(key, headSeq)
	stopKey := db.lEncodeListKey(key, tailSeq)

	rit := openkv.NewRangeIterator(it, &openkv.Range{
		Min:  startKey,
		Max:  stopKey,
		Type: openkv.RangeClose})
	for ; rit.Valid(); rit.Next() {
		t.Delete(rit.RawKey())
		num++
	}

	t.Delete(mk)

	return num, nil
}

func (db *DBList) lGetMeta(it *openkv.Iterator, ek []byte) (headSeq int32, tailSeq int32, size int32, err error) {
	var v []byte
	if it != nil {
		v = it.Find(ek)
	} else {
		v, err = db.IKVStoreDB.Get(ek)
	}
	if err != nil {
		return
	} else if v == nil {
		headSeq = listInitialSeq
		tailSeq = listInitialSeq
		size = 0
		return
	} else {
		headSeq = int32(binary.LittleEndian.Uint32(v[0:4]))
		tailSeq = int32(binary.LittleEndian.Uint32(v[4:8]))
		size = tailSeq - headSeq + 1
	}
	return
}
