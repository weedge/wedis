package storager

import "sync"

type DBString struct {
	*DB
	batch *Batch
}

func NewDBString(db *DB) *DBString {
	batch := NewBatch(db.store, db.IKVStoreDB.NewWriteBatch(),
		&dbBatchLocker{
			l:      &sync.Mutex{},
			wrLock: &db.store.wLock,
		})
	return &DBString{DB: db, batch: batch}
}

func (db *DBString) delete(t *Batch, key []byte) int64 {
	key = db.encodeStringKey(key)
	t.Delete(key)
	return 1
}
