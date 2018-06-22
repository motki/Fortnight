package localstore // import "github.com/motki/fortnight/localstore"

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/coreos/bbolt"
	"github.com/pkg/errors"

	"bytes"

	"github.com/motki/core/evedb"
	"github.com/motki/core/model"
)

type Kind int

const (
	KindUnknown Kind = iota
	KindLocation
	KindItemType
	KindInventoryItem
)

func (s Kind) String() string {
	switch s {
	case KindLocation:
		return "location"
	case KindItemType:
		return "item_type"
	case KindInventoryItem:
		return "inventory_item"
	default:
		return "unknown"
	}
}

func (s Kind) bucket() []byte {
	return []byte(s.String())
}

func (s Kind) identity() Value {
	switch s {
	case KindLocation:
		return &model.Location{}
	case KindItemType:
		return &evedb.ItemType{}
	case KindInventoryItem:
		return &model.InventoryItem{}
	default:
		return nil
	}
}

type Value interface{}

type Store struct {
	db *bolt.DB

	tx *bolt.Tx
}

func New(dataDir string) (*Store, error) {
	db, err := bolt.Open(dataDir, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	return &Store{db, nil}, nil
}

func (s *Store) Acquire(kind Kind) (*Bucket, error) {
	var err error
	if s.tx == nil {
		if s.tx, err = s.db.Begin(true); err != nil {
			return nil, err
		}
	}
	b, err := s.tx.CreateBucketIfNotExists(kind.bucket())
	if err != nil {
		return nil, err
	}
	return &Bucket{Bucket: b, kind: kind, prototype: kind.identity}, nil
}

type withCallback func(*Store) error

func (s *Store) With(fn withCallback) (err error) {
	if err = fn(s); err != nil {
		return errors.Wrap(s.Rollback(), err.Error())
	}
	return s.Commit()
}

func (s *Store) Commit() error {
	if s.tx == nil {
		return errors.New("no tx")
	}
	tx := s.tx
	s.tx = nil
	return tx.Commit()
}

func (s *Store) Rollback() error {
	if s.tx == nil {
		return errors.New("no tx")
	}
	tx := s.tx
	s.tx = nil
	return tx.Rollback()
}

type Bucket struct {
	*bolt.Bucket

	kind Kind

	prototype prototype
}

type prototype func() Value

func (bkt *Bucket) SetPrototype(p func() Value) {
	bkt.prototype = p
}

func (bkt *Bucket) Kind() Kind {
	return bkt.kind
}

func (bkt *Bucket) key(id int) []byte {
	return append(append([]byte(bkt.Kind().String()), ':'), []byte(strconv.Itoa(id))...)
}

func (bkt *Bucket) All() ([]Value, error) {
	cur := bkt.Cursor()
	prefix := []byte(bkt.kind.String())
	var res []Value
	for k, v := cur.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cur.Next() {
		val := bkt.prototype()
		if err := json.Unmarshal(v, val); err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}

func (bkt *Bucket) Put(id int, data Value) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return bkt.Bucket.Put(bkt.key(id), b)
}

var ErrNotFound = errors.New("not found")

func (bkt *Bucket) Get(id int) (Value, error) {
	v := bkt.Bucket.Get(bkt.key(id))
	if v == nil {
		return nil, ErrNotFound
	}
	res := bkt.prototype()
	if err := json.Unmarshal(v, res); err != nil {
		return nil, err
	}
	return res, nil
}
