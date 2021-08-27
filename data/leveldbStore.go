package data

import (
	"encoding/json"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LeveldbStore struct {
	dataDB *leveldb.DB
}

func CreateLeveldbStore(path string) *LeveldbStore {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &LeveldbStore{
		dataDB: db,
	}
}

func (lvdb *LeveldbStore) Get(key string) (string, error) {
	reqByte, err := lvdb.dataDB.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return string(reqByte), nil
}

func (lvdb *LeveldbStore) Add(key string, value string) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return lvdb.dataDB.Put([]byte(key), b, nil)
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
