package PebbleDB

import (
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/dberr"
	"bytes"
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
)

type PebbleDatabase struct {
	DB      *pebble.DB
	Address string
}

func NewDataBasePebble(address string) dbpak.DBClient {
	return &PebbleDatabase{
		Address: address,
	}
}

func (p *PebbleDatabase) Delete(key []byte) error {
	err := p.DB.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *PebbleDatabase) Open() error {
	var err error
	p.DB, err = pebble.Open(p.Address, &pebble.Options{})
	return err
}

func (p *PebbleDatabase) Close() {
	p.DB.Close()
}

func (p *PebbleDatabase) Add(key, value []byte) error {
	return p.DB.Set(key, value, nil)
}

func (p *PebbleDatabase) Get(key []byte) ([]byte, error) {
	if p.DB == nil {
		return nil, dberr.ErrDBNil
	}

	value, closer, err := p.DB.Get([]byte(key))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, dberr.ErrKeyNotFound
		}
		return nil, err
	}
	defer closer.Close()

	data := make([]byte, len(value))
	copy(data, value)

	return data, nil
}

func (p *PebbleDatabase) Read(start, end *[]byte, count int) ([]dbpak.KVData, error) {
	var Item []dbpak.KVData

	iterOptions := &pebble.IterOptions{}
	if start != nil {
		iterOptions.LowerBound = *start
	}
	if end != nil {

		iterOptions.UpperBound = *end
	}

	iter, err := p.DB.NewIter(iterOptions)
	if err != nil {
		return Item, err
	}
	defer iter.Close()

	cnt := 0
	if end != nil && start == nil {
		iter.Last()

		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())

		value := make([]byte, len(iter.Value()))
		copy(value, iter.Value())

		Item = append(Item, dbpak.KVData{Key: key, Value: value})
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
		}

		for i := 0; i < len(Item)/2; i++ {
			j := len(Item) - i - 1
			temp := Item[i]
			Item[i] = Item[j]
			Item[j] = temp
		}
	} else {
		if start != nil {
			iter.SeekGE(*start)
			iter.Next()
		} else {
			iter.First()
		}

		for iter.Valid() {
			cnt++
			if cnt > count {
				break
			}
			key := make([]byte, len(iter.Key()))
			copy(key, iter.Key())

			value := make([]byte, len(iter.Value()))
			copy(value, iter.Value())

			Item = append(Item, dbpak.KVData{Key: key, Value: value})
			iter.Next()
		}
	}

	return Item, nil
}

func (l *PebbleDatabase) Search(valueEntry []byte) ([][]byte, error) {
	var data [][]byte

	Iterator, err := l.DB.NewIter(nil)
	if err != nil {
		return data, err
	}

	defer Iterator.Close()
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
