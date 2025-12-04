package dbpak

type KVData struct {
	Key   []byte
	Value []byte
}

type DBClient interface {
	Open() error
	Close()
	Add(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Read(start, end *[]byte, count int) (error, []KVData)
	Delete(key []byte) error
	Search(value []byte) (error, [][]byte)
}
