package badgerDB

import (
	dbpak "DatabaseDB/internal/Databaces"
	"bytes"

	"github.com/dgraph-io/badger/v4"
)

type badgerDatabase struct {
	db      *badger.DB
	Address string
}

func NewDataBaseBadger(address string) dbpak.DBClient {
	return &badgerDatabase{
		Address: address,
	}
}

func (b *badgerDatabase) Open() error {
	var err error
	b.db, err = badger.Open(badger.DefaultOptions(b.Address))
	return err
}

func (b *badgerDatabase) Add(key, value []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (b *badgerDatabase) Close() {
	b.db.Close()
}

func (b *badgerDatabase) Get(key []byte) ([]byte, error) {
	var valORG []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		valORG = val
		return nil
	})
	return valORG, err
}

func (b *badgerDatabase) Delete(key []byte) error {
	b.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (c *badgerDatabase) Read(start, end *[]byte, count int) (error, []dbpak.KVData) {
	var items []dbpak.KVData
	var opts badger.IteratorOptions
	opts.PrefetchSize = count

	if end != nil && start == nil {
		opts.Reverse = true
	}

	err := c.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opts)
		defer iter.Close()

		cnt := 0

		if end != nil && start == nil {
			iter.Seek(*end)
			iter.Next()
			item := iter.Item()
			key := item.Key()
			for iter.Seek(key); iter.Valid(); iter.Next() {
				cnt++
				if cnt > count {
					break
				}
				item := iter.Item()

				valCopy, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				key1 := make([]byte, len(item.Key()))
				copy(key1, item.Key())

				value1 := make([]byte, len(valCopy))
				copy(value1, valCopy)
				items = append(items, dbpak.KVData{Key: key1, Value: value1})
			}

			for i := 0; i < len(items)/2; i++ {
				j := len(items) - i - 1
				temp := items[i]
				items[i] = items[j]
				items[j] = temp
			}
		} else {

			if start != nil {
				iter.Seek(*start)
				iter.Next()
			} else {

				iter.Rewind()
			}

			for ; iter.Valid(); iter.Next() {
				cnt++
				if cnt > count {
					break
				}

				item := iter.Item()

				valCopy, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				key1 := make([]byte, len(item.Key()))
				copy(key1, item.Key())

				value1 := make([]byte, len(valCopy))
				copy(value1, valCopy)
				items = append(items, dbpak.KVData{Key: key1, Value: value1})
			}
		}
		return nil
	})
	if err != nil {
		return err, nil
	}

	return nil, items
}

func (l *badgerDatabase) Search(valueEntry []byte) (error, [][]byte) {
	var data [][]byte
	var opts badger.IteratorOptions

	err := l.db.View(func(txn *badger.Txn) error {
		Iterator := txn.NewIterator(opts)
		defer Iterator.Close()
		Iterator.Rewind()

		for Iterator.Valid() {

			if bytes.Contains(Iterator.Item().Key(), valueEntry) {

				key1 := make([]byte, len(Iterator.Item().Key()))
				copy(key1, Iterator.Item().Key())

				data = append(data, key1)

			}
			Iterator.Next()
		}
		return nil
	})
	if err != nil {
		return err, data
	}
	return nil, data
}
