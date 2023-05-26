package storager

func (db *DB) rmExpire(t *Batch, dataType byte, key []byte) (int64, error) {
	mk := db.expEncodeMetaKey(dataType, key)
	v, err := db.IKVStoreDB.Get(mk)
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	when, err2 := Int64(v, nil)
	if err2 != nil {
		return 0, err2
	}

	tk := db.expEncodeTimeKey(dataType, key, when)
	t.Delete(mk)
	t.Delete(tk)
	return 1, nil
}
