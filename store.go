package gospider

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

//key-value存储数据
type Store interface {
	Add(key string, value string) error
	BatchAdd(m map[string]string) error
	Get(key string) (string, error)
	Del(key string) error
	List(prefix string, limit ...string) ([]string, error)
	Clear(prefix string, limit ...string)
}

//Store的默认实现
type LeveldbStore struct {
	path   string
	dataDB *leveldb.DB
}

func CreateLeveldbStore(path string) *LeveldbStore {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return &LeveldbStore{
		dataDB: db,
		path:   path,
	}
}

func (lvdb *LeveldbStore) Get(key string) (string, error) {
	reqByte, err := lvdb.dataDB.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(reqByte), nil
}

func (lvdb *LeveldbStore) Add(key string, value string) error {
	return lvdb.dataDB.Put([]byte(key), []byte(value), nil)
}

func (lvdb *LeveldbStore) BatchAdd(m map[string]string) error {
	batch := new(leveldb.Batch)
	for k, v := range m {
		batch.Put([]byte(k), []byte(v))
	}
	if batch.Len() > 0 {
		return lvdb.dataDB.Write(batch, nil)
	}
	return nil
}

func (lvdb *LeveldbStore) Del(key string) error {
	return lvdb.dataDB.Delete([]byte(key), nil)
}

func (lvdb *LeveldbStore) List(prefix string, limit ...string) ([]string, error) {
	var iter iterator.Iterator = nil
	if len(limit) > 0 {
		iter = lvdb.dataDB.NewIterator(&util.Range{Start: []byte(prefix), Limit: []byte(limit[0])}, nil)
	} else {
		iter = lvdb.dataDB.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	}
	defer iter.Release()
	listReq := make([]string, 0, 10)
	for iter.Next() {
		listReq = append(listReq, string(iter.Value()))
	}
	return listReq, nil
}

func (lvdb *LeveldbStore) Clear(prefix string, limit ...string) {
	var iter iterator.Iterator = nil
	if len(limit) > 0 {
		iter = lvdb.dataDB.NewIterator(&util.Range{Start: []byte(prefix), Limit: []byte(limit[0])}, nil)
	} else {
		iter = lvdb.dataDB.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	}
	defer iter.Release()
	for iter.Next() {
		lvdb.dataDB.Delete(iter.Key(), nil)
	}
}
