package util

import (
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var ldb *leveldb.DB = New()

func New() *leveldb.DB {
	d, err := leveldb.OpenFile("./dbstore", nil)
	if err != nil {
		log.Errorf("error while loading db: %v", err)
		return nil
	}
	return d
}

func DPut(collection string, key string, value string) {
	err := ldb.Put([]byte(collection+"/"+key), []byte(value), nil)
	if err != nil {
		log.Errorf("db put error: %v", err)
	}
}

func DGet(collection string, key string) string {
	result, err := ldb.Get([]byte(collection+"/"+key), nil)
	if err != nil {
		log.Errorf("db get error: %v", err)
		return ""
	}
	return string(result)
}

func DList(collection string) map[string]string {
	result := map[string]string{}

	iter := ldb.NewIterator(util.BytesPrefix([]byte(collection+"/")), nil)
	for iter.Next() {
		result[string(iter.Key())] = string(iter.Value())
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Errorf("db list error: %v", err)
		return nil
	}
	return result
}

func DDelete(collection string, key string) {
	err := ldb.Delete([]byte(collection+"/"+key), nil)
	if err != nil {
		log.Errorf("db delete error: %v", err)
	}
}
