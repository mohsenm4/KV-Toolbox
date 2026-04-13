package leveldbb

import (
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/dberr"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
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
	opts := &opt.Options{
		BlockCacheCapacity: 0,
		WriteBuffer:        0,
	}
	l.DB, err = leveldb.OpenFile(l.Address, opts)
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
		return nil, dberr.ErrDBNil
	}
	data, err := l.DB.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, dberr.ErrKeyNotFound
		}
		return nil, err
	}
	return data, err
}

func (c *LeveldbDatabase) Read(start, end *[]byte, count int) ([]dbpak.KVData, error) {
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

		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())

		value := make([]byte, len(iter.Value()))
		copy(value, iter.Value())
		Item = append(Item, dbpak.KVData{Key: key, Value: value})
		value = nil
		cnt++

		for iter.Prev() {
			cnt++
			if cnt > count {
				break
			}
			key := make([]byte, len(iter.Key()))
			copy(key, iter.Key())

			value := make([]byte, len(iter.Value()))
			copy(value, iter.Value())
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
			value = nil
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

			key := make([]byte, len(iter.Key()))
			copy(key, iter.Key())

			value := make([]byte, len(iter.Value()))
			copy(value, iter.Value())

			Item = append(Item, dbpak.KVData{Key: key, Value: value})
		}
	}

	return Item, nil
}

func (l *LeveldbDatabase) Search(valueEntry []byte) ([][]byte, error) {
	var data [][]byte

	Iterator := l.DB.NewIterator(nil, nil)
	defer Iterator.Release()

	if Iterator == nil {
		return data, fmt.Errorf("iterator is nil")
	}

	if !Iterator.First() {
		return data, fmt.Errorf("iterator is empty")
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

	return data, nil
}

func FormatKeyValue(item dbpak.KVData) (string, string) {
	truncatedKey := TruncateString(string(item.Key), 20)

	typeValue := mimetype.Detect(item.Value)
	var truncatedValue string
	if typeValue.Extension() != ".txt" {
		truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
	} else {
		truncatedValue = TruncateString(string(item.Value), 30)
	}

	return truncatedKey, truncatedValue
}

func TruncateString(input string, length int) string {
	nameData := input
	if len(nameData) > length {
		nameData = nameData[:length] + ". . ."
	}
	parts := strings.Split(nameData, "\n")
	if len(parts) > 1 {

		nameData = parts[0] + " . . ."
	}

	return nameData
}
