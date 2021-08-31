package gospider

import "github.com/syndtr/goleveldb/leveldb"

//CreateDataDB 创建存储库
func CreateDataDB(path string) *leveldb.DB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return db
}
