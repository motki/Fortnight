package localstore // import "github.com/motki/fortnight/localstore"

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/coreos/bbolt"
	"github.com/pkg/errors"

	"fmt"

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

type Key []byte

func IntKey(id int) Key {
	return Key(strconv.Itoa(id))
}

func IntPairKey(id1, id2 int) Key {
	return Key(fmt.Sprintf("%d,%d", id1, id2))
}

type Value interface{}

type Store struct {
	db *bolt.DB

	tx *bolt.Tx
}

func New(dataDir string) (*Store, error) {
	db, err := bolt.Open(dataDir, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	return &Store{db, nil}, nil
}

type BucketOption func(b *Bucket) error

type prototype func() Value

func WithPrototype(p func() Value) BucketOption {
	return func(b *Bucket) error {
		b.prototype = p
		return nil
	}
}

func (s *Store) RemoveBucket(kind Kind) error {
	_, err := s.Acquire(kind)
	if err != nil {
		return err
	}
	// Use the tx that we know was created during Acquire.
	return s.tx.DeleteBucket(kind.bucket())
}

func (s *Store) Acquire(kind Kind, opts ...BucketOption) (*Bucket, error) {
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
	bucket := &Bucket{Bucket: b, kind: kind, prototype: kind.identity}
	for _, o := range opts {
		if err := o(bucket); err != nil {
			return nil, err
		}
	}
	return bucket, nil
}

type withCallback func(*Store) error

func (s *Store) With(fn withCallback) error {
	if err := fn(s); err != nil {
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

func (bkt *Bucket) Kind() Kind {
	return bkt.kind
}

func (bkt *Bucket) All() ([]Value, error) {
	cur := bkt.Cursor()
	var res []Value
	for k, v := cur.First(); k != nil; k, v = cur.Next() {
		val := bkt.prototype()
		if err := json.Unmarshal(v, val); err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}

func (bkt *Bucket) Put(k Key, data Value) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return bkt.Bucket.Put([]byte(k), b)
}

var ErrNotFound = errors.New("not found")

func (bkt *Bucket) Get(k Key) (Value, error) {
	v := bkt.Bucket.Get([]byte(k))
	if v == nil {
		return nil, ErrNotFound
	}
	res := bkt.prototype()
	if err := json.Unmarshal(v, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (bkt *Bucket) Delete(k Key) error {
	return bkt.Bucket.Delete([]byte(k))
}
