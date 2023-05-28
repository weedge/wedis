package storager

import (
	"sync"

	"github.com/weedge/wedis/internal/srv/openkv"
)

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

func (db *DBString) delete(t *Batch, key []byte) (int64, error) {
	key = db.encodeStringKey(key)
	t.Delete(key)
	return 1, nil
}

func checkKeySize(key []byte) error {
	if len(key) > MaxKeySize || len(key) == 0 {
		return ErrKeySize
	}
	return nil
}

func checkValueSize(value []byte) error {
	if len(value) > MaxValueSize {
		return ErrValueSize
	}
	return nil
}

// Set sets the data.
func (db *DBString) Set(key []byte, value []byte) error {
	if err := checkKeySize(key); err != nil {
		return err
	}
	if err := checkValueSize(value); err != nil {
		return err
	}

	var err error
	key = db.encodeStringKey(key)

	db.batch.Lock()
	defer db.batch.Unlock()

	db.batch.Put(key, value)

	err = db.batch.Commit()

	return err
}

// GetSlice gets the slice of the data
func (db *DB) GetSlice(key []byte) (openkv.Slice, error) {
	if err := checkKeySize(key); err != nil {
		return nil, err
	}

	key = db.encodeStringKey(key)

	return db.IKVStoreDB.GetSlice(key)
}
