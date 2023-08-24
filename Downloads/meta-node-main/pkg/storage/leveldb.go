package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	db     *leveldb.DB
	closed bool
	path   string
}

func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{db, false, path}, nil
}

func (ldb *LevelDB) Get(key []byte) ([]byte, error) {
	return ldb.db.Get(key, nil)
}

func (ldb *LevelDB) Put(key, value []byte) error {
	return ldb.db.Put(key, value, nil)
}

func (ldb *LevelDB) Has(key []byte) bool {
	has, _ := ldb.db.Has(key, nil)
	return has
}

func (ldb *LevelDB) Delete(key []byte) error {
	return ldb.db.Delete(key, nil)
}

func (ldb *LevelDB) BatchPut(kvs [][2][]byte) error {
	batch := new(leveldb.Batch)
	for i := range kvs {
		batch.Put(kvs[i][0], kvs[i][1])
	}
	return ldb.db.Write(batch, nil)
}

func (ldb *LevelDB) Open() error {
	var err error
	if ldb.closed {
		ldb.db, err = leveldb.OpenFile(ldb.path, nil)
		if err != nil {
			return err
		}
		ldb.closed = false
	}
	return nil
}

func (ldb *LevelDB) Close() error {
	if !ldb.closed {
		err := ldb.db.Close()
		if err != nil {
			return err
		}
		ldb.closed = true
	}
	return nil
}

func (ldb *LevelDB) GetIterator() IIterator {
	return ldb.db.NewIterator(nil, nil)
}

func (ldb *LevelDB) CopyToNewPath(newPath string) (IStorage, error) {
	if ldb.closed {
		ldb.Open()
		defer ldb.Close()
	}
	newLDB, err := NewLevelDB(newPath)
	if err != nil {
		return nil, err
	}
	batch := new(leveldb.Batch)
	iter := ldb.GetIterator()
	for iter.Next() {
		if iter.Error() != nil {
			return nil, iter.Error()
		}
		// need to copy because iter done save value in next call Next
		cValue := make([]byte, len(iter.Value()))
		cKey := make([]byte, len(iter.Key()))
		copy(cValue, iter.Value())
		copy(cKey, iter.Key())
		batch.Put(cKey, cValue)
	}
	newLDB.db.Write(batch, nil)

	return newLDB, nil
}
