package leveldbb

import (
	dbpak "DatabaseDB/internal/Databaces"
	"bytes"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LeveldbDatabase struct {
	DB      *leveldb.DB
	Address string
}

func NewDataBaseLeveldb(address string) dbpak.DBClient {
	return &LeveldbDatabase{
		Address: address,
	}
}

func (l *LeveldbDatabase) Delete(key []byte) error {
	err := l.DB.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *LeveldbDatabase) Open() error {
	var err error
	l.DB, err = leveldb.OpenFile(l.Address, nil)
	return err
}

func (l *LeveldbDatabase) Close() {
	l.DB.Close()
}

func (l *LeveldbDatabase) Add(key, value []byte) error {
	return l.DB.Put(key, value, nil)
}

func (l *LeveldbDatabase) Get(key []byte) ([]byte, error) {
	if l.DB == nil {
		return nil, nil
	}
	data, err := l.DB.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (c *LeveldbDatabase) Read(start, end *[]byte, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData

	readRange := &util.Range{}
	if start != nil {
		readRange.Start = *start
	}
	if end != nil {
		readRange.Limit = *end
	}
	iter := c.DB.NewIterator(readRange, nil)
	defer iter.Release()
	cnt := 0
	if end != nil && start == nil {
		iter.Last()

		key1 := make([]byte, len(iter.Key()))
		copy(key1, iter.Key())

		value1 := make([]byte, len(iter.Value()))
		copy(value1, iter.Value())
		Item = append(Item, dbpak.KVData{Key: key1, Value: value1})
		cnt++

		for iter.Prev() {
			cnt++
			if cnt > count {
				break
			}
			key1 := make([]byte, len(iter.Key()))
			copy(key1, iter.Key())

			value1 := make([]byte, len(iter.Value()))
			copy(value1, iter.Value())
			Item = append(Item, dbpak.KVData{Key: key1, Value: value1})
		}
		//reverse items
		for i := 0; i < len(Item)/2; i++ {
			j := len(Item) - i - 1
			temp := Item[i]
			Item[i] = Item[j]
			Item[j] = temp
		}
	} else {
		if start != nil {

			iter.Next()
		}
		for iter.Next() {
			cnt++
			if cnt > count {
				break
			}

			key1 := make([]byte, len(iter.Key()))
			copy(key1, iter.Key())

			value1 := make([]byte, len(iter.Value()))
			copy(value1, iter.Value())

			Item = append(Item, dbpak.KVData{Key: key1, Value: value1})
		}
	}

	return nil, Item
}

func (l *LeveldbDatabase) Search(valueEntry []byte) (error, [][]byte) {
	var data [][]byte

	Iterator := l.DB.NewIterator(nil, nil)
	if Iterator == nil {

		return fmt.Errorf("Iterator is nil"), data
	}

	if !Iterator.First() {
		return fmt.Errorf("iterator is empty"), data
	}

	for Iterator.Valid() {

		if bytes.Contains(Iterator.Key(), valueEntry) {

			key1 := make([]byte, len(Iterator.Key()))
			copy(key1, Iterator.Key())

			data = append(data, key1)

		}
		if !Iterator.Next() {
			break
		}
	}

	return nil, data
}
