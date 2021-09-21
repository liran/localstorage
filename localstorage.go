package localstorage

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v3"
)

type LocalStorage struct {
	db *badger.DB
}

func NewLocalStorage(directory string) (*LocalStorage, error) {
	options := badger.DefaultOptions(directory)
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	return &LocalStorage{db: db}, nil
}

func (obj *LocalStorage) SetItem(key string, data interface{}) error {
	var value []byte
	switch v := data.(type) {
	case string: // Prevent repeated double quotes in the string
		value = []byte(v)
	default:
		bytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		value = bytes
	}

	return obj.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func (obj *LocalStorage) GetItem(key string) ([]byte, error) {
	var valCopy []byte
	err := obj.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err == badger.ErrKeyNotFound {
		err = nil
	}
	return valCopy, err
}

func (obj *LocalStorage) RemoveItem(key string) error {
	err := obj.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	if err == badger.ErrKeyNotFound {
		err = nil
	}
	return err
}

func (obj *LocalStorage) Exists(key string) bool {
	err := obj.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	return err == nil
}

func (obj *LocalStorage) Close() error {
	return obj.db.Close()
}
