package storager

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/weedge/pkg/utils"
	"github.com/weedge/wedis/internal/srv/storager/driver"
)

// DB core sturct
// impl like redis string, list, hash, set, zset, bitmap struct store db op
type DB struct {
	store *Storager
	// database index
	index int
	// database index to varint buffer
	indexVarBuf []byte
	// IKVStoreDB impl
	driver.IKVStoreDB

	string *DBString
	list   *DBList
	hash   *DBHash
	set    *DBSet
	zset   *DBZSet
	bitmap *DBBitmap

	lbKeys *lBlockKeys

	ttlChecker *TTLChecker
}

func NewDB(store *Storager, idx int) *DB {
	db := &DB{store: store}
	db.SetIndex(idx)
	db.IKVStoreDB = store.odb

	db.string = NewDBString(db)
	db.list = NewDBList(db)
	db.hash = NewDBHash(db)
	db.set = NewDBSet(db)
	db.zset = NewDBZSet(db)
	db.bitmap = NewDBBitmap(db)

	db.lbKeys = newLBlockKeys()

	db.ttlChecker = NewTTLChecker(db)

	return db
}

func (m *DB) DBString() *DBString {
	return m.string
}
func (m *DB) DBList() *DBList {
	return m.list
}
func (m *DB) DBHash() *DBHash {
	return m.hash
}
func (m *DB) DBSet() *DBSet {
	return m.set
}
func (m *DB) DBZSet() *DBZSet {
	return m.zset
}
func (m *DB) DBBitmap() *DBBitmap {
	return m.bitmap
}

func (m *DB) Close() (err error) {
	if utils.IsNil(m.IKVStoreDB) {
		return
	}

	return m.IKVStoreDB.Close()
}

// Index gets the index of database.
func (db *DB) Index() int {
	return db.index
}

// IndexVarBuf gets the index varint buf of database.
func (db *DB) IndexVarBuf() []byte {
	return db.indexVarBuf
}

// SetIndex set the index of database.
func (db *DB) SetIndex(index int) {
	db.index = index
	// the most size for varint is 10 bytes
	buf := make([]byte, 10)
	n := binary.PutUvarint(buf, uint64(index))

	db.indexVarBuf = buf[0:n]
}

func (db *DB) checkKeyIndex(buf []byte) (int, error) {
	if len(buf) < len(db.indexVarBuf) {
		return 0, fmt.Errorf("key is too small")
	}
	if !bytes.Equal(db.indexVarBuf, buf[0:len(db.indexVarBuf)]) {
		return 0, fmt.Errorf("invalid db index")
	}

	return len(db.indexVarBuf), nil
}
