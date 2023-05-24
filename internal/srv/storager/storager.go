package storager

// core store struct
type DB struct {
}

func New() *DB {
	return &DB{}
}

func (m *DB) Close() (err error) {
	return
}
